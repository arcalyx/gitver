package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	preCommitFlag  bool
	prePushFlag    bool
	postMergeFlag  bool
	installAllFlag bool
	removeFlag     bool
	keywordFlag    string
)

// hooksCmd represents the hooks command
var hooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "Manage Git hooks for automatic version management",
	Long: `Install or remove Git hooks for automatic version management.
	
This command allows you to set up Git hooks that automate version management tasks:

- pre-commit: Validate version consistency before commits
- pre-push: Ensure version is properly updated before pushing
- post-merge: Update package versions after merging

You can specify a keyword that must be present in commit messages to trigger version updates.

Example:
  gitver hooks --install-all                # Install all supported hooks
  gitver hooks --pre-commit                 # Install only the pre-commit hook
  gitver hooks --pre-commit --keyword=bump  # Install hook that triggers on "bump" keyword
  gitver hooks --remove                     # Remove all gitver hooks`,
	Run: func(cmd *cobra.Command, args []string) {
		if !preCommitFlag && !prePushFlag && !postMergeFlag && !installAllFlag && !removeFlag {
			fmt.Println("Please specify which hooks to install or use --remove to uninstall hooks")
			fmt.Println("Run 'gitver hooks --help' for usage information")
			return
		}

		gitDir, err := findGitDir()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		hooksDir := filepath.Join(gitDir, "hooks")
		if _, err := os.Stat(hooksDir); os.IsNotExist(err) {
			if err := os.MkdirAll(hooksDir, 0755); err != nil {
				fmt.Printf("Error creating hooks directory: %v\n", err)
				return
			}
		}

		if removeFlag {
			removeHooks(hooksDir)
			return
		}

		if installAllFlag {
			preCommitFlag = true
			prePushFlag = true
			postMergeFlag = true
		}

		// Save keyword to a config file if specified
		if keywordFlag != "" {
			if err := saveKeywordConfig(keywordFlag); err != nil {
				fmt.Printf("Error saving keyword config: %v\n", err)
				return
			}
		}

		if preCommitFlag {
			installHook(hooksDir, "pre-commit", generatePreCommitHook(keywordFlag))
		}

		if prePushFlag {
			installHook(hooksDir, "pre-push", generatePrePushHook(keywordFlag))
		}

		if postMergeFlag {
			installHook(hooksDir, "post-merge", generatePostMergeHook(keywordFlag))
		}
	},
}

func init() {
	rootCmd.AddCommand(hooksCmd)

	hooksCmd.Flags().BoolVar(&preCommitFlag, "pre-commit", false, "Install pre-commit hook")
	hooksCmd.Flags().BoolVar(&prePushFlag, "pre-push", false, "Install pre-push hook")
	hooksCmd.Flags().BoolVar(&postMergeFlag, "post-merge", false, "Install post-merge hook")
	hooksCmd.Flags().BoolVar(&installAllFlag, "install-all", false, "Install all supported hooks")
	hooksCmd.Flags().BoolVar(&removeFlag, "remove", false, "Remove all gitver hooks")
	hooksCmd.Flags().StringVar(&keywordFlag, "keyword", "", "Keyword in commit message that triggers version updates")
}

func saveKeywordConfig(keyword string) error {
	configDir := filepath.Join(projectDir, ".gitver")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}
	}

	configPath := filepath.Join(configDir, "hooks.config")
	content := fmt.Sprintf("KEYWORD=%s\n", keyword)

	return os.WriteFile(configPath, []byte(content), 0644)
}

func findGitDir() (string, error) {
	// Start from current directory and look for .git
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			return gitDir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root
			return "", fmt.Errorf("not a git repository (or any of the parent directories)")
		}
		dir = parent
	}
}

func installHook(hooksDir, hookName, content string) {
	hookPath := filepath.Join(hooksDir, hookName)

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		// Read existing hook
		existingContent, err := os.ReadFile(hookPath)
		if err != nil {
			fmt.Printf("Error reading existing %s hook: %v\n", hookName, err)
			return
		}

		// Check if our hook is already installed
		if string(existingContent) == content {
			fmt.Printf("%s hook already installed\n", hookName)
			return
		}

		// Backup existing hook
		backupPath := hookPath + ".backup"
		if err := os.WriteFile(backupPath, existingContent, 0755); err != nil {
			fmt.Printf("Error backing up existing %s hook: %v\n", hookName, err)
			return
		}
		fmt.Printf("Backed up existing %s hook to %s\n", hookName, backupPath)
	}

	// Write new hook
	if err := os.WriteFile(hookPath, []byte(content), 0755); err != nil {
		fmt.Printf("Error installing %s hook: %v\n", hookName, err)
		return
	}

	fmt.Printf("Successfully installed %s hook\n", hookName)
}

func removeHooks(hooksDir string) {
	hooks := []string{"pre-commit", "pre-push", "post-merge"}

	for _, hook := range hooks {
		hookPath := filepath.Join(hooksDir, hook)

		// Check if hook exists
		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			continue
		}

		// Read hook content
		content, err := os.ReadFile(hookPath)
		if err != nil {
			fmt.Printf("Error reading %s hook: %v\n", hook, err)
			continue
		}

		// Check if it's our hook
		if strings.Contains(string(content), "# gitver") {

			// Check for backup
			backupPath := hookPath + ".backup"
			if _, err := os.Stat(backupPath); err == nil {
				// Restore backup
				backupContent, err := os.ReadFile(backupPath)
				if err != nil {
					fmt.Printf("Error reading backup of %s hook: %v\n", hook, err)
					continue
				}

				if err := os.WriteFile(hookPath, backupContent, 0755); err != nil {
					fmt.Printf("Error restoring backup of %s hook: %v\n", hook, err)
					continue
				}

				if err := os.Remove(backupPath); err != nil {
					fmt.Printf("Error removing backup of %s hook: %v\n", hook, err)
				}

				fmt.Printf("Restored original %s hook\n", hook)
			} else {
				// No backup, just remove
				if err := os.Remove(hookPath); err != nil {
					fmt.Printf("Error removing %s hook: %v\n", hook, err)
					continue
				}
				fmt.Printf("Removed %s hook\n", hook)
			}
		} else {
			fmt.Printf("%s hook was not installed by gitver, skipping\n", hook)
		}
	}
}

func generatePreCommitHook(keyword string) string {
	hook := `#!/bin/sh
# gitver pre-commit hook
# Validates version consistency before commit

# Get the gitver executable
GITVER=$(which gitver 2>/dev/null)
if [ -z "$GITVER" ]; then
    echo "Warning: gitver not found in PATH, skipping version validation"
    exit 0
fi

# Get commit message
COMMIT_MSG_FILE=$1
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Check for keyword in commit message if specified
KEYWORD="%s"
if [ -n "$KEYWORD" ]; then
    if ! echo "$COMMIT_MSG" | grep -q "$KEYWORD"; then
        # If keyword is not found, just exit with success
        # This allows commits without the keyword to proceed without version validation
        exit 0
    fi
    # Keyword found, proceed with version validation
    echo "Found keyword '$KEYWORD' in commit message, validating version..."
fi

# Run version validation
$GITVER validate 2>&1
if [ $? -ne 0 ]; then
    echo "Error: Version validation failed. Please fix version inconsistencies before committing."
    echo "Run 'gitver validate' for details."
    exit 1
fi

exit 0
`

	if keyword != "" {
		return fmt.Sprintf(hook, keyword)
	} else {
		return strings.ReplaceAll(hook, "KEYWORD=\"%s\"\n", "KEYWORD=\"\"\n")
	}
}

func generatePrePushHook(keyword string) string {
	hook := `#!/bin/sh
# gitver pre-push hook
# Ensures version is properly updated before pushing

# Get the gitver executable
GITVER=$(which gitver 2>/dev/null)
if [ -z "$GITVER" ]; then
    echo "Warning: gitver not found in PATH, skipping version check"
    exit 0
fi

# Check for keyword in last commit message
KEYWORD="%s"
if [ -n "$KEYWORD" ]; then
    LAST_COMMIT_MSG=$(git log -1 --pretty=%B)
    if ! echo "$LAST_COMMIT_MSG" | grep -q "$KEYWORD"; then
        # If keyword is not found, just exit with success
        # This allows pushes without the keyword to proceed without version checks
        exit 0
    fi
    # Keyword found, proceed with version check
    echo "Found keyword '$KEYWORD' in last commit message, checking version..."
fi

# Check if we need to bump version based on commits
$GITVER bump --auto --dry-run 2>&1
if [ $? -eq 0 ]; then
    echo "Warning: Commits suggest a version bump may be needed."
    echo "Consider running 'gitver bump --auto' before pushing."
    # Not exiting with error to allow push to continue
fi

exit 0
`

	if keyword != "" {
		return fmt.Sprintf(hook, keyword)
	} else {
		return strings.ReplaceAll(hook, "KEYWORD=\"%s\"\n", "KEYWORD=\"\"\n")
	}
}

func generatePostMergeHook(keyword string) string {
	hook := `#!/bin/sh
# gitver post-merge hook
# Updates package versions after merging

# Get the gitver executable
GITVER=$(which gitver 2>/dev/null)
if [ -z "$GITVER" ]; then
    echo "Warning: gitver not found in PATH, skipping version update"
    exit 0
fi

# Check for keyword in merge commit message
KEYWORD="%s"
if [ -n "$KEYWORD" ]; then
    MERGE_MSG=$(git log -1 --pretty=%B)
    if ! echo "$MERGE_MSG" | grep -q "$KEYWORD"; then
        # If keyword is not found, just exit with success
        # This allows merges without the keyword to proceed without version updates
        exit 0
    fi
    # Keyword found, proceed with version update
    echo "Found keyword '$KEYWORD' in merge commit message, updating package versions..."
fi

# Check if version file was updated in the merge
if git diff-tree -r --name-only --no-commit-id ORIG_HEAD HEAD | grep -q ".gitver/.version"; then
    echo "Version file was updated in merge, updating package versions..."
    $GITVER sync-packages
fi

exit 0
`

	if keyword != "" {
		return fmt.Sprintf(hook, keyword)
	} else {
		return strings.ReplaceAll(hook, "KEYWORD=\"%s\"\n", "KEYWORD=\"\"\n")
	}
}
