package cmd

import (
	"github.com/arcalyx/gitver/internal/constants"
	"github.com/arcalyx/gitver/internal/gitops"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/viper"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	exitWithError bool   // Used to indicate if the program should exit with an error code
	projectDir    string // Project directory path, used across the cmd package
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.ProgramName,
	Short: "A Git-based semantic versioning tool",
	Long: `Gitver is a command-line tool for managing semantic versioning in Git repositories.
	
It helps you automate version bumping based on conventional commits, create version tags,
and manage releases. Gitver follows the Semantic Versioning 2.0.0 specification.

Examples:
  gitver init                       # Initialize gitver in your project
  gitver bump --auto                # Automatically determine version bump based on commit messages
  gitver bump --infer               # Infer version bump from commit keywords in config
  gitver bump --major               # Bump the major version
  gitver bump --minor               # Bump the minor version
  gitver bump --patch               # Bump the patch version
  gitver bump --changelog           # Generate a changelog based on commits
  gitver bump --update-packages     # Update version in package manager files
  gitver version                    # Display the current version
  gitver config init                # Initialize configuration
  gitver hooks install              # Install git hooks`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil || exitWithError {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	var err error
	projectDir, err = version.GetProjectDirectory()
	if err != nil {
		projectDir, err = os.Getwd()
		if err != nil {
			log.Fatal("Error determining current directory:", err)
		}
	}

	viper.SetConfigName(constants.ConfigName)
	viper.SetConfigType(constants.ConfigType)
	viper.AddConfigPath(projectDir + "/" + constants.ConfigFolderName)

	// Set default configuration values
	viper.SetDefault("Version", "0.0.0")
	viper.SetDefault("changelog.format", "markdown")
	viper.SetDefault("changelog.output_file", "CHANGELOG.md")
	viper.SetDefault("versioning.commit_keywords.major", []string{"BREAKING CHANGE", "major"})
	viper.SetDefault("versioning.commit_keywords.minor", []string{"feat", "feature", "minor"})
	viper.SetDefault("versioning.commit_keywords.patch", []string{"fix", "patch"})

	version.SetFilePath(projectDir + "/" + constants.ConfigFolderName)
	version.SetFileName(constants.VersionFileName)

	gitops.SetRepositoryPath(projectDir)
}
