IMAGE_TAG ?= $(shell date -u +%Y%m%d-%H%M%S)
LINT_CFG=.golangci.yml
TARGET=gwt
TARGETARCH ?= $(shell go env GOARCH)
GO=go
DOCKER=docker
DOCKER_COMPOSE=docker compose

TARGET:
	GOARCH=$(TARGETARCH) $(GO) build -o $(TARGET)

.PHONY: build
build:
	GOARCH=$(TARGETARCH) $(GO) build -o $(TARGET)

.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed or not on PATH"; \
		exit 1; \
	fi

.PHONY: docker_build
docker_build:
	$(DOCKER) build \
		--platform $(DOCKER_PLATFORM) \
		--build-arg TARGETARCH=$(TARGETARCH) \
		-t gwt_comptest:$(IMAGE_TAG) .

.PHONY: test
test:
	$(GO) test $$(go list ./... | grep -v '/comp_tests')

.PHONY: comptest
comptest:
	$(DOCKER_COMPOSE) run --rm --build --remove-orphans test pytest -q /comp_tests

.PHONY: clean
clean:
	rm -f $(TARGET)
	$(DOCKER_COMPOSE) down --remove-orphans
	$(DOCKER) builder prune -af