# crud-api

## TODOs

## Note

⚠️ **Warning**  
Please do not commit generated protobuf files, please generate them at build time..

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

#### Run with Docker

`Dockerfile.crud` refers to just running the

Make sure to create a network for this work since we need every other service to be accessible hence
we place them in the same network. 

```bash
docker network create crud-network
```

```bash
docker build -t crud-service -f Dockerfile.crud .
```

```bash
docker run -d \
  --name crud-service \
  --network crud-network \
  -p 50051:50051 \
  -e NEO4J_URI=bolt://neo4j-local:7687 \
  -e NEO4J_USER=${NEO4J_USER} \
  -e NEO4J_PASSWORD=${NEO4J_PASSWORD} \
  -e MONGO_URI=${MONGO_URI} \
  crud-service
```

#### Validate 

```bash
brew install grpcurl
```

```bash
grpcurl -plaintext localhost:50051 list
```

Output

```bash
crud.CrudService
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
```

### Run Tests (Please Validate this work)

```bash
# Build the test image
docker build -t crud-service-test -f Dockerfile.test .

# Run the tests
docker run --rm \
  --add-host=host.docker.internal:host-gateway \
  crud-service-test
```
