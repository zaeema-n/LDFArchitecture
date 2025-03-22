// import ballerina/io;
// import ballerina/protobuf.types.'any as pbAny;


// CrudServiceClient ep = check new ("http://localhost:50051");

// function unwrapAny(pbAny:Any anyValue) returns string|error {
//     // Extract the base64-encoded value from the Any object
//     // Unpack the value as a string
//     string|error result = pbAny:unpack(anyValue, string);
//     return result;
// }

// public function main() returns error? {
//     // Create the metadata array
//     record {| string key; pbAny:Any value; |}[] metadataArray = [];

//     // Pack string values into protobuf.Any directly
//     pbAny:Any packedValue1 = check pbAny:pack("value1");
//     pbAny:Any packedValue2 = check pbAny:pack("value2");

//     // Add packed values to the metadata array
//     metadataArray.push({key: "key1", value: packedValue1});
//     metadataArray.push({key: "key2", value: packedValue2});

//     Entity createEntityRequest = {
//         id: "ballerina-1",
//         kind: {
//             major: "ballerina",
//             minor: "ballerina"
//         },
//         created: "ballerina",
//         terminated: "ballerina",
//         name: {
//             startTime: "ballerina",
//             endTime: "ballerina",
//             value: check pbAny:pack("ballerina")
//         },
//         metadata: metadataArray
//     };

//     Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
//     io:println(createEntityResponse);
//     io:println("--------------------------------");
//     EntityId readEntityRequest = {id: "ballerina-1"};
//     Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
//     io:println("--------------------------------");
//     io:println(readEntityResponse);

//     foreach var item in readEntityResponse.metadata {
//         // Unpack the value as a StringValue
//         var unwrapped = unwrapAny(item.value);
//         if unwrapped is string {
//             io:println("Unwrapped value: ", unwrapped);
//         } else {
//             io:println("Error: ", unwrapped);
//         }
//     }

//     io:println("--------------------------------");

//     Entity updateEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: check 'any:pack("ballerina")}};
//     Entity updateEntityResponse = check ep->UpdateEntity(updateEntityRequest);
//     io:println(updateEntityResponse);

//     EntityId deleteEntityRequest = {id: "ballerina"};
//     Empty deleteEntityResponse = check ep->DeleteEntity(deleteEntityRequest);
//     io:println(deleteEntityResponse);
// }


// import ballerina/io;

// CrudServiceClient ep = check new ("http://localhost:9090");

// public function main() returns error? {
//     Entity createEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: check 'any:pack("ballerina")}};
//     Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
//     io:println(createEntityResponse);

//     EntityId readEntityRequest = {id: "ballerina"};
//     Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
//     io:println(readEntityResponse);

//     UpdateEntityRequest updateEntityRequest = {id: "ballerina", entity: {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: check 'any:pack("ballerina")}}};
//     Entity updateEntityResponse = check ep->UpdateEntity(updateEntityRequest);
//     io:println(updateEntityResponse);

//     EntityId deleteEntityRequest = {id: "ballerina"};
//     Empty deleteEntityResponse = check ep->DeleteEntity(deleteEntityRequest);
//     io:println(deleteEntityResponse);
// }
