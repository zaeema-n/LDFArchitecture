# Update API

## Generate Stubs

```bash
bal grpc --mode client --input ../crud-api/protos/types.proto --output client/
```

> 💡 **Note**  
> At the generation make sure to remove any sample code generated to show how to use the API. Because that might add an unnecessary main file. 

## Run Service

```bash
cd update-api
bal run
```

At the moment the port is hardcoded to 8080. This must be configurable via a config file.

