package tools

import (
	"testing"
)

const (
	OWNER = "geropl"
	REPO  = "github-mcp-go-test"
	BRANCH = "test/feature-branch-1"
)

func TestPullRequest(t *testing.T) {
	testCases := []TestCase{
		{
			Name: "SuccessfulCreation",
			Tool: "create_pull_request",
			Input: map[string]interface{}{
				"owner": OWNER,
				"repo":  REPO,
				"title": "Test PR",
				"body":  "Test PR body",
				"head":  BRANCH,
				"base":  "main",
				"draft": false,
			},
		},
		// {
		// 	Name: "MissingOwner",
		// 	Tool: "create_pull_request",
		// 	Input: map[string]interface{}{
		// 		"repo":  REPO,
		// 		"title": "Test PR",
		// 		"body":  "Test PR body",
		// 		"head":  BRANCH,
		// 		"base":  "main",
		// 		"draft": false,
		// 	},
		// },
		// {
		// 	Name: "MissingRepo",
		// 	Tool: "create_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner": OWNER,
		// 		"title": "Test PR",
		// 		"body":  "Test PR body",
		// 		"head":  BRANCH,
		// 		"base":  "main",
		// 		"draft": false,
		// 	},
		// },
		// {
		// 	Name: "MissingTitle",
		// 	Tool: "create_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner": OWNER,
		// 		"repo":  REPO,
		// 		"body":  "Test PR body",
		// 		"head":  BRANCH,
		// 		"base":  "main",
		// 		"draft": false,
		// 	},
		// },
		// {
		// 	Name: "MissingHead",
		// 	Tool: "create_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner": OWNER,
		// 		"repo":  REPO,
		// 		"title": "Test PR",
		// 		"body":  "Test PR body",
		// 		"base":  "main",
		// 		"draft": false,
		// 	},
		// },
		// {
		// 	Name: "MissingBase",
		// 	Tool: "create_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner": OWNER,
		// 		"repo":  REPO,
		// 		"title": "Test PR",
		// 		"body":  "Test PR body",
		// 		"head":  BRANCH,
		// 		"draft": false,
		// 	},
		// },
		// {
		// 	Name: "OptionalBodyOmitted",
		// 	Tool: "create_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner": OWNER,
		// 		"repo":  REPO,
		// 		"title": "Test PR",
		// 		"head":  BRANCH,
		// 		"base":  "main",
		// 		"draft": false,
		// 	},
		// },
		// {
		// 	Name: "OptionalDraftOmitted",
		// 	Tool: "create_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner": OWNER,
		// 		"repo":  REPO,
		// 		"title": "Test PR",
		// 		"body":  "Test PR body",
		// 		"head":  BRANCH,
		// 		"base":  "main",
		// 	},
		// },
		// // get_pull_request
		// {
		// 	Name: "SuccessfulRetrieval",
		// 	Tool: "get_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner":  "testowner",
		// 		"repo":   "testrepo",
		// 		"number": float64(1),
		// 	},
		// },
		// {
		// 	Name: "MissingOwner",
		// 	Tool: "get_pull_request",
		// 	Input: map[string]interface{}{
		// 		"repo":   "testrepo",
		// 		"number": float64(1),
		// 	},
		// },
		// {
		// 	Name: "MissingRepo",
		// 	Tool: "get_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner":  "testowner",
		// 		"number": float64(1),
		// 	},
		// },
		// {
		// 	Name: "MissingNumber",
		// 	Tool: "get_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner": OWNER,
		// 		"repo":  REPO,
		// 	},
		// },
		// {
		// 	Name: "InvalidNumberType",
		// 	Tool: "get_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner":  "testowner",
		// 		"repo":   "testrepo",
		// 		"number": "1", // String instead of number
		// 	},
		// },
		// {
		// 	Name: "NonExistentPullRequest",
		// 	Tool: "get_pull_request",
		// 	Input: map[string]interface{}{
		// 		"owner":  "testowner",
		// 		"repo":   "testrepo",
		// 		"number": float64(999),
		// 	},
		// },
		// // get_pull_request_diff
		// {
		// 	Name: "SuccessfulDiffRetrieval",
		// 	Tool: "get_pull_request_diff",
		// 	Input: map[string]interface{}{
		// 		"owner":  "testowner",
		// 		"repo":   "testrepo",
		// 		"number": float64(1),
		// 	},
		// },
		// {
		// 	Name: "MissingOwner",
		// 	Tool: "get_pull_request_diff",
		// 	Input: map[string]interface{}{
		// 		"repo":   "testrepo",
		// 		"number": float64(1),
		// 	},
		// },
		// {
		// 	Name: "MissingRepo",
		// 	Tool: "get_pull_request_diff",
		// 	Input: map[string]interface{}{
		// 		"owner":  "testowner",
		// 		"number": float64(1),
		// 	},
		// },
		// {
		// 	Name: "MissingNumber",
		// 	Tool: "get_pull_request_diff",
		// 	Input: map[string]interface{}{
		// 		"owner": OWNER,
		// 		"repo":  REPO,
		// 	},
		// },
		// {
		// 	Name: "InvalidNumberType",
		// 	Tool: "get_pull_request_diff",
		// 	Input: map[string]interface{}{
		// 		"owner":  "testowner",
		// 		"repo":   "testrepo",
		// 		"number": "1", // String instead of number
		// 	},
		// },
		// {
		// 	Name: "NonExistentPullRequest",
		// 	Tool: "get_pull_request_diff",
		// 	Input: map[string]interface{}{
		// 		"owner":  "testowner",
		// 		"repo":   "testrepo",
		// 		"number": float64(999),
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			RunTest(t, tc)
		})
	}
}
