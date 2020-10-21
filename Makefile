.PHONY: all build test-unit

all: build

build: main.go
	go build

test-unit:
	go test -race -coverprofile=coverage.out
