#!/bin/bash

set -e

echo "Setting up pre-commit hooks..."

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "pre-commit is not installed. Installing..."

    # Try to install using pip
    if command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    elif command -v pip &> /dev/null; then
        pip install pre-commit
    else
        echo "Error: pip/pip3 not found. Please install Python and pip first."
        echo "Visit: https://pre-commit.com/#install"
        exit 1
    fi
fi

# Install the git hook scripts
echo "Installing git hook scripts..."
pre-commit install --hook-type pre-commit
pre-commit install --hook-type commit-msg

echo ""
echo "Pre-commit hooks installed successfully!"
echo ""
echo "The following checks will run automatically:"
echo "  - Conventional commit message format validation"
echo "  - Go formatting (go fmt)"
echo "  - Go static analysis (go vet)"
echo "  - Go tests (go test)"
echo "  - Trailing whitespace removal"
echo "  - YAML validation"
echo "  - Large file detection"
echo "  - Private key detection"
echo ""
echo "To run checks manually: pre-commit run --all-files"
echo ""
echo "Commit message format examples:"
echo "  feat: add new feature"
echo "  fix: resolve bug in authentication"
echo "  docs: update README"
echo "  test: add unit tests for config"
echo "  refactor: simplify error handling"
echo "  chore: update dependencies"
echo ""
