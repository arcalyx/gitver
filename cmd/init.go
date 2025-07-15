/*
Copyright 2024 NAME HERE EMAIL ADDRESS
*/
package cmd

import (
	"github.com/arcalyx/gitver/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	versionFlag string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gitver in your project",
	Long: `Initialize gitver in your project by creating the necessary configuration files.

This command creates a .gitver directory in your project with the required configuration files.
You can specify an initial version with the --version flag (defaults to 0.0.1).

Example:
  gitver init                # Initialize with default version 0.0.1
  gitver init --version 1.0.0  # Initialize with specific version`,
	Run: func(cmd *cobra.Command, args []string) {
		err := version.FromString(versionFlag)
		if err != nil {
			log.Fatalf(err.Error())
		}

		err = version.SafeWriteVersion()
		if err != nil {
			log.Fatalf(err.Error())
			return
		}

		err = viper.SafeWriteConfig()
		if err != nil {
			log.Fatalf(err.Error())
			return
		}

		log.Println("Gitver initialized for the project.")
	},
}

func init() {
	configCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&versionFlag, "version", "v", "0.0.1", "Initial version to set (default: 0.0.1)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
