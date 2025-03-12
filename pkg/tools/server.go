package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"

	ghclient "github.com/geropl/github-mcp-go/pkg/github"
)

// Server wraps the MCP server and provides additional functionality
type Server struct {
	server *server.MCPServer
	client *ghclient.Client
	logger *logrus.Logger
}

// NewServer creates a new MCP server
func NewServer(name, version string, client *ghclient.Client, logger *logrus.Logger) *Server {
	s := server.NewMCPServer(
		name,
		version,
	)

	return &Server{
		server: s,
		client: client,
		logger: logger,
	}
}

// GetClient returns the GitHub client
func (s *Server) GetClient() *ghclient.Client {
	return s.client
}

// GetLogger returns the logger
func (s *Server) GetLogger() *logrus.Logger {
	return s.logger
}

// RegisterTool registers a tool with the server
func (s *Server) RegisterTool(tool mcp.Tool, handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	s.server.AddTool(tool, handler)
}

// Serve starts the server using stdio
func (s *Server) Serve() error {
	return server.ServeStdio(s.server)
}

func RegisterTools(s *Server) {
	RegisterRepositoryTools(s)
	RegisterPullRequestTools(s)
	RegisterFileTools(s)
	RegisterIssueTools(s)
	RegisterCommitTools(s)
	RegisterBranchTools(s)
	RegisterSearchTools(s)
}

// GetReadOnlyToolNames returns a map of tool names that are read-only
// These tools do not modify any state and are safe to auto-approve
func GetReadOnlyToolNames() map[string]bool {
	return map[string]bool{
		"search_repositories":   true,
		"search_code":           true,
		"search_issues":         true,
		"search_commits":        true,
		"get_file_contents":     true,
		"get_issue":             true,
		"list_issues":           true,
		"list_issue_comments":   true,
		"get_pull_request":      true,
		"get_pull_request_diff": true,
		"get_commit":            true,
		"list_commits":          true,
		"compare_commits":       true,
		"get_commit_status":     true,
		"list_commit_comments":  true,
		"list_branches":         true,
		"get_branch":            true,
	}
}
