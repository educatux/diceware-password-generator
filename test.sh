#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "Starting Diceware Generator tests..."

# Function to run a test case
run_test() {
    local test_name="$1"
    local test_cmd="$2"
    echo -n "Running $test_name... "
    if eval "$test_cmd" > /dev/null 2>&1; then
        echo -e "${GREEN}PASS${NC}"
        return 0
    else
        echo -e "${RED}FAIL${NC}"
        return 1
    fi
}

# Test build
run_test "Build test" "make build"

# Test container
run_test "Container build" "podman build -t diceware-test ."

# Test wordlist loading
run_test "Wordlist presence" "[ -f frenchdiceware.txt ] && [ -f diceware-fr-alt.txt ]"

# Test automated passphrase generation
run_test "Auto generation" "echo -e '1\n4\n1\nN\nN\n' | podman run -i diceware-test"

# Test QR code generation
run_test "QR generation" "echo -e '1\n4\n1\nN\nO\n' | podman run -i diceware-test"

# Cleanup
podman rmi diceware-test || true

echo "Tests completed."
