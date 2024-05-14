ifeq ($(OS),Windows_NT)
	GOOS=windows
	BINARY_NAME=services.exe
	AIR_COMMAND=air -c .air.win.toml
else
	GOOS=linux
	BINARY_NAME=services
	AIR_COMMAND=air
endif

run: 
	@go run cmd/main.go

templ:
	@templ generate -v

build: templ
	@GOOS="$(GOOS)" GOARCH="amd64" go build -o bin/$(BINARY_NAME) cmd/main.go

build-win: templ
	@GOOS="windows" GOARCH="amd64" go build -o bin/services.exe cmd/main.go

air: 
	@$(AIR_COMMAND)

