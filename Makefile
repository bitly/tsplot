SHELL := /bin/bash
GOBIN := $(shell which go 2> /dev/null)
ifeq ($(GOBIN),)
GOBIN := /usr/local/go/bin/go
endif

.PHONY: test
test:
	$(GOBIN) test -v ./...

