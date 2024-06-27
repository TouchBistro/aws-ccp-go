.DEFAULT_GOAL = build

# Get all dependencies
setup:
# Only install if missing
ifeq (,$(wildcard bin/golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.59.1
endif

	go mod tidy
.PHONY: setup

# clean all build artifacts
clean:
	rm -rf dist
	rm -rf coverage
	rm -rf ./codegen/clients
.PHONY: clean

# build source
build:
	go build ./...
.PHONY: build

# Run the linter
lint:
	./bin/golangci-lint run ./...
.PHONY: lint

# Run tests and collect coverage data
test:
	mkdir -p coverage
	go test -coverpkg=./... -coverprofile=coverage/coverage.txt ./...
	go tool cover -html=coverage/coverage.txt -o coverage/coverage.html
.PHONY: test

# Run tests and print coverage data to stdout
test-ci:
	mkdir -p coverage
	go test -coverpkg=./... -coverprofile=coverage/coverage.txt ./...
	go tool cover -func=coverage/coverage.txt
.PHONY: test-ci