#!/bin/bash

# Test for the version command
# Usage: test_01_version.sh <gitver_binary> <test_dir>

GITVER_BIN=$1
TEST_DIR=$2
TEST_NAME="Version Command Test"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to display test results
print_result() {
  if [ $1 -eq 0 ]; then
    echo -e "${GREEN}✓ $2${NC}"
    return 0
  else
    echo -e "${RED}✗ $2${NC}"
    return 1
  fi
}

# Create test directory
mkdir -p "$TEST_DIR/version_test"
cd "$TEST_DIR/version_test" || exit 1

# Initialize Git repository
git init > /dev/null 2>&1
git config user.email "test@example.com"
git config user.name "Test User"

# Test 1: version command should execute without errors
echo "Test 1: Executing the version command"
$GITVER_BIN version > /dev/null 2>&1
print_result $? "Version command executed"

# Test 2: version command should output "Current version:"
echo "Test 2: Checking output format"
VERSION_OUTPUT=$($GITVER_BIN version)
if [[ $VERSION_OUTPUT == *"Current version:"* ]]; then
  print_result 0 "Version has expected format: $VERSION_OUTPUT"
else
  print_result 1 "Version does not have the expected format: $VERSION_OUTPUT"
fi

# Test 3: Initialization with a specific version
echo "Test 3: Initialization with a specific version"
$GITVER_BIN init --version 1.2.3 > /dev/null 2>&1
VERSION_OUTPUT=$($GITVER_BIN version)
if [[ "$VERSION_OUTPUT" == *"Current version: 1.2.3"* ]]; then
  print_result 0 "Version was correctly set to 1.2.3"
else
  print_result 1 "Version was not set correctly. Expected: 1.2.3, Received: $VERSION_OUTPUT"
fi

# Test 4: Version with --format flag
echo "Test 4: Format output"
FORMAT_OUTPUT=$($GITVER_BIN version --format semver)
if [[ -n "$FORMAT_OUTPUT" ]]; then
  print_result 0 "Format output was generated: $FORMAT_OUTPUT"
else
  print_result 1 "Format output could not be generated"
fi

exit 0
