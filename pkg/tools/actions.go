package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/geropl/github-mcp-go/pkg/errors"
	"github.com/geropl/github-mcp-go/pkg/github"
	gh "github.com/google/go-github/v69/github"
)

// RegisterActionsTools registers GitHub Actions-related tools
func RegisterActionsTools(s *Server) {
	client := s.GetClient()
	logger := s.GetLogger()
	actionsOps := github.NewActionsOperations(client, logger)

	// Register list_workflows tool
	listWorkflowsTool := mcp.NewTool("list_workflows",
		mcp.WithDescription("List all workflows in a repository"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner (username or organization)"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default: 1)"),
		),
		mcp.WithNumber("perPage",
			mcp.Description("Number of results per page (default: 30, max: 100)"),
		),
	)

	s.RegisterTool(listWorkflowsTool, true, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters
		owner, ok := request.Params.Arguments["owner"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("owner must be a string"))), nil
		}

		repo, ok := request.Params.Arguments["repo"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("repo must be a string"))), nil
		}

		// Extract optional parameters with defaults
		page := 1
		if pageVal, ok := request.Params.Arguments["page"]; ok {
			if pageFloat, ok := pageVal.(float64); ok {
				page = int(pageFloat)
			}
		}

		perPage := 30
		if perPageVal, ok := request.Params.Arguments["perPage"]; ok {
			if perPageFloat, ok := perPageVal.(float64); ok {
				perPage = int(perPageFloat)
			}
		}

		// Call the operation
		workflows, err := actionsOps.ListWorkflows(ctx, owner, repo, page, perPage)
		if err != nil {
			if ghErr, ok := err.(*errors.GitHubError); ok {
				return mcp.NewToolResultError(errors.FormatGitHubError(ghErr)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Error listing workflows: %v", err)), nil
		}

		// Format the result as markdown
		markdown := formatWorkflowsToMarkdown(workflows)
		return mcp.NewToolResultText(markdown), nil
	})
}

// formatWorkflowsToMarkdown converts GitHub workflows to markdown
func formatWorkflowsToMarkdown(workflows *gh.Workflows) string {
	md := "# Workflows\n\n"

	if workflows.GetTotalCount() == 0 {
		md += "No workflows found.\n"
		return md
	}

	md += fmt.Sprintf("Found %d workflows.\n\n", workflows.GetTotalCount())

	for i, workflow := range workflows.Workflows {
		md += fmt.Sprintf("## %d. %s\n\n", i+1, workflow.GetName())
		md += fmt.Sprintf("**ID:** %d  \n", workflow.GetID())
		md += fmt.Sprintf("**Path:** %s  \n", workflow.GetPath())
		md += fmt.Sprintf("**State:** %s  \n", workflow.GetState())

		// Format dates if available
		createdAt := workflow.GetCreatedAt()
		if !createdAt.IsZero() {
			md += fmt.Sprintf("**Created:** %s  \n", createdAt.Format(time.RFC1123))
		}

		updatedAt := workflow.GetUpdatedAt()
		if !updatedAt.IsZero() {
			md += fmt.Sprintf("**Updated:** %s  \n", updatedAt.Format(time.RFC1123))
		}

		md += fmt.Sprintf("**URL:** %s  \n\n", workflow.GetHTMLURL())
	}

	return md
}
