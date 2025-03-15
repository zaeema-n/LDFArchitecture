# crud-api

## Pre-requisites

```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   go get google.golang.org/grpc
```

## Generate Go Code

```bash
go mod init github.com/zaeema-n/LDFArchitecture/tree/main/design/crud-api
```

```bash
   protoc --go_out=. --go-grpc_out=. --proto_path=protos protos/types.proto
```