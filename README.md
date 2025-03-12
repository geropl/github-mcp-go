# GitHub MCP Server

A Model Context Protocol (MCP) server for GitHub, implemented in Go. This server allows LLMs to interact with GitHub repositories, issues, pull requests, and more through a standardized interface.

## Features

- **Comprehensive GitHub API Access**: See [Tools](#available-tools) below
- **Setup Command**: Easy automatic setup for AI assistants
  - `--auto-approve` allows to pre-fill the auto-approval checkbox
  - `--write-access` enables write access for remote operations
  - Copies the binary into a stable location
  - Updates the tool MCP server config (cline and claude-desktop are supported)
- **Security Controls**:
  - Write access is disabled by default for safety
  - Fine-grained auto-approval options


## Installation

### Prerequisites

- GitHub Personal Access Token with appropriate permissions

### Using Pre-built Binaries

You can download pre-built binaries for your platform from the [GitHub Releases](https://github.com/geropl/github-mcp-go/releases) page.

```bash
# Download the latest release for your platform
# Example for Linux (amd64):
curl -L https://github.com/geropl/github-mcp-go/releases/latest/download/github-mcp-go_<version>_linux_amd64 -o github-mcp-go
chmod +x github-mcp-go
```

### Building from Source

If you prefer to build from source:

```bash
# Prerequisites: Go 1.23 or later

# Clone the repository
git clone https://github.com/geropl/github-mcp-go.git
cd github-mcp-go

# Build the server
go build -o github-mcp-go .
```

## Usage

### Environment Variables

The server requires a GitHub Personal Access Token to authenticate with the GitHub API:

```bash
export GITHUB_PERSONAL_ACCESS_TOKEN=your_token_here
```

### Running the Server

```bash
# Run the server directly (read-only mode)
./github-mcp-go serve

# Run the server with write access enabled
./github-mcp-go serve --write-access

# Show help
./github-mcp-go --help
```

### Setup Command

The server includes a convenient setup command to install and configure the MCP server for use with AI assistants:

```bash
# Set up for Cline with auto-approval for read-only tools
./github-mcp-go setup --auto-approve allow-read-only --tool cline

# Set up for Claude Desktop with auto-approval for specific tools
./github-mcp-go setup --auto-approve search_repositories,get_file_contents --tool claude-desktop

# Set up with write access enabled
./github-mcp-go setup --write-access --tool cline

# Show setup help
./github-mcp-go setup --help
```

#### Auto-Approval Options

The `--auto-approve` flag can be used to specify which tools should be auto-approved as a comma-separated list. `allow-read-only` is a special value to add all read-only tools to the auto-approve list (safe, no state changes).

The `--write-access` flag enables write access for remote operations. This allows tools that modify remote repositories to be used. By default, write access is disabled for safety.


## Available Tools

### Repository Tools

- `search_repositories`: Search for GitHub repositories
- `create_repository`: Create a new GitHub repository
- `fork_repository`: Fork a GitHub repository

### Pull Request Tools

- `create_pull_request`: Create a new pull request
- `get_pull_request`: Get detailed information about a pull request
- `get_pull_request_diff`: Get the diff of a pull request

### File Tools

- `get_file_contents`: Get the contents of a file or directory
- `create_or_update_file`: Create or update a file
- `push_files`: Push multiple files in a single commit

### Issue Tools

- `create_issue`: Create a new issue
- `list_issues`: List issues with filtering options
- `update_issue`: Update an existing issue
- `add_issue_comment`: Add a comment to an issue
- `get_issue`: Get details of a specific issue
- `list_issue_comments`: List comments on an issue

### Branch Tools

- `list_branches`: List branches in a repository
- `get_branch`: Get details about a specific branch
- `create_branch`: Create a new branch
- `merge_branches`: Merge one branch into another
- `delete_branch`: Delete a branch

### Commit Tools

- `get_commit`: Get details of a specific commit
- `list_commits`: List commits in a repository
- `compare_commits`: Compare two commits or branches
- `get_commit_status`: Get the combined status for a specific commit
- `create_commit_comment`: Add a comment to a specific commit
- `list_commit_comments`: List comments for a specific commit
- `create_commit`: Create a new commit directly

### Search Tools

- `search_code`: Search for code across repositories
- `search_issues`: Search for issues and pull requests
- `search_commits`: Search for commits across repositories

## Releases

The project follows [Semantic Versioning](https://semver.org/). New releases are automatically built and published to GitHub Releases when a new tag is pushed to the repository.

### Release Process

1. Update the CHANGELOG.md file with the changes in the new version
2. Create and push a new tag:
   ```bash
   git tag -a v0.1.1 -m "Release v0.1.1"
   git push origin v0.1.1
   ```
3. The GitHub Actions workflow will automatically build binaries for multiple platforms and create a release

### Available Platforms

Pre-built binaries are available for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Development

### Testing

The project uses table-driven tests with go-vcr for recording HTTP interactions:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
