build:
	@go build -o bin/FileSystem cmd/main.go

run: build
	@./bin/FileSystem

test:
	@go test -v ./...