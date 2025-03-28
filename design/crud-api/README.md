# crud-api

## Pre-requisites

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go get google.golang.org/grpc
go get go.mongodb.org/mongo-driver/mongo
```

## Development

### Generate Go Code (Developer)

```bash
cd crud-api
```

## Go Module Setup

Initialize a Go module. 

```bash
go mod init github.com/zaeema-n/LDFArchitecture/design/crud-api
```

Generating the protobuf stubs for Go-lang
This will be used to write the data pipeline from CRUD server to the API server
via Grpc. 

```bash
protoc --go_out=. --go-grpc_out=. --proto_path=protos protos/types_v1.proto
```

### Build

```bash
go build ./...
go build -o crud-service cmd/server/service.go cmd/server/utils.go
```

## Usage

### Run Service

```bash
./crud-service
```

The service will be running in port `50051` and it is hard coded. This needs to be configurable. 