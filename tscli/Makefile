BINDIR ?= ../bin
SHELL := /bin/bash
GOBIN := $(shell which go 2> /dev/null)
ifeq ($(GOBIN),)
GOBIN := /usr/local/go/bin/go
endif

default: build

.PHONY: build
build:
	@mkdir -p $(BINDIR)
	$(GOBIN) build -o $(BINDIR)/tscli
