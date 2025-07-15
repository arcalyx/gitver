/*
Copyright 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/constants"
	"github.com/arcalyx/gitver/internal/gitops"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/cobra"
	"log"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Create a release tag for the current version",
	Long: `Create a release tag for the current version of your project.

This command creates a Git tag with the format 'release-vX.Y.Z' for the current version.
Use this command when you want to mark a specific version as a release.

Example:
  gitver release  # Creates a tag like 'release-v1.2.3' for the current version`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()
		err := prepareGitOperation()
		if err != nil {
			log.Fatalf("Failed to prepare Git operation: %v", err)
		}

		tagName := fmt.Sprintf(constants.ReleaseTag, version.ToString())
		if err := gitops.CreateTag(tagName, constants.TagMessage); err != nil {
			log.Fatalf("Failed to create release tag: %v", err)
		}

		log.Printf("Tag '%s' created successfully", tagName)

		if pushFlag {
			if err := gitops.Push(); err != nil {
				log.Fatalf("Failed to push tag: %v", err)
			}
			log.Printf("Tag '%s' pushed successfully", tagName)
		}
	},
}

var pushFlag bool

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Add a flag to push the tag after creation
	releaseCmd.Flags().BoolVarP(&pushFlag, "push", "p", false, "Push the tag after creation")
}
