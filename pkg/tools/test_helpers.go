package tools

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	"github.com/geropl/github-mcp-go/pkg/github"
)

var (
	// UpdateGolden is a flag to update golden files
	golden = flag.Bool("golden", false, "update golden files")
	// RecordMode is a flag to control VCR recording
	record = flag.Bool("record", false, "record new HTTP interactions")
)

type TestCase struct {
	Name  string
	Tool  string
	Input map[string]interface{}
}

// TestResult represents the expected result of a test
type TestResult struct {
	Output string `json:"output"`
	Err    string `json:"err"`
}

func RunTest(t *testing.T, tc TestCase) {
	s := createTestServer(t, *record)

	testCtx := context.Background()
	actual, testErr := executeTestTool(testCtx, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		rawMessage, err := json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %v", err)
		}
		rawResponse := s.server.HandleMessage(ctx, rawMessage)
		jsonRpcResponse, ok := rawResponse.(mcp.JSONRPCResponse);
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T", rawResponse)
		}
		callToolResult, ok := jsonRpcResponse.Result.(mcp.CallToolResult)
		if !ok {
			return nil, fmt.Errorf("unexpected result type: %T", jsonRpcResponse.Result)
		}
		return &callToolResult, nil
	}, tc.Tool, tc.Input)
	if testErr != nil {
		t.Fatalf("Failed to execute test tool: %v", testErr)
	}

	goldenFile := filepath.Join("testdata", "golden", t.Name()+".golden")

	// If -golden flag is set, update the golden file
	if *golden {
		writeErr := writeGoldenFile(goldenFile, actual)
		if writeErr != nil {
			t.Fatalf("Failed to write golden file: %v", writeErr)
		}
		return
	}

	expected, readErr := readGoldenFile(goldenFile)
	if readErr != nil {
		t.Fatalf("Failed to read golden file: %v", readErr)
	}

	// Compare actual and expected results
	if actual.Output != expected.Output {
		t.Errorf("Output mismatch:\nExpected: %s\nActual: %s", expected.Output, actual.Output)
	}
	if actual.Err != expected.Err {
		t.Errorf("Error mismatch:\nExpected: %s\nActual: %s", expected.Err, actual.Err)
	}
}

// createTestServer creates a test server with a VCR recorder
func createTestServer(t *testing.T, doRecord bool) *Server {
	logger := logrus.New()

	options := []recorder.Option{
		// filter "Authorization" header
		recorder.WithHook(func(i *cassette.Interaction) error {
			i.Request.Headers.Del("Authorization")
			return nil
		}, recorder.AfterCaptureHook),
		recorder.WithMatcher(cassette.NewDefaultMatcher(cassette.WithIgnoreAuthorization())),
	}
	if doRecord {
		options = append(options, recorder.WithMode(recorder.ModeRecordOnly))
	} else {
		options = append(options, recorder.WithMode(recorder.ModeReplayOnly))
	}

	cassetteName := path.Join("testdata", "cassette", t.Name()+".yaml")
	r, err := recorder.New(cassetteName, options...)
	if err != nil {
		t.Fatalf("Failed to create recorder: %v", err)
	}

	// Create a GitHub client with the recorder's transport
	httpClient := &http.Client{Transport: r}

	// Only require token in record mode
	token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	githubClient := github.NewClientWithHTTPClient(token, httpClient, logger)

	// Create a server
	s := NewServer("test-server", "0.1.0", githubClient, logger)
	RegisterTools(s)
	return s
}

// executeTestTool executes a test tool with the given input and returns the result
func executeTestTool(ctx context.Context, handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error), toolName string, input map[string]interface{}) (*TestResult, error) {
	// Create a request
	request := mcp.CallToolRequest{}
	request.Params.Name = toolName
	request.Params.Arguments = input

	// Execute the tool handler directly
	result, err := handler(ctx, request)
	if err != nil {
		return &TestResult{
			Output: "",
			Err:    err.Error(),
		}, nil
	}

	// Process the result
	if result.IsError {
		// For error results, the content is the error message
		if len(result.Content) > 0 {
			// Try to extract the text content
			content := result.Content[0]
			if textContent, ok := content.(mcp.TextContent); ok {
				return &TestResult{
					Output: "",
					Err:    textContent.Text,
				}, nil
			}
			// If we can't extract the text, return a generic error
			return &TestResult{
				Output: "",
				Err:    "Unknown error",
			}, nil
		}
		return &TestResult{
			Output: "",
			Err:    "Empty error content",
		}, nil
	}

	// For success results, the content is the output
	if len(result.Content) > 0 {
		// Try to extract the text content
		content := result.Content[0]
		if textContent, ok := content.(mcp.TextContent); ok {
			return &TestResult{
				Output: textContent.Text,
				Err:    "",
			}, nil
		}
		// If we can't extract the text, return a generic output
		return &TestResult{
			Output: "Non-text content",
			Err:    "",
		}, nil
	}
	return &TestResult{
		Output: "Empty content",
		Err:    "",
	}, nil
}

func writeGoldenFile(goldenFile string, actual *TestResult) error {
	if err := os.MkdirAll(filepath.Dir(goldenFile), 0755); err != nil {
		return fmt.Errorf("failed to create golden directory: %v", err)
	}
	data, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal actual result: %v", err)
	}
	if err := os.WriteFile(goldenFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write golden file: %v", err)
	}
	return nil
}

func readGoldenFile(goldenFile string) (*TestResult, error) {
	// Read expected result from golden file
	data, err := os.ReadFile(goldenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read golden file: %v", err)
	}

	var expected TestResult
	if err := json.Unmarshal(data, &expected); err != nil {
		return nil, fmt.Errorf("failed to unmarshal expected result: %v", err)
	}
	return &expected, nil
}
