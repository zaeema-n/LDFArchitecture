import ballerina/grpc;
import ballerina/protobuf;

public const string TYPES_DESC = "0A0B74797065732E70726F746F12046372756422320A044B696E6412140A056D616A6F7218012001280952056D616A6F7212140A056D696E6F7218022001280952056D696E6F7222C1040A06456E74697479120E0A02696418012001280952026964121E0A046B696E6418022001280B320A2E637275642E4B696E6452046B696E6412180A0763726561746564180320012809520763726561746564121E0A0A7465726D696E61746564180420012809520A7465726D696E6174656412280A046E616D6518052001280B32142E637275642E54696D65426173656456616C756552046E616D6512360A086D6574616461746118062003280B321A2E637275642E456E746974792E4D65746164617461456E74727952086D65746164617461123C0A0A6174747269627574657318072003280B321C2E637275642E456E746974792E41747472696275746573456E747279520A6174747269627574657312450A0D72656C6174696F6E736869707318082003280B321F2E637275642E456E746974792E52656C6174696F6E7368697073456E747279520D72656C6174696F6E73686970731A3B0A0D4D65746164617461456E74727912100A036B657918012001280952036B657912140A0576616C7565180220012809520576616C75653A0238011A530A0F41747472696275746573456E74727912100A036B657918012001280952036B6579122A0A0576616C756518022001280B32142E637275642E54696D65426173656456616C7565520576616C75653A0238011A540A1252656C6174696F6E7368697073456E74727912100A036B657918012001280952036B657912280A0576616C756518022001280B32122E637275642E52656C6174696F6E73686970520576616C75653A023801225E0A0E54696D65426173656456616C7565121C0A09737461727454696D651801200128095209737461727454696D6512180A07656E6454696D651802200128095207656E6454696D6512140A0576616C7565180320012809520576616C756522700A0C52656C6174696F6E7368697012280A0F72656C61746564456E746974794964180120012809520F72656C61746564456E746974794964121C0A09737461727454696D651802200128095209737461727454696D6512180A07656E6454696D651803200128095207656E6454696D6532BB010A0B4372756453657276696365122A0A0C437265617465456E74697479120C2E637275642E456E746974791A0C2E637275642E456E7469747912280A0A52656164456E74697479120C2E637275642E456E746974791A0C2E637275642E456E74697479122A0A0C557064617465456E74697479120C2E637275642E456E746974791A0C2E637275642E456E74697479122A0A0C44656C657465456E74697479120C2E637275642E456E746974791A0C2E637275642E456E74697479421C5A1A6C6B2F64617461666F756E646174696F6E2F637275642D617069620670726F746F33";

public isolated client class CrudServiceClient {
    *grpc:AbstractClientEndpoint;

    private final grpc:Client grpcClient;

    public isolated function init(string url, *grpc:ClientConfiguration config) returns grpc:Error? {
        self.grpcClient = check new (url, config);
        check self.grpcClient.initStub(self, TYPES_DESC);
    }

    isolated remote function CreateEntity(Entity|ContextEntity req) returns Entity|grpc:Error {
        map<string|string[]> headers = {};
        Entity message;
        if req is ContextEntity {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("crud.CrudService/CreateEntity", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <Entity>result;
    }

    isolated remote function CreateEntityContext(Entity|ContextEntity req) returns ContextEntity|grpc:Error {
        map<string|string[]> headers = {};
        Entity message;
        if req is ContextEntity {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("crud.CrudService/CreateEntity", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <Entity>result, headers: respHeaders};
    }

    isolated remote function ReadEntity(Entity|ContextEntity req) returns Entity|grpc:Error {
        map<string|string[]> headers = {};
        Entity message;
        if req is ContextEntity {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("crud.CrudService/ReadEntity", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <Entity>result;
    }

    isolated remote function ReadEntityContext(Entity|ContextEntity req) returns ContextEntity|grpc:Error {
        map<string|string[]> headers = {};
        Entity message;
        if req is ContextEntity {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("crud.CrudService/ReadEntity", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <Entity>result, headers: respHeaders};
    }

    isolated remote function UpdateEntity(Entity|ContextEntity req) returns Entity|grpc:Error {
        map<string|string[]> headers = {};
        Entity message;
        if req is ContextEntity {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("crud.CrudService/UpdateEntity", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <Entity>result;
    }

    isolated remote function UpdateEntityContext(Entity|ContextEntity req) returns ContextEntity|grpc:Error {
        map<string|string[]> headers = {};
        Entity message;
        if req is ContextEntity {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("crud.CrudService/UpdateEntity", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <Entity>result, headers: respHeaders};
    }

    isolated remote function DeleteEntity(Entity|ContextEntity req) returns Entity|grpc:Error {
        map<string|string[]> headers = {};
        Entity message;
        if req is ContextEntity {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("crud.CrudService/DeleteEntity", message, headers);
        [anydata, map<string|string[]>] [result, _] = payload;
        return <Entity>result;
    }

    isolated remote function DeleteEntityContext(Entity|ContextEntity req) returns ContextEntity|grpc:Error {
        map<string|string[]> headers = {};
        Entity message;
        if req is ContextEntity {
            message = req.content;
            headers = req.headers;
        } else {
            message = req;
        }
        var payload = check self.grpcClient->executeSimpleRPC("crud.CrudService/DeleteEntity", message, headers);
        [anydata, map<string|string[]>] [result, respHeaders] = payload;
        return {content: <Entity>result, headers: respHeaders};
    }
}

public type ContextEntity record {|
    Entity content;
    map<string|string[]> headers;
|};

@protobuf:Descriptor {value: TYPES_DESC}
public type Entity record {|
    string id = "";
    Kind kind = {};
    string created = "";
    string terminated = "";
    TimeBasedValue name = {};
    record {|string key; string value;|}[] metadata = [];
    record {|string key; TimeBasedValue value;|}[] attributes = [];
    record {|string key; Relationship value;|}[] relationships = [];
|};

@protobuf:Descriptor {value: TYPES_DESC}
public type TimeBasedValue record {|
    string startTime = "";
    string endTime = "";
    string value = "";
|};

@protobuf:Descriptor {value: TYPES_DESC}
public type Kind record {|
    string major = "";
    string minor = "";
|};

@protobuf:Descriptor {value: TYPES_DESC}
public type Relationship record {|
    string relatedEntityId = "";
    string startTime = "";
    string endTime = "";
|};
