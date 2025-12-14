build:
	go build -v
.PHONY: build

test:
	go test -v -cover -timeout=15m
.PHONY: test

clean:
	go clean
.PHONY: clean

all: clean build test
.PHONY: all
