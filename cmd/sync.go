package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/arcalyx/gitver/internal/packagemanagers"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync-packages",
	Short: "Synchronize version across all package files",
	Long: `Update all package manager files with the current version from .gitver/.version.

This command updates the version in package.json, pom.xml, build.gradle, go.mod, 
setup.py, and pyproject.toml to match the version in .gitver/.version.

Example:
  gitver sync-packages`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get current version
		v := version.New()
		v.SetFilePath(filepath.Join(projectDir, ".gitver"))
		v.SetFileName(".version")
		err := v.ReadVersion()
		if err != nil {
			fmt.Printf("Error reading version: %v\n", err)
			return
		}

		versionString := v.ToString()
		fmt.Printf("Current version: %s\n", versionString)

		// Update all package manager files
		errors := packagemanagers.UpdateAllPackageManagers(projectDir, versionString)

		if len(errors) > 0 {
			fmt.Println("\nErrors occurred while updating package versions:")
			for _, err := range errors {
				fmt.Printf("  - %v\n", err)
			}
			cmd.SilenceUsage = true
			exitWithError = true
			return
		}

		fmt.Println("All package versions have been synchronized successfully!")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
