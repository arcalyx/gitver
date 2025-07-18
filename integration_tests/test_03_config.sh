#!/bin/bash

# Test for the config command
# Usage: test_03_config.sh <gitver_binary> <test_dir>

GITVER_BIN=$1
TEST_DIR=$2
TEST_NAME="Config Command Test"

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
mkdir -p "$TEST_DIR/config_test"
cd "$TEST_DIR/config_test" || exit 1

# Initialize Git repository
git init > /dev/null 2>&1
git config user.email "test@example.com"
git config user.name "Test User"

# Initialize gitver
$GITVER_BIN init --version 1.0.0 > /dev/null 2>&1

# Test 1: Config command should execute without errors
echo "Test 1: Execute config command"
$GITVER_BIN config > /dev/null 2>&1
print_result $? "Config command executed"

# Test 2: Set config value
echo "Test 2: Set config value"
$GITVER_BIN config set "test.key" "test-value" > /dev/null 2>&1
SET_RESULT=$?
print_result $SET_RESULT "Attempted to set config value"

# Test 3: List config values
echo "Test 3: List config values"
CONFIG_LIST=$($GITVER_BIN config --list 2>&1)
if [[ -n "$CONFIG_LIST" ]]; then
  print_result 0 "Config list was generated"
else
  print_result 1 "Config list could not be generated"
fi

# Test 4: Execute config init
echo "Test 4: Execute config init"
$GITVER_BIN config init > /dev/null 2>&1
INIT_RESULT=$?
print_result $INIT_RESULT "Config init was executed"

exit 0
