module github.com/arvindnama/golang-microservices/product-images-service

go 1.22.3

require (
	github.com/arvindnama/golang-microservices/libs/utils v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/hashicorp/go-hclog v1.6.3
	github.com/joho/godotenv v1.5.1
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028
)

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.17.0 // indirect
)

replace github.com/arvindnama/golang-microservices/libs/utils => ../../libs/utils
