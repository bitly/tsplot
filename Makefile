BINDIR ?= ./bin
SCRIPTDIR := ./scripts
SHELL := /bin/bash
GOBIN := $(shell which go 2> /dev/null)
ifeq ($(GOBIN),)
GOBIN := /usr/local/go/bin/go
endif

.PHONY: test
test:
	$(GOBIN) test -v ./...

.PHONY: build
build:
	make -C tscli

.PHONY: install
install: build
	cp $(BINDIR)/tscli /usr/local/bin/

.PHONY: clean
clean:
	rm -rf $(BINDIR)

