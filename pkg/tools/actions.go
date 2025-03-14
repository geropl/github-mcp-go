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

	// Register get_workflow tool
	getWorkflowTool := mcp.NewTool("get_workflow",
		mcp.WithDescription("Get detailed information about a specific workflow"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner (username or organization)"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("workflow_id",
			mcp.Required(),
			mcp.Description("The ID or filename of the workflow (can be a numeric ID or a filename)"),
		),
	)

	s.RegisterTool(getWorkflowTool, true, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters
		owner, ok := request.Params.Arguments["owner"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("owner must be a string"))), nil
		}

		repo, ok := request.Params.Arguments["repo"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("repo must be a string"))), nil
		}

		workflowID, ok := request.Params.Arguments["workflow_id"]
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("workflow_id is required"))), nil
		}

		// Call the operation
		workflow, err := actionsOps.GetWorkflow(ctx, owner, repo, workflowID)
		if err != nil {
			if ghErr, ok := err.(*errors.GitHubError); ok {
				return mcp.NewToolResultError(errors.FormatGitHubError(ghErr)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Error getting workflow: %v", err)), nil
		}

		// Format the result as markdown
		markdown := formatWorkflowToMarkdown(workflow)
		return mcp.NewToolResultText(markdown), nil
	})

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
	
	// Register list_workflow_runs tool
	listWorkflowRunsTool := mcp.NewTool("list_workflow_runs",
		mcp.WithDescription("Lists workflow runs for a repository or a specific workflow"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner (username or organization)"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("workflow_id",
			mcp.Description("The ID or filename of the workflow to filter runs by"),
		),
		mcp.WithString("branch",
			mcp.Description("Filter by branch name"),
		),
		mcp.WithString("status",
			mcp.Description("Filter by workflow run status (completed, action_required, cancelled, failure, neutral, skipped, stale, success, timed_out, in_progress, queued, requested, waiting)"),
		),
		mcp.WithString("event",
			mcp.Description("Filter by event type (push, pull_request, etc.)"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (default: 1)"),
		),
		mcp.WithNumber("perPage",
			mcp.Description("Number of results per page (default: 30, max: 100)"),
		),
	)
	
	// Register get_workflow_run tool
	getWorkflowRunTool := mcp.NewTool("get_workflow_run",
		mcp.WithDescription("Gets detailed information about a specific workflow run"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner (username or organization)"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("run_id",
			mcp.Required(),
			mcp.Description("The ID of the workflow run"),
		),
	)
	
	// Register download_workflow_run_logs tool
	downloadWorkflowRunLogsTool := mcp.NewTool("download_workflow_run_logs",
		mcp.WithDescription("Downloads and extracts logs for a workflow run"),
		mcp.WithString("owner",
			mcp.Required(),
			mcp.Description("Repository owner (username or organization)"),
		),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Repository name"),
		),
		mcp.WithString("run_id",
			mcp.Required(),
			mcp.Description("The ID of the workflow run"),
		),
	)
	
	s.RegisterTool(getWorkflowRunTool, true, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters
		owner, ok := request.Params.Arguments["owner"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("owner must be a string"))), nil
		}

		repo, ok := request.Params.Arguments["repo"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("repo must be a string"))), nil
		}

		runID, ok := request.Params.Arguments["run_id"]
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("run_id is required"))), nil
		}

		// Call the operation
		run, err := actionsOps.GetWorkflowRun(ctx, owner, repo, runID)
		if err != nil {
			if ghErr, ok := err.(*errors.GitHubError); ok {
				return mcp.NewToolResultError(errors.FormatGitHubError(ghErr)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Error getting workflow run: %v", err)), nil
		}

		// Format the result as markdown
		markdown := formatWorkflowRunToMarkdown(run)
		return mcp.NewToolResultText(markdown), nil
	})

	// Register the download_workflow_run_logs tool
	s.RegisterTool(downloadWorkflowRunLogsTool, true, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract parameters
		owner, ok := request.Params.Arguments["owner"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("owner must be a string"))), nil
		}

		repo, ok := request.Params.Arguments["repo"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("repo must be a string"))), nil
		}

		runID, ok := request.Params.Arguments["run_id"]
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("run_id is required"))), nil
		}

		// Call the operation
		result, err := actionsOps.DownloadWorkflowRunLogs(ctx, owner, repo, runID)
		if err != nil {
			if ghErr, ok := err.(*errors.GitHubError); ok {
				return mcp.NewToolResultError(errors.FormatGitHubError(ghErr)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Error downloading workflow run logs: %v", err)), nil
		}

		// Format the result as markdown
		markdown := formatLogsResultToMarkdown(result)
		return mcp.NewToolResultText(markdown), nil
	})

	s.RegisterTool(listWorkflowRunsTool, true, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract required parameters
		owner, ok := request.Params.Arguments["owner"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("owner must be a string"))), nil
		}

		repo, ok := request.Params.Arguments["repo"].(string)
		if !ok {
			return mcp.NewToolResultError(errors.FormatGitHubError(errors.NewInvalidArgumentError("repo must be a string"))), nil
		}

		// Extract optional parameters
		var workflowID interface{}
		if val, ok := request.Params.Arguments["workflow_id"]; ok {
			workflowID = val
		}

		branch := ""
		if val, ok := request.Params.Arguments["branch"].(string); ok {
			branch = val
		}

		status := ""
		if val, ok := request.Params.Arguments["status"].(string); ok {
			status = val
		}

		event := ""
		if val, ok := request.Params.Arguments["event"].(string); ok {
			event = val
		}

		// Extract pagination parameters with defaults
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
		runs, err := actionsOps.ListWorkflowRuns(ctx, owner, repo, workflowID, branch, status, event, page, perPage)
		if err != nil {
			if ghErr, ok := err.(*errors.GitHubError); ok {
				return mcp.NewToolResultError(errors.FormatGitHubError(ghErr)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Error listing workflow runs: %v", err)), nil
		}

		// Format the result as markdown
		markdown := formatWorkflowRunsToMarkdown(runs)
		return mcp.NewToolResultText(markdown), nil
	})
}

// formatWorkflowToMarkdown converts a GitHub workflow to markdown
func formatWorkflowToMarkdown(workflow *gh.Workflow) string {
	md := fmt.Sprintf("# Workflow: %s\n\n", workflow.GetName())
	
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
	
	// Add badge URL if available
	badgeURL := workflow.GetBadgeURL()
	if badgeURL != "" {
		md += fmt.Sprintf("**Badge URL:** %s  \n\n", badgeURL)
	}
	
	return md
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

// formatWorkflowRunToMarkdown converts a single GitHub workflow run to markdown
func formatWorkflowRunToMarkdown(run *gh.WorkflowRun) string {
	md := fmt.Sprintf("# Workflow Run: %s\n\n", run.GetName())
	
	md += fmt.Sprintf("**ID:** %d  \n", run.GetID())
	md += fmt.Sprintf("**Run Number:** %d  \n", run.GetRunNumber())
	md += fmt.Sprintf("**Workflow ID:** %d  \n", run.GetWorkflowID())
	md += fmt.Sprintf("**Branch:** %s  \n", run.GetHeadBranch())
	md += fmt.Sprintf("**Commit SHA:** %s  \n", run.GetHeadSHA())
	md += fmt.Sprintf("**Event:** %s  \n", run.GetEvent())
	md += fmt.Sprintf("**Status:** %s  \n", run.GetStatus())
	
	// Add conclusion if available
	conclusion := run.GetConclusion()
	if conclusion != "" {
		md += fmt.Sprintf("**Conclusion:** %s  \n", conclusion)
	}
	
	// Format dates if available
	createdAt := run.GetCreatedAt()
	if !createdAt.IsZero() {
		md += fmt.Sprintf("**Created:** %s  \n", createdAt.Format(time.RFC1123))
	}

	updatedAt := run.GetUpdatedAt()
	if !updatedAt.IsZero() {
		md += fmt.Sprintf("**Updated:** %s  \n", updatedAt.Format(time.RFC1123))
	}
	
	startedAt := run.GetRunStartedAt()
	if !startedAt.IsZero() {
		md += fmt.Sprintf("**Started:** %s  \n", startedAt.Format(time.RFC1123))
	}
	
	// Add actor information if available
	actor := run.GetActor()
	if actor != nil {
		md += fmt.Sprintf("**Triggered by:** %s  \n", actor.GetLogin())
	}
	
	// Add triggering actor if available and different from actor
	triggeringActor := run.GetTriggeringActor()
	if triggeringActor != nil && (actor == nil || triggeringActor.GetLogin() != actor.GetLogin()) {
		md += fmt.Sprintf("**Triggering Actor:** %s  \n", triggeringActor.GetLogin())
	}
	
	md += fmt.Sprintf("**Run Attempt:** %d  \n", run.GetRunAttempt())
	
	// Add URL
	md += fmt.Sprintf("**URL:** %s  \n\n", run.GetHTMLURL())
	
	return md
}

// formatWorkflowRunsToMarkdown converts GitHub workflow runs to markdown
func formatWorkflowRunsToMarkdown(runs *gh.WorkflowRuns) string {
	md := "# Workflow Runs\n\n"

	if runs.GetTotalCount() == 0 {
		md += "No workflow runs found.\n"
		return md
	}

	md += fmt.Sprintf("Found %d workflow runs.\n\n", runs.GetTotalCount())

	for i, run := range runs.WorkflowRuns {
		md += fmt.Sprintf("## %d. Run #%d: %s\n\n", i+1, run.GetRunNumber(), run.GetName())
		md += fmt.Sprintf("**ID:** %d  \n", run.GetID())
		md += fmt.Sprintf("**Workflow ID:** %d  \n", run.GetWorkflowID())
		md += fmt.Sprintf("**Branch:** %s  \n", run.GetHeadBranch())
		md += fmt.Sprintf("**Commit SHA:** %s  \n", run.GetHeadSHA())
		md += fmt.Sprintf("**Event:** %s  \n", run.GetEvent())
		md += fmt.Sprintf("**Status:** %s  \n", run.GetStatus())
		
		// Add conclusion if available
		conclusion := run.GetConclusion()
		if conclusion != "" {
			md += fmt.Sprintf("**Conclusion:** %s  \n", conclusion)
		}
		
		// Format dates if available
		createdAt := run.GetCreatedAt()
		if !createdAt.IsZero() {
			md += fmt.Sprintf("**Created:** %s  \n", createdAt.Format(time.RFC1123))
		}

		updatedAt := run.GetUpdatedAt()
		if !updatedAt.IsZero() {
			md += fmt.Sprintf("**Updated:** %s  \n", updatedAt.Format(time.RFC1123))
		}
		
		startedAt := run.GetRunStartedAt()
		if !startedAt.IsZero() {
			md += fmt.Sprintf("**Started:** %s  \n", startedAt.Format(time.RFC1123))
		}
		
		// Add actor information if available
		actor := run.GetActor()
		if actor != nil {
			md += fmt.Sprintf("**Triggered by:** %s  \n", actor.GetLogin())
		}
		
		md += fmt.Sprintf("**Run Attempt:** %d  \n", run.GetRunAttempt())
		md += fmt.Sprintf("**URL:** %s  \n\n", run.GetHTMLURL())
	}

	return md
}

// formatLogsResultToMarkdown converts a logs result to markdown
func formatLogsResultToMarkdown(result *github.LogsResult) string {
	md := fmt.Sprintf("# Workflow Run Logs: %s\n\n", result.WorkflowName)
	
	md += fmt.Sprintf("**Run ID:** %d  \n", result.RunID)
	
	// Use a fixed timestamp in test environment
	downloadTime := result.DownloadTime
	md += fmt.Sprintf("**Downloaded:** %s  \n", downloadTime.Format(time.RFC1123))
	
	md += fmt.Sprintf("**Size:** %d bytes  \n", result.Size)
	md += fmt.Sprintf("**Files:** %d  \n\n", result.FileCount)
	
	// Use a placeholder for the logs directory in test environment
	logsDir := result.LogsDir
	
	md += fmt.Sprintf("**Logs Directory:** %s  \n\n", logsDir)
	
	// List the extracted files (limit to 20 to avoid excessive output)
	if len(result.Files) > 0 {
		md += "## Extracted Log Files\n\n"
		
		maxFiles := 20
		if len(result.Files) <= maxFiles {
			for _, file := range result.Files {
				md += fmt.Sprintf("- %s\n", file)
			}
		} else {
			for i := 0; i < maxFiles; i++ {
				md += fmt.Sprintf("- %s\n", result.Files[i])
			}
			md += fmt.Sprintf("\n... and %d more files\n", len(result.Files)-maxFiles)
		}
	}
	
	md += "\n## Usage Instructions\n\n"
	md += "The logs have been extracted to the *Logs Directory* mentioned above. You can access the log files directly to view their contents.\n\n"
	md += "**Note:** The logs directory is temporary and will be cleaned up when the server restarts.\n"
	
	return md
}
