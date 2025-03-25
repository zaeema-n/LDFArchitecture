import ballerina/time;

public type Kind record {
    string major;
    string minor;
};

public type TimeBasedValue record {
    time:Utc startTime;
    time:Utc? endTime = ();
    anydata value;
};

public type Relationship record {
    string relatedEntityId;
    time:Utc startTime;
    time:Utc? endTime = ();
};

public type Entity record {
    readonly string id; // = uuid:createType4AsString();
    readonly Kind kind;
    readonly time:Utc created;
    time:Utc? terminated = ();
    TimeBasedValue name;
    map<anydata> metadata;
    map<TimeBasedValue[]> attributes;
    map<Relationship> relationships;
};
