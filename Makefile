build:
	@go build -o bin/FileSystem main.go

run: build
	@./bin/FileSystem

test:
	@go test -v ./...