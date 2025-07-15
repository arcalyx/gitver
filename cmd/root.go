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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.ProgrammName,
	Short: "A Git-based semantic versioning tool",
	Long: `Gitver is a command-line tool for managing semantic versioning in Git repositories.
	
It helps you automate version bumping based on conventional commits, create version tags,
and manage releases. Gitver follows the Semantic Versioning 2.0.0 specification.

Examples:
  gitver init                  # Initialize gitver in your project
  gitver bump --auto           # Automatically bump version based on commit messages
  gitver bump --major          # Bump the major version
  gitver bump --minor          # Bump the minor version
  gitver bump --patch          # Bump the patch version
  gitver release               # Create a release tag`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotver.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	projectDir, err := version.GetProjectDirectory()
	if err != nil {
		projectDir, err = os.Getwd()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	viper.SetConfigName(constants.ConfigName)
	viper.SetConfigType(constants.ConfigType)
	viper.AddConfigPath(projectDir + "/" + constants.ConfigFolderName)
	viper.SetDefault("Version", "0.0.0")

	version.SetFilePath(projectDir + "/" + constants.ConfigFolderName)
	version.SetFileName(constants.VersionFileName)

	gitops.SetRepositoryPath(projectDir)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
