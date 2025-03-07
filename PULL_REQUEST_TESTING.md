# Pull Request Testing Guide

This document provides instructions for implementing and running the pull request test cases in an iterative manner.

## Test Case Overview

The `pulls_test.go` file contains test cases for three pull request tools:

1. `create_pull_request` - Create a new pull request
2. `get_pull_request` - Get detailed information about a pull request
3. `get_pull_request_diff` - Get the diff of a pull request

Each tool has both happy path and error test cases, all of which are commented out except for the initial "SuccessfulCreation" test case.

## Iterative Testing Approach

Following the project's testing patterns, we'll implement one test case at a time, making sure each test case works completely before moving to the next. The general workflow is:

1. Uncomment one test case in `pulls_test.go`
2. Run the test with the `-record` flag to record HTTP interactions
3. Run the test with the `-golden` flag to create golden files
4. Run the test without flags to verify it passes
5. Move to the next test case

## Prerequisites

Before running the tests, ensure you have:

1. A GitHub personal access token with appropriate permissions set as an environment variable:
   ```bash
   export GITHUB_PERSONAL_ACCESS_TOKEN=your_token_here
   ```

2. Access to the test repository (`geropl/github-mcp-go-test`)

## Testing Workflow

### Step 1: Verify the existing test case

First, make sure the existing "SuccessfulCreation" test case works:

```bash
go test -v ./pkg/tools -run TestPullRequest/SuccessfulCreation
```

### Step 2: Implement the next test case

For each test case, follow these steps:

1. Uncomment the test case in `pulls_test.go`
2. Run the test with the `-record` flag to record HTTP interactions:
   ```bash
   go test -v ./pkg/tools -run TestPullRequest/[TestCaseName] -record
   ```
3. Run the test with the `-golden` flag to create golden files:
   ```bash
   go test -v ./pkg/tools -run TestPullRequest/[TestCaseName] -golden
   ```
4. Run the test without flags to verify it passes:
   ```bash
   go test -v ./pkg/tools -run TestPullRequest/[TestCaseName]
   ```

### Step 3: Repeat for each test case

Follow the implementation order below, implementing one test case at a time:

## Implementation Order

### create_pull_request - Happy Path

1. ✅ SuccessfulCreation (already implemented)
2. CreateDraftPR
3. CreatePRWithLabels
4. CreatePRWithAssignees
5. CreatePRWithReviewers

### create_pull_request - Error Cases

6. InvalidOwner
7. InvalidRepo
8. InvalidBranch
9. SameBranches
10. MissingRequiredFields

### get_pull_request - Happy Path

11. GetExistingPR
12. GetMergedPR
13. GetClosedPR

### get_pull_request - Error Cases

14. GetNonExistentPR
15. InvalidOwnerGetPR
16. InvalidRepoGetPR
17. InvalidPRNumber

### get_pull_request_diff - Happy Path

18. GetDiffForOpenPR
19. GetDiffForMergedPR
20. GetDiffForClosedPR
21. GetDiffWithLargeChanges

### get_pull_request_diff - Error Cases

22. GetDiffForNonExistentPR
23. InvalidOwnerGetDiff
24. InvalidRepoGetDiff
25. InvalidPRNumberGetDiff

## Example: Implementing the "CreateDraftPR" Test Case

1. Uncomment the "CreateDraftPR" test case in `pulls_test.go`:
   ```go
   {
       Name: "CreateDraftPR",
       Tool: "create_pull_request",
       Input: map[string]interface{}{
           "owner": OWNER,
           "repo":  REPO,
           "title": "Draft PR",
           "body":  "This is a draft PR",
           "head":  uniqueBranch + "-draft",
           "base":  "main",
           "draft": true,
       },
       Before: func(ctx context.Context, client *ghclient.Client) error {
           return createBranch(ctx, client, OWNER, REPO, uniqueBranch + "-draft", "main")
       },
       After: func(ctx context.Context, client *ghclient.Client) error {
           return deleteBranch(ctx, client, OWNER, REPO, uniqueBranch + "-draft")
       },
   },
   ```

2. Record HTTP interactions:
   ```bash
   go test -v ./pkg/tools -run TestPullRequest/CreateDraftPR -record
   ```

3. Create golden files:
   ```bash
   go test -v ./pkg/tools -run TestPullRequest/CreateDraftPR -golden
   ```

4. Verify the test passes:
   ```bash
   go test -v ./pkg/tools -run TestPullRequest/CreateDraftPR
   ```

## Notes on Test Cases

### PR Numbers

The test cases for `get_pull_request` and `get_pull_request_diff` assume certain PR numbers exist in the test repository:

- `PR_NUMBER` (1): An open PR
- `PR_NUMBER + 1` (2): A merged PR
- `PR_NUMBER + 2` (3): A closed PR
- `PR_NUMBER + 3` (4): A PR with large changes

You may need to adjust these numbers based on the actual PRs in the test repository.

### Error Cases

For error cases, no `Before` or `After` functions are needed since we're testing error conditions that don't require setup or cleanup.

### Test Repository State

Some tests may require specific conditions in the test repository. Ensure the repository has:

1. An open PR
2. A merged PR
3. A closed PR
4. A PR with large changes

## Running All Tests

After implementing all test cases, you can run all tests with:

```bash
go test -v ./pkg/tools -run TestPullRequest
```

This will run all the test cases in sequence.
