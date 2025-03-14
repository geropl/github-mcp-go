package github

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v69/github"
	"github.com/sirupsen/logrus"

	"github.com/geropl/github-mcp-go/pkg/errors"
)

// ActionsOperations handles GitHub Actions-related operations
type ActionsOperations struct {
	client *Client
	logger *logrus.Logger
}

// NewActionsOperations creates a new ActionsOperations
func NewActionsOperations(client *Client, logger *logrus.Logger) *ActionsOperations {
	return &ActionsOperations{
		client: client,
		logger: logger,
	}
}

// ListWorkflows lists all workflows in a repository
func (a *ActionsOperations) ListWorkflows(ctx context.Context, owner, repo string, page, perPage int) (*github.Workflows, error) {
	// Validate owner and repo
	if owner == "" {
		return nil, errors.NewValidationError("owner cannot be empty")
	}
	if repo == "" {
		return nil, errors.NewValidationError("repository name cannot be empty")
	}

	// Set default pagination values if not provided
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 30
	} else if perPage > 100 {
		perPage = 100
	}

	// Create list options for pagination
	opts := &github.ListOptions{
		Page:    page,
		PerPage: perPage,
	}

	// Call the GitHub API
	workflows, _, err := a.client.GetClient().Actions.ListWorkflows(ctx, owner, repo, opts)
	if err != nil {
		return nil, a.client.HandleError(err)
	}

	return workflows, nil
}

// GetWorkflow gets detailed information about a specific workflow
func (a *ActionsOperations) GetWorkflow(ctx context.Context, owner, repo string, workflowID interface{}) (*github.Workflow, error) {
	// Validate owner and repo
	if owner == "" {
		return nil, errors.NewValidationError("owner cannot be empty")
	}
	if repo == "" {
		return nil, errors.NewValidationError("repository name cannot be empty")
	}
	
	// Handle different types of workflowID (can be int64 or string)
	var id int64
	var name string
	
	switch v := workflowID.(type) {
	case int64:
		id = v
	case float64:
		id = int64(v)
	case int:
		id = int64(v)
	case string:
		if v == "" {
			return nil, errors.NewValidationError("workflow_id cannot be empty")
		}
		
		// Try to convert string to int64 if it looks like a number
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			id = i
		} else {
			name = v
		}
	default:
		return nil, errors.NewValidationError(fmt.Sprintf("workflow_id must be a string or number, got %T", workflowID))
	}
	
	// Call the GitHub API
	var workflow *github.Workflow
	var err error
	
	if id != 0 {
		workflow, _, err = a.client.GetClient().Actions.GetWorkflowByID(ctx, owner, repo, id)
	} else {
		workflow, _, err = a.client.GetClient().Actions.GetWorkflowByFileName(ctx, owner, repo, name)
	}
	
	if err != nil {
		return nil, a.client.HandleError(err)
	}
	
	return workflow, nil
}

// ListWorkflowRuns lists workflow runs for a repository or a specific workflow
func (a *ActionsOperations) ListWorkflowRuns(
	ctx context.Context,
	owner, repo string,
	workflowID interface{},
	branch, status, event string,
	page, perPage int,
) (*github.WorkflowRuns, error) {
	// Validate owner and repo
	if owner == "" {
		return nil, errors.NewValidationError("owner cannot be empty")
	}
	if repo == "" {
		return nil, errors.NewValidationError("repository name cannot be empty")
	}

	// Set default pagination values if not provided
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 30
	} else if perPage > 100 {
		perPage = 100
	}

	// Create options for the API call
	opts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}

	// Add optional filters if provided
	if branch != "" {
		opts.Branch = branch
	}
	if status != "" {
		opts.Status = status
	}
	if event != "" {
		opts.Event = event
	}

	// Call the GitHub API
	var runs *github.WorkflowRuns
	var err error

	// If workflowID is provided, list runs for that specific workflow
	if workflowID != nil {
		// Handle different types of workflowID (can be int64 or string)
		var id int64
		var name string

		switch v := workflowID.(type) {
		case int64:
			id = v
		case float64:
			id = int64(v)
		case int:
			id = int64(v)
		case string:
			if v == "" {
				// Empty string means no workflow filter, list all runs
				runs, _, err = a.client.GetClient().Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
				if err != nil {
					return nil, a.client.HandleError(err)
				}
				return runs, nil
			}

			// Try to convert string to int64 if it looks like a number
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				id = i
			} else {
				name = v
			}
		default:
			return nil, errors.NewValidationError(fmt.Sprintf("workflow_id must be a string or number, got %T", workflowID))
		}

		// Call the appropriate API method based on the workflow identifier
		if id != 0 {
			runs, _, err = a.client.GetClient().Actions.ListWorkflowRunsByID(ctx, owner, repo, id, opts)
		} else if name != "" {
			runs, _, err = a.client.GetClient().Actions.ListWorkflowRunsByFileName(ctx, owner, repo, name, opts)
		}
	} else {
		// No workflow ID provided, list all runs for the repository
		runs, _, err = a.client.GetClient().Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
	}

	if err != nil {
		return nil, a.client.HandleError(err)
	}

	return runs, nil
}

// GetWorkflowRun gets detailed information about a specific workflow run
func (a *ActionsOperations) GetWorkflowRun(ctx context.Context, owner, repo string, runID interface{}) (*github.WorkflowRun, error) {
	// Validate owner and repo
	if owner == "" {
		return nil, errors.NewValidationError("owner cannot be empty")
	}
	if repo == "" {
		return nil, errors.NewValidationError("repository name cannot be empty")
	}
	
	// Handle different types of runID (can be int64 or string)
	var id int64
	
	switch v := runID.(type) {
	case int64:
		id = v
	case float64:
		id = int64(v)
	case int:
		id = int64(v)
	case string:
		if v == "" {
			return nil, errors.NewValidationError("run_id cannot be empty")
		}
		
		// Try to convert string to int64
		var err error
		id, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, errors.NewValidationError(fmt.Sprintf("run_id must be a valid number, got %s", v))
		}
	default:
		return nil, errors.NewValidationError(fmt.Sprintf("run_id must be a string or number, got %T", runID))
	}
	
	// Call the GitHub API
	run, _, err := a.client.GetClient().Actions.GetWorkflowRunByID(ctx, owner, repo, id)
	if err != nil {
		return nil, a.client.HandleError(err)
	}
	
	return run, nil
}
