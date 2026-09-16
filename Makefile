LINT_CFG=.golangci.yml
TARGET=gwt

TARGET:
	go build -o $(TARGET)

build: $(TARGET)

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed or not on PATH"; \
		exit 1; \
	fi

docker_build:
	@echo "Docker build..."
	@exit 1

test:
	go test $$(go list ./... | grep -v '/spec_tests')

comptest:
	@echo "Component testing"
