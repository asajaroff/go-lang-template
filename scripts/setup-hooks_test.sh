#!/bin/bash

# Test script for setup-hooks.sh

set -e

echo "Testing setup-hooks.sh..."

# Test 1: Check if .pre-commit-config.yaml exists
echo "Test 1: Checking if .pre-commit-config.yaml exists..."
if [ -f ".pre-commit-config.yaml" ]; then
    echo "✓ .pre-commit-config.yaml exists"
else
    echo "✗ .pre-commit-config.yaml does not exist"
    exit 1
fi

# Test 2: Validate YAML syntax
echo "Test 2: Validating YAML syntax..."
if command -v python3 &> /dev/null; then
    python3 -c "
import yaml
import sys
try:
    with open('.pre-commit-config.yaml', 'r') as f:
        yaml.safe_load(f)
    print('✓ YAML syntax is valid')
except Exception as e:
    print(f'✗ YAML syntax error: {e}')
    sys.exit(1)
"
else
    echo "⚠ Python3 not found, skipping YAML validation"
fi

# Test 3: Check if setup-hooks.sh is executable
echo "Test 3: Checking if setup-hooks.sh is executable..."
if [ -x "scripts/setup-hooks.sh" ]; then
    echo "✓ setup-hooks.sh is executable"
else
    echo "✗ setup-hooks.sh is not executable"
    exit 1
fi

# Test 4: Check if Makefile exists and has required targets
echo "Test 4: Checking Makefile targets..."
if [ -f "Makefile" ]; then
    required_targets=("setup-hooks" "test" "lint" "build" "help")
    for target in "${required_targets[@]}"; do
        if grep -q "^${target}:" Makefile; then
            echo "✓ Makefile has ${target} target"
        else
            echo "✗ Makefile missing ${target} target"
            exit 1
        fi
    done
else
    echo "✗ Makefile does not exist"
    exit 1
fi

# Test 5: Check if CONTRIBUTING.md exists
echo "Test 5: Checking if CONTRIBUTING.md exists..."
if [ -f "CONTRIBUTING.md" ]; then
    echo "✓ CONTRIBUTING.md exists"

    # Check for key sections
    if grep -q "Conventional Commits" CONTRIBUTING.md; then
        echo "✓ CONTRIBUTING.md mentions Conventional Commits"
    else
        echo "✗ CONTRIBUTING.md missing Conventional Commits section"
        exit 1
    fi

    if grep -q "Semantic Versioning" CONTRIBUTING.md; then
        echo "✓ CONTRIBUTING.md mentions Semantic Versioning"
    else
        echo "✗ CONTRIBUTING.md missing Semantic Versioning section"
        exit 1
    fi
else
    echo "✗ CONTRIBUTING.md does not exist"
    exit 1
fi

# Test 6: Check pre-commit config for required hooks
echo "Test 6: Checking pre-commit config for required hooks..."
required_hooks=("conventional-pre-commit" "go-fmt" "go-vet" "go-test")
for hook in "${required_hooks[@]}"; do
    if grep -q "${hook}" .pre-commit-config.yaml; then
        echo "✓ .pre-commit-config.yaml includes ${hook} hook"
    else
        echo "✗ .pre-commit-config.yaml missing ${hook} hook"
        exit 1
    fi
done

echo ""
echo "All tests passed! ✓"
echo ""
echo "Note: To fully test the pre-commit hooks, you need to:"
echo "1. Install pre-commit: pip install pre-commit"
echo "2. Run: make setup-hooks"
echo "3. Make a test commit with conventional format"
