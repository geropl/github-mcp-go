# Pull Request Testing Implementation

## Overview

This project implements comprehensive testing for the GitHub MCP server's pull request operations. The tests cover three main tools:

1. `create_pull_request` - Create a new pull request
2. `get_pull_request` - Get detailed information about a pull request
3. `get_pull_request_diff` - Get the diff of a pull request

## Test Structure

The tests are implemented in `pkg/tools/pulls_test.go` using a table-driven approach with the following components:

- **TestCase struct**: Defines the test name, tool, input parameters, and setup/cleanup functions
- **go-vcr**: Records and replays HTTP interactions with the GitHub API
- **Golden files**: Store expected test results for comparison

## Test Cases

A total of 25 test cases have been defined across the three pull request tools:

- **create_pull_request**: 10 test cases (5 happy path, 5 error cases)
- **get_pull_request**: 7 test cases (3 happy path, 4 error cases)
- **get_pull_request_diff**: 8 test cases (4 happy path, 4 error cases)

All test cases except the initial "SuccessfulCreation" case are commented out, ready to be implemented one by one following the project's iterative testing approach.

## Implementation Approach

Following the project's testing patterns, we take an iterative approach:

1. Implement one test case at a time
2. Make each test case work completely before moving to the next
3. Start with "happy path" test cases before error cases
4. Use the `-record` flag to record new HTTP interactions
5. Use the `-golden` flag to create golden files
6. Run tests without flags to verify functionality

## Detailed Testing Guide

For detailed instructions on implementing and running the test cases, refer to [PULL_REQUEST_TESTING.md](./PULL_REQUEST_TESTING.md).

## Running Tests

To run a specific test case:

```bash
go test -v ./pkg/tools -run TestPullRequest/[TestCaseName]
```

To record HTTP interactions:

```bash
go test -v ./pkg/tools -run TestPullRequest/[TestCaseName] -record
```

To update golden files:

```bash
go test -v ./pkg/tools -run TestPullRequest/[TestCaseName] -golden
```

To run all pull request tests:

```bash
go test -v ./pkg/tools -run TestPullRequest
```

## Prerequisites

Before running the tests, ensure you have:

1. A GitHub personal access token with appropriate permissions
2. Access to the test repository (`geropl/github-mcp-go-test`)
3. The token set as an environment variable: `GITHUB_PERSONAL_ACCESS_TOKEN`
