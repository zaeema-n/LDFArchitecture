import ballerina/io;
import ballerina/time;
import ./types;

public function main() {
    io:println("Hello, World!");

    types:PersonEntity person = {
        majorType: "Person",
        minorType: "Citizen",
        dateOfCreation: time:currentTime(),
        name: "John Doe"
    };

    types:BusinessEntity business = {
        majorType: "Business",
        minorType: "Private",
        dateOfCreation: time:currentTime(),
        name: "Tech Corp"
    };

    types:GovernmentEntity government = {
        majorType: "Government",
        minorType: "TBD",
        dateOfCreation: time:currentTime(),
        name: "Government Entity"
    };

    types:LandParcelEntity landParcel = {
        majorType: "LandParcel",
        minorType: "Education",
        dateOfCreation: time:currentTime(),
        name: "School Land"
    };

    types:AdministrativeEntity administrative = {
        majorType: "Administrative",
        minorType: "Country",
        dateOfCreation: time:currentTime(),
        name: "Country Admin"
    };

    io:println("Person Entity: ", person);
    io:println("Business Entity: ", business);
    io:println("Government Entity: ", government);
    io:println("Land Parcel Entity: ", landParcel);
    io:println("Administrative Entity: ", administrative);
}
