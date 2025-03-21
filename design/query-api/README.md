# Query API

## Implement OpenAPI Contract

```bash
bal openapi -i ../contracts/rest/query_api.yaml --mode service
```

## Generate GRPC Stubs

The client stub generated here will be sending and receiving values via Grpc. 
This will send requests to the corresponding CRUD server endpoint. 

```bash
bal grpc --mode client --input ../crud-api/protos/types_v1.proto --output .
```
