# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Run go fmt against code
fmt: gci addlicense
	go mod tidy
	go fmt ./...
	find . -type f -name '*.go' -a ! -name '*zz_generated*' -exec $(GCI) write -s standard -s default -s "prefix(github.com/liqotech/liqo)" {} \;
	find . -type f -name '*.go' -exec $(ADDLICENSE) -l apache -c "cheina97" -y "2023-$(shell date +%Y)" {} \;

lint: golangci-lint
	$(GOLANGCILINT) run --new

# Install addlicense if not available
addlicense:
ifeq (, $(shell which addlicense))
	@go install github.com/google/addlicense@v1.0.0
ADDLICENSE=$(GOBIN)/addlicense
else
ADDLICENSE=$(shell which addlicense)
endif

# Install gci if not available
gci:
ifeq (, $(shell which gci))
	@go install github.com/daixiang0/gci@v0.7.0
GCI=$(GOBIN)/gci
else
GCI=$(shell which gci)
endif

# Install golangci-lint if not available
golangci-lint:
ifeq (, $(shell which golangci-lint))
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.0
GOLANGCILINT=$(GOBIN)/golangci-lint
else
GOLANGCILINT=$(shell which golangci-lint)
endif

build: 
	$(eval GIT_COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown"))
	go build -ldflags="-s -w -X 'main.version=$(GIT_COMMIT)'" -o gowatch cmd/main.go 