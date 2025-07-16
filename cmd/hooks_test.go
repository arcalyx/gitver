package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveKeywordConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gitver-hooks-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Temporarily set projectDir to our test directory
	originalProjectDir := projectDir
	projectDir = tempDir
	defer func() {
		projectDir = originalProjectDir // Restore original value after test
	}()

	// Test saving keyword config
	keyword := "version-bump"
	err = saveKeywordConfig(keyword)
	assert.NoError(t, err, "Saving keyword config should succeed")

	// Verify config file was created with correct content
	configDir := filepath.Join(tempDir, ".gitver")
	configPath := filepath.Join(configDir, "hooks.config")

	// Check if config directory was created
	_, err = os.Stat(configDir)
	assert.NoError(t, err, "Config directory should be created")

	// Check if config file was created
	_, err = os.Stat(configPath)
	assert.NoError(t, err, "Config file should be created")

	// Check file content
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "KEYWORD=version-bump", "Config file should contain the keyword")
}

func TestFindGitDir(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "gitver-git-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a nested directory structure
	nestedDir := filepath.Join(tempDir, "level1", "level2")
	err = os.MkdirAll(nestedDir, 0755)
	require.NoError(t, err)

	// Create a .git directory at the top level
	gitDir := filepath.Join(tempDir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	require.NoError(t, err)

	// Save current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(cwd) // Restore original directory after test

	// Change to the nested directory for testing
	err = os.Chdir(nestedDir)
	require.NoError(t, err)

	// Test finding git directory from nested location
	foundGitDir, err := findGitDir()
	assert.NoError(t, err, "Finding git directory should succeed")

	// On macOS, /var is often a symlink to /private/var, so we need to evaluate the paths
	// to ensure they're comparing the same thing
	expectedPath, err := filepath.EvalSymlinks(gitDir)
	require.NoError(t, err)
	actualPath, err := filepath.EvalSymlinks(foundGitDir)
	require.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath, "Should find the correct git directory")

	// Test with non-git directory
	nonGitDir, err := os.MkdirTemp("", "non-git-dir")
	require.NoError(t, err)
	defer os.RemoveAll(nonGitDir)

	err = os.Chdir(nonGitDir)
	require.NoError(t, err)

	_, err = findGitDir()
	assert.Error(t, err, "Should return error when not in a git repository")
	assert.Contains(t, err.Error(), "not a git repository", "Error message should indicate not a git repository")
}

func TestInstallHook(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-hooks-install-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create hooks directory
	hooksDir := filepath.Join(tempDir, "hooks")
	err = os.MkdirAll(hooksDir, 0755)
	require.NoError(t, err)

	// Test installing a new hook
	hookName := "pre-commit"
	hookContent := "#!/bin/sh\necho 'Test hook'"

	installHook(hooksDir, hookName, hookContent)

	// Verify hook was installed
	hookPath := filepath.Join(hooksDir, hookName)
	_, err = os.Stat(hookPath)
	assert.NoError(t, err, "Hook file should be created")

	content, err := os.ReadFile(hookPath)
	require.NoError(t, err)
	assert.Equal(t, hookContent, string(content), "Hook content should match")

	// Test installing over existing hook (should create backup)
	existingHookContent := "#!/bin/sh\necho 'Existing hook'"
	err = os.WriteFile(hookPath, []byte(existingHookContent), 0755)
	require.NoError(t, err)

	newHookContent := "#!/bin/sh\necho 'New hook'"
	installHook(hooksDir, hookName, newHookContent)

	// Verify backup was created
	backupPath := hookPath + ".backup"
	_, err = os.Stat(backupPath)
	assert.NoError(t, err, "Backup file should be created")

	backupContent, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, existingHookContent, string(backupContent), "Backup content should match original hook")

	// Verify hook was updated
	updatedContent, err := os.ReadFile(hookPath)
	require.NoError(t, err)
	assert.Equal(t, newHookContent, string(updatedContent), "Hook should be updated with new content")
}

func TestRemoveHooks(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-hooks-remove-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create hooks directory
	hooksDir := filepath.Join(tempDir, "hooks")
	err = os.MkdirAll(hooksDir, 0755)
	require.NoError(t, err)

	// Create gitver hooks
	hooks := []string{"pre-commit", "pre-push", "post-merge"}
	for _, hook := range hooks {
		hookPath := filepath.Join(hooksDir, hook)
		hookContent := "#!/bin/sh\n# gitver\necho 'Test hook'"
		err = os.WriteFile(hookPath, []byte(hookContent), 0755)
		require.NoError(t, err)
	}

	// Create a non-gitver hook
	nonGitverHook := filepath.Join(hooksDir, "pre-rebase")
	err = os.WriteFile(nonGitverHook, []byte("#!/bin/sh\necho 'Non-gitver hook'"), 0755)
	require.NoError(t, err)

	// Create backups for some hooks
	backupHook := hooks[0]
	backupPath := filepath.Join(hooksDir, backupHook+".backup")
	backupContent := "#!/bin/sh\necho 'Original hook'"
	err = os.WriteFile(backupPath, []byte(backupContent), 0755)
	require.NoError(t, err)

	// Test removing hooks
	removeHooks(hooksDir)

	// Verify gitver hooks were removed or restored
	for _, hook := range hooks {
		hookPath := filepath.Join(hooksDir, hook)

		if hook == backupHook {
			// This hook should be restored from backup
			_, err = os.Stat(hookPath)
			assert.NoError(t, err, "Hook should exist (restored from backup)")

			content, err := os.ReadFile(hookPath)
			require.NoError(t, err)
			assert.Equal(t, backupContent, string(content), "Hook should be restored from backup")

			// Backup should be deleted
			_, err = os.Stat(backupPath)
			assert.True(t, os.IsNotExist(err), "Backup file should be deleted")
		} else {
			// Other hooks should be removed
			_, err = os.Stat(hookPath)
			assert.True(t, os.IsNotExist(err), "Hook should be removed")
		}
	}

	// Verify non-gitver hook was not removed
	_, err = os.Stat(nonGitverHook)
	assert.NoError(t, err, "Non-gitver hook should not be removed")
}

func TestGenerateHooks(t *testing.T) {
	// Test generating hooks with and without keywords

	// Test pre-commit hook
	preCommitWithKeyword := generatePreCommitHook("version")
	assert.Contains(t, preCommitWithKeyword, "KEYWORD=\"version\"", "Pre-commit hook should contain keyword")

	preCommitWithoutKeyword := generatePreCommitHook("")
	assert.Contains(t, preCommitWithoutKeyword, "KEYWORD=\"\"", "Pre-commit hook should have empty keyword")

	// Test pre-push hook
	prePushWithKeyword := generatePrePushHook("bump")
	assert.Contains(t, prePushWithKeyword, "KEYWORD=\"bump\"", "Pre-push hook should contain keyword")

	prePushWithoutKeyword := generatePrePushHook("")
	assert.Contains(t, prePushWithoutKeyword, "KEYWORD=\"\"", "Pre-push hook should have empty keyword")

	// Test post-merge hook
	postMergeWithKeyword := generatePostMergeHook("release")
	assert.Contains(t, postMergeWithKeyword, "KEYWORD=\"release\"", "Post-merge hook should contain keyword")

	postMergeWithoutKeyword := generatePostMergeHook("")
	assert.Contains(t, postMergeWithoutKeyword, "KEYWORD=\"\"", "Post-merge hook should have empty keyword")

	// Verify all hooks contain the gitver marker
	hooks := []string{
		preCommitWithKeyword, preCommitWithoutKeyword,
		prePushWithKeyword, prePushWithoutKeyword,
		postMergeWithKeyword, postMergeWithoutKeyword,
	}

	for _, hook := range hooks {
		assert.Contains(t, hook, "# gitver", "Hook should contain gitver marker")
		assert.True(t, strings.HasPrefix(hook, "#!/bin/sh"), "Hook should start with shebang")
	}
}
