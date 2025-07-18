/*
Copyright 2024 NAME HERE EMAIL ADDRESS
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/arcalyx/gitver/internal/gitops"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var (
	historyFlag bool
	countFlag   int
	formatFlag  string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long: `Display the current version of your project and optionally show version history.

This command shows the current version of your project. With the --history flag,
it also displays a history of version tags from your Git repository.

Examples:
  gitver version              # Show current version
  gitver version --format semver  # Show only the semver string
  gitver version --format json    # Show version info as JSON
  gitver version --format full    # Show detailed version information
  gitver version --history    # Show version history
  gitver version --history --count 5  # Show last 5 versions`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		// Display current version based on format
		switch formatFlag {
		case "semver":
			// Just output the raw semver string
			fmt.Print(version.ToString())
			return
		case "json":
			// Output version information as JSON
			versionInfo := map[string]interface{}{
				"version":     version.ToString(),
				"major":       version.GetMajor(),
				"minor":       version.GetMinor(),
				"patch":       version.GetPatch(),
				"prerelease":  version.GetPreRelease(),
				"buildmeta":   version.GetBuildMetadata(),
				"lastVersion": version.GetLastVersion(),
			}
			jsonData, err := json.MarshalIndent(versionInfo, "", "  ")
			if err != nil {
				log.Fatalf("Failed to create JSON output: %v", err)
			}
			fmt.Println(string(jsonData))
			return
		case "full":
			// Output detailed version information
			fmt.Printf("Current version: %s\n", version.ToString())
			fmt.Printf("Major: %d\n", version.GetMajor())
			fmt.Printf("Minor: %d\n", version.GetMinor())
			fmt.Printf("Patch: %d\n", version.GetPatch())

			preRelease := version.GetPreRelease()
			if preRelease != "" {
				fmt.Printf("Pre-release: %s\n", preRelease)
			}

			buildMeta := version.GetBuildMetadata()
			if buildMeta != "" {
				fmt.Printf("Build metadata: %s\n", buildMeta)
			}

			lastVersion := version.GetLastVersion()
			if lastVersion != "" && lastVersion != version.ToString() {
				fmt.Printf("Previous version: %s\n", lastVersion)
			}
			return
		default:
			// Default format - simple version display
			fmt.Printf("Current version: %s\n", version.ToString())
		}

		if historyFlag {
			// Get version tags from Git
			tags, err := gitops.GetTags()
			if err != nil {
				log.Fatalf("Failed to get version history: %v", err)
			}

			if len(tags) == 0 {
				fmt.Println("No version tags found in repository.")
				return
			}

			// Limit the number of tags to display if count is specified
			if countFlag > 0 && countFlag < len(tags) {
				tags = tags[:countFlag]
			}

			fmt.Println("\nVersion history:")
			fmt.Println("----------------")

			for _, tag := range tags {
				// Get tag creation time
				tagTime, err := gitops.GetTagTime(tag.Name)
				timeStr := ""
				if err == nil {
					timeStr = tagTime.Format("2006-01-02 15:04:05")
				}

				// Format the output
				if strings.HasPrefix(tag.Name, "v") {
					fmt.Printf("%-15s %-20s %s\n", tag.Name, timeStr, tag.Message)
				} else if strings.HasPrefix(tag.Name, "release-v") {
					fmt.Printf("%-15s %-20s %s [RELEASE]\n", tag.Name, timeStr, tag.Message)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&historyFlag, "history", false, "Show version history")
	versionCmd.Flags().IntVar(&countFlag, "count", 0, "Number of versions to show in history (0 = all)")
	versionCmd.Flags().StringVar(&formatFlag, "format", "", "Output format: semver, json, or full")
}
