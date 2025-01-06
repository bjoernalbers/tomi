release:
	goreleaser release --clean
.PHONY: release

test:
	go test ./...
.PHONY: test
