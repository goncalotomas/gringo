MAIN=.
TAG_VERSION=$(shell git describe --tags --exact-match 2>/dev/null || git rev-parse --short HEAD)
VERSION=$(subst /,_,$(TAG_VERSION))

all: build
	@echo all built

build:
	 go build $(MAIN)

clean:
	@rm -rf gringo target
	@echo all removed

lint:
	golangci-lint run ./...
	@echo all code is linted

format:
	gofmt -w -s .

format-check:
	gofmt -l .

run:
	go run $(MAIN) -loglevel DEBUG

test:
	go test -v ./...

sbom: build
	@mkdir -p target/sbom
	cyclonedx-gomod bin -json -output ./target/sbom/gringo-$(VERSION).bom.json ./gringo

.PHONY: all build clean lint format format-check run test sbom
