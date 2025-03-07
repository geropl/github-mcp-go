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
}
