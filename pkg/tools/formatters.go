package tools

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v69/github"
)

// formatPullRequestToMarkdown converts a GitHub PullRequest to markdown
func formatPullRequestToMarkdown(pr *github.PullRequest) string {
	md := fmt.Sprintf("# Pull Request: %s\n\n", pr.GetTitle())
	md += fmt.Sprintf("**Number:** #%d  \n", pr.GetNumber())
	md += fmt.Sprintf("**State:** %s  \n", pr.GetState())
	md += fmt.Sprintf("**Created:** %s  \n", pr.GetCreatedAt().Format(time.RFC1123))
	md += fmt.Sprintf("**URL:** %s  \n\n", pr.GetHTMLURL())
	
	if pr.GetBody() != "" {
		md += fmt.Sprintf("## Description\n\n%s\n\n", pr.GetBody())
	}
	
	md += fmt.Sprintf("## Details\n\n")
	md += fmt.Sprintf("- **Head:** %s  \n", pr.GetHead().GetRef())
	md += fmt.Sprintf("- **Base:** %s  \n", pr.GetBase().GetRef())
	// Mergeable state might not be available immediately
	if pr.Mergeable != nil {
		md += fmt.Sprintf("- **Mergeable:** %t  \n", *pr.Mergeable)
	}
	md += fmt.Sprintf("- **Draft:** %t  \n", pr.GetDraft())
	md += fmt.Sprintf("- **Changes:** +%d/-%d in %d files  \n", 
		pr.GetAdditions(), pr.GetDeletions(), pr.GetChangedFiles())
	
	return md
}

// formatRepositoryToMarkdown converts a GitHub Repository to markdown
func formatRepositoryToMarkdown(repo *github.Repository) string {
	md := fmt.Sprintf("# Repository: %s\n\n", repo.GetFullName())
	md += fmt.Sprintf("**URL:** %s  \n", repo.GetHTMLURL())
	
	if repo.GetDescription() != "" {
		md += fmt.Sprintf("**Description:** %s  \n", repo.GetDescription())
	}
	
	md += fmt.Sprintf("**Created:** %s  \n\n", repo.GetCreatedAt().Format(time.RFC1123))
	
	md += fmt.Sprintf("## Details\n\n")
	md += fmt.Sprintf("- **Default Branch:** %s  \n", repo.GetDefaultBranch())
	md += fmt.Sprintf("- **Stars:** %d  \n", repo.GetStargazersCount())
	md += fmt.Sprintf("- **Forks:** %d  \n", repo.GetForksCount())
	md += fmt.Sprintf("- **Open Issues:** %d  \n", repo.GetOpenIssuesCount())
	md += fmt.Sprintf("- **Private:** %t  \n", repo.GetPrivate())
	
	return md
}

// formatRepositorySearchToMarkdown converts GitHub repository search results to markdown
func formatRepositorySearchToMarkdown(result *github.RepositoriesSearchResult) string {
	md := fmt.Sprintf("# Repository Search Results\n\n")
	md += fmt.Sprintf("**Total Results:** %d\n\n", result.GetTotal())
	
	if len(result.Repositories) == 0 {
		md += "No repositories found.\n"
		return md
	}
	
	md += "## Repositories\n\n"
	
	for i, repo := range result.Repositories {
		md += fmt.Sprintf("### %d. [%s](%s)\n\n", i+1, repo.GetFullName(), repo.GetHTMLURL())
		
		if repo.GetDescription() != "" {
			md += fmt.Sprintf("%s\n\n", repo.GetDescription())
		}
		
		md += fmt.Sprintf("- **Stars:** %d\n", repo.GetStargazersCount())
		md += fmt.Sprintf("- **Language:** %s\n", repo.GetLanguage())
		md += fmt.Sprintf("- **Updated:** %s\n\n", repo.GetUpdatedAt().Format(time.RFC1123))
	}
	
	return md
}

// formatFileContentToMarkdown converts GitHub file content to markdown
func formatFileContentToMarkdown(content map[string]interface{}) string {
	md := fmt.Sprintf("# File: %s\n\n", content["path"])
	
	md += fmt.Sprintf("**Path:** %s  \n", content["path"])
	md += fmt.Sprintf("**Size:** %v bytes  \n", content["size"])
	md += fmt.Sprintf("**URL:** %s  \n\n", content["html_url"])
	
	if fileContent, ok := content["content"].(string); ok && fileContent != "" {
		// Determine if we should use a code block based on file extension
		path := content["path"].(string)
		extension := ""
		if lastDot := strings.LastIndex(path, "."); lastDot >= 0 {
			extension = path[lastDot+1:]
		}
		
		// Add content with appropriate formatting
		md += "## Content\n\n"
		if extension != "" {
			md += fmt.Sprintf("```%s\n%s\n```\n", extension, fileContent)
		} else {
			md += fmt.Sprintf("```\n%s\n```\n", fileContent)
		}
	}
	
	return md
}

// formatDirectoryContentToMarkdown converts GitHub directory content to markdown
func formatDirectoryContentToMarkdown(content map[string]interface{}) string {
	md := fmt.Sprintf("# Directory: %s\n\n", content["path"])
	
	contents, ok := content["contents"].([]map[string]interface{})
	if !ok {
		return md + "No contents found."
	}
	
	md += "## Contents\n\n"
	
	for i, item := range contents {
		itemType := item["type"].(string)
		name := item["name"].(string)
		htmlURL := item["html_url"].(string)
		
		if itemType == "dir" {
			md += fmt.Sprintf("%d. 📁 [%s](%s)  \n", i+1, name, htmlURL)
		} else {
			md += fmt.Sprintf("%d. 📄 [%s](%s)  \n", i+1, name, htmlURL)
		}
	}
	
	return md
}

// formatFileUpdateToMarkdown converts GitHub file update result to markdown
func formatFileUpdateToMarkdown(result *github.RepositoryContentResponse) string {
	md := fmt.Sprintf("# File Update: %s\n\n", result.GetContent().GetPath())
	
	md += fmt.Sprintf("**Path:** %s  \n", result.GetContent().GetPath())
	md += fmt.Sprintf("**SHA:** %s  \n", result.GetContent().GetSHA())
	md += fmt.Sprintf("**URL:** %s  \n\n", result.GetContent().GetHTMLURL())
	
	// Commit information is not directly accessible in RepositoryContentResponse
	// We'll just show the content information
	
	return md
}

// formatCommitToMarkdown converts GitHub commit result to markdown
func formatCommitToMarkdown(result *github.Commit) string {
	md := fmt.Sprintf("# Commit\n\n")
	
	md += fmt.Sprintf("**SHA:** %s  \n", result.GetSHA())
	md += fmt.Sprintf("**Message:** %s  \n", result.GetMessage())
	
	// Get author information if available
	if author := result.GetAuthor(); author != nil {
		md += fmt.Sprintf("**Author:** %s  \n", author.GetLogin())
	}
	
	// Get committer information if available
	if committer := result.GetCommitter(); committer != nil {
		md += fmt.Sprintf("**Committer:** %s  \n", committer.GetLogin())
	}
	
	md += fmt.Sprintf("**URL:** %s  \n\n", result.GetHTMLURL())
	
	return md
}

// formatPullRequestDiffToMarkdown formats a pull request diff as markdown
func formatPullRequestDiffToMarkdown(number int, diff string) string {
	md := fmt.Sprintf("# Pull Request Diff (#%d)\n\n", number)
	
	// Truncate diff if it's too long
	const maxDiffLength = 10000
	if len(diff) > maxDiffLength {
		diff = diff[:maxDiffLength] + "\n\n... (diff truncated due to size)"
	}
	
	md += "```diff\n" + diff + "\n```\n"
	
	return md
}
