BINARY_NAME=order-api-service
.DEFAULT_GOAL := run-dev

build:
	 GOOS=linux GOARCH=amd64 go build -o ./.bin/${BINARY_NAME}-linux main.go
	 GOOS=darwin GOARCH=arm64 go build -o ./.bin/${BINARY_NAME}-darwin main.go

build-migration-tool:
	 GOOS=linux GOARCH=amd64 go build -o ./.bin/${BINARY_NAME}-migration-linux migration-main.go
	 GOOS=darwin GOARCH=arm64 go build -o ./.bin/${BINARY_NAME}-migration-darwin migration-main.go
	 
build-docker-image: build build-migration-tool
	docker build -t order-api-service:latest .

run-dev:
	ls **/*.go | entr -r go run main.go

run: build
	./.bin/${BINARY_NAME}-darwin

clean:
	go clean
	rm -rf ./.bin/${BINARY_NAME}-darwin
	rm -rf ./.bin/${BINARY_NAME}-linux


check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	swagger generate spec -o ./swagger.yaml --scan-models

generate_client:
	cd sdk && swagger generate client -f ../swagger.yaml -A order-api

migration:
	@migrate create -ext sql -dir migrate/migrations $(filter-out $@,$(MAKECMDGOALS)) 

migrate-up:
	@go run migration-main.go up

migrate-down:
	@go run migration-main.go down