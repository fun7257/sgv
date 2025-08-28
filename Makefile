build:
	go build -o sgv .

clean:
	rm -f sgv

install:
	go install .

local: build
	sudo cp sgv /usr/local/bin/sgv

.PHONY: build clean install
