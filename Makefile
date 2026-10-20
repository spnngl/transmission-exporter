GO ?= go
export CGO_ENABLED = 0

.PHONY: all
all: install

# Every check CI runs; lint includes gofmt.
.PHONY: check
check: lint tidy-check vet test build

.PHONY: clean
clean:
	$(GO) clean -i ./...

.PHONY: install
install:
	$(GO) install -v ./cmd/transmission-exporter

.PHONY: build
build:
	$(GO) build -v ./cmd/transmission-exporter

.PHONY: test
test:
	CGO_ENABLED=1 $(GO) test -race ./...

.PHONY: tidy-check
tidy-check:
	$(GO) mod tidy -diff

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: dashboards
dashboards:
	jsonnet fmt -i dashboards/transmission.jsonnet
	jsonnet -J dashboards/vendor -m dashboards -e "(import 'dashboards/transmission.jsonnet').grafanaDashboards"
