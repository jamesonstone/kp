BINARY=kp

.PHONY: build test install release

build:
go build ./...

test:
go test ./...

install:
go install ./cmd/kp

release:
goreleaser release --snapshot --clean
