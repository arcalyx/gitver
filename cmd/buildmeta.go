/*
Copyright 2024 NAME HERE EMAIL ADDRESS
*/
package cmd

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/cobra"
	"log"
)

var (
	buildMetadataValue string
)

// buildMetaCmd represents the build-meta command
var buildMetaCmd = &cobra.Command{
	Use:   "build-meta",
	Short: "Manage build metadata for semantic versions",
	Long: `Manage build metadata for semantic versions.

Build metadata is specified by appending a plus sign and a series of dot-separated
identifiers immediately following the patch or pre-release version.
Build metadata does NOT affect version precedence.

Examples:
  gitver build-meta set build.123      # Set build metadata to "build.123"
  gitver build-meta set sha.abc123     # Set build metadata to "sha.abc123"
  gitver build-meta clear              # Clear build metadata`,
}

// setCmd represents the set command for build metadata
var setMetaCmd = &cobra.Command{
	Use:   "set [metadata]",
	Short: "Set build metadata for the current version",
	Long: `Set build metadata for the current version.

Build metadata must consist of alphanumeric characters and hyphens [0-9A-Za-z-]
and can be separated by dots.

Examples:
  gitver build-meta set build.123      # Results in version like 1.2.3+build.123
  gitver build-meta set sha.abc123     # Results in version like 1.2.3+sha.abc123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		buildMetadata := args[0]
		err := version.SetBuildMetadata(buildMetadata)
		if err != nil {
			log.Fatalf("Failed to set build metadata: %v", err)
		}

		fmt.Printf("Build metadata set to '%s'\n", buildMetadata)
		fmt.Printf("New version: %s\n", version.ToString())
	},
}

// clearMetaCmd represents the clear command for build metadata
var clearMetaCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear build metadata from the current version",
	Long: `Clear build metadata from the current version.

Example:
  gitver build-meta clear      # Removes build metadata from version`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		err := version.ClearBuildMetadata()
		if err != nil {
			log.Fatalf("Failed to clear build metadata: %v", err)
		}

		fmt.Printf("Build metadata cleared\n")
		fmt.Printf("New version: %s\n", version.ToString())
	},
}

// showMetaCmd represents the show command for build metadata
var showMetaCmd = &cobra.Command{
	Use:   "show",
	Short: "Show build metadata of the current version",
	Long: `Show build metadata of the current version.

Example:
  gitver build-meta show      # Shows current build metadata`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		buildMetadata := version.GetBuildMetadata()
		if buildMetadata == "" {
			fmt.Println("No build metadata is set")
		} else {
			fmt.Printf("Current build metadata: %s\n", buildMetadata)
		}
		fmt.Printf("Full version: %s\n", version.ToString())
	},
}

func init() {
	rootCmd.AddCommand(buildMetaCmd)
	buildMetaCmd.AddCommand(setMetaCmd)
	buildMetaCmd.AddCommand(clearMetaCmd)
	buildMetaCmd.AddCommand(showMetaCmd)
}
