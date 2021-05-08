BINDIR ?= ../bin

default: build

.PHONY: build
build:
	@mkdir -p $(BINDIR)
	GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o $(BINDIR)/tscli
