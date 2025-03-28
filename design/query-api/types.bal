// AUTO-GENERATED FILE.
// This file is auto-generated by the Ballerina OpenAPI tool.

import ballerina/http;

public type entities_search_body record {
    string kind?;
    string created?;
    string terminated?;
    # Attribute filters. Example: {"height":180, "eyeColor":"blue"}
    record {} attributes?;
};

public type EntitiesEntityIdMetadataResponse record {
};

public type InlineResponse2002ArrayOk record {|
    *http:Ok;
    inline_response_200_2[] body;
|};

public type InlineResponse200Ok record {|
    *http:Ok;
    inline_response_200 body;
|};

public type inline_response_200_1 record {string 'start?; string? end?; string value?;}|record {string 'start?; string? end?; string value?;}[]|string?;

public type inline_response_200 record {
    string[] body?;
};

public type inline_response_200_2 record {
    string relatedEntityId?;
    string startTime?;
    string endTime?;
    string id?;
    string name?;
};

public type entityId_relations_body record {
    string relatedEntityId;
    string startTime;
    string endTime;
    string id;
    string name;
};
