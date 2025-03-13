import ballerina/graphql;
import ballerina/log;
import ballerina/time;

enum MajorType {
    Person,
    Business
};

type Entity record {
    string id;
    MajorType kind;
    int birthDate?;
    int deathDate?;
};

type EntityMetadata record {
    string kind;
    int birthDate;
    int deathDate;
};

// Define the new record type
type AttributeValue record {
    time:Utc startTime;
    time:Utc end;
    anydata value;
};

service /query on new graphql:Listener(9090) {

    resource function get getEntryEntityIds(MajorType kind) returns string[] {
        // Implement your logic to fetch the entity IDs of the given kind
        log:printInfo("Fetching entity IDs for kind: " + kind.toString());
        // Example data, replace with actual logic to fetch entity IDs
        string[] entityIds = ["id1", "id2", "id3"];
        return entityIds;
    }

    resource function get findEntities(MajorType kind, int? birthDate, int? deathDate, map<string,anydata>? attributes) returns string[] {
        // Implement your logic to fetch the entities based on the given criteria
        log:printInfo("Finding entities with criteria - kind: " + kind.toString() + ", birthDate: " + birthDate.toString() + ", deathDate: " + deathDate.toString());

        // Example data, replace with actual logic to fetch entity IDs
        string[] entityIds = ["id4", "id5", "id6"];
        return entityIds;
    }

    resource function get getEntityMetadata(string entityId) returns EntityMetadata {
        // Implement your logic to fetch the entity metadata based on the given entityId
        log:printInfo("Fetching metadata for entity ID: " + entityId);
        // Example data, replace with actual logic to fetch entity metadata
        EntityMetadata metadata = {
            kind: "Person",
            birthDate: 1990,
            deathDate: 2020
        };
        return metadata;
    }

    resource function get getEntityAttribute(string entityId, string attributeName, time:Utc? ts) returns AttributeValue|AttributeValue[]|() {
        // Implement your logic to fetch the attribute value based on the given entityId, attributeName, and optional timestamp
        log:printInfo("Fetching attribute for entity ID: " + entityId + ", attribute: " + attributeName + ", timestamp: " + ts.toString());

        // Example data, replace with actual logic to fetch attribute values
        if (ts is time:Utc) {
            // Return a single value for the given timestamp
            AttributeValue value = {
                startTime: time:utcNow(),
                end: time:utcNow(),
                value: "exampleValue"
            };
            return value;
        } else {
            // Return all values for all time ranges
            AttributeValue[] values = [
                {
                    startTime: time:utcNow(),
                    end: time:utcNow(),
                    value: "exampleValue1"
                },
                {
                    startTime: time:utcNow(),
                    end: time:utcNow(),
                    value: "exampleValue2"
                }
            ];
            return values;
        }
    }

    // New endpoint for getRelatedEntityIds
    resource function get getRelatedEntityIds(string entityId, string relationship, time:Utc ts) returns string[] {
        // Implement your logic to fetch the related entity IDs based on the given entityId, relationship, and timestamp
        log:printInfo("Fetching related entity IDs for entity ID: " + entityId + ", relationship: " + relationship + ", timestamp: " + ts.toString());

        // Example data, replace with actual logic to fetch related entity IDs
        string[] relatedEntityIds = ["relatedId1", "relatedId2", "relatedId3"];
        return relatedEntityIds;
    }
}
