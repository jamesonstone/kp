BIN_DIR := ./bin
BINARY := $(BIN_DIR)/kp
CMD := ./cmd/kp
PREFIX ?= /usr/local
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || printf unknown)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)

.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

.PHONY: test
test:
	go test ./...

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" $(CMD)

.PHONY: fmt
fmt:
	gofmt -w prompts.go cmd internal

.PHONY: clean
clean:
	rm -f $(BINARY)
