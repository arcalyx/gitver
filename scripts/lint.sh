#!/bin/bash
set -e

# Run golangci-lint if available
if command -v golangci-lint &> /dev/null; then
  echo "Running golangci-lint..."
  golangci-lint run --timeout=5m
else
  echo "golangci-lint not found, skipping linting"
  echo "::warning::Linting skipped - golangci-lint not available"
fi

# Success
echo "Lint checks completed"
