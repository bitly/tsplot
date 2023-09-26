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
build: $(BINDIR)/codegen
	make -C tscli

$(BINDIR)/codegen: $(SCRIPTDIR)/codegen.go
	@mkdir -p $(BINDIR)
	GOOS=linux GOARCH=amd64 $(GOBIN) build -o $(BINDIR)/codegen $<

.PHONY: install
install: build
	cp $(BINDIR)/tscli /usr/local/bin/

.PHONY: clean
clean:
	rm -rf $(BINDIR)

CODEGENFILE="set_aggregation_opts.go"
.PHONY: codegen
codegen:
	$(BINDIR)/codegen -output ./tsplot/$(CODEGENFILE)
	go fmt ./tsplot/$(CODEGENFILE)

.PHONY: codegendiff
codegendiff:
	diff ./tsplot/$(CODEGENFILE) <($(BINDIR)/codegen -stdout | gofmt)
