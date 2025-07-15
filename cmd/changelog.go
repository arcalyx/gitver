/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/gitops"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	fromTag      string
	toTag        string
	outputFile   string
	markdownFlag bool
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate a changelog from commit messages",
	Long: `Generate a changelog from commit messages between two tags or from the latest tag to HEAD.

This command analyzes commit messages following the Conventional Commits specification
and generates a structured changelog. You can specify the starting and ending tags,
or let the command automatically use the latest tag and HEAD.

Examples:
  gitver changelog                      # Generate changelog from latest tag to HEAD
  gitver changelog --from v1.0.0        # Generate changelog from v1.0.0 to HEAD
  gitver changelog --from v1.0.0 --to v2.0.0  # Generate changelog between two tags
  gitver changelog --output CHANGELOG.md      # Save changelog to a file
  gitver changelog --markdown                 # Format output as markdown`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		// If no fromTag is specified, use the latest tag
		if fromTag == "" {
			var err error
			fromTag, err = gitops.GetLastTag()
			if err != nil {
				log.Fatalf("Failed to get latest tag: %v", err)
			}
			if fromTag == "" {
				log.Fatalf("No tags found in repository. Please specify a starting point with --from")
			}
		}

		// Get commits between tags
		var commits []*object.Commit
		var err error
		if toTag == "" {
			commits, err = gitops.GetCommits(fromTag)
		} else {
			commits, err = gitops.GetCommitsBetweenTags(toTag, fromTag)
		}
		if err != nil {
			log.Fatalf("Failed to get commits: %v", err)
		}

		if len(commits) == 0 {
			log.Println("No commits found between the specified tags.")
			return
		}

		// Generate changelog
		changelog := generateChangelog(commits, markdownFlag)

		// Output changelog
		if outputFile != "" {
			err := os.WriteFile(outputFile, []byte(changelog), 0644)
			if err != nil {
				log.Fatalf("Failed to write changelog to file: %v", err)
			}
			log.Printf("Changelog written to %s", outputFile)
		} else {
			fmt.Println(changelog)
		}
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)
	changelogCmd.Flags().StringVar(&fromTag, "from", "", "Starting tag (default: latest tag)")
	changelogCmd.Flags().StringVar(&toTag, "to", "", "Ending tag (default: HEAD)")
	changelogCmd.Flags().StringVar(&outputFile, "output", "", "Output file for changelog")
	changelogCmd.Flags().BoolVar(&markdownFlag, "markdown", false, "Format changelog as markdown")
}

// generateChangelog creates a structured changelog from commit messages
func generateChangelog(commits []*object.Commit, markdown bool) string {
	var sb strings.Builder

	// Add header
	currentVersion := version.ToString()
	date := time.Now().Format("2006-01-02")

	if markdown {
		sb.WriteString(fmt.Sprintf("# Changelog for v%s (%s)\n\n", currentVersion, date))
	} else {
		sb.WriteString(fmt.Sprintf("Changelog for v%s (%s)\n", currentVersion, date))
		sb.WriteString(strings.Repeat("=", 40) + "\n\n")
	}

	// Categorize commits
	var features, fixes, breakingChanges, others []string

	// Regular expressions for parsing conventional commits
	featureRegex := regexp.MustCompile(`^feat(\([^)]+\))?:\s*(.+)`)
	fixRegex := regexp.MustCompile(`^fix(\([^)]+\))?:\s*(.+)`)
	breakingChangeRegex := regexp.MustCompile(`BREAKING CHANGE:\s*(.+)`)

	for _, commit := range commits {
		message := commit.Message

		// Skip merge commits
		if strings.HasPrefix(message, "Merge") {
			continue
		}

		// Check for breaking changes
		if breakingChangeRegex.MatchString(message) {
			matches := breakingChangeRegex.FindStringSubmatch(message)
			breakingChanges = append(breakingChanges, matches[1])
			continue
		}

		// Check for features
		if featureRegex.MatchString(message) {
			matches := featureRegex.FindStringSubmatch(message)
			scope := ""
			if matches[1] != "" {
				scope = matches[1][1:len(matches[1])-1] + ": " // Remove parentheses
			}
			features = append(features, scope+matches[2])
			continue
		}

		// Check for fixes
		if fixRegex.MatchString(message) {
			matches := fixRegex.FindStringSubmatch(message)
			scope := ""
			if matches[1] != "" {
				scope = matches[1][1:len(matches[1])-1] + ": " // Remove parentheses
			}
			fixes = append(fixes, scope+matches[2])
			continue
		}

		// Other commits
		firstLine := strings.Split(message, "\n")[0]
		others = append(others, firstLine)
	}

	// Add sections to changelog
	if len(breakingChanges) > 0 {
		if markdown {
			sb.WriteString("## ⚠️ BREAKING CHANGES\n\n")
			for _, change := range breakingChanges {
				sb.WriteString(fmt.Sprintf("* %s\n", change))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("BREAKING CHANGES:\n")
			sb.WriteString(strings.Repeat("-", 15) + "\n")
			for _, change := range breakingChanges {
				sb.WriteString(fmt.Sprintf("* %s\n", change))
			}
			sb.WriteString("\n")
		}
	}

	if len(features) > 0 {
		if markdown {
			sb.WriteString("## ✨ Features\n\n")
			for _, feature := range features {
				sb.WriteString(fmt.Sprintf("* %s\n", feature))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("Features:\n")
			sb.WriteString(strings.Repeat("-", 9) + "\n")
			for _, feature := range features {
				sb.WriteString(fmt.Sprintf("* %s\n", feature))
			}
			sb.WriteString("\n")
		}
	}

	if len(fixes) > 0 {
		if markdown {
			sb.WriteString("## 🐛 Bug Fixes\n\n")
			for _, fix := range fixes {
				sb.WriteString(fmt.Sprintf("* %s\n", fix))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString("Bug Fixes:\n")
			sb.WriteString(strings.Repeat("-", 10) + "\n")
			for _, fix := range fixes {
				sb.WriteString(fmt.Sprintf("* %s\n", fix))
			}
			sb.WriteString("\n")
		}
	}

	if len(others) > 0 {
		if markdown {
			sb.WriteString("## 🔄 Other Changes\n\n")
			for _, other := range others {
				sb.WriteString(fmt.Sprintf("* %s\n", other))
			}
		} else {
			sb.WriteString("Other Changes:\n")
			sb.WriteString(strings.Repeat("-", 14) + "\n")
			for _, other := range others {
				sb.WriteString(fmt.Sprintf("* %s\n", other))
			}
		}
	}

	return sb.String()
}
