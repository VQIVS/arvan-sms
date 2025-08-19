#!/bin/bash

# Simple test runner script for SMS service tests

set -e

echo "SMS Service Test Runner"
echo "======================="

# Change to project directory
cd "$(dirname "$0")/.."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed or not in PATH"
    exit 1
fi

# Function to run basic tests
run_basic_tests() {
    echo "Running basic unit tests..."
    if go test ./tests/... -v; then
        echo "Basic tests passed!"
    else
        echo "Basic tests failed!"
        exit 1
    fi
}

# Function to run tests with coverage
run_coverage_tests() {
    echo "Running tests with coverage..."
    if go test ./tests/... -cover -coverprofile=coverage.out; then
        echo "Coverage tests completed!"
        
        # Generate HTML coverage report
        echo "Generating HTML coverage report..."
        go tool cover -html=coverage.out -o coverage.html
        echo "Coverage report generated: coverage.html"
    else
        echo "Coverage tests failed!"
        exit 1
    fi
}

# Function to run benchmark tests
run_benchmarks() {
    echo "Running benchmark tests..."
    if go test -bench=. -benchmem ./tests/; then
        echo "Benchmark tests completed!"
    else
        echo "Benchmark tests failed!"
        exit 1
    fi
}

# Function to run specific test
run_specific_test() {
    local test_name="$1"
    echo "Running specific test: $test_name"
    if go test -run "$test_name" ./tests/... -v; then
        echo "Test '$test_name' passed!"
    else
        echo "Test '$test_name' failed!"
        exit 1
    fi
}

# Function to run race condition tests
run_race_tests() {
    echo "Running tests with race detection..."
    if go test -race ./tests/...; then
        echo "Race condition tests passed!"
    else
        echo "Race condition detected!"
        exit 1
    fi
}

# Function to clean test artifacts
clean_artifacts() {
    echo "Cleaning test artifacts..."
    rm -f coverage.out coverage.html
    echo "Artifacts cleaned!"
}

# Main menu
show_help() {
    echo "Usage: $0 [OPTION]"
    echo 
    echo "Options:"
    echo "  basic      Run basic unit tests"
    echo "  coverage   Run tests with coverage report"
    echo "  bench      Run benchmark tests"
    echo "  race       Run tests with race detection"
    echo "  all        Run all test types"
    echo "  clean      Clean test artifacts"
    echo "  specific   Run specific test (example: $0 specific TestSMSService_SendSMS_Success)"
    echo "  help       Show this help message"
}

# Parse command line arguments
case "${1:-help}" in
    "basic")
        run_basic_tests
        ;;
    "coverage")
        run_coverage_tests
        ;;
    "bench")
        run_benchmarks
        ;;
    "race")
        run_race_tests
        ;;
    "all")
        echo "Running complete test suite..."
        run_basic_tests
        echo
        run_coverage_tests
        echo
        run_benchmarks
        echo
        run_race_tests
        echo
        echo "All tests completed successfully!"
        ;;
    "clean")
        clean_artifacts
        ;;
    "specific")
        if [ -z "$2" ]; then
            echo "Error: Please provide test name as second argument"
            echo "Example: $0 specific TestSMSService_SendSMS_Success"
            exit 1
        fi
        run_specific_test "$2"
        ;;
    "help"|*)
        show_help
        ;;
esac
