// Example client code to test the CRUD service

// import ballerina/io;

// CrudServiceClient ep = check new ("http://localhost:9090");

// public function main() returns error? {
//     Entity createEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: "ballerina"}};
//     Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
//     io:println(createEntityResponse);

//     Entity readEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: "ballerina"}};
//     Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
//     io:println(readEntityResponse);

//     Entity updateEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: "ballerina"}};
//     Entity updateEntityResponse = check ep->UpdateEntity(updateEntityRequest);
//     io:println(updateEntityResponse);

//     Entity deleteEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: "ballerina"}};
//     Entity deleteEntityResponse = check ep->DeleteEntity(deleteEntityRequest);
//     io:println(deleteEntityResponse);
// }
