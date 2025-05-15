# How It Works: End-to-End Data Flow

This document describes the complete workflow of how data flows through the system, from initial JSON input to final storage in the databases.

## 1. Data Entry Point (Update API)

The system receives data through a REST API built with Ballerina. The API accepts JSON payloads for entity creation and updates.

### Example JSON Input
```json
{
    "id": "entity123",
    "kind": {
        "major": "Person",
        "minor": "Employee"
    },
    "name": {
        "startTime": "2024-01-01T00:00:00Z",
        "endTime": "",
        "value": "John Doe"
    },
    "metadata": {
        "department": "Engineering",
        "role": "Software Engineer"
    },
    "attributes": {
        "salary": {
            "values": [
                {
                    "startTime": "2024-01-01T00:00:00Z",
                    "endTime": "",
                    "value": "100000"
                }
            ]
        }
    },
    "relationships": {
        "reports_to": {
            "relatedEntityId": "manager123",
            "startTime": "2024-01-01T00:00:00Z",
            "endTime": "",
            "id": "rel123",
            "name": "reports_to"
        }
    }
}
```

## 2. Data Transformation (Update API → CRUD Service)

### 2.1 JSON to Protobuf Conversion
The Update API converts the JSON payload into a protobuf Entity message. This conversion happens in the `convertJsonToEntity` function:

```ballerina
function convertJsonToEntity(json jsonPayload) returns Entity|error {
    // Convert metadata to protobuf Any type
    record {| string key; pbAny:Any value; |}[] metadataArray = [];
    foreach var [key, val] in metadataMap.entries() {
        pbAny:Any packedValue = check pbAny:pack(val.toString());
        metadataArray.push({key: key, value: packedValue});
    }

    // Convert attributes to TimeBasedValueList
    record {| string key; TimeBasedValueList value; |}[] attributesArray = [];
    foreach var [key, val] in attributesMap.entries() {
        TimeBasedValue[] timeBasedValues = [];
        // Process each time-based value
        TimeBasedValue tbv = {
            startTime: val.startTime,
            endTime: val.endTime,
            value: check pbAny:pack(val.value.toString())
        };
        timeBasedValues.push(tbv);
    }

    // Convert relationships
    record {| string key; Relationship value; |}[] relationshipsArray = [];
    foreach var [key, val] in relationshipsMap.entries() {
        Relationship relationship = {
            relatedEntityId: val.relatedEntityId,
            startTime: val.startTime,
            endTime: val.endTime,
            id: val.id,
            name: val.name
        };
        relationshipsArray.push({key: key, value: relationship});
    }

    // Create final Entity
    return Entity {
        id: jsonPayload.id,
        kind: {major: kindJson.major, minor: kindJson.minor},
        name: {startTime: startTimeValue, endTime: endTimeValue, value: namePackedValue},
        metadata: metadataArray,
        attributes: attributesArray,
        relationships: relationshipsArray
    };
}
```

### 2.2 gRPC Communication
The converted protobuf message is sent to the CRUD service via gRPC. The communication happens on port 50051.

## 3. CRUD Service Processing

The CRUD service receives the protobuf message and processes it through multiple steps:

### 3.1 Create Entity Flow
```go
func (s *Server) CreateEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
    // 1. Save metadata in MongoDB
    err := s.mongoRepo.HandleMetadata(ctx, req.Id, req)
    
    // 2. Save entity in Neo4j
    success, err := s.neo4jRepo.HandleGraphEntityCreation(ctx, req)
    
    // 3. Handle relationships in Neo4j
    err = s.neo4jRepo.HandleGraphRelationshipsCreate(ctx, req)
    
    return req, nil
}
```

### 3.2 Data Storage

#### MongoDB Storage (Metadata)
The metadata is stored in MongoDB as a document:
```json
{
    "_id": "entity123",
    "metadata": {
        "department": "Engineering",
        "role": "Software Engineer"
    }
}
```

#### Neo4j Storage (Entity and Relationships)
The entity and its relationships are stored in Neo4j as nodes and relationships:

```cypher
// Entity Node
CREATE (e:Entity {
    id: 'entity123',
    kind_major: 'Person',
    kind_minor: 'Employee',
    name: 'John Doe',
    created: '2024-01-01T00:00:00Z',
    terminated: null
})

// Relationship
CREATE (e)-[r:REPORTS_TO {
    id: 'rel123',
    startTime: '2024-01-01T00:00:00Z',
    endTime: null
}]->(m:Entity {id: 'manager123'})
```

## 4. Data Retrieval Flow

### 4.1 Read Entity Flow
```go
func (s *Server) ReadEntity(ctx context.Context, req *pb.Entity) (*pb.Entity, error) {
    // 1. Get metadata from MongoDB
    metadata, _ := s.mongoRepo.GetMetadata(ctx, req.Id)

    // 2. Get entity info from Neo4j
    kind, name, created, terminated, _ := s.neo4jRepo.GetGraphEntity(ctx, req.Id)

    // 3. Get relationships from Neo4j
    relationships, _ := s.neo4jRepo.GetGraphRelationships(ctx, req.Id)

    // 4. Return combined entity
    return &pb.Entity{
        Id:            req.Id,
        Kind:          kind,
        Name:          name,
        Created:       created,
        Terminated:    terminated,
        Metadata:      metadata,
        Attributes:    make(map[string]*pb.TimeBasedValueList),
        Relationships: relationships,
    }, nil
}
```

### 4.2 Data Transformation (CRUD Service → Update API)
The retrieved data is converted back to JSON in the Update API before being sent to the client.

## 5. Error Handling

The system implements error handling at multiple levels:

1. **Update API Level**
   - JSON validation
   - Protobuf conversion errors
   - gRPC communication errors

2. **CRUD Service Level**
   - Database connection errors
   - Data validation errors
   - Transaction errors

3. **Repository Level**
   - Database-specific errors
   - Query execution errors
   - Data consistency errors

## 6. Data Consistency

The system maintains data consistency through:

1. **Atomic Operations**
   - MongoDB transactions for metadata
   - Neo4j transactions for entity and relationships

2. **Error Recovery(TODO)**
   - Rollback mechanisms
   - Error logging and monitoring
   - Retry mechanisms for failed operations

## 7. Performance Considerations

1. **Connection Pooling(TODO)**
   - MongoDB connection pool
   - Neo4j connection pool
   - gRPC connection management

2. **Caching(TODO)**
   - Metadata caching
   - Entity relationship caching

3. **Query Optimization(TODO)**
   - Indexed queries
   - Efficient relationship traversal
   - Batch operations
