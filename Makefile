VERSION = $(shell gobump show -r .)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X main.revision=$(CURRENT_REVISION)"
DIST_DIR = dist
u := $(if $(update),-u)

.PHONY: default
default: test

.PHONY: test
test:
	go test

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) .

.PHONY: install
install:
	go install -ldflags=$(BUILD_LDFLAGS) .

.PHONY: deps
deps:
	go get ${u}
	go mod tidy

.PHONY: devel-deps
devel-deps:
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	go install github.com/Songmu/goxz/cmd/goxz@latest
	go install github.com/x-motemen/gobump/cmd/gobump@latest
	go install github.com/tcnksm/ghr@latest

.PHONY: credits
credits: go.sum deps devel-deps
	gocredits -skip-missing -w

.PHONY: cross-build
cross-build: credits
	rm -rf $(DIST_DIR)
	goxz -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) -static -d $(DIST_DIR) .

.PHONY: release
release: cross-build
	ghr v$(VERSION) $(DIST_DIR)
