/*
Copyright 2024 NAME HERE EMAIL ADDRESS
*/
package cmd

import (
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/cobra"
	"log"
)

var (
	alphaFlag         bool
	betaFlag          bool
	rcFlag            bool
	clearFlag         bool
	promoteFlag       bool
	bumpFlag          bool
	versionNumberFlag int
)

// prereleaseCmd represents the prerelease command
var prereleaseCmd = &cobra.Command{
	Use:   "prerelease",
	Short: "Manage pre-release versions",
	Long: `Manage pre-release versions (alpha, beta, rc) for your project.

This command allows you to set, bump, promote, or clear pre-release versions.
Pre-release versions follow the format: MAJOR.MINOR.PATCH-TYPE.NUMBER
Where TYPE is one of: alpha, beta, rc

Examples:
  gitver prerelease --alpha           # Set version to alpha.1
  gitver prerelease --beta --version 2 # Set version to beta.2
  gitver prerelease --rc              # Set version to rc.1
  gitver prerelease --bump            # Increment the pre-release number
  gitver prerelease --promote         # Promote alpha to beta, beta to rc, or rc to release
  gitver prerelease --clear           # Remove pre-release info (convert to stable release)`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		// Count the number of flags set
		flagCount := 0
		if alphaFlag {
			flagCount++
		}
		if betaFlag {
			flagCount++
		}
		if rcFlag {
			flagCount++
		}
		if clearFlag {
			flagCount++
		}
		if promoteFlag {
			flagCount++
		}
		if bumpFlag {
			flagCount++
		}

		// Ensure only one action flag is set
		if flagCount != 1 {
			log.Fatalf("Please specify exactly one action: --alpha, --beta, --rc, --bump, --promote, or --clear")
		}

		var err error

		switch {
		case alphaFlag:
			if cmd.Flags().Changed("version") {
				err = version.SetPreRelease("alpha", versionNumberFlag)
			} else {
				err = version.SetPreRelease("alpha", 1)
			}
		case betaFlag:
			if cmd.Flags().Changed("version") {
				err = version.SetPreRelease("beta", versionNumberFlag)
			} else {
				err = version.SetPreRelease("beta", 1)
			}
		case rcFlag:
			if cmd.Flags().Changed("version") {
				err = version.SetPreRelease("rc", versionNumberFlag)
			} else {
				err = version.SetPreRelease("rc", 1)
			}
		case bumpFlag:
			err = version.BumpPreRelease()
		case promoteFlag:
			err = version.PromotePreRelease()
		case clearFlag:
			err = version.ClearPreRelease()
		}

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		log.Printf("Version updated: %s -> %s", version.GetLastVersion(), version.ToString())

		// Execute Git operations if requested
		if gitFlag == CommitTag || gitFlag == CommitTagPush {
			executeGitOperations()
		}
	},
}

func init() {
	rootCmd.AddCommand(prereleaseCmd)

	prereleaseCmd.Flags().BoolVar(&alphaFlag, "alpha", false, "Set version to alpha")
	prereleaseCmd.Flags().BoolVar(&betaFlag, "beta", false, "Set version to beta")
	prereleaseCmd.Flags().BoolVar(&rcFlag, "rc", false, "Set version to release candidate")
	prereleaseCmd.Flags().BoolVar(&bumpFlag, "bump", false, "Bump the pre-release version number")
	prereleaseCmd.Flags().BoolVar(&promoteFlag, "promote", false, "Promote to next pre-release level (alpha->beta->rc->release)")
	prereleaseCmd.Flags().BoolVar(&clearFlag, "clear", false, "Clear pre-release information (convert to stable release)")
	prereleaseCmd.Flags().IntVar(&versionNumberFlag, "version", 0, "Specify pre-release version number (default: 1)")
	prereleaseCmd.Flags().StringVar(&gitFlag, "git", "", "Git operations to perform after updating. Valid values: COMMIT_TAG, COMMIT_TAG_PUSH")
}
