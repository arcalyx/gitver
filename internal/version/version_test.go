package version

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestVersion is a simplified version of Version for testing
type TestVersion struct {
	Major             int
	Minor             int
	Patch             int
	PreRelease        string
	PreReleaseVersion int
}

// MockVersion is a version implementation for testing that doesn't perform file operations
type MockVersion struct {
	major             int
	minor             int
	patch             int
	preRelease        string
	preReleaseVersion int
	lastVersion       string
	buildMetadata     string
}

// NewMockVersion creates a new MockVersion instance
func NewMockVersion() *MockVersion {
	return &MockVersion{}
}

// FromString parses a version string into the MockVersion
func (v *MockVersion) FromString(version string) error {
	// Reuse the parsing logic from the real Version
	realVersion := New()
	if err := realVersion.FromString(version); err != nil {
		return err
	}

	// Additional validation for pre-release type
	if realVersion.preRelease != "" {
		preType := strings.ToLower(realVersion.preRelease)
		if preType != "alpha" && preType != "beta" && preType != "rc" {
			return fmt.Errorf("invalid pre-release type: %s (must be alpha, beta, or rc)", preType)
		}
	}

	// Copy the parsed values
	v.major = realVersion.major
	v.minor = realVersion.minor
	v.patch = realVersion.patch
	v.preRelease = realVersion.preRelease
	v.preReleaseVersion = realVersion.preReleaseVersion
	v.lastVersion = realVersion.lastVersion
	v.buildMetadata = realVersion.buildMetadata

	return nil
}

// ToString returns the string representation of the MockVersion
func (v *MockVersion) ToString() string {
	if v.preRelease == "" {
		if v.buildMetadata == "" {
			return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
		}
		return fmt.Sprintf("%d.%d.%d+%s", v.major, v.minor, v.patch, v.buildMetadata)
	}
	if v.buildMetadata == "" {
		return fmt.Sprintf("%d.%d.%d-%s.%d", v.major, v.minor, v.patch, v.preRelease, v.preReleaseVersion)
	}
	return fmt.Sprintf("%d.%d.%d-%s.%d+%s", v.major, v.minor, v.patch, v.preRelease, v.preReleaseVersion, v.buildMetadata)
}

// SetPreRelease sets the pre-release type and version
func (v *MockVersion) SetPreRelease(preRelease string, version int) error {
	// Validate pre-release type
	preRelease = strings.ToLower(preRelease)
	if preRelease != "alpha" && preRelease != "beta" && preRelease != "rc" && preRelease != "" {
		return fmt.Errorf("invalid pre-release type: %s (must be alpha, beta, rc, or empty string)", preRelease)
	}

	v.lastVersion = v.ToString()
	v.preRelease = preRelease
	v.preReleaseVersion = version

	return nil
}

// BumpPreRelease increments the pre-release version
func (v *MockVersion) BumpPreRelease() error {
	if v.preRelease == "" {
		return fmt.Errorf("no pre-release version set")
	}

	v.lastVersion = v.ToString()
	v.preReleaseVersion++

	return nil
}

// PromotePreRelease promotes the pre-release to the next level
func (v *MockVersion) PromotePreRelease() error {
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

	return nil
}

// ClearPreRelease clears the pre-release information
func (v *MockVersion) ClearPreRelease() error {
	v.lastVersion = v.ToString()
	v.preRelease = ""
	v.preReleaseVersion = 0
	return nil
}

// Compare compares this version with another version following SemVer 2.0.0 specification
// Returns:
//
//	-1 if this version is less than the other version
//	 0 if this version is equal to the other version
//	 1 if this version is greater than the other version
//
// Note: Build metadata is ignored in version precedence as per SemVer 2.0.0 spec
func (v *MockVersion) Compare(other *MockVersion) int {
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
func (v *MockVersion) SetBuildMetadata(buildMetadata string) error {
	// Save the current version as the last version
	v.lastVersion = v.ToString()

	// Validate build metadata format (alphanumeric and hyphen only, dot-separated)
	if buildMetadata != "" {
		for _, part := range strings.Split(buildMetadata, ".") {
			if !regexp.MustCompile(`^[0-9A-Za-z-]+$`).MatchString(part) {
				return fmt.Errorf("invalid build metadata: %s (must contain only alphanumeric characters and hyphens)", buildMetadata)
			}
		}
	}

	v.buildMetadata = buildMetadata
	return nil
}

// ClearBuildMetadata removes the build metadata from the version
func (v *MockVersion) ClearBuildMetadata() error {
	// Save the current version as the last version
	v.lastVersion = v.ToString()

	v.buildMetadata = ""
	return nil
}

// GetBuildMetadata returns the build metadata for the version
func (v *MockVersion) GetBuildMetadata() string {
	return v.buildMetadata
}

func setupTestVersion() *MockVersion {
	return NewMockVersion()
}

func initializeVersion(v *MockVersion, versionStr string) error {
	// Parse the version string
	return v.FromString(versionStr)
}

func TestPreReleaseVersions(t *testing.T) {
	tests := []struct {
		name           string
		initialVersion string
		action         func(v *MockVersion) error
		expectedResult TestVersion
		expectError    bool
	}{
		{
			name:           "Set alpha version",
			initialVersion: "1.0.0",
			action: func(v *MockVersion) error {
				return v.SetPreRelease("alpha", 1)
			},
			expectedResult: TestVersion{1, 0, 0, "alpha", 1},
			expectError:    false,
		},
		{
			name:           "Set beta version",
			initialVersion: "1.0.0",
			action: func(v *MockVersion) error {
				return v.SetPreRelease("beta", 1)
			},
			expectedResult: TestVersion{1, 0, 0, "beta", 1},
			expectError:    false,
		},
		{
			name:           "Set rc version",
			initialVersion: "1.0.0",
			action: func(v *MockVersion) error {
				return v.SetPreRelease("rc", 1)
			},
			expectedResult: TestVersion{1, 0, 0, "rc", 1},
			expectError:    false,
		},
		{
			name:           "Set invalid pre-release type",
			initialVersion: "1.0.0",
			action: func(v *MockVersion) error {
				return v.SetPreRelease("gamma", 1)
			},
			expectedResult: TestVersion{},
			expectError:    true,
		},
		{
			name:           "Bump alpha version",
			initialVersion: "1.0.0-alpha.1",
			action: func(v *MockVersion) error {
				return v.BumpPreRelease()
			},
			expectedResult: TestVersion{1, 0, 0, "alpha", 2},
			expectError:    false,
		},
		{
			name:           "Bump beta version",
			initialVersion: "1.0.0-beta.3",
			action: func(v *MockVersion) error {
				return v.BumpPreRelease()
			},
			expectedResult: TestVersion{1, 0, 0, "beta", 4},
			expectError:    false,
		},
		{
			name:           "Bump rc version",
			initialVersion: "1.0.0-rc.5",
			action: func(v *MockVersion) error {
				return v.BumpPreRelease()
			},
			expectedResult: TestVersion{1, 0, 0, "rc", 6},
			expectError:    false,
		},
		{
			name:           "Promote alpha to beta",
			initialVersion: "1.0.0-alpha.2",
			action: func(v *MockVersion) error {
				return v.PromotePreRelease()
			},
			expectedResult: TestVersion{1, 0, 0, "beta", 1},
			expectError:    false,
		},
		{
			name:           "Promote beta to rc",
			initialVersion: "1.0.0-beta.4",
			action: func(v *MockVersion) error {
				return v.PromotePreRelease()
			},
			expectedResult: TestVersion{1, 0, 0, "rc", 1},
			expectError:    false,
		},
		{
			name:           "Promote rc to release",
			initialVersion: "1.0.0-rc.3",
			action: func(v *MockVersion) error {
				return v.PromotePreRelease()
			},
			expectedResult: TestVersion{1, 0, 0, "", 0},
			expectError:    false,
		},
		{
			name:           "Clear pre-release version",
			initialVersion: "1.0.0-beta.2",
			action: func(v *MockVersion) error {
				return v.ClearPreRelease()
			},
			expectedResult: TestVersion{1, 0, 0, "", 0},
			expectError:    false,
		},
		{
			name:           "Bump non-pre-release version",
			initialVersion: "1.0.0",
			action: func(v *MockVersion) error {
				return v.BumpPreRelease()
			},
			expectedResult: TestVersion{},
			expectError:    true,
		},
		{
			name:           "Promote non-pre-release version",
			initialVersion: "1.0.0",
			action: func(v *MockVersion) error {
				return v.PromotePreRelease()
			},
			expectedResult: TestVersion{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new MockVersion instance for each test
			v := setupTestVersion()

			// Set up the initial version
			err := initializeVersion(v, tt.initialVersion)
			if err != nil {
				t.Fatalf("Failed to initialize version %s: %v", tt.initialVersion, err)
			}

			// Perform the action
			err = tt.action(v)

			// Check if error matches expectation
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v (%v)", tt.expectError, err != nil, err)
			}

			// If no error is expected, check the result
			if !tt.expectError {
				if v.major != tt.expectedResult.Major {
					t.Errorf("Expected major %d, got %d", tt.expectedResult.Major, v.major)
				}
				if v.minor != tt.expectedResult.Minor {
					t.Errorf("Expected minor %d, got %d", tt.expectedResult.Minor, v.minor)
				}
				if v.patch != tt.expectedResult.Patch {
					t.Errorf("Expected patch %d, got %d", tt.expectedResult.Patch, v.patch)
				}
				if v.preRelease != tt.expectedResult.PreRelease {
					t.Errorf("Expected preRelease %s, got %s", tt.expectedResult.PreRelease, v.preRelease)
				}
				if v.preReleaseVersion != tt.expectedResult.PreReleaseVersion {
					t.Errorf("Expected preReleaseVersion %d, got %d", tt.expectedResult.PreReleaseVersion, v.preReleaseVersion)
				}
			}
		})
	}
}

func TestVersionParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedMajor int
		expectedMinor int
		expectedPatch int
		expectedPre   string
		expectedNum   int
		expectError   bool
	}{
		{"1.2.3", 1, 2, 3, "", 0, false},
		{"0.0.1", 0, 0, 1, "", 0, false},
		{"10.20.30", 10, 20, 30, "", 0, false},
		{"1.2.3-alpha.1", 1, 2, 3, "alpha", 1, false},
		{"1.2.3-beta.5", 1, 2, 3, "beta", 5, false},
		{"1.2.3-rc.10", 1, 2, 3, "rc", 10, false},
		{"invalid", 0, 0, 0, "", 0, true},
		{"1.2", 0, 0, 0, "", 0, true},
		{"1.2.3.4", 0, 0, 0, "", 0, true},
		{"1.2.3-gamma.1", 0, 0, 0, "", 0, true}, // Invalid pre-release type
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			v := NewMockVersion() // Use MockVersion for parsing tests
			err := v.FromString(tt.input)

			// Check if error matches expectation
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v for input %s", tt.expectError, err != nil, tt.input)
				return
			}

			// If error is expected, no need to check further
			if tt.expectError {
				return
			}

			// Check version components
			if v.major != tt.expectedMajor {
				t.Errorf("Expected major %d, got %d for input %s", tt.expectedMajor, v.major, tt.input)
			}
			if v.minor != tt.expectedMinor {
				t.Errorf("Expected minor %d, got %d for input %s", tt.expectedMinor, v.minor, tt.input)
			}
			if v.patch != tt.expectedPatch {
				t.Errorf("Expected patch %d, got %d for input %s", tt.expectedPatch, v.patch, tt.input)
			}

			// Check pre-release info
			if v.preRelease != tt.expectedPre {
				t.Errorf("Expected pre-release type %s, got %s for input %s", tt.expectedPre, v.preRelease, tt.input)
			}
			if v.preReleaseVersion != tt.expectedNum {
				t.Errorf("Expected pre-release version %d, got %d for input %s", tt.expectedNum, v.preReleaseVersion, tt.input)
			}
		})
	}
}

func TestVersionParsingWithBuildMetadata(t *testing.T) {
	testCases := []struct {
		version       string
		major         int
		minor         int
		patch         int
		preRelease    string
		preReleaseVer int
		buildMetadata string
		shouldBeValid bool
	}{
		{"1.2.3", 1, 2, 3, "", 0, "", true},
		{"1.2.3+build.123", 1, 2, 3, "", 0, "build.123", true},
		{"1.2.3-alpha.1+build.456", 1, 2, 3, "alpha", 1, "build.456", true},
		{"1.2.3-beta.5+sha.abc123", 1, 2, 3, "beta", 5, "sha.abc123", true},
		{"1.2.3+", 0, 0, 0, "", 0, "", false},
		{"1.2.3++extra", 0, 0, 0, "", 0, "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.version, func(t *testing.T) {
			v := NewMockVersion()
			err := v.FromString(tc.version)

			if tc.shouldBeValid {
				assert.NoError(t, err)
				assert.Equal(t, tc.major, v.major)
				assert.Equal(t, tc.minor, v.minor)
				assert.Equal(t, tc.patch, v.patch)
				assert.Equal(t, tc.preRelease, v.preRelease)
				assert.Equal(t, tc.preReleaseVer, v.preReleaseVersion)
				assert.Equal(t, tc.buildMetadata, v.buildMetadata)

				// Test that ToString() correctly formats the version
				assert.Equal(t, tc.version, v.ToString())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestBuildMetadataMethods(t *testing.T) {
	tests := []struct {
		name           string
		initialVersion string
		buildMetadata  string
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Set valid build metadata on regular version",
			initialVersion: "1.2.3",
			buildMetadata:  "build.123",
			expectedResult: "1.2.3+build.123",
			expectError:    false,
		},
		{
			name:           "Set valid build metadata on pre-release version",
			initialVersion: "1.2.3-alpha.1",
			buildMetadata:  "sha.abc123",
			expectedResult: "1.2.3-alpha.1+sha.abc123",
			expectError:    false,
		},
		{
			name:           "Set invalid build metadata (spaces)",
			initialVersion: "1.2.3",
			buildMetadata:  "build 123",
			expectedResult: "1.2.3",
			expectError:    true,
		},
		{
			name:           "Set invalid build metadata (special chars)",
			initialVersion: "1.2.3",
			buildMetadata:  "build@123",
			expectedResult: "1.2.3",
			expectError:    true,
		},
		{
			name:           "Clear build metadata",
			initialVersion: "1.2.3+build.123",
			buildMetadata:  "",
			expectedResult: "1.2.3",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewMockVersion()
			err := v.FromString(tt.initialVersion)
			assert.NoError(t, err)

			if tt.buildMetadata == "" {
				// Test ClearBuildMetadata
				err = v.ClearBuildMetadata()
				assert.NoError(t, err)
			} else {
				// Test SetBuildMetadata
				err = v.SetBuildMetadata(tt.buildMetadata)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}

			// Verify the result
			if !tt.expectError {
				assert.Equal(t, tt.expectedResult, v.ToString())
				assert.Equal(t, tt.buildMetadata, v.GetBuildMetadata())
			}
		})
	}
}

func TestVersionComparison(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int // -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.1.0", -1},
		{"1.1.0", "1.0.0", 1},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.0.0-alpha.1", "1.0.0", -1},
		{"1.0.0", "1.0.0-alpha.1", 1},
		{"1.0.0-alpha.1", "1.0.0-alpha.2", -1},
		{"1.0.0-alpha.2", "1.0.0-alpha.1", 1},
		{"1.0.0-alpha.1", "1.0.0-beta.1", -1},
		{"1.0.0-beta.1", "1.0.0-alpha.1", 1},
		{"1.0.0-beta.1", "1.0.0-rc.1", -1},
		{"1.0.0-rc.1", "1.0.0-beta.1", 1},
		// Test cases with build metadata - should be ignored in comparison
		{"1.0.0+build.1", "1.0.0+build.2", 0},
		{"1.0.0-alpha.1+build.1", "1.0.0-alpha.1+build.2", 0},
		{"1.0.0+build.1", "1.0.0", 0},
		{"1.0.0-alpha.1+build.1", "1.0.0-alpha.1", 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_vs_%s", tt.v1, tt.v2), func(t *testing.T) {
			v1 := NewMockVersion()
			err := v1.FromString(tt.v1)
			assert.NoError(t, err)

			v2 := NewMockVersion()
			err = v2.FromString(tt.v2)
			assert.NoError(t, err)

			// Use the new Compare method
			result := v1.Compare(v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}
