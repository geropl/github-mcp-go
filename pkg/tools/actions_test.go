package tools

import (
	"testing"
)

func TestActions(t *testing.T) {
	testCases := []*TestCase{
		// list_workflows - Happy Path
		{
			Name: "ListWorkflows",
			Tool: "list_workflows",
			Input: map[string]interface{}{
				"owner": "geropl",
				"repo":  "github-mcp-go-test",
			},
		},
		
		// list_workflows - Error Cases
		{
			Name: "ListWorkflowsInvalidOwner",
			Tool: "list_workflows",
			Input: map[string]interface{}{
				"owner": "",
				"repo":  "github-mcp-go-test",
			},
		},
		{
			Name: "ListWorkflowsInvalidRepo",
			Tool: "list_workflows",
			Input: map[string]interface{}{
				"owner": "geropl",
				"repo":  "",
			},
		},
		{
			Name: "ListWorkflowsNonExistentRepo",
			Tool: "list_workflows",
			Input: map[string]interface{}{
				"owner": "geropl",
				"repo":  "non-existent-repo",
			},
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			RunTest(t, tc)
		})
	}
}
