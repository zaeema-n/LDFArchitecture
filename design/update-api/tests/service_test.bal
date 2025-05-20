import ballerina/io;
import ballerina/test;
import ballerina/protobuf.types.'any as pbAny;
import ballerina/http;
import ballerina/os;

// Get environment variables without fallback values
string testCrudHostname = os:getEnv("CRUD_SERVICE_HOST");
string testCrudPort = os:getEnv("CRUD_SERVICE_PORT");
string testUpdateHostname = os:getEnv("UPDATE_SERVICE_HOST");
string testUpdatePort = os:getEnv("UPDATE_SERVICE_PORT");

// Construct URLs using string concatenation
string testCrudServiceUrl = "http://" + testCrudHostname + ":" + testCrudPort;
string testUpdateServiceUrl = "http://" + testUpdateHostname + ":" + testUpdatePort;

// Before Suite Function
@test:BeforeSuite
function beforeSuiteFunc() {
    io:println("I'm the before suite function!");
    io:println("CRUD Service URL: " + testCrudServiceUrl);
    io:println("Update Service URL: " + testUpdateServiceUrl);
}

// After Suite Function
@test:AfterSuite
function afterSuiteFunc() {
    io:println("I'm the after suite function!");
}

// Helper function to unpack Any values to strings
function unwrapAny(pbAny:Any anyValue) returns string|error {
    return pbAny:unpack(anyValue, string);
}

@test:Config {}
function testMetadataHandling() returns error? {
    // Initialize the client
    CrudServiceClient ep = check new (testCrudServiceUrl);
    
    // Test data setup
    string testId = "test-entity-1";
    string expectedValue1 = "value1";
    string expectedValue2 = "value2";
    
    // Create the metadata array
    record {| string key; pbAny:Any value; |}[] metadataArray = [];

    // Pack string values into protobuf.Any directly
    pbAny:Any packedValue1 = check pbAny:pack(expectedValue1);
    pbAny:Any packedValue2 = check pbAny:pack(expectedValue2);

    // Add packed values to the metadata array
    metadataArray.push({key: "key1", value: packedValue1});
    metadataArray.push({key: "key2", value: packedValue2});

    // Create entity request
    Entity createEntityRequest = {
        id: testId,
        kind: {
            major: "test",
            minor: "metadata"
        },
        created: "2023-01-01",
        terminated: "",
        name: {
            startTime: "2023-01-01",
            endTime: "",
            value: check pbAny:pack("test-entity")
        },
        metadata: metadataArray
    };

    // Create entity
    Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
    io:println("Entity created with ID: " + createEntityResponse.id);
    io:println("Created entity metadata: ", createEntityResponse.metadata);
    
    // Read entity
    ReadEntityRequest readEntityRequest = {
        id: testId,
        entity: {
            id: testId,
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["metadata"]
    };
    io:println("ReadEntityRequest: ", readEntityRequest);
    Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
    io:println("Entity retrieved, verifying data...");
    io:println("Retrieved entity: ", readEntityResponse);
    io:println("Retrieved entity metadata: ", readEntityResponse.metadata);
    
    // Verify metadata values
    map<string> actualValues = {};
    foreach var item in readEntityResponse.metadata {
        string|error unwrapped = unwrapAny(item.value);
        if unwrapped is string {
            actualValues[item.key] = unwrapped.trim();
        } else {
            test:assertFail("Failed to unpack metadata value for key: " + item.key);
        }
    }
    
    // Assert the values match
    test:assertEquals(actualValues["key1"], expectedValue1, "Metadata value for key1 doesn't match");
    test:assertEquals(actualValues["key2"], expectedValue2, "Metadata value for key2 doesn't match");
    
    // Clean up
    Empty _ = check ep->DeleteEntity({id: testId});
    Empty _ = check ep->DeleteEntity({id: testId});
    io:println("Test entity deleted");
    
    return;
}

// TODO: Re-enable once the Result type response handling is added
// See: https://github.com/zaeema-n/LDFArchitecture/issues/23
@test:Config {
    enable: false
}
function testMetadataUnpackError() returns error? {
    // Test case to verify handling of non-existent entities
    CrudServiceClient ep = check new (testCrudServiceUrl);
    
    // Try to read a non-existent entity
    ReadEntityRequest readEntityRequest = {
        id: "non-existent-entity",
        entity: {
            id: "non-existent-entity",
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["metadata"]
    };
    Entity|error response = ep->ReadEntity(readEntityRequest);
    
    // Assert that we get an error for non-existent entity
    test:assertTrue(response is error, "Expected error for non-existent entity");
    
    return;
}

@test:Config {}
function testMetadataUpdating() returns error? {
    // Initialize the client
    CrudServiceClient ep = check new (testCrudServiceUrl);
    
    // Test data setup
    string testId = "test-entity-update";
    
    // Initial metadata values
    string initialValue1 = "initial-value1";
    string initialValue2 = "initial-value2";
    
    // Updated metadata values
    string updatedValue1 = "updated-value1";
    string updatedValue2 = "updated-value2";
    string newValue3 = "new-value3";
    
    // Create the initial metadata array
    record {| string key; pbAny:Any value; |}[] initialMetadataArray = [];
    pbAny:Any packedInitialValue1 = check pbAny:pack(initialValue1);
    pbAny:Any packedInitialValue2 = check pbAny:pack(initialValue2);
    initialMetadataArray.push({key: "key1", value: packedInitialValue1});
    initialMetadataArray.push({key: "key2", value: packedInitialValue2});

    // Create initial entity request
    Entity createEntityRequest = {
        id: testId,
        kind: {
            major: "test",
            minor: "update-metadata"
        },
        created: "2023-01-01",
        terminated: "",
        name: {
            startTime: "2023-01-01",
            endTime: "",
            value: check pbAny:pack("test-update-entity")
        },
        metadata: initialMetadataArray
    };

    // Create entity
    Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
    io:println("Entity created with ID: " + createEntityResponse.id);
    
    // Verify initial metadata
    ReadEntityRequest readEntityRequest = {
        id: testId,
        entity: {
            id: testId,
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["metadata"]
    };
    Entity initialReadResponse = check ep->ReadEntity(readEntityRequest);
    verifyMetadata(initialReadResponse.metadata, {"key1": initialValue1, "key2": initialValue2});
    io:println("Initial metadata verified");
    
    // Create updated metadata array
    record {| string key; pbAny:Any value; |}[] updatedMetadataArray = [];
    pbAny:Any packedUpdatedValue1 = check pbAny:pack(updatedValue1);
    pbAny:Any packedUpdatedValue2 = check pbAny:pack(updatedValue2);
    pbAny:Any packedNewValue3 = check pbAny:pack(newValue3);
    updatedMetadataArray.push({key: "key1", value: packedUpdatedValue1});
    updatedMetadataArray.push({key: "key2", value: packedUpdatedValue2});
    updatedMetadataArray.push({key: "key3", value: packedNewValue3});

    // Update entity request
    Entity updateEntityRequest = {
        id: testId,
        kind: {
            major: "test",
            minor: "update-metadata"
        },
        created: "2023-01-01",
        terminated: "",
        name: {
            startTime: "2023-01-01",
            endTime: "",
            value: check pbAny:pack("test-update-entity")
        },
        metadata: updatedMetadataArray
    };
    
    // Update entity
    UpdateEntityRequest updateRequest = {
        id: testId,
        entity: updateEntityRequest
    };
    Entity updateEntityResponse = check ep->UpdateEntity(updateRequest);
    io:println("Entity updated with ID: " + updateEntityResponse.id);
    
    // Verify updated metadata
    ReadEntityRequest updatedReadRequest = {
        id: testId,
        entity: {
            id: testId,
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["metadata"]
    };
    Entity updatedReadResponse = check ep->ReadEntity(updatedReadRequest);
    verifyMetadata(updatedReadResponse.metadata, {
        "key1": updatedValue1, 
        "key2": updatedValue2,
        "key3": newValue3
    });
    io:println("Updated metadata verified");
    
    // Clean up
    Empty _ = check ep->DeleteEntity({id: testId});
    Empty _ = check ep->DeleteEntity({id: testId});
    io:println("Test entity deleted");
    
    return;
}

// Helper function to verify metadata
function verifyMetadata(record {| string key; pbAny:Any value; |}[] metadata, map<string> expected) {
    map<string> actual = {};
    foreach var item in metadata {
        string|error unwrapped = unwrapAny(item.value);
        if unwrapped is string {
            actual[item.key] = unwrapped.trim();
        }
    }
    
    // Verify all expected key-value pairs exist in the actual metadata
    foreach var [key, expectedValue] in expected.entries() {
        test:assertTrue(actual.hasKey(key), "Metadata key not found: " + key);
        test:assertEquals(actual[key] ?: "", expectedValue, 
            string `Metadata value for ${key} doesn't match: expected ${expectedValue}, got ${actual[key] ?: ""}`);
    }
}

@test:Config {}
function testEntityReading() returns error? {
    // Initialize the client
    CrudServiceClient ep = check new (testCrudServiceUrl);
    
    // Test data setup
    string testId = "test-entity-read";
    string metadataKey = "readTest";
    string metadataValue = "read-test-value";
    
    // Create a test entity first
    record {| string key; pbAny:Any value; |}[] metadataArray = [];
    pbAny:Any packedValue = check pbAny:pack(metadataValue);
    metadataArray.push({key: metadataKey, value: packedValue});
    
    Entity createEntityRequest = {
        id: testId,
        kind: {
            major: "test",
            minor: "read-entity"
        },
        created: "2023-01-01",
        terminated: "",
        name: {
            startTime: "2023-01-01",
            endTime: "",
            value: check pbAny:pack("test-read-entity")
        },
        metadata: metadataArray,
        attributes: [],
        relationships: []
    };
    
    // Create entity
    Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
    io:println("Test entity created with ID: " + createEntityResponse.id);
    io:println("Created entity metadata: ", createEntityResponse.metadata);
    
    // Read the entity
    ReadEntityRequest readEntityRequest = {
        id: testId,
        entity: {
            id: testId,
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["metadata"]
    };
    io:println("ReadEntityRequest: ", readEntityRequest);
    Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
    io:println("Entity retrieved, verifying data...");
    io:println("Retrieved entity: ", readEntityResponse);
    io:println("Retrieved entity metadata: ", readEntityResponse.metadata);
    
    // Verify entity fields
    test:assertEquals(readEntityResponse.id, testId, "Entity ID mismatch");
    
    // Verify metadata
    boolean foundMetadata = false;
    foreach var item in readEntityResponse.metadata {
        if item.key == metadataKey {
            string|error unwrapped = unwrapAny(item.value);
            if unwrapped is string {
                test:assertEquals(unwrapped.trim(), metadataValue, 
                    string `Metadata value mismatch: expected ${metadataValue}, got ${unwrapped}`);
                foundMetadata = true;
            }
        }
    }
    
    test:assertTrue(foundMetadata, "Expected metadata key not found");
    
    // Test reading non-existent entity
    string nonExistentId = "non-existent-entity-" + testId;
    ReadEntityRequest nonExistentRequest = {
        id: nonExistentId,
        entity: {
            id: nonExistentId,
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["metadata"]
    };
    Entity nonExistentEntity = check ep->ReadEntity(nonExistentRequest);
    io:println("Non-existent entity: " + nonExistentEntity.id);
    io:println("Non-existent entity metadata: ", nonExistentEntity.metadata);
    
    // Validate that metadata for non-existent entity is empty
    test:assertEquals(nonExistentEntity.metadata.length(), 0, "Non-existent entity should have empty metadata");
    
    // Assert that we get an error for non-existent entity
    // For non-existence entities, we send a response with an empty data
    // But once the Result API is integrated this can be tested. 
    // FIXME: https://github.com/zaeema-n/LDFArchitecture/issues/23
    // test:assertTrue(nonExistentResponse is error, "Expected error for non-existent entity ID");
    
    // Clean up
    Empty _ = check ep->DeleteEntity({id: testId});
    Empty _ = check ep->DeleteEntity({id: testId});
    io:println("Test entity deleted");
    
    return;
}

@test:Config {}
function testCreateMinimalGraphEntity() returns error? {
    // Initialize the client
    CrudServiceClient ep = check new (testCrudServiceUrl);
    
    // Test data setup - minimal entity with just required fields
    string testId = "test-minimal-entity";
    
    // Create entity request with only required fields - no metadata, attributes, or relationships
    Entity createEntityRequest = {
        id: testId,
        kind: {
            major: "test",
            minor: "minimal"
        },
        created: "2023-01-01",
        terminated: "",
        name: {
            startTime: "2023-01-01",
            endTime: "",
            value: check pbAny:pack("minimal-test-entity")
        },
        metadata: [],
        attributes: [],
        relationships: []
    };

    // Create entity
    Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
    io:println("Minimal entity created with ID: " + createEntityResponse.id);
    
    // Verify entity was created correctly
    ReadEntityRequest readEntityRequest = {
        id: testId,
        entity: {
            id: testId,
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["metadata", "attributes", "relationships"]
    };
    Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
    
    // Basic entity verification
    test:assertEquals(readEntityResponse.id, testId, "Entity ID doesn't match");
    test:assertEquals(readEntityResponse.kind.major, "test", "Entity kind.major doesn't match");
    test:assertEquals(readEntityResponse.kind.minor, "minimal", "Entity kind.minor doesn't match");
    
    // Verify empty collections
    test:assertEquals(readEntityResponse.metadata.length(), 0, "Metadata should be empty");
    test:assertEquals(readEntityResponse.attributes.length(), 0, "Attributes should be empty");
    test:assertEquals(readEntityResponse.relationships.length(), 0, "Relationships should be empty");
    
    // Clean up
    Empty _ = check ep->DeleteEntity({id: testId});
    Empty _ = check ep->DeleteEntity({id: testId});
    io:println("Test minimal entity deleted");
    
    return;
}

@test:Config {}
function testCreateMinimalGraphEntityViaRest() returns error? {
    // Initialize an HTTP client for the REST API
    http:Client restClient = check new (testUpdateServiceUrl);
    
    // Test data setup - minimal JSON entity
    string testId = "test-minimal-json-entity";
    
    // Minimal JSON payload with required fields matching the Entity structure
    json minimalEntityJson = {
        "id": testId,
        "kind": {
            "major": "test",
            "minor": "minimal-json"
        },
        "created": "2023-01-01",
        "terminated": "",
        "name": {
            "startTime": "2023-01-01",
            "endTime": "",
            "value": "minimal-json-test-entity"
        },
        "metadata": [],
        "attributes": [],
        "relationships": []
    };

    // Create entity via REST API
    http:Response|error response = restClient->post("/entities", minimalEntityJson);
    
    // Verify HTTP request was successful
    if response is error {
        test:assertFail("Failed to create entity via REST API: " + response.message());
    }
    
    http:Response httpResponse = <http:Response>response;
    test:assertEquals(httpResponse.statusCode, 201, "Expected 201 OK status code");
    
    // Parse response JSON
    json responseJson = check httpResponse.getJsonPayload();
    test:assertEquals(check responseJson.id, testId, "Entity ID in response doesn't match");
    
    // Initialize the gRPC client to verify entity was properly created
    CrudServiceClient ep = check new (testCrudServiceUrl);
    
    // Verify entity data
    ReadEntityRequest readEntityRequest = {
        id: testId,
        entity: {
            id: testId,
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["metadata","attributes", "relationships"]
    };
    Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
    
    // Basic entity verification
    test:assertEquals(readEntityResponse.id, testId, "Entity ID doesn't match");
    test:assertEquals(readEntityResponse.kind.major, "test", "Entity kind.major doesn't match");
    test:assertEquals(readEntityResponse.kind.minor, "minimal-json", "Entity kind.minor doesn't match");
    
    // Verify empty collections
    test:assertEquals(readEntityResponse.metadata.length(), 0, "Metadata should be empty");
    test:assertEquals(readEntityResponse.attributes.length(), 0, "Attributes should be empty");
    test:assertEquals(readEntityResponse.relationships.length(), 0, "Relationships should be empty");
    
    // Clean up
    Empty _ = check ep->DeleteEntity({id: testId});
    Empty _ = check ep->DeleteEntity({id: testId});
    io:println("Test minimal JSON entity deleted");
    
    return;
}

@test:Config {
    groups: ["entity", "relationship"]
}
function testEntityWithRelationship() returns error? {
    // Test IDs for entities
    string sourceEntityId = "test-entity-with-relationship-source";
    string targetEntityId = "test-entity-with-relationship-target";
    
    // Initialize REST client
    http:Client restClient = check new (testUpdateServiceUrl);
    
    // Create source entity
    json sourceEntityJson = {
        "id": sourceEntityId,
        "kind": {
            "major": "test",
            "minor": "relationship-source"
        },
        "created": "2023-01-01",
        "terminated": "",
        "name": {
            "startTime": "2023-01-01",
            "endTime": "",
            "value": "source-entity"
        },
        "metadata": [],
        "attributes": [],
        "relationships": []
    };
    
    // Create target entity
    json targetEntityJson = {
        "id": targetEntityId,
        "kind": {
            "major": "test",
            "minor": "relationship-target"
        },
        "created": "2023-01-01", 
        "terminated": "",
        "name": {
            "startTime": "2023-01-01",
            "endTime": "",
            "value": "target-entity"
        },
        "metadata": [],
        "attributes": [],
        "relationships": []
    };
    
    // Create both entities via REST API
    http:Response|error sourceResponse = restClient->post("/entities", sourceEntityJson);
    http:Response|error targetResponse = restClient->post("/entities", targetEntityJson);
    
    // Verify HTTP requests were successful
    if sourceResponse is error {
        test:assertFail("Failed to create source entity: " + sourceResponse.message());
    }
    if targetResponse is error {
        test:assertFail("Failed to create target entity: " + targetResponse.message());
    }
    
    http:Response sourceHttpResponse = <http:Response>sourceResponse;
    http:Response targetHttpResponse = <http:Response>targetResponse;
    test:assertEquals(sourceHttpResponse.statusCode, 201, "Expected 201 status code for source entity");
    test:assertEquals(targetHttpResponse.statusCode, 201, "Expected 201 status code for target entity");
    
    // Create relationship between entities - include full entity structure
    string relationshipId = "rel-" + sourceEntityId + "-" + targetEntityId;
    json relationshipJson = {
        "id": sourceEntityId,
        "kind": {
        },
        "created": "",
        "terminated": "",
        "name": {
        },
        "metadata": [],
        "attributes": [],
        "relationships": {
            relationshipId: {
                "relatedEntityId": targetEntityId,
                "startTime": "2023-01-01",
                "endTime": "",
                "id": relationshipId,
                "name": "CONNECTS_TO"
            }
        }
    };
    
    // Update source entity with relationship
    http:Response|error updateResponse = restClient->put("/entities/" + sourceEntityId, relationshipJson);
    
    // Verify update was successful
    if updateResponse is error {
        test:assertFail("Failed to update entity with relationship: " + updateResponse.message());
    }
    
    http:Response updateHttpResponse = <http:Response>updateResponse;
    test:assertEquals(updateHttpResponse.statusCode, 200, "Expected 200 status code for relationship update");
    
    // Initialize the gRPC client to verify relationship was properly created
    CrudServiceClient ep = check new (testCrudServiceUrl);
    
    // Read source entity to verify relationship
    ReadEntityRequest readEntityRequest = {
        id: sourceEntityId,
        entity: {
            id: sourceEntityId,
            kind: {},
            created: "",
            terminated: "",
            name: {
                startTime: "",
                endTime: "",
                value: check pbAny:pack("")
            },
            metadata: [],
            attributes: [],
            relationships: []
        },
        output: ["relationships"]
    };
    Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);
    
    // Verify relationship data
    test:assertEquals(readEntityResponse.relationships.length(), 1, "Entity should have one relationship");
    
    // Find the relationship by iterating through the array
    Relationship? targetRelationship = ();
    foreach var rel in readEntityResponse.relationships {
        if rel.key == relationshipId {
            targetRelationship = rel.value;
            break;
        }
    }
    
    io:println("Target relationship: " + targetRelationship.toJsonString());
    test:assertFalse(targetRelationship is (), "Relationship with key 'CONNECTS_TO' not found");
    Relationship relationship = <Relationship>targetRelationship;
    test:assertEquals(relationship.relatedEntityId, targetEntityId, "Related entity ID doesn't match");
    test:assertEquals(relationship.name, "CONNECTS_TO", "Relationship name doesn't match");
    test:assertEquals(relationship.startTime, "2023-01-01T00:00:00Z", "Relationship start time doesn't match");
    test:assertEquals(relationship.id, relationshipId, "Relationship ID doesn't match");
    
    // Clean up
    Empty _ = check ep->DeleteEntity({id: sourceEntityId});
    Empty _ = check ep->DeleteEntity({id: sourceEntityId});
    Empty _ = check ep->DeleteEntity({id: targetEntityId});
    io:println("Test entities with relationship deleted");
    
    return;
}

@test:Config {
    groups: ["entity", "attributes"]
}
function testEntityWithTabularAttributes() returns error? {
    // Test ID for entity
    string testId = "test-entity-tabular";
    
    // Initialize the gRPC client to verify entity
    CrudServiceClient ep = check new (testCrudServiceUrl);
    
    // Create entity with tabular data in attributes
    json salaryGraph = {
        "nodes": [
            {"id": "salary_2024", "type": "salary_record", "properties": {"year": "2024", "amount": "100000", "bonus": "10000", "department": "Engineering"}},
            {"id": "salary_2023", "type": "salary_record", "properties": {"year": "2023", "amount": "90000", "bonus": "8000", "department": "Engineering"}},
            {"id": "salary_2022", "type": "salary_record", "properties": {"year": "2022", "amount": "80000", "bonus": "5000", "department": "Engineering"}},
            {"id": "dept_eng", "type": "department", "properties": {"name": "Engineering", "location": "HQ"}}
        ],
        "edges": [
            {"source": "salary_2024", "target": "dept_eng", "type": "belongs_to", "properties": {"effective_date": "2024-01-01"}},
            {"source": "salary_2023", "target": "dept_eng", "type": "belongs_to", "properties": {"effective_date": "2023-01-01"}},
            {"source": "salary_2022", "target": "dept_eng", "type": "belongs_to", "properties": {"effective_date": "2022-01-01"}},
            {"source": "salary_2024", "target": "salary_2023", "type": "promotion", "properties": {"date": "2024-01-01", "increase": "10000"}},
            {"source": "salary_2023", "target": "salary_2022", "type": "promotion", "properties": {"date": "2023-01-01", "increase": "10000"}}
        ]
    };

    json projectGraph = {
        "nodes": [
            {"id": "proj_redesign", "type": "project", "properties": {"id": "P001", "name": "System Redesign", "status": "active"}},
            {"id": "proj_migration", "type": "project", "properties": {"id": "P002", "name": "API Migration", "status": "completed"}},
            {"id": "proj_audit", "type": "project", "properties": {"id": "P003", "name": "Security Audit", "status": "completed"}},
            {"id": "role_lead", "type": "role", "properties": {"title": "Lead Developer", "level": "senior"}},
            {"id": "role_dev", "type": "role", "properties": {"title": "Developer", "level": "mid"}}
        ],
        "edges": [
            {"source": "proj_redesign", "target": "role_lead", "type": "has_role", "properties": {"start_date": "2024-01-01", "end_date": ""}},
            {"source": "proj_migration", "target": "role_dev", "type": "has_role", "properties": {"start_date": "2023-06-01", "end_date": "2023-12-31"}},
            {"source": "proj_audit", "target": "role_dev", "type": "has_role", "properties": {"start_date": "2023-01-01", "end_date": "2023-05-31"}},
            {"source": "proj_redesign", "target": "proj_migration", "type": "follows", "properties": {"transition_date": "2023-12-31"}},
            {"source": "proj_migration", "target": "proj_audit", "type": "follows", "properties": {"transition_date": "2023-05-31"}}
        ]
    };

    Entity createEntityRequest = {
        id: testId,
        kind: {
            major: "test",
            minor: "graph"
        },
        created: "2024-01-01T00:00:00Z",
        terminated: "",
        name: {
            startTime: "2024-01-01T00:00:00Z",
            endTime: "",
            value: check pbAny:pack("test-entity")
        },
        metadata: [
            {
                key: "test_metadata",
                value: check pbAny:pack("test_value")
            }
        ],
        attributes: [
            {
                key: "employee_salary_history",
                value: {
                    values: [
                        {
                            startTime: "2024-01-01T00:00:00Z",
                            endTime: "",
                            value: check pbAny:pack(salaryGraph.toJsonString())
                        }
                    ]
                }
            },
            {
                key: "project_assignments",
                value: {
                    values: [
                        {
                            startTime: "2024-01-01T00:00:00Z",
                            endTime: "",
                            value: check pbAny:pack(projectGraph.toJsonString())
                        }
                    ]
                }
            }
        ],
        relationships: []
    };

    // Create entity via gRPC
    Entity createEntityResponse = check ep->CreateEntity(createEntityRequest);
    io:println("Entity created with ID: " + createEntityResponse.id);
    
    // Read entity to verify attributes
    EntityId readEntityRequest = {id: testId};
    Entity readEntityResponse = check ep->ReadEntity(readEntityRequest);

    io:println("Entity: " + readEntityResponse.toJsonString());
    
    // Verify basic entity data
    test:assertEquals(readEntityResponse.id, testId, "Entity ID doesn't match");
    test:assertEquals(readEntityResponse.kind.major, "test", "Entity kind.major doesn't match");
    test:assertEquals(readEntityResponse.kind.minor, "graph", "Entity kind.minor doesn't match");
    
    // Verify metadata
    test:assertTrue(readEntityResponse.metadata.length() > 0, "Entity should have metadata");
    boolean foundMetadata = false;
    foreach var item in readEntityResponse.metadata {
        if item.key == "test_metadata" {
            string|error unwrapped = unwrapAny(item.value);
            if unwrapped is string {
                test:assertEquals(unwrapped.trim(), "test_value", "Metadata value doesn't match");
                foundMetadata = true;
            }
        }
    }
    test:assertTrue(foundMetadata, "Expected metadata key not found");
    
    // Verify attributes
    test:assertTrue(readEntityResponse.attributes.length() > 0, "Entity should have attributes");
    boolean foundSalaryHistory = false;
    boolean foundProjectAssignments = false;
    
    foreach var attr in readEntityResponse.attributes {
        if attr.key == "employee_salary_history" {
            foundSalaryHistory = true;
            test:assertTrue(attr.value.values.length() > 0, "Salary history should have values");
            
            // Parse the graph data
            json|error graphData = check unwrapAny(attr.value.values[0].value);
            if graphData is json {
                json parsedData = <json>graphData;
                test:assertTrue(parsedData is map<json>, "Graph data should be a map");
                
                // Verify nodes
                map<json> dataMap = <map<json>>parsedData;
                json nodesJson = dataMap["nodes"];
                json[] nodes = <json[]>nodesJson;
                test:assertEquals(nodes.length(), 4, "Should have 4 nodes");
                
                // Verify first node data
                map<json> firstNode = <map<json>>nodes[0];
                test:assertEquals(firstNode["id"], "salary_2024", "First node id should be 'salary_2024'");
                test:assertEquals(firstNode["type"], "salary_record", "First node type should be 'salary_record'");
                
                // Verify node properties
                map<json> properties = <map<json>>firstNode["properties"];
                test:assertEquals(properties["year"], "2024", "Year should be 2024");
                test:assertEquals(properties["amount"], "100000", "Amount should be 100000");
            }
        }
        if attr.key == "project_assignments" {
            foundProjectAssignments = true;
            test:assertTrue(attr.value.values.length() > 0, "Project assignments should have values");
            
            // Parse the graph data
            json|error graphData = check unwrapAny(attr.value.values[0].value);
            if graphData is json {
                json parsedData = <json>graphData;
                test:assertTrue(parsedData is map<json>, "Graph data should be a map");
                
                // Verify nodes
                map<json> dataMap = <map<json>>parsedData;
                json nodesJson = dataMap["nodes"];
                json[] nodes = <json[]>nodesJson;
                test:assertEquals(nodes.length(), 5, "Should have 5 nodes");
                
                // Verify first node data
                map<json> firstNode = <map<json>>nodes[0];
                test:assertEquals(firstNode["id"], "proj_redesign", "First node id should be 'proj_redesign'");
                test:assertEquals(firstNode["type"], "project", "First node type should be 'project'");
                
                // Verify node properties
                map<json> properties = <map<json>>firstNode["properties"];
                test:assertEquals(properties["id"], "P001", "Project ID should be P001");
                test:assertEquals(properties["name"], "System Redesign", "Project name should be System Redesign");
            }
        }
    }
    
    test:assertTrue(foundSalaryHistory, "Salary history attribute not found");
    test:assertTrue(foundProjectAssignments, "Project assignments attribute not found");
    
    // Verify relationships
    test:assertTrue(readEntityResponse.relationships.length() > 0, "Entity should have relationships");
    boolean foundRelationship = false;
    foreach var rel in readEntityResponse.relationships {
        if rel.key == "reports_to" {
            foundRelationship = true;
            test:assertEquals(rel.value.relatedEntityId, "manager123", "Related entity ID doesn't match");
            test:assertEquals(rel.value.name, "reports_to", "Relationship name doesn't match");
            test:assertEquals(rel.value.startTime, "2024-01-01T00:00:00Z", "Relationship start time doesn't match");
            test:assertEquals(rel.value.id, "rel123", "Relationship ID doesn't match");
        }
    }
    test:assertTrue(foundRelationship, "Expected relationship not found");
    
    // Clean up
    Empty _ = check ep->DeleteEntity(readEntityRequest);
    Empty _ = check ep->DeleteEntity({id: testId});
    io:println("Test entity with graph attributes deleted");
    
    return;
}


