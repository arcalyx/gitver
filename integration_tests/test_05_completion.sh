#!/bin/bash

# Test for the completion command
# Usage: test_05_completion.sh <gitver_binary> <test_dir>

GITVER_BIN=$1
TEST_DIR=$2
TEST_NAME="Completion Command Test"

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
mkdir -p "$TEST_DIR/completion_test"
cd "$TEST_DIR/completion_test" || exit 1

# Test 1: Generate Bash completion
echo "Test 1: Generate Bash completion"
BASH_COMPLETION=$($GITVER_BIN completion bash 2>&1)
if [[ -n "$BASH_COMPLETION" ]]; then
  print_result 0 "Bash completion was successfully generated"
else
  print_result 1 "Bash completion could not be generated"
fi

# Test 2: Generate Zsh completion
echo "Test 2: Generate Zsh completion"
ZSH_COMPLETION=$($GITVER_BIN completion zsh 2>&1)
if [[ -n "$ZSH_COMPLETION" ]]; then
  print_result 0 "Zsh completion was successfully generated"
else
  print_result 1 "Zsh completion could not be generated"
fi

# Test 3: Generate Fish completion
echo "Test 3: Generate Fish completion"
FISH_COMPLETION=$($GITVER_BIN completion fish 2>&1)
if [[ -n "$FISH_COMPLETION" ]]; then
  print_result 0 "Fish completion was successfully generated"
else
  print_result 1 "Fish completion could not be generated"
fi

# Test 4: Invalid shell
echo "Test 4: Invalid shell"
INVALID_SHELL=$($GITVER_BIN completion invalid 2>&1)
if [[ "$INVALID_SHELL" == *"error"* || "$INVALID_SHELL" == *"Error"* || "$INVALID_SHELL" == *"invalid"* || "$INVALID_SHELL" == *"Invalid"* || $? -ne 0 ]]; then
  print_result 0 "Invalid shell was correctly rejected"
else
  print_result 1 "Invalid shell was not correctly rejected"
fi

exit 0
