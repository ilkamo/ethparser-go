GOBIN    := $(GOPATH)/bin
PATH     := $(GOBIN):$(PATH)

GOLANGCI_VERSION=v1.57.2

TOOLS += golang.org/x/tools/cmd/goimports

.PHONY: tools
tools: $(TOOLS) golangci_lint

.PHONY: $(TOOLS)
$(TOOLS): %:
	cd /tmp && GOBIN=$(GOBIN) go get -u $*

.PHONY: golangci_lint
golangci_lint:
	cd /tmp && wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_VERSION}

.PHONY: lint
lint:
	$(info Running Go code checkers and linters)
	golangci-lint --version
	golangci-lint $(V) run -v -E goimports --timeout=5m

.PHONY: fmt
fmt:
	golangci-lint run -E gofumpt --fix ./...

test:
	go test --race -v ./...

test-e2e:
	go test --race -tags=e2e -v ./...

.PHONY: display_coverage
display_coverage:
	go test --race --cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
