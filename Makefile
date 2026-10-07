GO ?= go
MOCKERY_VERSION := v3.8.0

.PHONY: build test mock-gen

build:
	mkdir -p bin
	$(GO) build -o bin/shipping-rate ./cmd/cli

test:
	$(GO) test -race -cover ./...

mock-gen:
	$(GO) run github.com/vektra/mockery/v3@$(MOCKERY_VERSION) --config .mockery.yml
