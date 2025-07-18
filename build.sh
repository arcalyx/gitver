#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display headers
print_header() {
  echo -e "\n${BLUE}=== $1 ===${NC}\n"
}

# Function to display success/failure
print_result() {
  if [ $1 -eq 0 ]; then
    echo -e "${GREEN}✓ $2 successful${NC}"
  else
    echo -e "${RED}✗ $2 failed${NC}"
    if [ "$3" != "continue" ]; then
      exit 1
    fi
  fi
}

# Check if we are in the correct directory
if [ ! -d ".git" ] || [ ! -f "go.mod" ]; then
  echo -e "${RED}Error: This script must be executed in the main directory of the gitver project.${NC}"
  exit 1
fi

# Display help
show_help() {
  echo "Usage: $0 [options]"
  echo
  echo "Options:"
  echo "  --help          Show this help"
  echo "  --lint          Only perform linting"
  echo "  --test          Only run tests"
  echo "  --integration   Only run integration tests"
  echo "  --build         Only build the binary"
  echo "  --install       Install gitver after building"
  echo "  --all           Run everything (without installation)"
  echo "  --verbose       Verbose output"
  echo
}

# Process parameters
DO_LINT=false
DO_TEST=false
DO_BUILD=false
DO_INTEGRATION=false
DO_INSTALL=false
VERBOSE=""

if [ $# -eq 0 ]; then
  # Run everything by default
  DO_LINT=true
  DO_TEST=true
  DO_BUILD=true
  DO_INTEGRATION=true
else
  for arg in "$@"; do
    case $arg in
      --help)
        show_help
        exit 0
        ;;
      --lint)
        DO_LINT=true
        ;;
      --test)
        DO_TEST=true
        ;;
      --integration)
        DO_INTEGRATION=true
        ;;
      --build)
        DO_BUILD=true
        ;;
      --install)
        DO_BUILD=true  # Installation requires build
        DO_INSTALL=true
        ;;
      --all)
        DO_LINT=true
        DO_TEST=true
        DO_BUILD=true
        DO_INTEGRATION=true
        ;;
      --verbose)
        VERBOSE="-v"
        ;;
      *)
        echo -e "${RED}Unknown option: $arg${NC}"
        show_help
        exit 1
        ;;
    esac
  done
fi

# Main functions

run_lint() {
  print_header "Go Linting"
  
  # Check if golangci-lint is installed
  if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}golangci-lint is not installed. Attempting installation...${NC}"
    
    # Install golangci-lint
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    # Add Go bin path to PATH if command is not found
    GOPATH=$(go env GOPATH)
    export PATH=$PATH:$GOPATH/bin
    
    # Check again if golangci-lint is now available
    if ! command -v golangci-lint &> /dev/null; then
      echo -e "${RED}golangci-lint could not be installed or is not in PATH.${NC}"
      echo -e "${YELLOW}Please install golangci-lint manually: https://golangci-lint.run/usage/install/${NC}"
      return 1
    fi
    
    print_result $? "Installation of golangci-lint" "continue"
  fi
  
  echo "Running golangci-lint..."
  golangci-lint run $VERBOSE ./...
  print_result $? "Linting"
}

run_tests() {
  print_header "Go Tests"
  
  # Run tests with coverage
  echo "Running tests with coverage..."
  go test $VERBOSE -race -coverprofile=coverage.out ./...
  TEST_RESULT=$?
  print_result $TEST_RESULT "Tests"
  
  # Show coverage report if tests were successful
  if [ $TEST_RESULT -eq 0 ]; then
    echo -e "\n${BLUE}Coverage report:${NC}"
    go tool cover -func=coverage.out
  fi
}

run_integration() {
  print_header "CLI Integration Tests"
  
  # Run integration tests
  echo "Running CLI integration tests..."
  bash ./integration_tests/run_tests.sh "$(pwd)/build/gitver"
  print_result $? "CLI Integration Tests"
}

run_build() {
  print_header "Go Build"
  
  # Platform-specific builds
  build_for_platform() {
    local os=$1
    local arch=$2
    local output="build/gitver"
    
    # Add .exe for Windows
    if [ "$os" = "windows" ]; then
      output="${output}.exe"
    fi
    
    echo "Building for $os/$arch..."
    GOOS=$os GOARCH=$arch go build $VERBOSE -o $output
    print_result $? "Build for $os/$arch" "continue"
  }
  
  # Create build directory
  mkdir -p build
  
  # Build for current platform
  echo "Building for current platform..."
  go build $VERBOSE -o build/gitver
  print_result $? "Build for current platform"
  
  # Ask if builds for other platforms should be created
  echo -e "\n${YELLOW}Do you want to build for other platforms as well? (y/n)${NC}"
  read -r build_all
  
  if [[ $build_all =~ ^[yY] ]]; then
    build_for_platform "linux" "amd64"
    build_for_platform "windows" "amd64"
    build_for_platform "darwin" "amd64"
    build_for_platform "darwin" "arm64"
  fi
  
  echo -e "\n${GREEN}Build completed. Binaries are located in the 'build' directory.${NC}"
}

run_install() {
  print_header "Installation"
  
  # Check if Go is installed
  if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
  fi
  
  # Check if binary exists
  if [ ! -f "build/gitver" ]; then
    echo -e "${RED}Error: Binary not found. Please run the build first.${NC}"
    exit 1
  fi
  
  # Determine installation path based on platform
  INSTALL_DIR=""
  if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    INSTALL_DIR="/usr/local/bin"
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    INSTALL_DIR="/usr/local/bin"
  elif [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo -e "${YELLOW}Windows detected. Installing only to ./build.${NC}"
    echo -e "${GREEN}Installation successful! Binary is located at ./build/gitver${NC}"
    echo "Add this directory to your PATH to use gitver from anywhere."
    return 0
  else
    echo -e "${YELLOW}Unknown operating system. Installing only to ./build.${NC}"
    echo -e "${GREEN}Installation successful! Binary is located at ./build/gitver${NC}"
    echo "Add this directory to your PATH to use gitver from anywhere."
    return 0
  fi
  
  # Install the binary
  echo "Installing gitver to $INSTALL_DIR..."
  if [ -w "$INSTALL_DIR" ]; then
    # User has write permissions to the installation directory
    cp build/gitver "$INSTALL_DIR/gitver"
  else
    # Sudo is required for installation
    echo "Administrator privileges are required to install to $INSTALL_DIR"
    sudo cp build/gitver "$INSTALL_DIR/gitver"
  fi
  
  if [ $? -ne 0 ]; then
    echo -e "${RED}Installation failed!${NC}"
    echo "You can manually copy the binary from ./build/gitver to a directory in your PATH."
    return 1
  fi
  
  echo -e "${GREEN}Installation successful!${NC}"
  echo "You can now use gitver from anywhere."
  echo "Run 'gitver --help' to get started."
}

# Run main program

if [ "$DO_LINT" = true ]; then
  run_lint
fi

if [ "$DO_TEST" = true ]; then
  run_tests
fi

if [ "$DO_BUILD" = true ]; then
  run_build
fi

if [ "$DO_INTEGRATION" = true ]; then
  run_integration
fi

if [ "$DO_INSTALL" = true ]; then
  run_install
fi

# Show summary
echo -e "\n${GREEN}All selected tasks have been completed.${NC}"
