all: build

build:
	go build ./cmd...

image: 
	docker build . -t logconv

clean: 
	rm -f logconv	
	rm -f test-artifacts/*

test: unit-test

unit-test:
	go test -v ./...

