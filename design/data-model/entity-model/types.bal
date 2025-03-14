import ballerina/time;
import ballerina/uuid;

type Kind record {
    string major;
    string minor;
};

type TimeBasedValue record {
    time:Utc startTime;
    time:Utc? endTime = ();
    anydata value;
};

type Relationship record {
    string relatedEntityId;
    time:Utc startTime;
    time:Utc? endTime = ();
};

type Entity record {
    readonly string id; // = uuid:createType4AsString();
    readonly Kind kind;
    readonly time:Utc created;
    time:Utc? terminated = ();
    TimeBasedValue name;
    map<anydata> metadata;
    map<TimeBasedValue[]> attributes;
    map<Relationship> relationships;
};

function foo() {
  Entity entity1 = {
    id: uuid:createType4AsString(),
    kind: {major: "Person", minor: "Citizen"},
    created: time:utcNow(),
    name: {startTime: time:utcNow(), value: "John Doe"},
    metadata: {},
    attributes: {},
    relationships: {}
};

entity1.attributes["age"] = [{startTime: time:utcNow(), value: 25}];
entity1.attributes["address"] = [{startTime: time:utcNow(), value: "Colombo"}];
entity1.relationships["father"] = {relatedEntityId: uuid:createType4AsString(), startTime: time:utcNow()};
entity1.relationships["mother"] = {relatedEntityId: uuid:createType4AsString(), startTime: time:utcNow()};
entity1.relationships["spouse"] = {relatedEntityId: uuid:createType4AsString(), startTime: time:utcNow()};
entity1.relationships["child"] = {relatedEntityId: uuid:createType4AsString(), startTime: time:utcNow()};
entity1.relationships["friend"] = {relatedEntityId: uuid:createType4AsString(), startTime: time:utcNow()};


}


