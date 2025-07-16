package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/arcalyx/gitver/internal/packagemanagers"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate version consistency across package files",
	Long: `Validate that the version in .gitver/.version matches all package manager files.

This command checks that the version in .gitver/.version is consistent with 
versions in package.json, pom.xml, build.gradle, go.mod, setup.py, and pyproject.toml.

Example:
  gitver validate`,
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

		// Check package manager files
		inconsistencies := validatePackageVersions(projectDir, versionString)

		if len(inconsistencies) > 0 {
			fmt.Println("\nVersion inconsistencies found:")
			for _, inc := range inconsistencies {
				fmt.Printf("  - %s\n", inc)
			}
			fmt.Println("\nRun 'gitver sync-packages' to update all package files.")
			cmd.SilenceUsage = true
			exitWithError = true
			return
		}

		fmt.Println("All package versions are consistent! ")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func validatePackageVersions(projectDir, versionString string) []string {
	var inconsistencies []string

	// Check each package manager
	for _, pm := range packagemanagers.PackageManagers {
		if pm.Detect(projectDir) {
			// Get the file path based on package manager
			var filePath string
			switch pm.Name() {
			case "Maven":
				filePath = filepath.Join(projectDir, "pom.xml")
			case "NPM":
				filePath = filepath.Join(projectDir, "package.json")
			case "Gradle":
				filePath = filepath.Join(projectDir, "build.gradle")
			case "Kotlin Gradle":
				filePath = filepath.Join(projectDir, "build.gradle.kts")
			case "Go":
				filePath = filepath.Join(projectDir, "go.mod")
			case "Python setup.py":
				filePath = filepath.Join(projectDir, "setup.py")
			case "Python pyproject.toml":
				filePath = filepath.Join(projectDir, "pyproject.toml")
			default:
				continue
			}

			// Check if version in file matches current version
			// This is a simplified check - in a real implementation,
			// we would need to parse each file format properly
			if !packagemanagers.CheckVersionInFile(filePath, versionString) {
				inconsistencies = append(inconsistencies, fmt.Sprintf("%s (%s)", pm.Name(), filepath.Base(filePath)))
			}
		}
	}

	return inconsistencies
}
