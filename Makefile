LINT_CFG=.golangci.yml
TARGET=gwt
TARGETARCH ?= $(shell go env GOARCH)

TARGET:
	GOARCH=$(TARGETARCH) go build -o $(TARGET)

.PHONY: build
build:
	GOARCH=$(TARGETARCH) go build -o $(TARGET)

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

.PHONY: comptest
comptest:
	docker compose run --rm test pytest -q /spec_tests
