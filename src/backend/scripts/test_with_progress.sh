#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   Running Tests with Progress${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Count total test packages
TOTAL_PACKAGES=$(go list ./test/... 2>/dev/null | wc -l | tr -d ' ')
CURRENT=0

echo -e "${YELLOW}Found ${TOTAL_PACKAGES} test packages${NC}"
echo ""

# Function to run tests for a package with progress
run_test_package() {
    local package=$1
    CURRENT=$((CURRENT + 1))
    
    echo -e "${BLUE}[${CURRENT}/${TOTAL_PACKAGES}]${NC} Testing: ${package}"
    
    # Run tests and capture output
    OUTPUT=$(go test -v "$package" 2>&1)
    EXIT_CODE=$?
    
    # Count PASS and FAIL
    PASS_COUNT=$(echo "$OUTPUT" | grep -c "--- PASS:" || echo "0")
    FAIL_COUNT=$(echo "$OUTPUT" | grep -c "--- FAIL:" || echo "0")
    
    if [ $EXIT_CODE -eq 0 ]; then
        echo -e "  ${GREEN}✓ PASS${NC} - ${PASS_COUNT} test(s) passed"
    else
        echo -e "  ${RED}✗ FAIL${NC} - ${FAIL_COUNT} test(s) failed"
        echo "$OUTPUT" | grep -E "(FAIL|Error|panic)" | head -5
    fi
    echo ""
}

# Run tests for each package
for package in $(go list ./test/... 2>/dev/null); do
    run_test_package "$package"
done

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   Test Summary${NC}"
echo -e "${BLUE}========================================${NC}"

# Final summary
go test ./test/... 2>&1 | grep -E "(ok|FAIL)" | while read line; do
    if echo "$line" | grep -q "ok"; then
        echo -e "${GREEN}✓ $line${NC}"
    else
        echo -e "${RED}✗ $line${NC}"
    fi
done

