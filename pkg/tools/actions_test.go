package tools

import (
	"testing"
)

func TestActions(t *testing.T) {
	testCases := []*TestCase{
		// get_workflow - Happy Path
		{
			Name: "GetWorkflowByID",
			Tool: "get_workflow",
			Input: map[string]interface{}{
				"owner":       "geropl",
				"repo":        "github-mcp-go-test",
				"workflow_id": "149593961", // Actual workflow ID
			},
		},
		{
			Name: "GetWorkflowByFileName",
			Tool: "get_workflow",
			Input: map[string]interface{}{
				"owner":       "geropl",
				"repo":        "github-mcp-go-test",
				"workflow_id": "test-workflow.yml", // Just the filename, not the full path
			},
		},

		// get_workflow - Error Cases
		{
			Name: "GetWorkflowInvalidOwner",
			Tool: "get_workflow",
			Input: map[string]interface{}{
				"owner":       "",
				"repo":        "github-mcp-go-test",
				"workflow_id": "12345678",
			},
		},
		{
			Name: "GetWorkflowInvalidRepo",
			Tool: "get_workflow",
			Input: map[string]interface{}{
				"owner":       "geropl",
				"repo":        "",
				"workflow_id": "12345678",
			},
		},
		{
			Name: "GetWorkflowInvalidID",
			Tool: "get_workflow",
			Input: map[string]interface{}{
				"owner":       "geropl",
				"repo":        "github-mcp-go-test",
				"workflow_id": "",
			},
		},
		{
			Name: "GetWorkflowNonExistentID",
			Tool: "get_workflow",
			Input: map[string]interface{}{
				"owner":       "geropl",
				"repo":        "github-mcp-go-test",
				"workflow_id": "99999999", // Non-existent workflow ID
			},
		},
		{
			Name: "GetWorkflowNonExistentFileName",
			Tool: "get_workflow",
			Input: map[string]interface{}{
				"owner":       "geropl",
				"repo":        "github-mcp-go-test",
				"workflow_id": "non-existent.yml",
			},
		},
		
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
