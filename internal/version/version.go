package version

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/constants"
	"github.com/arcalyx/gitver/internal/packagemanagers"
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var v *Version
var updatePackagesFlag = true

// Version represents a semantic version
type Version struct {
	major               int
	minor               int
	patch               int
	preRelease          string
	preReleaseVersion   int
	buildMetadata       string
	versionFilePath     string
	versionFileName     string
	lastVersionFileName string
	lastVersion         string
	fs                  afero.Fs
}

func init() {
	v = New()
}

func New() *Version {
	v := new(Version)
	v.major = 0
	v.minor = 0
	v.patch = 0
	v.preRelease = ""
	v.preReleaseVersion = 0
	v.buildMetadata = ""
	v.versionFilePath = ""
	v.versionFileName = ""
	v.lastVersionFileName = ".lastversion"
	v.lastVersion = "0.0.0"
	v.fs = afero.NewOsFs()
	return v
}

// GetProjectDirectory returns the directory containing the .gitver folder.
func GetProjectDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", UnhandledError(err.Error())
	}

	for {
		// Check if the current directory contains the .gitver folder
		if _, err := os.Stat(filepath.Join(dir, constants.ConfigFolderName)); !os.IsNotExist(err) {
			return filepath.Clean(dir), nil
		}

		// Get the parent directory
		parentDir := filepath.Dir(dir)

		// If the current directory is the same as the parent directory,
		// we've reached the root directory
		if parentDir == dir {
			break
		}

		dir = parentDir
	}

	return "", ProjectDirectoryNotFoundError(dir)
}

func ReadVersion() error {
	return v.ReadVersion()
}

func (v *Version) ReadVersion() error {
	versionFilePath := filepath.Join(v.versionFilePath, v.versionFileName)
	lastVersionFilePath := filepath.Join(v.versionFilePath, v.lastVersionFileName)

	data, err := os.ReadFile(versionFilePath)
	if err != nil {
		return FileNotFoundError(versionFilePath)
	}

	// Use FromString to parse the version with proper regex validation
	versionStr := strings.TrimSpace(string(data))
	if err := v.FromString(versionStr); err != nil {
		return FileFormatError(versionFilePath)
	}

	if _, err := os.Stat(lastVersionFilePath); os.IsNotExist(err) {
		return nil
	}

	data, err = os.ReadFile(lastVersionFilePath)
	if err != nil {
		return FileNotFoundError(lastVersionFilePath)
	}

	_, err = fmt.Sscan(string(data), &v.lastVersion)
	if err != nil {
		return FileFormatError(versionFilePath)
	}

	return nil
}

func SetFilePath(filePath string) {
	v.SetFilePath(filePath)
}

func (v *Version) SetFilePath(filePath string) {
	v.versionFilePath = filePath
}

func SetFileName(fileName string) {
	v.SetFileName(fileName)
}

func (v *Version) SetFileName(fileName string) {
	v.versionFileName = fileName
}

func SafeWriteVersion() error {
	return v.SafeWriteVersion()
}

func (v *Version) SafeWriteVersion() error {
	dir := filepath.Join(v.versionFilePath)
	versionFilePath := filepath.Join(v.versionFilePath, v.versionFileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	alreadyExists, err := afero.Exists(v.fs, versionFilePath)
	if alreadyExists && err == nil {
		return FileAlreadyExistsError(versionFilePath)
	}

	if err := v.WriteVersion(); err != nil {
		return err
	}

	return nil
}

func ToString() string {
	return v.ToString()
}

func (v *Version) ToString() string {
	if v.preRelease == "" && v.buildMetadata == "" {
		return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	} else if v.preRelease != "" && v.buildMetadata == "" {
		return fmt.Sprintf("%d.%d.%d-%s.%d", v.major, v.minor, v.patch, v.preRelease, v.preReleaseVersion)
	} else if v.preRelease == "" && v.buildMetadata != "" {
		return fmt.Sprintf("%d.%d.%d+%s", v.major, v.minor, v.patch, v.buildMetadata)
	} else {
		return fmt.Sprintf("%d.%d.%d-%s.%d+%s", v.major, v.minor, v.patch, v.preRelease, v.preReleaseVersion, v.buildMetadata)
	}
}

func GetLastVersion() string {
	return v.GetLastVersion()
}

func (v *Version) GetLastVersion() string {
	return v.lastVersion
}

func FromString(version string) error {
	return v.FromString(version)
}

func (v *Version) FromString(version string) error {
	// Regular expression for semantic versioning with optional pre-release and build metadata
	semverRegex := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z]+)\.(\d+))?(?:\+([0-9A-Za-z-.]+))?$`)
	matches := semverRegex.FindStringSubmatch(version)

	if matches == nil {
		return InputValueError(version)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return InputValueError(version)
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return InputValueError(version)
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return InputValueError(version)
	}

	v.major = major
	v.minor = minor
	v.patch = patch

	// Check if pre-release info is present
	if len(matches) > 4 && matches[4] != "" {
		v.preRelease = strings.ToLower(matches[4])
		preReleaseVersion, err := strconv.Atoi(matches[5])
		if err != nil {
			return InputValueError(version)
		}
		v.preReleaseVersion = preReleaseVersion
	} else {
		v.preRelease = ""
		v.preReleaseVersion = 0
	}

	// Check if build metadata is present
	if len(matches) > 6 && matches[6] != "" {
		v.buildMetadata = matches[6]
	} else {
		v.buildMetadata = ""
	}

	return nil
}

func WriteVersion() error {
	return v.WriteVersion()
}

func (v *Version) WriteVersion() error {
	versionFilePath := filepath.Join(v.versionFilePath, v.versionFileName)
	lastVersionFilePath := filepath.Join(v.versionFilePath, v.lastVersionFileName)

	// Create the version file if it doesn't exist
	if _, err := v.fs.Stat(versionFilePath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := v.fs.MkdirAll(filepath.Dir(versionFilePath), 0755); err != nil {
			return err
		}

		// Create the file
		file, err := v.fs.Create(versionFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	// Write the version to the file
	file, err := v.fs.OpenFile(versionFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s", v.ToString())
	if err != nil {
		return err
	}

	// Create the last version file if it doesn't exist
	if _, err := v.fs.Stat(lastVersionFilePath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := v.fs.MkdirAll(filepath.Dir(lastVersionFilePath), 0755); err != nil {
			return err
		}

		// Create the file
		file, err := v.fs.Create(lastVersionFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	// Write the last version to the file
	file, err = v.fs.OpenFile(lastVersionFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s", v.lastVersion)
	if err != nil {
		return err
	}

	// Update versions in package manager files if the flag is set
	if updatePackagesFlag {
		versionString := v.ToString()
		projectDir, err := GetProjectDirectory()
		if err == nil {
			errors := packagemanagers.UpdateAllPackageManagers(projectDir, versionString)
			if len(errors) > 0 {
				// Log errors but don't fail the version update
				for _, err := range errors {
					log.Printf("Warning: Failed to update package manager file: %v", err)
				}
			}
		}
	}

	return nil
}

func BumpMajor() error {
	return v.BumpMajor()
}

func (v *Version) BumpMajor() error {
	v.lastVersion = v.ToString()
	v.major++
	v.minor = 0
	v.patch = 0

	return v.WriteVersion()
}

func BumpMinor() error {
	return v.BumpMinor()
}

func (v *Version) BumpMinor() error {
	v.lastVersion = v.ToString()
	v.minor++
	v.patch = 0

	return v.WriteVersion()
}

func BumpPatch() error {
	return v.BumpPatch()
}

func (v *Version) BumpPatch() error {
	v.lastVersion = v.ToString()
	v.patch++

	return v.WriteVersion()
}

// Pre-release version management functions
func SetPreRelease(preRelease string, version int) error {
	return v.SetPreRelease(preRelease, version)
}

func (v *Version) SetPreRelease(preRelease string, version int) error {
	// Validate pre-release type
	preRelease = strings.ToLower(preRelease)
	if preRelease != "alpha" && preRelease != "beta" && preRelease != "rc" && preRelease != "" {
		return fmt.Errorf("invalid pre-release type: %s (must be alpha, beta, rc, or empty string)", preRelease)
	}

	v.lastVersion = v.ToString()
	v.preRelease = preRelease
	v.preReleaseVersion = version

	return v.WriteVersion()
}

func BumpPreRelease() error {
	return v.BumpPreRelease()
}

func (v *Version) BumpPreRelease() error {
	if v.preRelease == "" {
		return fmt.Errorf("no pre-release version set")
	}

	v.lastVersion = v.ToString()
	v.preReleaseVersion++

	return v.WriteVersion()
}

func ClearPreRelease() error {
	return v.ClearPreRelease()
}

func (v *Version) ClearPreRelease() error {
	if v.preRelease == "" {
		return nil // Already cleared
	}

	v.lastVersion = v.ToString()
	v.preRelease = ""
	v.preReleaseVersion = 0

	return v.WriteVersion()
}

func PromotePreRelease() error {
	return v.PromotePreRelease()
}

func (v *Version) PromotePreRelease() error {
	if v.preRelease == "" {
		return fmt.Errorf("no pre-release version set")
	}

	v.lastVersion = v.ToString()

	// Promote from alpha to beta, beta to rc, or rc to release
	switch v.preRelease {
	case "alpha":
		v.preRelease = "beta"
		v.preReleaseVersion = 1
	case "beta":
		v.preRelease = "rc"
		v.preReleaseVersion = 1
	case "rc":
		v.preRelease = ""
		v.preReleaseVersion = 0
	}

	return v.WriteVersion()
}

func SetUpdatePackages(update bool) {
	updatePackagesFlag = update
}

func GetUpdatePackages() bool {
	return updatePackagesFlag
}

// Compare compares this version with another version following SemVer 2.0.0 specification
// Returns:
//
//	-1 if this version is less than the other version
//	 0 if this version is equal to the other version
//	 1 if this version is greater than the other version
//
// Note: Build metadata is ignored in version precedence as per SemVer 2.0.0 spec
func (v *Version) Compare(other *Version) int {
	// Compare major version
	if v.major != other.major {
		if v.major > other.major {
			return 1
		}
		return -1
	}

	// Compare minor version
	if v.minor != other.minor {
		if v.minor > other.minor {
			return 1
		}
		return -1
	}

	// Compare patch version
	if v.patch != other.patch {
		if v.patch > other.patch {
			return 1
		}
		return -1
	}

	// If one has a pre-release and the other doesn't, the one without is greater
	if v.preRelease == "" && other.preRelease != "" {
		return 1
	}
	if v.preRelease != "" && other.preRelease == "" {
		return -1
	}

	// If both have pre-releases, compare them
	if v.preRelease != "" && other.preRelease != "" {
		// Compare pre-release types (alpha < beta < rc)
		if v.preRelease != other.preRelease {
			if v.preRelease == "alpha" {
				return -1
			}
			if other.preRelease == "alpha" {
				return 1
			}
			if v.preRelease == "beta" {
				return -1
			}
			if other.preRelease == "beta" {
				return 1
			}
		}

		// Compare pre-release versions
		if v.preReleaseVersion != other.preReleaseVersion {
			if v.preReleaseVersion > other.preReleaseVersion {
				return 1
			}
			return -1
		}
	}

	// Versions are equal (build metadata doesn't affect precedence)
	return 0
}

// SetBuildMetadata sets the build metadata for the version
// Build metadata is specified by appending a plus sign and a series of dot-separated identifiers
// immediately following the patch or pre-release version
func (v *Version) SetBuildMetadata(buildMetadata string) error {
	// Save the current version as the last version
	v.lastVersion = v.ToString()

	// Validate build metadata format (alphanumeric and hyphen only, dot-separated)
	if buildMetadata != "" {
		for _, part := range strings.Split(buildMetadata, ".") {
			if !regexp.MustCompile(`^[0-9A-Za-z-]+$`).MatchString(part) {
				return InputValueError(fmt.Sprintf("invalid build metadata: %s (must contain only alphanumeric characters and hyphens)", buildMetadata))
			}
		}
	}

	v.buildMetadata = buildMetadata
	return v.WriteVersion()
}

// ClearBuildMetadata removes the build metadata from the version
func (v *Version) ClearBuildMetadata() error {
	// Save the current version as the last version
	v.lastVersion = v.ToString()

	v.buildMetadata = ""
	return v.WriteVersion()
}

// GetBuildMetadata returns the build metadata for the version
func (v *Version) GetBuildMetadata() string {
	return v.buildMetadata
}

func SetBuildMetadata(buildMetadata string) error {
	return v.SetBuildMetadata(buildMetadata)
}

func ClearBuildMetadata() error {
	return v.ClearBuildMetadata()
}

func GetBuildMetadata() string {
	return v.GetBuildMetadata()
}

// GetFileSystem returns the current file system used by the version package
func GetFileSystem() afero.Fs {
	return v.fs
}

// SetFileSystem sets the file system to use for the version package
func SetFileSystem(fs afero.Fs) {
	v.fs = fs
}
