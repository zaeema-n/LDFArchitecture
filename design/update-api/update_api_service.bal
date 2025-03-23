// AUTO-GENERATED FILE.
// This file is auto-generated by the Ballerina OpenAPI tool.

import ballerina/http;
import ballerina/protobuf.types.'any as pbAny;
import ballerina/io;


listener http:Listener ep0 = new (8080, config = {host: "localhost"});

CrudServiceClient ep = check new ("http://localhost:50051");

service / on ep0 {
    # Delete an entity
    #
    # + return - Entity deleted 
    resource function delete entities/[string id]() returns http:NoContent|error {
        var result = ep->DeleteEntity({id: id});
        if result is error {
            return result;
        }
        return http:NO_CONTENT;
    }

    # Create a new entity
    #
    # + return - Entity created 
    resource function post entities(@http:Payload json jsonPayload) returns Entity|error {
        // Convert JSON to Entity with custom mapping
        Entity payload = check convertJsonToEntity(jsonPayload);
        
        io:println(payload);
        var result = ep->CreateEntity(payload);
        if result is error {
            return result;
        }
        return result;
    }

    # Update an existing entity
    #
    # + return - Entity updated 
    resource function put entities/[string id](@http:Payload json jsonPayload) returns Entity|error {
        // Convert JSON to Entity with custom mapping
        Entity payload = check convertJsonToEntity(jsonPayload);

        // Create UpdateEntityRequest with both id from URL and entity from payload
        UpdateEntityRequest updateRequest = {
            id: id,
            entity: payload
        };
        
        var result = ep->UpdateEntity(updateRequest);
        if result is error {
            return result;
        }
        return result;
    }

    # TODO: Remove/Don't expose this endpoint from Ingest API (has been only added for testing purposes)
    # Read an entity by ID
    #
    # + id - The ID of the entity to retrieve
    # + return - The entity or an error
    resource function get entities/[string id]() returns Entity|error {
        // Call the ReadEntity function with the ID
        EntityId readEntityRequest = {id: id};
        Entity|error result = ep->ReadEntity(readEntityRequest);
        
        if result is error {
            return result;
        }
        
        // Successfully retrieved the entity
        io:println("Retrieved entity with ID: ", id);
        return result;
    }
}

// Helper function to convert JSON to Entity
function convertJsonToEntity(json jsonPayload) returns Entity|error {
    // Check if metadata is present and handle accordingly
    record {| string key; pbAny:Any value; |}[] metadataArray = [];

    if jsonPayload.metadata != () {
        if jsonPayload?.metadata is json[] {
            json[] metadataJsonArray = <json[]>check jsonPayload.metadata;
            foreach json item in metadataJsonArray {
                string key = (check item.key).toString();
                pbAny:Any packedValue = check pbAny:pack((check item.value).toString());
                metadataArray.push({key: key, value: packedValue});
            }
        } else if jsonPayload?.metadata is map<json> {
            map<json> metadataMap = <map<json>>check jsonPayload.metadata;
            foreach var [key, val] in metadataMap.entries() {
                pbAny:Any packedValue = check pbAny:pack(val.toString());
                metadataArray.push({key: key, value: packedValue});
            }
        }
    }
    
    // Process attributes if present
    record {| string key; TimeBasedValueList value; |}[] attributesArray = [];
    
    if jsonPayload.attributes != () {
        if jsonPayload?.attributes is json[] {
            json[] attributesJsonArray = <json[]>check jsonPayload.attributes;
            foreach json item in attributesJsonArray {
                string key = (check item.key).toString();
                
                // Add safe type checking for value
                json valueJson = check item.value;
                TimeBasedValue[] timeBasedValues = [];
                
                if valueJson is json[] {
                    json[] valuesJson = <json[]>valueJson;
                    foreach json valueItem in valuesJson {
                        TimeBasedValue tbv = {
                            startTime: (check valueItem.startTime).toString(),
                            endTime: valueItem?.endTime is () ? "" : (check valueItem.endTime).toString(),
                            value: check pbAny:pack((check valueItem.value).toString())
                        };
                        timeBasedValues.push(tbv);
                    }
                } else {
                    // Handle the case when value is not an array (could be a single object)
                    // Create a single TimeBasedValue object
                    TimeBasedValue tbv = {
                        startTime: valueJson?.startTime is () ? "" : (check valueJson.startTime).toString(),
                        endTime: valueJson?.endTime is () ? "" : (check valueJson.endTime).toString(),
                        value: check pbAny:pack(valueJson?.value is () ? "" : (check valueJson.value).toString())
                    };
                    timeBasedValues.push(tbv);
                }
                
                TimeBasedValueList tbvList = {
                    values: timeBasedValues
                };
                
                attributesArray.push({key: key, value: tbvList});
            }
        } else if jsonPayload?.attributes is map<json> {
            map<json> attributesMap = <map<json>>check jsonPayload.attributes;
            foreach var [key, val] in attributesMap.entries() {
                TimeBasedValue[] timeBasedValues = [];
                
                // Add safe type checking for val
                if val is json[] {
                    json[] valuesJson = <json[]>val;
                    foreach json valueItem in valuesJson {
                        TimeBasedValue tbv = {
                            startTime: (check valueItem.startTime).toString(),
                            endTime: valueItem?.endTime is () ? "" : (check valueItem.endTime).toString(),
                            value: check pbAny:pack((check valueItem.value).toString())
                        };
                        timeBasedValues.push(tbv);
                    }
                } else {
                    // Handle the case when val is not an array
                    TimeBasedValue tbv = {
                        startTime: val?.startTime is () ? "" : (check val.startTime).toString(),
                        endTime: val?.endTime is () ? "" : (check val.endTime).toString(),
                        value: check pbAny:pack(val?.value is () ? "" : (check val.value).toString())
                    };
                    timeBasedValues.push(tbv);
                }
                
                TimeBasedValueList tbvList = {
                    values: timeBasedValues
                };
                
                attributesArray.push({key: key, value: tbvList});
            }
        }
    }
    
    // Process relationships if present
    record {| string key; Relationship value; |}[] relationshipsArray = [];
    
    if jsonPayload.relationships != () {
        if jsonPayload?.relationships is json[] {
            json[] relationshipsJsonArray = <json[]>check jsonPayload.relationships;
            foreach json item in relationshipsJsonArray {
                string key = (check item.key).toString();
                json relJson = check item.value;
                
                Relationship relationship = {
                    relatedEntityId: (check relJson.relatedEntityId).toString(),
                    startTime: (check relJson.startTime).toString(),
                    endTime: relJson?.endTime is () ? "" : (check relJson.endTime).toString(),
                    id: (check relJson.id).toString(),
                    name: (check relJson.name).toString()
                };
                
                relationshipsArray.push({key: key, value: relationship});
            }
        } else if jsonPayload?.relationships is map<json> {
            map<json> relationshipsMap = <map<json>>check jsonPayload.relationships;
            foreach var [key, val] in relationshipsMap.entries() {
                Relationship relationship = {
                    relatedEntityId: (check val.relatedEntityId).toString(),
                    startTime: (check val.startTime).toString(),
                    endTime: val?.endTime is () ? "" : (check val.endTime).toString(),
                    id: (check val.id).toString(),
                    name: (check val.name).toString()
                };
                
                relationshipsArray.push({key: key, value: relationship});
            }
        }
    }
    
    // Create the entity with proper type reference
    Entity entity = {
        id: (check jsonPayload.id).toString(),
        kind: {
            major: (check jsonPayload.kind.major).toString(),
            minor: (check jsonPayload.kind.minor).toString()
        },
        created: (check jsonPayload.created).toString(),
        terminated: jsonPayload?.terminated is () ? "" : (check jsonPayload.terminated).toString(),
        name: {
            startTime: (check jsonPayload.name.startTime).toString(),
            endTime: jsonPayload?.name?.endTime is () ? "" : (check jsonPayload.name.endTime).toString(),
            value: check pbAny:pack((check jsonPayload.name.value).toString())
        },
        metadata: metadataArray,
        attributes: attributesArray,
        relationships: relationshipsArray
    };
    
    return entity;
}
