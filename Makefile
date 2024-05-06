run: 
	@go run cmd/main.go

templ:
	@templ generate -watch -proxy=http://localhost:42069

build:
	@go build -o bin/app cmd/main.go

