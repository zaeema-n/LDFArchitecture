import ballerina/graphql;
import ballerina/log;

enum MajorType {
    Person,
    Business
};

// type Entity record {
//     string id;
//     MajorType kind;
//     int birthDate?;
//     int deathDate?;
// };

type EntityMetadata record {
    string kind;
    int birthDate;
    int deathDate;
};

// Record for getfindEntities function as graphql function does not accept maps (check this)
type AttributeInput record {
    string name;
    anydata value;
};

// Define the new record type
type AttributeValue record {
    string startTime;
    string end;
    string value;
};

service /query on new graphql:Listener(9090) {

    resource function get getEntryEntityIds(MajorType kind) returns string[] {
        // Implement your logic to fetch the entity IDs of the given kind
        log:printInfo("Fetching entity IDs for kind: " + kind.toString());
        // Example data, replace with actual logic to fetch entity IDs
        string[] entityIds = ["id1", "id2", "id3"];
        return entityIds;
    }

    resource function get getfindEntities(MajorType kind, int? birthDate, AttributeInput[]? attributes, int? deathDate) returns string[] {
        // Implement your logic to fetch the entities based on the given criteria
        log:printInfo("Finding entities with criteria - kind: " + kind.toString() + ", birthDate: " + birthDate.toString() + ", deathDate: " + deathDate.toString() + ", attributes: " + attributes.toString());

        if attributes != () {
            foreach AttributeInput attribute in attributes {
                log:printInfo("Attribute name: " + attribute.name + ", value: " + attribute.value.toString());
            }
        }
        // foreach AttributeInput attribute in attributes {
        //     log:printInfo("Attribute name: " + attribute.name + ", value: " + attribute.value.toString());
        // }

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

    resource function get getEntityAttribute(string entityId, string attributeName, string ts) returns AttributeValue[] {
        log:printInfo("Fetching attribute for entity ID: " + entityId + ", attribute: " + attributeName + ", timestamp: " + ts);

        AttributeValue[] values;

        if (ts != "") { // If timestamp is provided, return a single value
            values = [
                {
                    startTime: "2025-03-13T12:00:00Z",
                    end: "2025-03-13T12:30:00Z",
                    value: "exampleValue"
                }
            ];
        } else { // Return multiple values for all time ranges
            values = [
                {
                    startTime: "2025-03-13T12:00:00Z",
                    end: "2025-03-13T12:30:00Z",
                    value: "exampleValue1"
                },
                {
                    startTime: "2025-03-13T13:00:00Z",
                    end: "2025-03-13T13:30:00Z",
                    value: "exampleValue2"
                }
            ];
        }

        return values;
    }

    // New endpoint for getRelatedEntityIds
    resource function get getRelatedEntityIds(string entityId, string relationship, string ts) returns string[] {
        // Implement your logic to fetch the related entity IDs based on the given entityId, relationship, and timestamp
        log:printInfo("Fetching related entity IDs for entity ID: " + entityId + ", relationship: " + relationship + ", timestamp: " + ts);

        // Example data, replace with actual logic to fetch related entity IDs
        string[] relatedEntityIds = ["relatedId1", "relatedId2", "relatedId3"];
        return relatedEntityIds;
    }
}
