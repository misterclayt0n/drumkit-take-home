.PHONY: help run build test fmt tidy clean

BINARY := bin/api

help:
	@echo "Available targets:"
	@echo "  make run    - start the Go API server"
	@echo "  make build  - compile the Go API server to $(BINARY)"
	@echo "  make test   - run Go tests"
	@echo "  make fmt    - format Go code"
	@echo "  make tidy   - tidy Go modules"
	@echo "  make clean  - remove compiled artifacts"

run:
	go run ./cmd/api

build:
	@mkdir -p bin
	go build -o $(BINARY) ./cmd/api

test:
	go test ./...

fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')

tidy:
	go mod tidy

clean:
	rm -rf bin
