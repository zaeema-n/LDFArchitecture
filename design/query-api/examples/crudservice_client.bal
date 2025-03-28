// import ballerina/io;

// CrudServiceClient ep = check new ("http://localhost:9090");

// public function main() returns error? {
//     Entity createEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: check 'any:pack("ballerina")}};
//     Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
//     io:println(createEntityResponse);

//     Entity readEntityRequest = {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: check 'any:pack("ballerina")}};
//     Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
//     io:println(readEntityResponse);

//     UpdateEntityRequest updateEntityRequest = {id: "ballerina", entity: {id: "ballerina", kind: {major: "ballerina", minor: "ballerina"}, created: "ballerina", terminated: "ballerina", name: {startTime: "ballerina", endTime: "ballerina", value: check 'any:pack("ballerina")}}};
//     Entity updateEntityResponse = check ep->UpdateEntity(updateEntityRequest);
//     io:println(updateEntityResponse);

//     EntityId deleteEntityRequest = {id: "ballerina"};
//     Empty deleteEntityResponse = check ep->DeleteEntity(deleteEntityRequest);
//     io:println(deleteEntityResponse);
// }
