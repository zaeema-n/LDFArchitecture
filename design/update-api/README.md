# Update API

## Generate Open API 

```bash
bal openapi -i ../contracts/rest/update_api.yaml  --mode service
```

## Generate GRPC Stubs

```bash
bal grpc --mode client --input ../crud-api/protos/types_v1.proto --output .
```

> 💡 **Note**  
> At the generation make sure to remove any sample code generated to show how to use the API. Because that might add an unnecessary main file. 

## Run Service

```bash
cd update-api
bal run
```

At the moment the port is hardcoded to 8080. This must be configurable via a config file.

