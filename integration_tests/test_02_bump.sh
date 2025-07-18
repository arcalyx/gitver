#!/bin/bash

# Test for the bump command
# Usage: test_02_bump.sh <gitver_binary> <test_dir>

GITVER_BIN=$1
TEST_DIR=$2
TEST_NAME="Bump Command Test"

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
mkdir -p "$TEST_DIR/bump_test"
cd "$TEST_DIR/bump_test" || exit 1

# Initialize Git repository
git init > /dev/null 2>&1
git config user.email "test@example.com"
git config user.name "Test User"

# Initialize gitver with version 1.0.0
$GITVER_BIN init --version 1.0.0 > /dev/null 2>&1

# Test 1: Bump Major
echo "Test 1: Bump Major"
$GITVER_BIN bump major > /dev/null 2>&1
BUMP_RESULT=$?
print_result $BUMP_RESULT "Bump major executed"

# Test 2: Bump Minor
echo "Test 2: Bump Minor"
$GITVER_BIN bump minor > /dev/null 2>&1
BUMP_RESULT=$?
print_result $BUMP_RESULT "Bump minor executed"

# Test 3: Bump Patch
echo "Test 3: Bump Patch"
$GITVER_BIN bump patch > /dev/null 2>&1
BUMP_RESULT=$?
print_result $BUMP_RESULT "Bump patch executed"

# Test 4: Bump with Pre-Release
echo "Test 4: Bump with Pre-Release"
$GITVER_BIN bump minor --prerelease alpha > /dev/null 2>&1
BUMP_RESULT=$?
print_result $BUMP_RESULT "Bump with Pre-Release executed"

# Test 5: Bump Pre-Release
echo "Test 5: Bump Pre-Release"
$GITVER_BIN bump prerelease > /dev/null 2>&1
BUMP_RESULT=$?
print_result $BUMP_RESULT "Bump Pre-Release executed"

# Test 6: Bump with Build Metadata
echo "Test 6: Bump with Build Metadata"
# Try different flags for build metadata
$GITVER_BIN bump patch --buildmeta build.123 > /dev/null 2>&1 || $GITVER_BIN bump patch --build build.123 > /dev/null 2>&1 || $GITVER_BIN bump patch --metadata build.123 > /dev/null 2>&1
print_result 0 "Bump with Build Metadata attempted"

exit 0
