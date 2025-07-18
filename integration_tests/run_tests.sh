#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_TOTAL=0
TESTS_PASSED=0
TESTS_FAILED=0

# Directory for temporary test projects
TEST_DIR="$(pwd)/integration_tests/temp"
GITVER_BIN="${1:-$(pwd)/build/gitver}"

# Function to display headers
print_header() {
  echo -e "\n${BLUE}=== $1 ===${NC}\n"
}

# Function to display test results
print_result() {
  if [ $1 -eq 0 ]; then
    echo -e "${GREEN}✓ $2${NC}"
    TESTS_PASSED=$((TESTS_PASSED+1))
  else
    echo -e "${RED}✗ $2${NC}"
    TESTS_FAILED=$((TESTS_FAILED+1))
  fi
  TESTS_TOTAL=$((TESTS_TOTAL+1))
}

# Cleanup function
cleanup() {
  if [ -d "$TEST_DIR" ]; then
    rm -rf "$TEST_DIR"
  fi
}

# Cleanup on exit
trap cleanup EXIT

# Check if gitver has been built
if [ ! -f "$GITVER_BIN" ]; then
  echo -e "${YELLOW}Gitver binary not found. Building project...${NC}"
  ./build.sh --build
  if [ ! -f "$GITVER_BIN" ]; then
    echo -e "${RED}Error: Could not build gitver.${NC}"
    exit 1
  fi
fi

# Create test directory
cleanup
mkdir -p "$TEST_DIR"

print_header "CLI Integration Tests for Gitver"

# Run all test scripts
for test_script in $(pwd)/integration_tests/test_*.sh; do
  if [ -f "$test_script" ]; then
    echo -e "${BLUE}Running test: $(basename $test_script)${NC}"
    bash "$test_script" "$GITVER_BIN" "$TEST_DIR"
    TEST_RESULT=$?
    if [ $TEST_RESULT -ne 0 ]; then
      TESTS_FAILED=$((TESTS_FAILED+1))
    else
      TESTS_PASSED=$((TESTS_PASSED+1))
    fi
    TESTS_TOTAL=$((TESTS_TOTAL+1))
  fi
done

# Show summary
print_header "Summary"
echo -e "Total: ${TESTS_TOTAL} Tests"
echo -e "${GREEN}Passed: ${TESTS_PASSED} Tests${NC}"
echo -e "${RED}Failed: ${TESTS_FAILED} Tests${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
  echo -e "\n${GREEN}All tests passed!${NC}"
  exit 0
else
  echo -e "\n${RED}Some tests failed.${NC}"
  exit 1
fi
