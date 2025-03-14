import ballerina/time;

type PersonMinorType "Citizen" | "Resident" | "Visitor" | "Foreigner";
type BusinessMinorType "Private" | "Public" | "NonProfit";
type GovernmentMinorType "TBD";
type LandParcelMinorType "Plain" | "Police" | "Education" | "Health" | "Election";
type AdministrativeMinorType "World" | "Country" | "Ocean" | "Continent";

type MajorType "Person" | "Business" | "Government" | "LandParcel" | "Administrative";

type MinorType PersonMinorType | BusinessMinorType | GovernmentMinorType | LandParcelMinorType | AdministrativeMinorType;

type Entity record {
    MajorType majorType;
    MinorType minorType;
    map<anydata> parameters?;
    time:Utc dateOfCreation;
    time:Utc? dateOfTermination;
    string name;
};

type PersonEntity record {|
    *Entity;
    MajorType majorType = "Person";
    PersonMinorType minorType;
|};

type BusinessEntity record {|
    *Entity;
    MajorType majorType = "Business";
    BusinessMinorType minorType;
|};

type GovernmentEntity record {|
    *Entity;
    MajorType majorType = "Government";
    GovernmentMinorType minorType;
|};

type LandParcelEntity record {|
    *Entity;
    MajorType majorType = "LandParcel";
    LandParcelMinorType minorType;
|};

type AdministrativeEntity record {|
    *Entity;
    MajorType majorType = "Administrative";
    AdministrativeMinorType minorType;
|};

