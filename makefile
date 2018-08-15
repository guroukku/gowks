default: test

build:
	go build -ldflags="-s -w" -v -o ./bin/gowks 

run: build
	./bin/gowks

test:
	go test -v ./...
