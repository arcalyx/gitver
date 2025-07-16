package cmd

import (
	"github.com/arcalyx/gitver/internal/version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildMetadataFunctionality(t *testing.T) {
	// Save original version
	originalVersion := version.ToString()
	defer func() {
		// Restore original version after tests
		err := version.FromString(originalVersion)
		if err != nil {
			t.Fatalf("Failed to restore original version: %v", err)
		}
	}()

	// Test SetBuildMetadata
	t.Run("SetBuildMetadata", func(t *testing.T) {
		// Setup
		err := version.FromString("1.2.3")
		if err != nil {
			t.Fatalf("Failed to set version: %v", err)
		}

		// Execute
		err = version.SetBuildMetadata("build.123")

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, "build.123", version.GetBuildMetadata())
		assert.Equal(t, "1.2.3+build.123", version.ToString())
	})

	// Test ClearBuildMetadata
	t.Run("ClearBuildMetadata", func(t *testing.T) {
		// Setup
		err := version.FromString("1.2.3+build.123")
		if err != nil {
			t.Fatalf("Failed to set version: %v", err)
		}

		// Execute
		err = version.ClearBuildMetadata()

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, "", version.GetBuildMetadata())
		assert.Equal(t, "1.2.3", version.ToString())
	})

	// Test GetBuildMetadata
	t.Run("GetBuildMetadata", func(t *testing.T) {
		// Setup
		err := version.FromString("1.2.3+build.123")
		if err != nil {
			t.Fatalf("Failed to set version: %v", err)
		}

		// Execute & Verify
		assert.Equal(t, "build.123", version.GetBuildMetadata())

		// Setup for no build metadata
		err = version.FromString("1.2.3")
		if err != nil {
			t.Fatalf("Failed to set version: %v", err)
		}

		// Execute & Verify
		assert.Equal(t, "", version.GetBuildMetadata())
	})
}
