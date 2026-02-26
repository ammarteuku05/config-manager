.PHONY: run-controller run-agent run-worker test test-coverage mock

# Run the controller server
run-controller:
	go run main.go controller

# Run the agent server
run-agent:
	go run main.go agent

# Run the worker server
run-worker:
	go run main.go worker

# Run all unit tests
test:
	go test -v ./...

# Run tests and generate coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	grep -v "mocks" coverage.out > coverage_filtered.out && mv coverage_filtered.out coverage.out
	go tool cover -func=coverage.out

# Generate mocks using mockery
mock:
	mockery
