include Makefile.Common

GO_MODULE_DIRS := $(shell find . -type f -name go.mod -not -path '*/vendor/*' -exec dirname {} \; | sort)

.PHONY: build ci for-all-target gofmt golint gotidy ktlo test vet

build:
	go build ./...

ci: checklicense fmt-check vet test build

fmt-check:
	@files="$$(find . -name '*.go' -not -path './.git/*' -print)"; \
	if [ -n "$$files" ]; then \
		unformatted="$$(gofmt -l $$files)"; \
		if [ -n "$$unformatted" ]; then \
			printf '%s\n' "$$unformatted"; \
			exit 1; \
		fi; \
	fi

for-all-target:
	@set -e; \
	for dir in $(GO_MODULE_DIRS); do \
		if [ ! -f "$$dir/Makefile" ]; then \
			echo "missing Makefile for Go module $$dir"; \
			exit 1; \
		fi; \
		echo "running $(TARGET) in $$dir"; \
		$(MAKE) --no-print-directory -C "$$dir" "$(TARGET)"; \
	done

gofmt:
	@$(MAKE) --no-print-directory for-all-target TARGET=fmt

golint:
	@$(MAKE) --no-print-directory for-all-target TARGET=lint

gotidy:
	@$(MAKE) --no-print-directory for-all-target TARGET=tidy

test:
	go test ./...

vet:
	go vet ./...

ktlo:
	mkdir -p ./bin
	go build -o ./bin/ktlo .
