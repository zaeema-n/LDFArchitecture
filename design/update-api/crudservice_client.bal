import ballerina/io;
import ballerina/protobuf.types.'any as pbAny;

CrudServiceClient ep = check new ("http://localhost:50051");

public function main() returns error? {
    // Create the metadata array
    record {| string key; pbAny:Any value; |}[] metadataArray = [];

    // Pack string values into protobuf.Any directly
    pbAny:Any packedValue1 = check pbAny:pack("value1");
    pbAny:Any packedValue2 = check pbAny:pack("value2");

    // Add packed values to the metadata array
    metadataArray.push({key: "key1", value: packedValue1});
    metadataArray.push({key: "key2", value: packedValue2});

    Entity createEntityRequest = {
        id: "ballerina-1",
        kind: {
            major: "ballerina",
            minor: "ballerina"
        },
        created: "ballerina",
        terminated: "ballerina",
        name: {
            startTime: "ballerina",
            endTime: "ballerina",
            value: check pbAny:pack("ballerina")
        },
        metadata: metadataArray
    };

    Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
    io:println(createEntityResponse);
    io:println("--------------------------------");
    EntityId readEntityRequest = {id: "ballerina-1"};
    Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
    io:println("--------------------------------");
    io:println(readEntityResponse);

    // Enumerate and unpack metadata
    foreach var item in readEntityResponse.metadata {
        // Unpack the value as a string
        io:println(item.value);
    }

    io:println("--------------------------------");

    // Entity updateEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: check 'any:pack("ballerina")}};
    // Entity updateEntityResponse = check ep->UpdateEntity(updateEntityRequest);
    // io:println(updateEntityResponse);

    // EntityId deleteEntityRequest = {id: "ballerina"};
    // Empty deleteEntityResponse = check ep->DeleteEntity(deleteEntityRequest);
    // io:println(deleteEntityResponse);
}