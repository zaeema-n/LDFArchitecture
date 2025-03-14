import ballerina/graphql;
import ballerina/log;

enum MajorType {
    Person,
    Business,
    Government,
    LandParcel,
    Administrative
};

type EntityMetadata record {
    string kind;
    int birthDate;
    int deathDate;
};

type AttributeInput record {
    string name;
    anydata value;
};

type AttributeValue record {
    string startTime;
    string end;
    string value;
};

service /update on new graphql:Listener(9091) {

    remote function updateEntityMetadata(string entityId, EntityMetadata metadata) returns boolean {
        // Implement your logic to update the entity metadata based on the given entityId
        log:printInfo("Updating metadata for entity ID: " + entityId + " with metadata: " + metadata.toString());
        // Example logic, replace with actual update logic
        boolean success = true;
        return success;
    }

    remote function updateEntityAttribute(string entityId, string attributeName, AttributeValue value) returns boolean {
        log:printInfo("Updating attribute for entity ID: " + entityId + ", attribute: " + attributeName + " with value: " + value.toString());
        // Example logic, replace with actual update logic
        boolean success = true;
        return success;
    }

    remote function addRelatedEntity(string entityId, string relatedEntityId, string relationship) returns boolean {
        log:printInfo("Adding related entity ID: " + relatedEntityId + " to entity ID: " + entityId + " with relationship: " + relationship);
        // Example logic, replace with actual update logic
        boolean success = true;
        return success;
    }

    remote function removeRelatedEntity(string entityId, string relatedEntityId, string relationship) returns boolean {
        log:printInfo("Removing related entity ID: " + relatedEntityId + " from entity ID: " + entityId + " with relationship: " + relationship);
        // Example logic, replace with actual update logic
        boolean success = true;
        return success;
    }
}