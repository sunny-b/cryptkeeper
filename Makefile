ROOT:=$(shell git rev-parse --show-toplevel)

.PHONY: build
build:
	@VERSION=$$(cat VERSION); \
	go build -ldflags "-X github.com/sunny-b/cryptkeeper/internal/version.Version=v$${VERSION}" -o ./bin/cryptkeeper ./cmd/cryptkeeper/*.go 

############################################################################
# Testing
############################################################################

tests = \
				test-go \
				test-go-lint \
				test-bash \
				test-fish \
				test-zsh

.PHONY: test-go
test-go:
	go test -coverprofile=coverage.out -covermode=atomic -v -race ./...

.PHONY: $(tests)
test: build $(tests)
	@echo
	@echo SUCCESS!

.PHONY: test-go-lint
test-go-lint:
	golangci-lint run

.PHONY: test-bash
test-bash:
	bash ${ROOT}/test/ck-test.bash

.PHONY: test-fish
test-fish:
	fish ${ROOT}/test/ck-test.fish

.PHONY: test-zsh
test-zsh:
	zsh ${ROOT}/test/ck-test.zsh

docker-test:
	docker build -t test && docker run --rm test