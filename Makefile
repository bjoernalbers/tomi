VERSION := $(shell git describe --tags | tr -d v)

all: test
.PHONY: all

tomi: ## Build tomi binary
	for arch in arm64 amd64; do GOARCH="$${arch}" go build -o "$@-$${arch}" -ldflags '-X main.version=$(VERSION)'; done
	lipo $@-* -create -output $@
	rm $@-*
.PHONY: tomi

release: ## Publish new release
	goreleaser release --clean
.PHONY: release

test: ## Run tests
	go test ./...
.PHONY: test

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help
