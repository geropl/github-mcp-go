package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

// SetupOptions contains the options for setting up the GitHub MCP server
type SetupOptions struct {
	BinaryPath  string
	Token       string
	AutoApprove string
	Tool        string
}

// Setup sets up the GitHub MCP server for use with an AI assistant
func Setup(options SetupOptions, readOnlyTools map[string]bool) error {
	// Find the config directory for the specified tool
	configDir, err := FindConfigDir(options.Tool)
	if err != nil {
		return fmt.Errorf("failed to find config directory: %w", err)
	}

	// Process auto-approve flag
	autoApproveTools := ProcessAutoApproveFlag(options.AutoApprove, readOnlyTools)

	// Create the server configuration
	serverArgs := []string{"serve"}

	// Create the server configuration
	serverConfig := map[string]interface{}{
		"command":     options.BinaryPath,
		"args":        serverArgs,
		"disabled":    false,
		"autoApprove": autoApproveTools,
	}

	// Add environment variables if token is provided
	if options.Token != "" {
		serverConfig["env"] = map[string]interface{}{
			"GITHUB_PERSONAL_ACCESS_TOKEN": options.Token,
		}
	}

	// Set up the tool-specific configuration
	var settingsPath string
	switch options.Tool {
	case "cline":
		settingsPath = filepath.Join(configDir, "cline_mcp_settings.json")
	case "claude-desktop":
		settingsPath = filepath.Join(configDir, "claude_desktop_config.json")
	default:
		return fmt.Errorf("unsupported tool: %s", options.Tool)
	}
	
	err = UpdateSettingsFile(settingsPath, serverConfig)
	if err != nil {
		return fmt.Errorf("failed to update settings file at '%s': %w", settingsPath, err)
	}
	fmt.Printf("%s MCP settings updated at %s\n", options.Tool, settingsPath)

	fmt.Printf("github-mcp-go successfully set up for %s\n", options.Tool)
	return nil
}

// InstallBinary installs the github-mcp-go binary to the specified directory
func InstallBinary() (string, error) {
	// Create the MCP servers directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	mcpServersDir := filepath.Join(homeDir, "mcp-servers")
	if err := os.MkdirAll(mcpServersDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create MCP servers directory: %w", err)
	}

	// Check if the github-mcp-go binary is already on the path
	binaryPath, found := CheckBinary(mcpServersDir)
	if !found {
		fmt.Printf("github-mcp-go binary not found on path, copying current binary to '%s'...\n", binaryPath)
		err := CopySelfToBinaryPath(binaryPath)
		if err != nil {
			return "", fmt.Errorf("failed to copy github-mcp-go binary: %w", err)
		}
	}

	return binaryPath, nil
}
