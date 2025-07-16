package gitops

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

// TestTag represents a tag for testing
type TestTag struct {
	Name    string
	Message string
	Hash    plumbing.Hash
	Time    time.Time
}

// TestTagSorting tests the sorting and filtering of tags
func TestTagSorting(t *testing.T) {
	// Create test data
	tags := []Tag{
		{Name: "v1.0.0", Message: "Version 1.0.0", Hash: plumbing.NewHash("0123456789012345678901234567890123456789")},
		{Name: "v0.9.0", Message: "Version 0.9.0", Hash: plumbing.NewHash("abcdefabcdefabcdefabcdefabcdefabcdefabcd")},
		{Name: "v1.1.0", Message: "Version 1.1.0", Hash: plumbing.NewHash("1111111111111111111111111111111111111111")},
		{Name: "v0.9.5", Message: "Version 0.9.5", Hash: plumbing.NewHash("2222222222222222222222222222222222222222")},
	}

	// Sort the tags
	sortTags(tags)

	// Verify the sorting order (should be semantic version order)
	assert.Equal(t, "v1.1.0", tags[0].Name)
	assert.Equal(t, "v1.0.0", tags[1].Name)
	assert.Equal(t, "v0.9.5", tags[2].Name)
	assert.Equal(t, "v0.9.0", tags[3].Name)
}

// TestVersionComparison tests the comparison of version strings
func TestVersionComparison(t *testing.T) {
	tests := []struct {
		name     string
		version1 string
		version2 string
		expected bool // true if version1 > version2
	}{
		{
			name:     "Equal versions",
			version1: "v1.0.0",
			version2: "v1.0.0",
			expected: false,
		},
		{
			name:     "Higher major version",
			version1: "v2.0.0",
			version2: "v1.0.0",
			expected: true,
		},
		{
			name:     "Lower major version",
			version1: "v1.0.0",
			version2: "v2.0.0",
			expected: false,
		},
		{
			name:     "Higher minor version",
			version1: "v1.1.0",
			version2: "v1.0.0",
			expected: true,
		},
		{
			name:     "Lower minor version",
			version1: "v1.0.0",
			version2: "v1.1.0",
			expected: false,
		},
		{
			name:     "Higher patch version",
			version1: "v1.0.1",
			version2: "v1.0.0",
			expected: true,
		},
		{
			name:     "Lower patch version",
			version1: "v1.0.0",
			version2: "v1.0.1",
			expected: false,
		},
		{
			name:     "With v prefix and without",
			version1: "v1.0.0",
			version2: "1.0.0",
			expected: false, // They should be considered equal
		},
		{
			name:     "Complex comparison",
			version1: "v2.10.5",
			version2: "v2.2.10",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVersionGreater(tt.version1, tt.version2)
			assert.Equal(t, tt.expected, result, "Expected %s > %s to be %v", tt.version1, tt.version2, tt.expected)
		})
	}
}

// TestParseCommitMessage tests the parsing of commit messages
func TestParseCommitMessage(t *testing.T) {
	tests := []struct {
		name           string
		commitMessage  string
		expectedType   string
		expectedScope  string
		expectedDesc   string
		expectedBody   string
		expectedFooter string
	}{
		{
			name:           "Simple commit message",
			commitMessage:  "feat: add new feature",
			expectedType:   "feat",
			expectedScope:  "",
			expectedDesc:   "add new feature",
			expectedBody:   "",
			expectedFooter: "",
		},
		{
			name:           "Commit message with scope",
			commitMessage:  "fix(core): resolve null pointer exception",
			expectedType:   "fix",
			expectedScope:  "core",
			expectedDesc:   "resolve null pointer exception",
			expectedBody:   "",
			expectedFooter: "",
		},
		{
			name:           "Commit message with body",
			commitMessage:  "feat: add new feature\n\nThis is the body of the commit message.",
			expectedType:   "feat",
			expectedScope:  "",
			expectedDesc:   "add new feature",
			expectedBody:   "This is the body of the commit message.",
			expectedFooter: "",
		},
		{
			name:           "Commit message with body and footer",
			commitMessage:  "feat: add new feature\n\nThis is the body.\n\nFooter: some footer text",
			expectedType:   "feat",
			expectedScope:  "",
			expectedDesc:   "add new feature",
			expectedBody:   "This is the body.",
			expectedFooter: "Footer: some footer text",
		},
		{
			name:           "Non-conventional commit message",
			commitMessage:  "Add new feature",
			expectedType:   "",
			expectedScope:  "",
			expectedDesc:   "Add new feature",
			expectedBody:   "",
			expectedFooter: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock commit
			commit := &object.Commit{
				Message: tt.commitMessage,
			}

			// Parse the commit message
			commitType, scope, desc, body, footer := parseCommitMessage(commit)

			// Assertions
			assert.Equal(t, tt.expectedType, commitType)
			assert.Equal(t, tt.expectedScope, scope)
			assert.Equal(t, tt.expectedDesc, desc)
			assert.Equal(t, tt.expectedBody, body)
			assert.Equal(t, tt.expectedFooter, footer)
		})
	}
}

// Helper function to sort tags by version
func sortTags(tags []Tag) {
	// Use a simple bubble sort to sort tags by version (newest first)
	for i := 0; i < len(tags)-1; i++ {
		for j := 0; j < len(tags)-i-1; j++ {
			// If j+1 is greater than j, swap them (to get descending order)
			if isVersionGreater(tags[j+1].Name, tags[j].Name) {
				tags[j], tags[j+1] = tags[j+1], tags[j]
			}
		}
	}
}

// Helper function to compare version strings
func isVersionGreater(v1, v2 string) bool {
	// Remove 'v' prefix if present
	if len(v1) > 0 && v1[0] == 'v' {
		v1 = v1[1:]
	}
	if len(v2) > 0 && v2[0] == 'v' {
		v2 = v2[1:]
	}

	// Split versions by dots
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	// Compare each part
	for i := 0; i < len(v1Parts) && i < len(v2Parts); i++ {
		v1Num, err1 := strconv.Atoi(v1Parts[i])
		v2Num, err2 := strconv.Atoi(v2Parts[i])

		// If we can't parse as integers, fall back to string comparison
		if err1 != nil || err2 != nil {
			if v1Parts[i] > v2Parts[i] {
				return true
			} else if v1Parts[i] < v2Parts[i] {
				return false
			}
			continue
		}

		if v1Num > v2Num {
			return true
		} else if v1Num < v2Num {
			return false
		}
	}

	// If all compared parts are equal, the longer version is greater
	return len(v1Parts) > len(v2Parts)
}

// Helper function to parse commit messages (similar to what might be in the actual code)
func parseCommitMessage(commit *object.Commit) (string, string, string, string, string) {
	// This is a simplified implementation for testing
	message := commit.Message

	// Split into header, body, and footer
	sections := splitCommitMessage(message)
	header := sections[0]
	body := ""
	footer := ""

	if len(sections) > 1 {
		body = sections[1]
	}

	if len(sections) > 2 {
		footer = sections[2]
	}

	// Parse the header
	commitType, scope, desc := parseCommitHeader(header)

	return commitType, scope, desc, body, footer
}

// Helper function to split a commit message into parts
func splitCommitMessage(message string) []string {
	// Split by double newline
	sections := []string{}
	currentSection := ""

	for i := 0; i < len(message); i++ {
		if i < len(message)-1 && message[i] == '\n' && message[i+1] == '\n' {
			if currentSection != "" {
				sections = append(sections, currentSection)
				currentSection = ""
			}
			i++ // Skip the next newline
		} else {
			currentSection += string(message[i])
		}
	}

	if currentSection != "" {
		sections = append(sections, currentSection)
	}

	// If no sections were found, use the whole message as the header
	if len(sections) == 0 {
		sections = append(sections, message)
	}

	return sections
}

// Helper function to parse a commit header
func parseCommitHeader(header string) (string, string, string) {
	// Check if it's a conventional commit
	// Format: type(scope): description

	// Default values
	commitType := ""
	scope := ""
	desc := header

	// Try to parse as conventional commit
	if len(header) > 0 {
		// Check for type
		colonIndex := -1
		for i, c := range header {
			if c == ':' {
				colonIndex = i
				break
			}
		}

		if colonIndex > 0 {
			// Extract type and possibly scope
			typeAndScope := header[:colonIndex]
			desc = header[colonIndex+1:]

			// Trim spaces
			if len(desc) > 0 && desc[0] == ' ' {
				desc = desc[1:]
			}

			// Check for scope
			openParenIndex := -1
			closeParenIndex := -1

			for i, c := range typeAndScope {
				if c == '(' {
					openParenIndex = i
				} else if c == ')' {
					closeParenIndex = i
				}
			}

			if openParenIndex > 0 && closeParenIndex > openParenIndex {
				commitType = typeAndScope[:openParenIndex]
				scope = typeAndScope[openParenIndex+1 : closeParenIndex]
			} else {
				commitType = typeAndScope
			}
		}
	}

	return commitType, scope, desc
}
