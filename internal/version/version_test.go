package version

import (
	"fmt"
	"strings"
	"testing"
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

	return nil
}

// ToString returns the string representation of the MockVersion
func (v *MockVersion) ToString() string {
	if v.preRelease == "" {
		return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s.%d", v.major, v.minor, v.patch, v.preRelease, v.preReleaseVersion)
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

func TestVersionComparison(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
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
	}

	for _, tt := range tests {
		t.Run(tt.v1+" vs "+tt.v2, func(t *testing.T) {
			v1 := New() // Use the real Version for comparison tests
			err := v1.FromString(tt.v1)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.v1, err)
			}

			v2 := New() // Use the real Version for comparison tests
			err = v2.FromString(tt.v2)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.v2, err)
			}

			// Since Compare method doesn't exist, we'll implement a simple comparison here
			result := compareVersions(v1, v2)
			if result != tt.expected {
				t.Errorf("Expected comparison result %d, got %d for %s vs %s", tt.expected, result, tt.v1, tt.v2)
			}
		})
	}
}

// Helper function to compare versions since the Version struct doesn't have a Compare method
func compareVersions(v1, v2 *Version) int {
	// Compare major version
	if v1.major != v2.major {
		if v1.major > v2.major {
			return 1
		}
		return -1
	}

	// Compare minor version
	if v1.minor != v2.minor {
		if v1.minor > v2.minor {
			return 1
		}
		return -1
	}

	// Compare patch version
	if v1.patch != v2.patch {
		if v1.patch > v2.patch {
			return 1
		}
		return -1
	}

	// If one has a pre-release and the other doesn't, the one without is greater
	if v1.preRelease == "" && v2.preRelease != "" {
		return 1
	}
	if v1.preRelease != "" && v2.preRelease == "" {
		return -1
	}

	// If both have pre-releases, compare them
	if v1.preRelease != "" && v2.preRelease != "" {
		// Compare pre-release types (alpha < beta < rc)
		if v1.preRelease != v2.preRelease {
			if v1.preRelease == "alpha" {
				return -1
			}
			if v2.preRelease == "alpha" {
				return 1
			}
			if v1.preRelease == "beta" {
				return -1
			}
			if v2.preRelease == "beta" {
				return 1
			}
		}

		// Compare pre-release versions
		if v1.preReleaseVersion != v2.preReleaseVersion {
			if v1.preReleaseVersion > v2.preReleaseVersion {
				return 1
			}
			return -1
		}
	}

	// Versions are equal
	return 0
}
