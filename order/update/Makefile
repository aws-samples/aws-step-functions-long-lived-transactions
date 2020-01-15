.PHONY: deps clean build

default: build

deps:
	go get .

clean:
	rm -rf ./build

build:
	GOOS=linux GOARCH=amd64 go build -o build/main .