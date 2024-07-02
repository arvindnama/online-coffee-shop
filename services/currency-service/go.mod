module github.com/arvindnama/golang-microservices/currency-service

go 1.22.3

require (
	github.com/arvindnama/golang-microservices/libs/grpc-protos v0.0.0-00010101000000-000000000000
	github.com/arvindnama/golang-microservices/libs/utils v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-hclog v1.6.3
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.64.0
)

require (
	github.com/fatih/color v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

replace github.com/arvindnama/golang-microservices/libs/utils => ../../libs/utils

replace github.com/arvindnama/golang-microservices/libs/grpc-protos => ../../libs/grpc-protos
