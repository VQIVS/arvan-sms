# SMS Service Unit Tests

This directory contains unit tests for the SMS service in the SMS Gateway application.

## Test Files

- **`sms_service_test.go`** - Basic unit tests for SMS service methods
- **`sms_service_table_test.go`** - Table-driven tests for comprehensive scenarios
- **`sms_service_benchmark_test.go`** - Performance benchmark tests
- **`sms_service_integration_test.go`** - Integration and edge case tests
- **`sms_service_race_test.go`** - Race condition and concurrency tests
- **`mock_service.go`** - Mock implementation of the service interface
- **`test_helper.go`** - Utility functions for testing

## What is Tested

### SendSMS Function
- Successful SMS creation
- Service layer errors (CreateSMS, UserBalanceUpdate)
- Empty recipient and message handling
- Long messages and international numbers
- Special characters and emojis
- Concurrent requests
- Context cancellation and timeout

### GetSMSMessage Function
- Successful SMS retrieval
- SMS not found errors
- Different SMS statuses (pending, delivered, failed)
- Edge cases with ID values

### Performance
- Benchmark tests for both functions
- Parallel execution performance
- Memory allocation tracking

### Race Conditions
- Concurrent SendSMS operations
- Concurrent GetSMSMessage operations
- Mixed operation scenarios
- Shared data access patterns
- Context cancellation under load
- Error handling in concurrent scenarios
- High load stress testing

## Running Tests

### Basic Commands
```bash
# Run all tests
go test ./tests/...

# Run with coverage
go test -cover ./tests/...

# Run benchmarks
go test -bench=. ./tests/

# Run race condition tests
go test -race ./tests/

# Run specific test
go test -run TestSMSService_SendSMS_Success ./tests/...
```

### Using Test Script
```bash
# Basic tests
./tests/run_tests.sh basic

# Coverage report
./tests/run_tests.sh coverage

# Benchmarks
./tests/run_tests.sh bench

# Race condition tests
./tests/run_tests.sh race

# All tests
./tests/run_tests.sh all
```

## Features

- No external dependencies (uses Go's built-in testing framework)
- Mock implementations for dependency isolation
- Table-driven tests for comprehensive coverage
- Benchmark tests for performance monitoring
- Context testing for timeout and cancellation
- Concurrent testing for race condition detection

## Adding New Tests

When adding new tests:
1. Follow existing naming conventions
2. Use the TestHelper for common operations
3. Add both positive and negative test cases
4. Include edge cases and boundary conditions
5. Update this README with new scenarios
