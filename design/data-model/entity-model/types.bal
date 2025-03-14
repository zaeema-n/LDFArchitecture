import ballerina/time;

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
