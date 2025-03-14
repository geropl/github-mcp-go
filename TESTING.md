# Testing Guide

## Overview

This document provides a comprehensive guide to the testing approach used in the GitHub MCP Go project. The tests are designed as integration tests over the tools offered, where we use go-vcr for recording HTTP interactions, to re-run the tests like unit tests.
All tests are specified as table-driven tests, and use golden files for persisting expectations.

## Testing Architecture

### Components

- **TestCase struct**: Defines test parameters including name, tool, input, and setup/cleanup functions
- **go-vcr**: Records and replays HTTP interactions with the GitHub API
- **Golden files**: Store expected test results for comparison

### Test Structure

Tests are implemented using a table-driven approach:

```go
func TestPullRequest(t *testing.T) {
    testCases := []TestCase{
        {
            Name: "SuccessfulCreation",
            Tool: "create_pull_request",
            Input: map[string]interface{}{
                // Test parameters
            },
            Before: func(ctx context.Context, client *ghclient.Client) error {
                // Setup code
            },
            After: func(ctx context.Context, client *ghclient.Client) error {
                // Cleanup code
            },
        },
        // Additional test cases
    }

    for _, tc := range testCases {
        t.Run(tc.Name, func(t *testing.T) {
            RunTest(t, tc)
        })
    }
}
```

## Testing Approach

### go-vcr Recording

The project uses go-vcr to record and replay HTTP interactions with the GitHub API. This approach has several benefits:

1. **Deterministic testing**: Tests can run without making actual API calls
2. **Faster execution**: Recorded tests run much faster than live API calls
3. **No API rate limiting**: Avoids GitHub API rate limits during testing
4. **Consistent results**: Tests produce the same results regardless of external factors

When running tests with the `-record` flag, go-vcr records all HTTP interactions and stores them in cassette files. These cassettes are then used for subsequent test runs.

### Golden Files

Golden files store the expected results of test runs in JSON format. They serve as the "ground truth" for test assertions. When running tests with the `-golden` flag, the current test results are written to the golden files, updating the expected results.

This approach allows for easy updating of expected results when the implementation changes intentionally, while still catching unintended changes.

## Running Tests

### Normal Mode

Run tests without recording or updating golden files:

```bash
go test -v ./pkg/tools -run TestPullRequest
```

This mode uses existing cassettes and golden files to verify that the implementation produces the expected results.

### Recording Mode

#### Prerequisites

Before running tests, ensure you have GitHub personal access token with the following permissions to the test repo `geropl/github-mcp-go-test`:
 - content r/w
 - issues r/w
 - pull reqeuests r/w
 - actions r/o

```bash
export GITHUB_PERSONAL_ACCESS_TOKEN=your_token_here
```

#### Run tests in recording mode

Record new HTTP interactions with the GitHub API:

```bash
go test -v ./pkg/tools -run TestPullRequest -record
```

Use this mode when:
- Implementing a new test case
- The GitHub API responses have changed
- The test implementation has changed significantly

### Golden Mode

Update the golden files with the current test results:

```bash
go test -v ./pkg/tools -run TestPullRequest -golden
```

Use this mode when:
- The expected output format has changed intentionally
- Adding new fields to the response
- Changing the formatting of the response

## Test Organization

All tests reside in `pkg/tools`. The test files are organized by GitHub resource type:

```
testdata/
├── TestBranches/
│   ├── CreateBranchFromAnotherBranch/
│   │   ├── create_branch-CreateBranchFromAnotherBranch.golden
│   │   └── create_branch-CreateBranchFromAnotherBranch.yaml
│   └── ...
├── TestCommits/
│   └── ...
├── TestFiles/
│   └── ...
├── TestIssues/
│   └── ...
├── TestPullRequest/
│   ├── SuccessfulCreation/
│   │   ├── create_pull_request-SuccessfulCreation.golden
│   │   └── create_pull_request-SuccessfulCreation.yaml
│   └── ...
├── TestRepository/
│   └── ...
└── TestSearch/
    └── ...
```

Each test case has its own directory containing:
- `.yaml` files: VCR cassettes with recorded HTTP interactions
- `.golden` files: Expected test results in JSON format

## Workflow Guide

The project follows an iterative testing approach:

1. Implement one test case at a time
2. Make each test case work completely before moving to the next
3. Start with "happy path" test cases before error cases

## Test Completion Requirements

Before moving to the next test case or feature, ensure:

1. ✅ Test passes in normal mode (not just recording/golden mode)
2. ✅ Error cases are implemented and passing
3. ✅ Edge cases are considered and tested
4. ✅ Test output is verified against expected behavior

A test is only "complete" when all these criteria are met.

### Step-by-Step Implementation

1. **Uncomment or add a test case** in the appropriate test file
2. **Record HTTP interactions**:
   ```bash
   go test -v ./pkg/tools -run TestCategory/TestCaseName -record
   ```
3. **Create golden files**:
   ```bash
   go test -v ./pkg/tools -run TestCategory/TestCaseName -golden
   ```
4. **Verify the test passes**:
   ```bash
   go test -v ./pkg/tools -run TestCategory/TestCaseName
   ```
   - Ensure the test passes without flags
   - Check that error cases fail appropriately
   - Verify the output matches expected behavior
5. **Move to the next test case**

## Troubleshooting Common Test Issues

### Test Passes in Recording Mode But Fails in Normal Mode
- Check for time-dependent values in responses
- Verify API responses haven't changed
- Ensure golden files match current implementation

### Inconsistent Test Results
- Check for race conditions or external dependencies
- Verify test isolation (tests shouldn't depend on each other)
- Ensure cleanup functions run properly

### VCR Recording Issues
- Check GitHub token permissions
- Verify network connectivity
- Ensure the test repository exists and is accessible
