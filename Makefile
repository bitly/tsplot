.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	make -C tscli

.PHONY: install
install: build
	cp ./bin/tscli /usr/local/bin/

.PHONY: clean
clean:
	rm -rf ./bin
