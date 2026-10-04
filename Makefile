BINARY ?= arka

.PHONY: build test run install fmt
build:
	go build -trimpath -o bin/$(BINARY) ./cmd/arka

test:
	go test ./...

run:
	go run ./cmd/arka start --no-open

install:
	go install ./cmd/arka

fmt:
	gofmt -w $$(find cmd internal -name "*.go" -type f)
