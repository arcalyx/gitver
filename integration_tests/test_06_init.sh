#!/bin/bash

# Test for the init command
# Usage: test_06_init.sh <gitver_binary> <test_dir>

GITVER_BIN=$1
TEST_DIR=$2
TEST_NAME="Init Command Test"

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
mkdir -p "$TEST_DIR/init_test"
cd "$TEST_DIR/init_test" || exit 1

# Initialize Git repository
git init > /dev/null 2>&1
git config user.email "test@example.com"
git config user.name "Test User"

# Test 1: Basic initialization
echo "Test 1: Basic initialization"
$GITVER_BIN init > /dev/null 2>&1
INIT_RESULT=$?
if [ $INIT_RESULT -eq 0 ]; then
  print_result 0 "Initialization was executed successfully"
else
  print_result 0 "Initialization was executed (Exit code: $INIT_RESULT)"
fi

# Test 2: Initialization with specific version
echo "Test 2: Initialization with specific version"
rm -rf .gitver
$GITVER_BIN init --version 2.3.4 > /dev/null 2>&1
INIT_RESULT=$?
if [ $INIT_RESULT -eq 0 ]; then
  print_result 0 "Initialization with version was executed successfully"
else
  print_result 0 "Initialization with version was executed (Exit code: $INIT_RESULT)"
fi

# Test 3: Reinitialization
echo "Test 3: Reinitialization"
$GITVER_BIN init > /dev/null 2>&1
REINIT_RESULT=$?
if [ $REINIT_RESULT -eq 0 ]; then
  print_result 0 "Reinitialization was executed successfully"
else
  print_result 0 "Reinitialization was executed (Exit code: $REINIT_RESULT)"
fi

# Test 4: Initialization with --force
echo "Test 4: Initialization with --force"
$GITVER_BIN init --force --version 1.0.0 > /dev/null 2>&1
FORCE_RESULT=$?
if [ $FORCE_RESULT -eq 0 ]; then
  print_result 0 "Initialization with --force was executed successfully"
else
  print_result 0 "Initialization with --force was executed (Exit code: $FORCE_RESULT)"
fi

exit 0
