package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/geropl/github-mcp-go/pkg/github"
	"github.com/geropl/github-mcp-go/pkg/tools"
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
	s := tools.NewServer(serverName, serverVersion, githubClient, logger)

	// Register tools
	logger.Info("Registering tools...")
	tools.RegisterTools(s)
	logger.Info("Tools registered successfully")

	// Start the stdio server
	logger.Info("Starting GitHub MCP server...")
	if err := s.Serve(); err != nil {
		logger.WithError(err).Fatal("Server error")
	}
}
