package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"

	"github.com/modelcontextprotocol/github-mcp-go/internal/errors"
	"github.com/modelcontextprotocol/github-mcp-go/internal/github"
	"github.com/modelcontextprotocol/github-mcp-go/internal/server"
	"github.com/modelcontextprotocol/github-mcp-go/pkg/tools"
)

const (
	serverName    = "github-mcp-server"
	serverVersion = "0.1.0"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Check for GitHub token
	token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	if token == "" {
		logger.Fatal("GITHUB_PERSONAL_ACCESS_TOKEN environment variable is required")
	}

	// Create GitHub client
	githubClient := github.NewClient(token, logger)

	// Create MCP server
	s := server.NewServer(serverName, serverVersion, githubClient, logger)

	// Register tools
	registerTools(s)

	// Start the stdio server
	logger.Info("Starting GitHub MCP server...")
	if err := s.Serve(); err != nil {
		logger.WithError(err).Fatal("Server error")
	}
}

func registerTools(s *server.Server) {
	logger := s.GetLogger()
	logger.Info("Registering tools...")

	// Register repository tools
	tools.RegisterRepositoryTools(s)

	// Register pull request tools
	tools.RegisterPullRequestTools(s)

	// Register file tools
	tools.RegisterFileTools(s)

	// Example tool registration
	helloTool := mcp.NewTool("hello_github",
		mcp.WithDescription("Say hello to GitHub"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	// Add tool handler
	s.RegisterTool(helloTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("name must be a string"))), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Hello, %s! Welcome to the GitHub MCP server.", name)), nil
	})

	logger.Info("Tools registered successfully")
}
