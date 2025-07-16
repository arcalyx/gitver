package packagemanagers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/arcalyx/gitver/internal/constants"
	"github.com/arcalyx/gitver/internal/xml"
	"github.com/pelletier/go-toml/v2"
)

// PackageManager is an interface for different package managers
type PackageManager interface {
	UpdateVersion(projectDir, newVersion string) error
	Detect(projectDir string) bool
	Name() string
}

// PackageManagers is a slice of all available package managers
var PackageManagers []PackageManager

// Initialize all package managers
func init() {
	PackageManagers = []PackageManager{
		&MavenPackageManager{},
		&NPMPackageManager{},
		&GradlePackageManager{},
		&KotlinGradlePackageManager{},
		&GoPackageManager{},
		&PythonSetupPackageManager{},
		&PyProjectTomlPackageManager{},
	}
}

// UpdateAllPackageManagers updates the version in all detected package managers
func UpdateAllPackageManagers(projectDir, newVersion string) []error {
	var errors []error
	var updatedManagers []string

	for _, pm := range PackageManagers {
		if pm.Detect(projectDir) {
			if err := pm.UpdateVersion(projectDir, newVersion); err != nil {
				errors = append(errors, fmt.Errorf(constants.PackageUpdateFailed, pm.Name(), err))
			} else {
				updatedManagers = append(updatedManagers, pm.Name())
			}
		}
	}

	// Print success messages for updated package managers
	for _, name := range updatedManagers {
		fmt.Printf(constants.PackageUpdateSuccess+"\n", name, newVersion)
	}

	return errors
}

// MavenPackageManager handles version updates for Maven projects
type MavenPackageManager struct{}

// Name returns the name of the package manager
func (m *MavenPackageManager) Name() string {
	return "Maven"
}

// Detect checks if the project is a Maven project
func (m *MavenPackageManager) Detect(projectDir string) bool {
	pomPath := filepath.Join(projectDir, constants.MavenPomFile)
	_, err := os.Stat(pomPath)
	return err == nil
}

// UpdateVersion updates the version in pom.xml
func (m *MavenPackageManager) UpdateVersion(projectDir, newVersion string) error {
	pomPath := filepath.Join(projectDir, constants.MavenPomFile)
	return xml.SetVersion(pomPath, "//project/version", newVersion)
}

// NPMPackageManager handles version updates for npm projects
type NPMPackageManager struct{}

// Name returns the name of the package manager
func (n *NPMPackageManager) Name() string {
	return "NPM"
}

// Detect checks if the project is an npm project
func (n *NPMPackageManager) Detect(projectDir string) bool {
	packagePath := filepath.Join(projectDir, constants.NPMPackageFile)
	_, err := os.Stat(packagePath)
	return err == nil
}

// UpdateVersion updates the version in package.json
func (n *NPMPackageManager) UpdateVersion(projectDir, newVersion string) error {
	packagePath := filepath.Join(projectDir, constants.NPMPackageFile)

	// Read package.json
	data, err := os.ReadFile(packagePath)
	if err != nil {
		return fmt.Errorf("error reading package.json: %w", err)
	}

	// Parse JSON
	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return fmt.Errorf("error parsing package.json: %w", err)
	}

	// Update version
	packageJSON["version"] = newVersion

	// Write back to file
	updatedData, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing package.json: %w", err)
	}

	if err := os.WriteFile(packagePath, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing package.json: %w", err)
	}

	return nil
}

// GradlePackageManager handles version updates for Gradle projects
type GradlePackageManager struct{}

// Name returns the name of the package manager
func (g *GradlePackageManager) Name() string {
	return "Gradle"
}

// Detect checks if the project is a Gradle project
func (g *GradlePackageManager) Detect(projectDir string) bool {
	buildGradlePath := filepath.Join(projectDir, constants.GradleBuildFile)
	_, err := os.Stat(buildGradlePath)
	return err == nil
}

// UpdateVersion updates the version in build.gradle
func (g *GradlePackageManager) UpdateVersion(projectDir, newVersion string) error {
	buildGradlePath := filepath.Join(projectDir, constants.GradleBuildFile)

	// Read build.gradle
	data, err := os.ReadFile(buildGradlePath)
	if err != nil {
		return fmt.Errorf("error reading build.gradle: %w", err)
	}

	content := string(data)

	// Use regular expressions for more reliable pattern matching
	versionPattern := regexp.MustCompile(`version\s*=\s*['"]([^'"]*)['"]`)
	versionSetPattern := regexp.MustCompile(`version\s+['"]([^'"]*)['"]`)

	if versionPattern.MatchString(content) {
		content = versionPattern.ReplaceAllString(content, fmt.Sprintf(`version = '%s'`, newVersion))
	} else if versionSetPattern.MatchString(content) {
		content = versionSetPattern.ReplaceAllString(content, fmt.Sprintf(`version '%s'`, newVersion))
	} else {
		// If no pattern matched, try to add version to the top of the file
		content = fmt.Sprintf("version = '%s'\n\n%s", newVersion, content)
	}

	// Write back to file
	if err := os.WriteFile(buildGradlePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing build.gradle: %w", err)
	}

	return nil
}

// KotlinGradlePackageManager handles version updates for Kotlin Gradle projects
type KotlinGradlePackageManager struct{}

// Name returns the name of the package manager
func (k *KotlinGradlePackageManager) Name() string {
	return "Kotlin Gradle"
}

// Detect checks if the project is a Kotlin Gradle project
func (k *KotlinGradlePackageManager) Detect(projectDir string) bool {
	buildGradlePath := filepath.Join(projectDir, constants.KotlinBuildFile)
	_, err := os.Stat(buildGradlePath)
	return err == nil
}

// UpdateVersion updates the version in build.gradle.kts
func (k *KotlinGradlePackageManager) UpdateVersion(projectDir, newVersion string) error {
	buildGradlePath := filepath.Join(projectDir, constants.KotlinBuildFile)

	// Read build.gradle.kts
	data, err := os.ReadFile(buildGradlePath)
	if err != nil {
		return fmt.Errorf("error reading build.gradle.kts: %w", err)
	}

	content := string(data)

	// Use regular expressions for more reliable pattern matching
	versionPattern := regexp.MustCompile(`version\s*=\s*"([^"]*)"`)
	versionFuncPattern := regexp.MustCompile(`version\("([^"]*)"\)`)

	if versionPattern.MatchString(content) {
		content = versionPattern.ReplaceAllString(content, fmt.Sprintf(`version = "%s"`, newVersion))
	} else if versionFuncPattern.MatchString(content) {
		content = versionFuncPattern.ReplaceAllString(content, fmt.Sprintf(`version("%s")`, newVersion))
	} else {
		// If no pattern matched, try to add version to the top of the file
		content = fmt.Sprintf("version = \"%s\"\n\n%s", newVersion, content)
	}

	// Write back to file
	if err := os.WriteFile(buildGradlePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing build.gradle.kts: %w", err)
	}

	return nil
}

// GoPackageManager handles version updates for Go modules
type GoPackageManager struct{}

// Name returns the name of the package manager
func (g *GoPackageManager) Name() string {
	return "Go"
}

// Detect checks if the project is a Go module
func (g *GoPackageManager) Detect(projectDir string) bool {
	goModPath := filepath.Join(projectDir, constants.GoModFile)
	_, err := os.Stat(goModPath)
	return err == nil
}

// UpdateVersion updates the version in go.mod
// Note: Go modules don't typically include version in go.mod, but we can add a comment
// or update the module path if it includes a version suffix
func (g *GoPackageManager) UpdateVersion(projectDir, newVersion string) error {
	goModPath := filepath.Join(projectDir, constants.GoModFile)

	// Read go.mod
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("error reading go.mod: %w", err)
	}

	content := string(data)

	// Try to find and update the module path with version suffix
	// Example: module github.com/user/repo/v1 -> module github.com/user/repo/v2
	modulePattern := regexp.MustCompile(`(?m)^module\s+([^\s]+)(?:/v\d+)?`)

	if modulePattern.MatchString(content) {
		// Extract major version from the new version
		versionPattern := regexp.MustCompile(`^(\d+)\.`)
		matches := versionPattern.FindStringSubmatch(newVersion)

		if len(matches) > 1 {
			majorVersion := matches[1]

			// Only update module path if major version > 1
			if majorVersion != "0" && majorVersion != "1" {
				// Get the base module path without version suffix
				moduleMatches := modulePattern.FindStringSubmatch(content)
				modulePath := moduleMatches[1]

				// Remove any existing version suffix
				moduleBase := regexp.MustCompile(`^(.*?)(?:/v\d+)?$`).FindStringSubmatch(modulePath)[1]

				// Replace module declaration with new version suffix
				newModuleLine := fmt.Sprintf("module %s/v%s", moduleBase, majorVersion)
				content = modulePattern.ReplaceAllString(content, newModuleLine)
			}
		}
	}

	// Add or update a version comment at the top of the file
	versionCommentPattern := regexp.MustCompile(`(?m)^// Version: .*$`)
	versionComment := fmt.Sprintf("// Version: %s", newVersion)

	if versionCommentPattern.MatchString(content) {
		content = versionCommentPattern.ReplaceAllString(content, versionComment)
	} else {
		// Add version comment at the top of the file
		content = versionComment + "\n" + content
	}

	// Write back to file
	if err := os.WriteFile(goModPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing go.mod: %w", err)
	}

	return nil
}

// PythonSetupPackageManager handles version updates for Python setup.py
type PythonSetupPackageManager struct{}

// Name returns the name of the package manager
func (p *PythonSetupPackageManager) Name() string {
	return "Python setup.py"
}

// Detect checks if the project has a setup.py file
func (p *PythonSetupPackageManager) Detect(projectDir string) bool {
	setupPath := filepath.Join(projectDir, constants.PythonSetupFile)
	_, err := os.Stat(setupPath)
	return err == nil
}

// UpdateVersion updates the version in setup.py
func (p *PythonSetupPackageManager) UpdateVersion(projectDir, newVersion string) error {
	setupPath := filepath.Join(projectDir, constants.PythonSetupFile)

	// Read setup.py
	data, err := os.ReadFile(setupPath)
	if err != nil {
		return fmt.Errorf("error reading setup.py: %w", err)
	}

	content := string(data)

	// Try to find and replace version in different formats
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)version\s*=\s*['"]([^'"]*)['"]`),
		regexp.MustCompile(`(?m)version\s*=\s*"([^"]*)"`),
		regexp.MustCompile(`(?m)version\s*=\s*'([^']*)'`),
	}

	replaced := false
	for _, pattern := range patterns {
		if pattern.MatchString(content) {
			content = pattern.ReplaceAllStringFunc(content, func(match string) string {
				// Preserve the original quote style
				if strings.Contains(match, "\"") {
					return fmt.Sprintf("version=\"%s\"", newVersion)
				} else {
					return fmt.Sprintf("version='%s'", newVersion)
				}
			})
			replaced = true
			break
		}
	}

	if !replaced {
		// If no pattern matched, try to add version after setup(
		setupPattern := regexp.MustCompile(`setup\(\s*`)
		if setupPattern.MatchString(content) {
			content = setupPattern.ReplaceAllString(content, fmt.Sprintf("setup(\n    version='%s',\n    ", newVersion))
		} else {
			return fmt.Errorf("could not find a suitable location to add version in setup.py")
		}
	}

	// Write back to file
	if err := os.WriteFile(setupPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing setup.py: %w", err)
	}

	return nil
}

// PyProjectTomlPackageManager handles version updates for Python pyproject.toml
type PyProjectTomlPackageManager struct{}

// Name returns the name of the package manager
func (p *PyProjectTomlPackageManager) Name() string {
	return "Python pyproject.toml"
}

// Detect checks if the project has a pyproject.toml file
func (p *PyProjectTomlPackageManager) Detect(projectDir string) bool {
	pyprojectPath := filepath.Join(projectDir, constants.PyProjectTomlFile)
	_, err := os.Stat(pyprojectPath)
	return err == nil
}

// UpdateVersion updates the version in pyproject.toml
func (p *PyProjectTomlPackageManager) UpdateVersion(projectDir, newVersion string) error {
	pyprojectPath := filepath.Join(projectDir, constants.PyProjectTomlFile)

	// Read pyproject.toml
	data, err := os.ReadFile(pyprojectPath)
	if err != nil {
		return fmt.Errorf("error reading pyproject.toml: %w", err)
	}

	// First try to update using regex for more reliability
	content := string(data)

	// Look for version in [project] section
	projectVersionPattern := regexp.MustCompile(`(?m)(\[project\]\s*(?:[^\[]*\s*)?version\s*=\s*)["']([^"']*)["']`)
	poetryVersionPattern := regexp.MustCompile(`(?m)(\[tool\.poetry\]\s*(?:[^\[]*\s*)?version\s*=\s*)["']([^"']*)["']`)

	if projectVersionPattern.MatchString(content) {
		content = projectVersionPattern.ReplaceAllString(content, fmt.Sprintf(`$1"%s"`, newVersion))
	} else if poetryVersionPattern.MatchString(content) {
		content = poetryVersionPattern.ReplaceAllString(content, fmt.Sprintf(`$1"%s"`, newVersion))
	} else {
		// If regex approach fails, try parsing the TOML
		var pyprojectData map[string]interface{}
		if err := toml.Unmarshal(data, &pyprojectData); err != nil {
			return fmt.Errorf("error parsing pyproject.toml: %w", err)
		}

		updated := false

		// Check for project section
		if project, ok := pyprojectData["project"].(map[string]interface{}); ok {
			project["version"] = newVersion
			updated = true
		}

		// Check for tool.poetry section
		if tool, ok := pyprojectData["tool"].(map[string]interface{}); ok {
			if poetry, ok := tool["poetry"].(map[string]interface{}); ok {
				poetry["version"] = newVersion
				updated = true
			}
		}

		if !updated {
			return fmt.Errorf("could not find version field in pyproject.toml")
		}

		// Marshal back to TOML
		updatedData, err := toml.Marshal(pyprojectData)
		if err != nil {
			return fmt.Errorf("error serializing pyproject.toml: %w", err)
		}

		content = string(updatedData)
	}

	// Write back to file
	if err := os.WriteFile(pyprojectPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing pyproject.toml: %w", err)
	}

	return nil
}

// CheckVersionInFile checks if the specified version exists in the given file
// This is a simplified check that looks for the version string in the file
func CheckVersionInFile(filePath, version string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	content := string(data)

	// For different file types, we might need different patterns
	// This is a simplified approach that just checks if the version string exists
	// A more robust implementation would parse each file format properly

	// Check for version with quotes (for JSON, TOML, etc.)
	if strings.Contains(content, fmt.Sprintf("\"%s\"", version)) {
		return true
	}

	// Check for version with single quotes (for some build files)
	if strings.Contains(content, fmt.Sprintf("'%s'", version)) {
		return true
	}

	// Check for version as a comment (for Go modules)
	if strings.Contains(content, fmt.Sprintf("// Version: %s", version)) {
		return true
	}

	// Check for raw version string (fallback)
	if strings.Contains(content, version) {
		return true
	}

	return false
}
