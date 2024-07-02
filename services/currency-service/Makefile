BINARY_NAME=currency-service
.DEFAULT_GOAL := run-dev

build:
	 GOOS=linux GOARCH=amd64 go build -o ./.bin/${BINARY_NAME}-linux main.go
	 GOOS=darwin GOARCH=arm64 go build -o ./.bin/${BINARY_NAME}-darwin main.go
	 
run-dev:
	ls **/*.go | entr -r go run .

run: build
	./.bin/${BINARY_NAME}-darwin

clean:
	go clean
	rm -rf ./.bin/${BINARY_NAME}-darwin
	rm -rf ./.bin/${BINARY_NAME}-linux
