package github

import (
	"context"

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
