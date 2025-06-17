build:
	go build -o sgv .

clean:
	rm -f sgv

install:
	go install -o sgv .

.PHONY: build clean install
