# crud-api

## Pre-requisites

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go get google.golang.org/grpc
```

## Development

### Generate Go Code (Developer)

```bash
cd crud-api
```

```bash
go mod init github.com/zaeema-n/LDFArchitecture/design/crud-api
```

```bash
protoc --go_out=. --go-grpc_out=.--proto_path=protos protos/types.proto
```

## Usage

### Run Service

```bash
./crud-service
```

The service will be running in port `50051` and it is hard coded. This needs to be configurable. 