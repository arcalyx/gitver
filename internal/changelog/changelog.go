package changelog

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/gitops"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/viper"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// CommitType represents the type of a conventional commit
type CommitType string

const (
	TypeFeat     CommitType = "feat"
	TypeFix      CommitType = "fix"
	TypeDocs     CommitType = "docs"
	TypeStyle    CommitType = "style"
	TypeRefactor CommitType = "refactor"
	TypePerf     CommitType = "perf"
	TypeTest     CommitType = "test"
	TypeBuild    CommitType = "build"
	TypeCI       CommitType = "ci"
	TypeChore    CommitType = "chore"
)

// CommitInfo represents a parsed conventional commit
type CommitInfo struct {
	Type        CommitType
	Scope       string
	Description string
	Body        string
	Footer      string
	Breaking    bool
	Hash        string
	Author      string
	Date        time.Time
}

// Format represents the output format for the changelog
type Format string

const (
	FormatMarkdown Format = "markdown"
	FormatPlain    Format = "plaintext"
)

// conventionalCommitPattern is a regex pattern for parsing conventional commits
var conventionalCommitPattern = regexp.MustCompile(`^(?:(\w+)(?:\(([^)]+)\))?: (.+))(?:\n\n(.+))?(?:\n\n(.+))?$`)

// ParseCommit parses a commit message according to the Conventional Commits specification
func ParseCommit(commit *object.Commit) *CommitInfo {
	message := commit.Message

	matches := conventionalCommitPattern.FindStringSubmatch(message)
	if len(matches) < 4 {
		return nil // Not a conventional commit
	}

	commitType := matches[1]
	scope := matches[2]
	description := matches[3]
	body := ""
	footer := ""

	if len(matches) > 4 {
		body = matches[4]
	}
	if len(matches) > 5 {
		footer = matches[5]
	}

	// Check for breaking changes
	breaking := strings.Contains(message, "BREAKING CHANGE:") ||
		strings.Contains(message, "!:") ||
		strings.HasSuffix(commitType, "!")

	return &CommitInfo{
		Type:        CommitType(commitType),
		Scope:       scope,
		Description: description,
		Body:        body,
		Footer:      footer,
		Breaking:    breaking,
		Hash:        commit.Hash.String()[:7],
		Author:      commit.Author.Name,
		Date:        commit.Author.When,
	}
}

// GenerateChangelog generates a changelog based on commits since the last tag
func GenerateChangelog(version string) (string, error) {
	// Get the format from config
	format := Format(viper.GetString("changelog.format"))
	if format == "" {
		format = FormatMarkdown
	}

	// Get commits since last tag
	tag, err := gitops.GetLastTag()
	if err != nil {
		return "", fmt.Errorf("failed to get last tag: %v", err)
	}

	commits, err := gitops.GetCommits(tag)
	if err != nil {
		return "", fmt.Errorf("failed to get commits: %v", err)
	}

	// Parse commits
	var conventionalCommits []*CommitInfo
	for _, commit := range commits {
		info := ParseCommit(commit)
		if info != nil {
			conventionalCommits = append(conventionalCommits, info)
		}
	}

	// Generate changelog content
	var content string
	if format == FormatMarkdown {
		content = generateMarkdownChangelog(version, conventionalCommits)
	} else {
		content = generatePlainTextChangelog(version, conventionalCommits)
	}

	return content, nil
}

// SaveChangelog saves the changelog to a file
func SaveChangelog(content string) error {
	outputFile := viper.GetString("changelog.output_file")
	if outputFile == "" {
		outputFile = "CHANGELOG.md"
	}

	// If the file exists, prepend the new content
	var finalContent string
	if _, err := os.Stat(outputFile); err == nil {
		existingContent, err := os.ReadFile(outputFile)
		if err != nil {
			return fmt.Errorf("failed to read existing changelog: %v", err)
		}
		finalContent = content + string(existingContent)
	} else {
		finalContent = content
	}

	// Write to file
	err := os.WriteFile(outputFile, []byte(finalContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write changelog: %v", err)
	}

	log.Printf("Changelog saved to %s", outputFile)
	return nil
}

// generateMarkdownChangelog generates a changelog in Markdown format
func generateMarkdownChangelog(version string, commits []*CommitInfo) string {
	var sb strings.Builder

	// Add header
	date := time.Now().Format("2006-01-02")
	sb.WriteString(fmt.Sprintf("# %s (%s)\n\n", version, date))

	// Group commits by type
	featCommits := filterCommitsByType(commits, TypeFeat)
	fixCommits := filterCommitsByType(commits, TypeFix)
	docsCommits := filterCommitsByType(commits, TypeDocs)
	styleCommits := filterCommitsByType(commits, TypeStyle)
	refactorCommits := filterCommitsByType(commits, TypeRefactor)
	perfCommits := filterCommitsByType(commits, TypePerf)
	testCommits := filterCommitsByType(commits, TypeTest)
	buildCommits := filterCommitsByType(commits, TypeBuild)
	ciCommits := filterCommitsByType(commits, TypeCI)
	choreCommits := filterCommitsByType(commits, TypeChore)

	// Breaking changes
	breakingCommits := filterBreakingCommits(commits)
	if len(breakingCommits) > 0 {
		sb.WriteString("## ⚠ BREAKING CHANGES\n\n")
		for _, commit := range breakingCommits {
			sb.WriteString(fmt.Sprintf("* **%s:** %s ([%s](commit/%s))\n",
				commit.Scope, commit.Description, commit.Hash, commit.Hash))
		}
		sb.WriteString("\n")
	}

	// Features
	if len(featCommits) > 0 {
		sb.WriteString("## ✨ Features\n\n")
		for _, commit := range featCommits {
			scope := ""
			if commit.Scope != "" {
				scope = "**" + commit.Scope + ":** "
			}
			sb.WriteString(fmt.Sprintf("* %s%s ([%s](commit/%s))\n",
				scope, commit.Description, commit.Hash, commit.Hash))
		}
		sb.WriteString("\n")
	}

	// Bug Fixes
	if len(fixCommits) > 0 {
		sb.WriteString("## 🐛 Bug Fixes\n\n")
		for _, commit := range fixCommits {
			scope := ""
			if commit.Scope != "" {
				scope = "**" + commit.Scope + ":** "
			}
			sb.WriteString(fmt.Sprintf("* %s%s ([%s](commit/%s))\n",
				scope, commit.Description, commit.Hash, commit.Hash))
		}
		sb.WriteString("\n")
	}

	// Documentation
	if len(docsCommits) > 0 {
		sb.WriteString("## 📚 Documentation\n\n")
		for _, commit := range docsCommits {
			scope := ""
			if commit.Scope != "" {
				scope = "**" + commit.Scope + ":** "
			}
			sb.WriteString(fmt.Sprintf("* %s%s ([%s](commit/%s))\n",
				scope, commit.Description, commit.Hash, commit.Hash))
		}
		sb.WriteString("\n")
	}

	// Other changes (refactor, style, perf, test, build, ci, chore)
	otherCommits := append(refactorCommits, styleCommits...)
	otherCommits = append(otherCommits, perfCommits...)
	otherCommits = append(otherCommits, testCommits...)
	otherCommits = append(otherCommits, buildCommits...)
	otherCommits = append(otherCommits, ciCommits...)
	otherCommits = append(otherCommits, choreCommits...)

	if len(otherCommits) > 0 {
		sb.WriteString("## 🔧 Other Changes\n\n")
		for _, commit := range otherCommits {
			scope := ""
			if commit.Scope != "" {
				scope = "**" + commit.Scope + ":** "
			}
			sb.WriteString(fmt.Sprintf("* **%s:** %s%s ([%s](commit/%s))\n",
				commit.Type, scope, commit.Description, commit.Hash, commit.Hash))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// generatePlainTextChangelog generates a changelog in plain text format
func generatePlainTextChangelog(version string, commits []*CommitInfo) string {
	var sb strings.Builder

	// Add header
	date := time.Now().Format("2006-01-02")
	sb.WriteString(fmt.Sprintf("%s (%s)\n\n", version, date))

	// Group commits by type
	featCommits := filterCommitsByType(commits, TypeFeat)
	fixCommits := filterCommitsByType(commits, TypeFix)
	docsCommits := filterCommitsByType(commits, TypeDocs)
	styleCommits := filterCommitsByType(commits, TypeStyle)
	refactorCommits := filterCommitsByType(commits, TypeRefactor)
	perfCommits := filterCommitsByType(commits, TypePerf)
	testCommits := filterCommitsByType(commits, TypeTest)
	buildCommits := filterCommitsByType(commits, TypeBuild)
	ciCommits := filterCommitsByType(commits, TypeCI)
	choreCommits := filterCommitsByType(commits, TypeChore)

	// Breaking changes
	breakingCommits := filterBreakingCommits(commits)
	if len(breakingCommits) > 0 {
		sb.WriteString("BREAKING CHANGES:\n\n")
		for _, commit := range breakingCommits {
			sb.WriteString(fmt.Sprintf("- %s: %s (%s)\n",
				commit.Scope, commit.Description, commit.Hash))
		}
		sb.WriteString("\n")
	}

	// Features
	if len(featCommits) > 0 {
		sb.WriteString("Features:\n\n")
		for _, commit := range featCommits {
			scope := ""
			if commit.Scope != "" {
				scope = commit.Scope + ": "
			}
			sb.WriteString(fmt.Sprintf("- %s%s (%s)\n",
				scope, commit.Description, commit.Hash))
		}
		sb.WriteString("\n")
	}

	// Bug Fixes
	if len(fixCommits) > 0 {
		sb.WriteString("Bug Fixes:\n\n")
		for _, commit := range fixCommits {
			scope := ""
			if commit.Scope != "" {
				scope = commit.Scope + ": "
			}
			sb.WriteString(fmt.Sprintf("- %s%s (%s)\n",
				scope, commit.Description, commit.Hash))
		}
		sb.WriteString("\n")
	}

	// Other changes
	otherCommits := append(docsCommits, styleCommits...)
	otherCommits = append(otherCommits, refactorCommits...)
	otherCommits = append(otherCommits, perfCommits...)
	otherCommits = append(otherCommits, testCommits...)
	otherCommits = append(otherCommits, buildCommits...)
	otherCommits = append(otherCommits, ciCommits...)
	otherCommits = append(otherCommits, choreCommits...)

	if len(otherCommits) > 0 {
		sb.WriteString("Other Changes:\n\n")
		for _, commit := range otherCommits {
			scope := ""
			if commit.Scope != "" {
				scope = commit.Scope + ": "
			}
			sb.WriteString(fmt.Sprintf("- %s: %s%s (%s)\n",
				commit.Type, scope, commit.Description, commit.Hash))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// filterCommitsByType filters commits by their type
func filterCommitsByType(commits []*CommitInfo, commitType CommitType) []*CommitInfo {
	var filtered []*CommitInfo
	for _, commit := range commits {
		if commit.Type == commitType {
			filtered = append(filtered, commit)
		}
	}
	return filtered
}

// filterBreakingCommits filters commits that are breaking changes
func filterBreakingCommits(commits []*CommitInfo) []*CommitInfo {
	var filtered []*CommitInfo
	for _, commit := range commits {
		if commit.Breaking {
			filtered = append(filtered, commit)
		}
	}
	return filtered
}
