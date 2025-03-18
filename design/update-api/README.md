# Update API

## Generate Open API 

```bash
bal openapi -i ../contracts/rest/update_api.yaml --mode service
```

## Generate GRPC Stubs

```bash
bal grpc --mode client --input ../crud-api/protos/types_v1.proto --output .
```

> 💡 **Note**  
> At the generation make sure to remove any sample code generated to show how to use the API. Because that might add an unnecessary main file. 

## Run Test

```bash
# Run all tests in the current package
bal test

# Run tests with verbose output
bal test --test-report

# Run a specific test file
bal test tests/service_test.bal

# Run a specific test function
bal test --tests testMetadataHandling

# Run tests and generate a coverage report
bal test --code-coverage
```

## Run Service

```bash
cd update-api
bal run
```

At the moment the port is hardcoded to 8080. This must be configurable via a config file.

