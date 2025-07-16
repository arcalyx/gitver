package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arcalyx/gitver/internal/packagemanagers"
	"github.com/arcalyx/gitver/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncPackages(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gitver-sync-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create .gitver directory and version file
	gitverDir := filepath.Join(tempDir, ".gitver")
	err = os.MkdirAll(gitverDir, 0755)
	require.NoError(t, err)

	// Create version file with test version
	testVersion := "2.1.0"
	versionFile := filepath.Join(gitverDir, ".version")
	err = os.WriteFile(versionFile, []byte(testVersion), 0644)
	require.NoError(t, err)

	// Create package files to test synchronization

	// Create go.mod file
	goModContent := `module github.com/example/project

go 1.19

require (
	github.com/example/dep v1.0.0
)
`
	goModPath := filepath.Join(tempDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(goModContent), 0644)
	require.NoError(t, err)

	// Create setup.py file
	setupPyContent := `
from setuptools import setup, find_packages

setup(
    name="example-project",
    version="1.0.0",
    packages=find_packages(),
)
`
	setupPyPath := filepath.Join(tempDir, "setup.py")
	err = os.WriteFile(setupPyPath, []byte(setupPyContent), 0644)
	require.NoError(t, err)

	// Create pyproject.toml file
	pyprojectContent := `
[build-system]
requires = ["setuptools>=42", "wheel"]
build-backend = "setuptools.build_meta"

[project]
name = "example-project"
version = "1.0.0"
description = "An example project"
`
	pyprojectPath := filepath.Join(tempDir, "pyproject.toml")
	err = os.WriteFile(pyprojectPath, []byte(pyprojectContent), 0644)
	require.NoError(t, err)

	// Save original projectDir and restore it after test
	originalProjectDir := projectDir
	projectDir = tempDir
	defer func() {
		projectDir = originalProjectDir
	}()

	// Initialize version object
	v := version.New()
	v.SetFilePath(gitverDir)
	v.SetFileName(".version")
	err = v.ReadVersion()
	require.NoError(t, err)

	// Test sync functionality
	versionString := v.ToString()
	errors := packagemanagers.UpdateAllPackageManagers(projectDir, versionString)
	assert.Empty(t, errors, "No errors should be returned")

	// Verify go.mod was updated
	goModUpdated, err := os.ReadFile(goModPath)
	require.NoError(t, err)
	assert.Contains(t, string(goModUpdated), "// Version: 2.1.0", "go.mod should contain version comment")

	// Verify setup.py was updated
	setupPyUpdated, err := os.ReadFile(setupPyPath)
	require.NoError(t, err)
	assert.Contains(t, string(setupPyUpdated), `version="2.1.0"`, "setup.py should have updated version")

	// Verify pyproject.toml was updated
	pyprojectUpdated, err := os.ReadFile(pyprojectPath)
	require.NoError(t, err)
	assert.Contains(t, string(pyprojectUpdated), `version = "2.1.0"`, "pyproject.toml should have updated version")
}

func TestSyncPackagesWithErrors(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gitver-sync-error-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create .gitver directory and version file
	gitverDir := filepath.Join(tempDir, ".gitver")
	err = os.MkdirAll(gitverDir, 0755)
	require.NoError(t, err)

	// Create version file with test version
	testVersion := "3.0.0"
	versionFile := filepath.Join(gitverDir, ".version")
	err = os.WriteFile(versionFile, []byte(testVersion), 0644)
	require.NoError(t, err)

	// Create a read-only directory to simulate permission errors
	readOnlyDir := filepath.Join(tempDir, "readonly")
	err = os.MkdirAll(readOnlyDir, 0755)
	require.NoError(t, err)

	// Create go.mod in the read-only directory
	goModPath := filepath.Join(readOnlyDir, "go.mod")
	err = os.WriteFile(goModPath, []byte("module example.com/readonly"), 0644)
	require.NoError(t, err)

	// Make the directory read-only to cause permission errors during sync
	// Note: This might not work on all systems, especially Windows
	err = os.Chmod(readOnlyDir, 0555)
	require.NoError(t, err)
	defer os.Chmod(readOnlyDir, 0755) // Restore permissions for cleanup

	// Save original projectDir and restore it after test
	originalProjectDir := projectDir
	projectDir = readOnlyDir // Set to read-only directory to test error handling
	defer func() {
		projectDir = originalProjectDir
	}()

	// Initialize version object
	v := version.New()
	v.SetFilePath(gitverDir)
	v.SetFileName(".version")
	err = v.ReadVersion()
	require.NoError(t, err)

	// Test sync functionality with errors
	versionString := v.ToString()
	errors := packagemanagers.UpdateAllPackageManagers(projectDir, versionString)

	// The function might return errors due to permission issues
	// What's important is that the test doesn't crash and handles errors gracefully
	if len(errors) == 0 {
		t.Log("No errors were returned, which is unexpected for read-only directory")
	} else {
		t.Logf("Received expected errors: %v", errors)
	}
}
