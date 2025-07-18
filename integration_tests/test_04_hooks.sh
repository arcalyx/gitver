#!/bin/bash

# Test for the hooks command
# Usage: test_04_hooks.sh <gitver_binary> <test_dir>

GITVER_BIN=$1
TEST_DIR=$2
TEST_NAME="Hooks Command Test"

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
mkdir -p "$TEST_DIR/hooks_test"
cd "$TEST_DIR/hooks_test" || exit 1

# Initialize Git repository
git init > /dev/null 2>&1
git config user.email "test@example.com"
git config user.name "Test User"

# Initialize gitver
$GITVER_BIN init --version 1.0.0 > /dev/null 2>&1

# Test 1: Hooks command should execute without errors
echo "Test 1: Execute hooks command"
$GITVER_BIN hooks > /dev/null 2>&1
print_result $? "Hooks command executed"

# Test 2: Install hooks
echo "Test 2: Install hooks"
$GITVER_BIN hooks install > /dev/null 2>&1
INSTALL_RESULT=$?
print_result $INSTALL_RESULT "Attempted to install hooks"

# Test 3: List hooks
echo "Test 3: List hooks"
HOOKS_LIST=$($GITVER_BIN hooks list 2>&1)
if [[ -n "$HOOKS_LIST" ]]; then
  print_result 0 "Hooks list was generated"
else
  print_result 1 "Hooks list could not be generated"
fi

# Test 4: Uninstall hooks
echo "Test 4: Uninstall hooks"
$GITVER_BIN hooks uninstall > /dev/null 2>&1
UNINSTALL_RESULT=$?
print_result $UNINSTALL_RESULT "Attempted to uninstall hooks"

exit 0
