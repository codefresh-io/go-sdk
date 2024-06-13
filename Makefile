VERSION=v1.3.0

ifndef GOBIN
ifndef GOPATH
$(error GOPATH is not set, please make sure you set your GOPATH correctly!)
endif
GOBIN=$(GOPATH)/bin
ifndef GOBIN
$(error GOBIN is not set, please make sure you set your GOBIN correctly!)
endif
endif

.PHONY: test-all
test-all: test test-fmt

.PHONY: test
test:
	@sh ./scripts/test.sh

.PHONY: test-fmt
test-fmt:
	@sh ./scripts/test-fmt.sh

.PHONY: lint
lint: $(GOBIN)/golangci-lint
	@echo linting go code...
	@$(GOBIN)/golangci-lint run --fix --timeout 10m


# Fix fmt errors in file
.PHONY: fmt
fmt:
	go fmt ./...

# Generate mock struct from interface
# example: make mock PKG=./pkg/runtime NAME=Runtime
.PHONY: mock
mock: $(GOBIN)/mockery
	@mockery

# Runs cript to upload codecov coverage data
.PHONY: upload-coverage
upload-coverage:
	@./scripts/codecov.sh -t $(CODECOV_TOKEN)

.PHONY: cur-version
cur-version:
	@echo -n $(VERSION)

$(GOBIN)/mockery:
	@go install github.com/vektra/mockery/v2@v2.42.0
	@mockery --version

$(GOBIN)/golangci-lint:
	@echo installing: golangci-lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.56.2
