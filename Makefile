test:
	@go test -v ./...

build:
	@go build -o bin/BlockC


run: build
	./bin/docker
