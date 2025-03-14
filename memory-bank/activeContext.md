# Active Context

## Current Focus

We are implementing GitHub Actions tools following PRD 001-action-tools, with current progress at 43%. The focus is on read-only tools for GitHub Actions workflows.

Current priorities:
1. **GitHub Actions tools implementation** (43% complete)
   - Implemented: list_workflows, get_workflow, list_workflow_runs
   - Next: get_workflow_run, download_workflow_run_logs, list_workflow_jobs, get_workflow_job

2. **Testing completion**
   - Repository operations tests need completion (50% done)
   - Actions operations tests in progress (43% done)
   - All other test categories are complete

3. **End-to-end testing preparation**

## Recent Changes

- Added ListWorkflowRuns operation and list_workflow_runs tool with tests
- Updated PRD 001-action-tools to reflect implementation progress (43%)
- Added formatter for workflow runs in markdown format

## Next Steps

1. Continue implementing GitHub Actions tools per PRD 001-action-tools
2. Complete repository operations tests for create_repository and fork_repository
3. Plan for end-to-end testing

## Active Decisions

1. **Feature Development Process**
   - Using PRDs to document requirements and track implementation progress
   - PRDs are numbered sequentially and stored in the `prds/` directory
   - Implementation status is tracked within each PRD

2. **Error Handling Strategy**
   - Using logrus for structured logging
   - Custom error types for different error categories
   - Detailed error messages for troubleshooting

3. **Authentication**
   - GitHub personal access token via environment variable
   - Validation for token presence

4. **Testing Approach**
   - Iterative: one test case at a time, starting with "happy path"
   - Using go-vcr for recording HTTP interactions
   - Golden files for expected results
   - Comprehensive test cases for all tools

5. **Response Formatting**
   - Markdown formatting for all tool responses
   - Formatters for each GitHub resource type
   - Clear, readable format with headers and sections

## Current Challenges

1. Ensuring proper error handling between GitHub API and MCP errors
2. Implementing proper validation for tool inputs
3. Managing test fixtures and cassettes for deterministic testing
4. Handling GitHub API permissions for tests requiring write access
5. Ensuring tests pass in normal mode, not just recording/golden mode
