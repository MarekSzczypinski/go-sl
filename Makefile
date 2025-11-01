.DEFAULT_GOAL := build

BINARY_NAME := sl
BUILD_DIR := ./bin

.PHONY: build fmt vet staticcheck govulncheck clean

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
ifeq (, $(shell which staticcheck))
	@echo "Installing staticcheck"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
endif
	staticcheck ./...

govulncheck:
ifeq (, $(shell which govulncheck))
	@echo "Installing govulncheck"
	@go install golang.org/x/vuln/cmd/govulncheck@latest
endif
	govulncheck ./...

build: fmt vet staticcheck govulncheck
	@mkdir -p $(BUILD_DIR)
	go build -v -o $(BUILD_DIR)/$(BINARY_NAME)

clean:
	@rm -fr $(BUILD_DIR)
	go clean
