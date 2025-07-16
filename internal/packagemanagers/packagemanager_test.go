package packagemanagers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNPMPackageManager(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-npm-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a package.json file
	packageJSON := map[string]interface{}{
		"name":    "test-package",
		"version": "1.0.0",
	}
	packageJSONBytes, err := json.MarshalIndent(packageJSON, "", "  ")
	require.NoError(t, err)

	packageJSONPath := filepath.Join(tempDir, "package.json")
	err = os.WriteFile(packageJSONPath, packageJSONBytes, 0644)
	require.NoError(t, err)

	// Create NPM package manager
	npm := &NPMPackageManager{}

	// Test detection
	detected := npm.Detect(tempDir)
	assert.True(t, detected, "NPM package manager should be detected")

	// Test version update
	err = npm.UpdateVersion(tempDir, "2.0.0")
	assert.NoError(t, err, "Version update should succeed")

	// Verify the version was updated
	updatedData, err := os.ReadFile(packageJSONPath)
	require.NoError(t, err)

	var updatedPackageJSON map[string]interface{}
	err = json.Unmarshal(updatedData, &updatedPackageJSON)
	require.NoError(t, err)

	assert.Equal(t, "2.0.0", updatedPackageJSON["version"], "Version should be updated to 2.0.0")
}

func TestGradlePackageManager(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-gradle-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a build.gradle file
	buildGradle := `
plugins {
    id 'java'
}

group = 'com.example'
version = '1.0.0'

repositories {
    mavenCentral()
}

dependencies {
    testImplementation 'junit:junit:4.13.2'
}
`
	buildGradlePath := filepath.Join(tempDir, "build.gradle")
	err = os.WriteFile(buildGradlePath, []byte(buildGradle), 0644)
	require.NoError(t, err)

	// Create Gradle package manager
	gradle := &GradlePackageManager{}

	// Test detection
	detected := gradle.Detect(tempDir)
	assert.True(t, detected, "Gradle package manager should be detected")

	// Test version update
	err = gradle.UpdateVersion(tempDir, "2.0.0")
	assert.NoError(t, err, "Version update should succeed")

	// Verify the version was updated
	updatedData, err := os.ReadFile(buildGradlePath)
	require.NoError(t, err)

	updatedContent := string(updatedData)
	assert.Contains(t, updatedContent, "version = '2.0.0'", "Version should be updated to 2.0.0")
}

func TestKotlinGradlePackageManager(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-kotlin-gradle-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a build.gradle.kts file
	buildGradleKts := `
plugins {
    kotlin("jvm") version "1.5.10"
}

group = "com.example"
version = "1.0.0"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(kotlin("test"))
}
`
	buildGradleKtsPath := filepath.Join(tempDir, "build.gradle.kts")
	err = os.WriteFile(buildGradleKtsPath, []byte(buildGradleKts), 0644)
	require.NoError(t, err)

	// Create Kotlin Gradle package manager
	kotlinGradle := &KotlinGradlePackageManager{}

	// Test detection
	detected := kotlinGradle.Detect(tempDir)
	assert.True(t, detected, "Kotlin Gradle package manager should be detected")

	// Test version update
	err = kotlinGradle.UpdateVersion(tempDir, "2.0.0")
	assert.NoError(t, err, "Version update should succeed")

	// Verify the version was updated
	updatedData, err := os.ReadFile(buildGradleKtsPath)
	require.NoError(t, err)

	updatedContent := string(updatedData)
	assert.Contains(t, updatedContent, "version = \"2.0.0\"", "Version should be updated to 2.0.0")
}

func TestGoPackageManager(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-go-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a go.mod file
	goMod := `module github.com/example/project

go 1.19

require (
	github.com/example/dep v1.0.0
)
`
	goModPath := filepath.Join(tempDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(goMod), 0644)
	require.NoError(t, err)

	// Create Go package manager
	goManager := &GoPackageManager{}

	// Test detection
	detected := goManager.Detect(tempDir)
	assert.True(t, detected, "Go package manager should be detected")

	// Test version update
	err = goManager.UpdateVersion(tempDir, "2.0.0")
	assert.NoError(t, err, "Version update should succeed")

	// Verify the version comment was added
	updatedData, err := os.ReadFile(goModPath)
	require.NoError(t, err)

	updatedContent := string(updatedData)
	assert.Contains(t, updatedContent, "// Version: 2.0.0", "Version comment should be added")

	// Test with major version update
	err = goManager.UpdateVersion(tempDir, "3.0.0")
	assert.NoError(t, err, "Version update should succeed")

	// Verify module path was updated for major version > 1
	updatedData, err = os.ReadFile(goModPath)
	require.NoError(t, err)

	updatedContent = string(updatedData)
	assert.Contains(t, updatedContent, "module github.com/example/project/v3", "Module path should be updated with major version")
	assert.Contains(t, updatedContent, "// Version: 3.0.0", "Version comment should be updated")
}

func TestPythonSetupPackageManager(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-python-setup-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a setup.py file
	setupPy := `
from setuptools import setup, find_packages

setup(
    name="example-project",
    version="1.0.0",
    packages=find_packages(),
)
`
	setupPyPath := filepath.Join(tempDir, "setup.py")
	err = os.WriteFile(setupPyPath, []byte(setupPy), 0644)
	require.NoError(t, err)

	// Create Python setup package manager
	pythonSetup := &PythonSetupPackageManager{}

	// Test detection
	detected := pythonSetup.Detect(tempDir)
	assert.True(t, detected, "Python setup.py package manager should be detected")

	// Test version update
	err = pythonSetup.UpdateVersion(tempDir, "2.0.0")
	assert.NoError(t, err, "Version update should succeed")

	// Verify the version was updated
	updatedData, err := os.ReadFile(setupPyPath)
	require.NoError(t, err)

	updatedContent := string(updatedData)
	assert.Contains(t, updatedContent, `version="2.0.0"`, "Version should be updated to 2.0.0")
}

func TestPyProjectTomlPackageManager(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-pyproject-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a pyproject.toml file with project section
	pyprojectToml := `
[build-system]
requires = ["setuptools>=42", "wheel"]
build-backend = "setuptools.build_meta"

[project]
name = "example-project"
version = "1.0.0"
description = "An example project"
`
	pyprojectPath := filepath.Join(tempDir, "pyproject.toml")
	err = os.WriteFile(pyprojectPath, []byte(pyprojectToml), 0644)
	require.NoError(t, err)

	// Create PyProject TOML package manager
	pyproject := &PyProjectTomlPackageManager{}

	// Test detection
	detected := pyproject.Detect(tempDir)
	assert.True(t, detected, "PyProject TOML package manager should be detected")

	// Test version update
	err = pyproject.UpdateVersion(tempDir, "2.0.0")
	assert.NoError(t, err, "Version update should succeed")

	// Verify the version was updated
	updatedData, err := os.ReadFile(pyprojectPath)
	require.NoError(t, err)

	updatedContent := string(updatedData)
	assert.Contains(t, updatedContent, `version = "2.0.0"`, "Version should be updated to 2.0.0")

	// Test with Poetry format
	poetryToml := `
[build-system]
requires = ["poetry-core>=1.0.0"]
build-backend = "poetry.core.masonry.api"

[tool.poetry]
name = "example-project"
version = "1.0.0"
description = "An example project"
`
	// Create a new temp dir for poetry test
	tempDirPoetry, err := os.MkdirTemp("", "gitver-poetry-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDirPoetry)

	poetryPath := filepath.Join(tempDirPoetry, "pyproject.toml")
	err = os.WriteFile(poetryPath, []byte(poetryToml), 0644)
	require.NoError(t, err)

	// Test detection
	detected = pyproject.Detect(tempDirPoetry)
	assert.True(t, detected, "PyProject TOML package manager should be detected for Poetry")

	// Test version update
	err = pyproject.UpdateVersion(tempDirPoetry, "3.0.0")
	assert.NoError(t, err, "Version update should succeed for Poetry")

	// Verify the version was updated
	updatedData, err = os.ReadFile(poetryPath)
	require.NoError(t, err)

	updatedContent = string(updatedData)
	assert.True(t, strings.Contains(updatedContent, `version = "3.0.0"`) ||
		strings.Contains(updatedContent, `version="3.0.0"`),
		"Version should be updated to 3.0.0 in Poetry section")
}

func TestUpdateAllPackageManagers(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gitver-all-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create package.json
	packageJSON := map[string]interface{}{
		"name":    "test-package",
		"version": "1.0.0",
	}
	packageJSONBytes, err := json.MarshalIndent(packageJSON, "", "  ")
	require.NoError(t, err)
	packageJSONPath := filepath.Join(tempDir, "package.json")
	err = os.WriteFile(packageJSONPath, packageJSONBytes, 0644)
	require.NoError(t, err)

	// Create build.gradle
	buildGradle := `version = '1.0.0'`
	buildGradlePath := filepath.Join(tempDir, "build.gradle")
	err = os.WriteFile(buildGradlePath, []byte(buildGradle), 0644)
	require.NoError(t, err)

	// Test updating all package managers
	errors := UpdateAllPackageManagers(tempDir, "3.0.0")
	assert.Empty(t, errors, "No errors should be returned")

	// Verify npm version was updated
	updatedNPMData, err := os.ReadFile(packageJSONPath)
	require.NoError(t, err)
	var updatedPackageJSON map[string]interface{}
	err = json.Unmarshal(updatedNPMData, &updatedPackageJSON)
	require.NoError(t, err)
	assert.Equal(t, "3.0.0", updatedPackageJSON["version"], "NPM version should be updated to 3.0.0")

	// Verify gradle version was updated
	updatedGradleData, err := os.ReadFile(buildGradlePath)
	require.NoError(t, err)
	updatedGradleContent := string(updatedGradleData)
	assert.Contains(t, updatedGradleContent, "version = '3.0.0'", "Gradle version should be updated to 3.0.0")
}
