GO_VERSION := $(shell go version | awk '{print $$3}')
COMMIT := $(shell git rev-parse HEAD)

build:
	go build -ldflags "-X main.goVersion=$(GO_VERSION) -X main.commit=$(COMMIT)" -o sgv .

clean:
	rm -f sgv

install:
	go install -ldflags "-X main.goVersion=$(GO_VERSION) -X main.commit=$(COMMIT)" .

.PHONY: build clean install
