// import ballerina/http;

// # A service representing a network-accessible API
// # bound to port `9090`.
// service / on new http:Listener(9090) {

//     # A resource for generating greetings
//     # + name - name as a string or nil
//     # + return - string name with hello message or error
//     resource function get greeting(string? name) returns string|error {
//         // Send a response back to the caller.
//         if name is () {
//             return error("name should not be empty!");
//         }
//         return string `Hello, ${name}`;
//     }
// }

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

    resource function get getEntityAttribute(string entityId, string attributeName, time:Utc? ts) returns record { time:Utc start; time:Utc end; anydata value }|record { time:Utc start; time:Utc end; anydata value }[]|() {
        // Implement your logic to fetch the attribute value based on the given entityId, attributeName, and optional timestamp
        log:printInfo("Fetching attribute for entity ID: " + entityId + ", attribute: " + attributeName + ", timestamp: " + ts.toString());

        // Example data, replace with actual logic to fetch attribute values
        if (ts is time:Utc) {
            // Return a single value for the given timestamp
            record { time:Utc start; time:Utc end; anydata value } value = {
                start: time:currentTime(),
                end: time:currentTime(),
                value: "exampleValue"
            };
            return value;
        } else {
            // Return all values for all time ranges
            record { time:Utc start; time:Utc end; anydata value }[] values = [
                {
                    start: time:currentTime(),
                    end: time:currentTime(),
                    value: "exampleValue1"
                },
                {
                    start: time:currentTime(),
                    end: time:currentTime(),
                    value: "exampleValue2"
                }
            ];
            return values;
        }
    }
}
