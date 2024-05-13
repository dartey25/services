ifeq ($(OS),Windows_NT)
	GOOS=windows
	BINARY_NAME=services.exe
else
	GOOS=linux
	BINARY_NAME=services
endif

run: 
	@go run cmd/main.go

templ:
	@templ generate -v
	@#templ generate -watch -proxy=http://localhost:42069

build: templ
	@GOOS="$(GOOS)" GOARCH="amd64" go build -o bin/$(BINARY_NAME) cmd/main.go

