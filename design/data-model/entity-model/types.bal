enum MajorType {
    Person,
    Business
};

enum MinorType {
    Citizen,
    Resident,
    Visitor,
    Foreigner
};

enum BusinessType {
    Private,
    Public,
    NonProfit
};

type PersonRecord record {
    MajorType kind; // Major type: Person
    MinorType subkind;
    string name;            // Additional parameter
    int bornDate;           // Additional parameter
    int deathDate; 
};

type CitizenRecord record {
    PersonRecord person;
    string citizenId;
};

type ResidentRecord record {
    PersonRecord person;
    string residentPermit;
};

type VisitorRecord record {
    PersonRecord person;
    string visaType;
};

type ForeignerRecord record {
    PersonRecord person;
    string passportNumber;
};

type PersonType PersonRecord|Citizen|Resident|Visitor|Foreigner;

