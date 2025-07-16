/*
Copyright 2024 NAME HERE EMAIL ADDRESS
*/
package cmd

import (
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
  gitver version --history    # Show version history
  gitver version --history --count 5  # Show last 5 versions
  gitver version --format=semver  # Show only the semver string`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		// Display current version based on format
		if formatFlag == "semver" {
			fmt.Print(version.ToString())
			return
		}

		fmt.Printf("Current version: %s\n", version.ToString())

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
	versionCmd.Flags().StringVar(&formatFlag, "format", "", "Output format (semver)")
}
