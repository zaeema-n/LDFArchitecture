public table<Kind> kindTable = table [];

function populateKindTable() {
    kindTable.add({major: "person", minor: "citizen"});
    kindTable.add({major: "person", minor: "resident"});
    kindTable.add({major: "person", minor: "visitor"});
    kindTable.add({major: "person", minor: "foreigner"});
    
    kindTable.add({major: "business", minor: "private"});
    kindTable.add({major: "business", minor: "public"});
    kindTable.add({major: "business", minor: "non-profit"});
    
    kindTable.add({major: "government", minor: nil});
    
    kindTable.add({major: "land-parcel", minor: "plain"});
    
    kindTable.add({major: "administrative", minor: "world"});
    kindTable.add({major: "administrative", minor: "country"});
    kindTable.add({major: "administrative", minor: "ocean"});
    kindTable.add({major: "administrative", minor: "continent"});
    
    kindTable.add({major: "asset", minor: "infrastructure"});
    kindTable.add({major: "asset", minor: "equipment"});
    kindTable.add({major: "asset", minor: "vehicle"});
    kindTable.add({major: "asset", minor: "digital"});
}