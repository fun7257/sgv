test:
	go test ./...

build: test
	go build -o sgv .
	sudo cp sgv /usr/local/bin/sgv

clean:
	rm -f sgv

.PHONY: build clean test
