#!/bin/bash
# Deployment Test Script for GitHub Actions
# This script validates deployment manifests and configuration

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
MANIFEST_DIR="$SCRIPT_DIR/manifest"

echo "=================================="
echo "Sensor Verwerker Deployment Tests"
echo "=================================="

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

test_count=0
passed_count=0
failed_count=0

# Test function
run_test() {
    local test_name=$1
    local test_command=$2
    
    test_count=$((test_count + 1))
    echo -e "\n${YELLOW}[Test $test_count]${NC} $test_name"
    
    if eval "$test_command"; then
        echo -e "${GREEN}✓ PASSED${NC}"
        passed_count=$((passed_count + 1))
    else
        echo -e "${RED}✗ FAILED${NC}"
        failed_count=$((failed_count + 1))
    fi
}

# Test 1: Check if manifest files exist
run_test "Manifest files exist" \
    "test -f $MANIFEST_DIR/deployment.yaml && test -f $MANIFEST_DIR/k8s-deployment.yaml && test -f $MANIFEST_DIR/serviceaccount.yaml"

# Test 2: Validate YAML syntax
run_test "deployment.yaml YAML syntax" \
    "kubectl apply --dry-run=client -f $MANIFEST_DIR/deployment.yaml > /dev/null 2>&1"

run_test "k8s-deployment.yaml YAML syntax" \
    "kubectl apply --dry-run=client -f $MANIFEST_DIR/k8s-deployment.yaml > /dev/null 2>&1"

run_test "serviceaccount.yaml YAML syntax" \
    "kubectl apply --dry-run=client -f $MANIFEST_DIR/serviceaccount.yaml > /dev/null 2>&1"

run_test "test-deployment.yaml YAML syntax" \
    "kubectl apply --dry-run=client -f $MANIFEST_DIR/test-deployment.yaml > /dev/null 2>&1"

# Test 3: Check Go code structure
run_test "main.go exists" \
    "test -f $SCRIPT_DIR/main.go"

run_test "main_test.go exists" \
    "test -f $SCRIPT_DIR/main_test.go"

run_test "go.mod exists" \
    "test -f $SCRIPT_DIR/go.mod"

# Test 4: Check for required fields in deployment
run_test "Deployment has replicas" \
    "grep -q 'replicas:' $MANIFEST_DIR/deployment.yaml"

run_test "Deployment has resource limits" \
    "grep -q 'limits:' $MANIFEST_DIR/deployment.yaml"

run_test "Deployment has security context" \
    "grep -q 'securityContext:' $MANIFEST_DIR/deployment.yaml"

# Test 5: Check for health probes
run_test "Deployment has liveness probe" \
    "grep -q 'livenessProbe:' $MANIFEST_DIR/deployment.yaml"

run_test "Deployment has startup probe" \
    "grep -q 'startupProbe:' $MANIFEST_DIR/deployment.yaml"

# Test 6: Verify image references
run_test "Image reference in deployment" \
    "grep -q 'roelandvdberg/sensor-verwerker' $MANIFEST_DIR/deployment.yaml"

# Test 7: Check for environment variables
run_test "Namespace configuration" \
    "grep -q 'namespace:' $MANIFEST_DIR/deployment.yaml"

# Test 8: Verify Go code constants
run_test "TruckVIN constant defined" \
    "grep -q 'TruckVIN.*=.*FC-TRUCK' $SCRIPT_DIR/main.go"

run_test "PostgrestAPI constant defined" \
    "grep -q 'PostgrestAPI' $SCRIPT_DIR/main.go"

run_test "SensorCount constant defined" \
    "grep -q 'SensorCount.*=' $SCRIPT_DIR/main.go"

# Test 9: Dockerfile validation
run_test "Dockerfile exists" \
    "test -f $SCRIPT_DIR/Dockerfile"

# Print summary
echo ""
echo "=================================="
echo "Test Summary"
echo "=================================="
echo -e "${GREEN}Passed: $passed_count${NC}"
echo -e "${RED}Failed: $failed_count${NC}"
echo "Total:  $test_count"
echo "=================================="

if [ $failed_count -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi
