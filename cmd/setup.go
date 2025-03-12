package cmd

import (
	"fmt"
	"os"

	"github.com/geropl/github-mcp-go/pkg/setup"
	"github.com/geropl/github-mcp-go/pkg/tools"
	"github.com/spf13/cobra"
)

var (
	autoApprove string
	tool        string
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up the GitHub MCP server for use with an AI assistant",
	Long: `Set up the GitHub MCP server for use with an AI assistant.

This command sets up the GitHub MCP server for use with an AI assistant by installing the binary and configuring the AI assistant to use it.

The --auto-approve flag can be used to specify which tools should be auto-approved. It takes a comma-separated list of tool names. "allow-read-only" is a special value to auto-approve all read-only tools`,
	Run: func(cmd *cobra.Command, args []string) {
		// Install the binary
		binaryPath, err := setup.InstallBinary()
		if err != nil {
			fmt.Printf("Error installing binary: %v\n", err)
			os.Exit(1)
		}

		token := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
		if token == "" {
			fmt.Println("GITHUB_PERSONAL_ACCESS_TOKEN environment variable is required")
			os.Exit(1)
		}

		// Set up the tool-specific configuration
		options := setup.SetupOptions{
			BinaryPath:  binaryPath,
			Token:       token,
			AutoApprove: autoApprove,
			Tool:        tool,
		}

		if err := setup.Setup(options, tools.GetReadOnlyToolNames()); err != nil {
			fmt.Printf("Error setting up %s: %v\n", tool, err)
			os.Exit(1)
		}

		fmt.Printf("github-mcp-go binary successfully set up for %s\n", tool)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Add flags to the setup command
	setupCmd.Flags().StringVar(&autoApprove, "auto-approve", "", "Comma-separated list of tools to auto-approve. 'allow-read-only' is a special value to auto-approve all read-only tools")
	setupCmd.Flags().StringVar(&tool, "tool", "cline", "The AI assistant tool to set up for (cline or claude-desktop)")
}
