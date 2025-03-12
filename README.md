# GitHub MCP Server

A Model Context Protocol (MCP) server for GitHub, implemented in Go. This server allows LLMs to interact with GitHub repositories, issues, pull requests, and more through a standardized interface.

## Features

- **Repository Operations**
  - Search repositories
  - Create repositories
  - Fork repositories

- **Pull Request Operations**
  - Create pull requests
  - Get pull request details
  - Get pull request diffs

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
./github-mcp-go
```

### Connecting to Claude for Desktop

To use the server with Claude for Desktop, add the following to your Claude for Desktop configuration file:

```json
{
  "mcpServers": {
    "github": {
      "command": "/path/to/github-mcp-go",
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "your_token_here"
      }
    }
  }
}
```

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

### Project Structure

```
github-mcp-go/
├── cmd/
│   └── github-mcp-go/
│       └── main.go
├── internal/
│   ├── server/
│   │   └── server.go
│   ├── github/
│   │   ├── client.go
│   │   ├── repository.go
│   │   ├── pulls.go
│   │   └── ...
│   └── errors/
│       └── errors.go
├── pkg/
│   └── tools/
│       ├── repository.go
│       ├── pulls.go
│       └── ...
└── test/
    └── ...
```

### Testing

The project uses table-driven tests with go-vcr for recording HTTP interactions:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
