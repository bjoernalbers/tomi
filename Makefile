all: unit
.PHONY: all

release:
	goreleaser release --clean
.PHONY: release

unit:
	go test ./...
.PHONY: unit
