build:
	go build -o sgv .

clean:
	rm -f sgv

install:
	go install .

.PHONY: build clean install
