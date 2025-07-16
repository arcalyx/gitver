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
		version.FromString(originalVersion)
	}()

	// Test SetBuildMetadata
	t.Run("SetBuildMetadata", func(t *testing.T) {
		// Setup
		version.FromString("1.2.3")

		// Execute
		err := version.SetBuildMetadata("build.123")

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, "build.123", version.GetBuildMetadata())
		assert.Equal(t, "1.2.3+build.123", version.ToString())
	})

	// Test ClearBuildMetadata
	t.Run("ClearBuildMetadata", func(t *testing.T) {
		// Setup
		version.FromString("1.2.3+build.123")

		// Execute
		err := version.ClearBuildMetadata()

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, "", version.GetBuildMetadata())
		assert.Equal(t, "1.2.3", version.ToString())
	})

	// Test GetBuildMetadata
	t.Run("GetBuildMetadata", func(t *testing.T) {
		// Setup
		version.FromString("1.2.3+build.123")

		// Execute & Verify
		assert.Equal(t, "build.123", version.GetBuildMetadata())

		// Setup for no build metadata
		version.FromString("1.2.3")

		// Execute & Verify
		assert.Equal(t, "", version.GetBuildMetadata())
	})
}
