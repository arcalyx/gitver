package cmd

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var (
	getFlag        bool
	listFlag       bool
	detectProjects bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "View or set configuration values",
	Long: `View or set configuration values for gitver.

This command allows you to get, list, or set configuration values that control
how gitver operates. Configuration is stored in the .gitver directory in your project.

Examples:
  gitver config --list                # List all configuration values
  gitver config --get user.name       # Get a specific configuration value
  gitver config user.name "Alex"      # Set a configuration value
  gitver config init                  # Initialize configuration directory
  gitver config init --detect-projects # Initialize and detect subprojects`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		// Handle --list flag
		if listFlag {
			listAllConfig()
			return
		}

		// Handle --get flag
		if getFlag {
			if len(args) != 1 {
				log.Fatalf("Error: --get requires exactly one key argument")
			}
			getConfigValue(args[0])
			return
		}

		// Handle set operation
		if len(args) == 2 {
			setConfigValue(args[0], args[1])
			return
		}

		// If no arguments provided, show help
		if len(args) == 0 && !listFlag && !getFlag {
			if err := cmd.Help(); err != nil {
				log.Fatalf("Error displaying help: %v", err)
			}
			return
		}

		fmt.Printf("Error: Invalid arguments. Run 'gitver config --help' for usage.\n")
	},
}

// initCmd represents the init subcommand
var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gitver configuration",
	Long: `Initialize gitver configuration directory and create default config file.

This command creates a .gitver directory in the project root and writes a
config.yaml with default values. With --detect-projects flag, it will also
scan the repository for subprojects and add them to the configuration.

Examples:
  gitver config init
  gitver config init --detect-projects`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create .gitver directory if it doesn't exist
		gitverDir := filepath.Join(projectDir, ".gitver")
		if _, err := os.Stat(gitverDir); os.IsNotExist(err) {
			if err := os.MkdirAll(gitverDir, 0755); err != nil {
				log.Fatalf("Error creating .gitver directory: %v", err)
			}
			fmt.Println("Created .gitver directory")
		}

		// Create default config file
		configPath := filepath.Join(gitverDir, "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			defaultConfig := `version: 1.0.0
versioning:
  auto_detect: true
  default_bump: patch
  commit_keywords:
    major: ["#breaking", "BREAKING CHANGE"]
    minor: ["feat", "feature"]
    patch: ["fix", "bugfix"]
changelog:
  format: markdown
  output_file: CHANGELOG.md
hooks:
  pre_commit: true
  pre_push: true
  post_merge: true
`
			if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
				log.Fatalf("Error writing default config: %v", err)
			}
			fmt.Println("Created default config.yaml")
		} else {
			fmt.Println("Config file already exists, skipping creation")
		}

		// Detect subprojects if requested
		if detectProjects {
			fmt.Println("Detecting subprojects...")
			detectSubprojects()
		}

		fmt.Println("Configuration initialized successfully!")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Add subcommands
	configCmd.AddCommand(initConfigCmd)

	// Add flags to config command
	configCmd.Flags().BoolVar(&getFlag, "get", false, "Get the value for a key")
	configCmd.Flags().BoolVar(&listFlag, "list", false, "List all config entries")

	// Add flags to init subcommand
	initConfigCmd.Flags().BoolVar(&detectProjects, "detect-projects", false, "Automatically detect subprojects")
}

// listAllConfig prints all configuration values
func listAllConfig() {
	settings := viper.AllSettings()
	if len(settings) == 0 {
		fmt.Println("No configuration values found")
		return
	}

	fmt.Println("Configuration values:")
	printConfigMap("", settings)
}

// printConfigMap recursively prints a map of configuration values
func printConfigMap(prefix string, configMap map[string]interface{}) {
	for key, value := range configMap {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if subMap, ok := value.(map[string]interface{}); ok {
			printConfigMap(fullKey, subMap)
		} else {
			fmt.Printf("%s = %v\n", fullKey, value)
		}
	}
}

// getConfigValue retrieves and prints a specific configuration value
func getConfigValue(key string) {
	value := viper.Get(key)
	if value == nil {
		fmt.Printf("No value found for key: %s\n", key)
		return
	}

	fmt.Printf("%s = %v\n", key, value)
}

// setConfigValue sets a configuration value and saves it
func setConfigValue(key, value string) {
	viper.Set(key, value)

	// Make sure the config directory exists
	configDir := filepath.Join(projectDir, constants.ConfigFolderName)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatalf("Error creating config directory: %v", err)
		}
	}

	// Save the configuration
	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("Error writing config: %v", err)
	}

	fmt.Printf("Set %s = %s\n", key, value)
}

// detectSubprojects scans the repository for known project types
func detectSubprojects() {
	// Get current subprojects configuration
	var subprojects []map[string]string
	if viper.IsSet("subprojects") {
		if err := viper.UnmarshalKey("subprojects", &subprojects); err != nil {
			log.Printf("Error reading existing subprojects: %v", err)
			subprojects = []map[string]string{}
		}
	}

	// Define project types to detect
	projectPatterns := map[string]string{
		"npm":    "package.json",
		"maven":  "pom.xml",
		"gradle": "build.gradle",
		"poetry": "pyproject.toml",
		"python": "setup.py",
		"go":     "go.mod",
	}

	// Scan for subprojects
	fmt.Println("Scanning for subprojects...")

	// Walk through the project directory
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Skip .gitver directory
		if info.IsDir() && info.Name() == constants.ConfigFolderName {
			return filepath.SkipDir
		}

		// Check if the file matches any project pattern
		for projectType, pattern := range projectPatterns {
			if filepath.Base(path) == pattern {
				// Convert to relative path
				relPath, err := filepath.Rel(projectDir, path)
				if err != nil {
					return nil
				}

				// Check if this project is already in the list
				exists := false
				for _, proj := range subprojects {
					if proj["path"] == relPath {
						exists = true
						break
					}
				}

				// Add if not exists
				if !exists {
					subprojects = append(subprojects, map[string]string{
						"path": relPath,
						"type": projectType,
					})
					fmt.Printf("Detected %s project: %s\n", projectType, relPath)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error scanning for subprojects: %v", err)
	}

	// Save detected subprojects to config
	viper.Set("subprojects", subprojects)
	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("Error writing config: %v", err)
	}
}
