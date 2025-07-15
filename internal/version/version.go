package version

import (
	"fmt"
	"github.com/arcalyx/gitver/internal/constants"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var v *Version

type Version struct {
	major               int
	minor               int
	patch               int
	preRelease          string
	preReleaseVersion   int
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

	_, err = fmt.Sscanf(string(data), "%d.%d.%d", &v.major, &v.minor, &v.patch)
	if err != nil {
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
	if v.preRelease == "" {
		return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s.%d", v.major, v.minor, v.patch, v.preRelease, v.preReleaseVersion)
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
	// Regular expression for semantic versioning with optional pre-release
	semverRegex := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z]+)\.(\d+))?$`)
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

	return nil
}

func WriteVersion() error {
	return v.WriteVersion()
}

func (v *Version) WriteVersion() error {
	versionFilePath := filepath.Join(v.versionFilePath, v.versionFileName)
	lastVersionFilePath := filepath.Join(v.versionFilePath, v.lastVersionFileName)

	if _, err := os.Stat(versionFilePath); !os.IsNotExist(err) {
		source, err := os.Open(versionFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer source.Close()

		// Create destination file
		destination, err := os.Create(lastVersionFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer destination.Close()

		// Copy content
		_, err = io.Copy(destination, source)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := os.WriteFile(versionFilePath, []byte(v.ToString()), os.ModePerm)
	if err != nil {
		return WriteOperationFailedError{versionFilePath, err}
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
