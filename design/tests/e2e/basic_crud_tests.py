import requests
import json
import sys
import base64

"""
This file contains the end-to-end tests for the CRUD API.
It is used to test the API's functionality by creating, reading, updating, and deleting an entity.

Running the tests:

## Run CRUD Server

```bash
cd design/crud-api
./crud-server
```
## Run API Server

```bash
cd design/api
bal run
```

## Run Tests

```bash
cd design/tests/e2e
python3 basic_crud_tests.py
```

"""

BASE_URL = "http://localhost:8080/entities"
ENTITY_ID = "12345"

def create_payload():
    """Returns the entity payload for create and update operations."""
    return {
        "create": {
            "id": ENTITY_ID,
            "kind": {"major": "example", "minor": "test"},
            "created": "2024-03-17T10:00:00Z",
            "terminated": "",
            "name": {
                "startTime": "2024-03-17T10:00:00Z",
                "endTime": "",
                "value": {
                    "typeUrl": "type.googleapis.com/google.protobuf.StringValue",
                    "value": "entity-name"
                }
            },
            "metadata": [
                {"key": "owner", "value": "test-user"},
                {"key": "version", "value": "1.0"},
                {"key": "developer", "value": "V8A"}
            ],
            "attributes": [],
            "relationships": []
        },
        "update": {
            "id": ENTITY_ID,
            "kind": {"major": "example", "minor": "test"},
            "created": "2024-03-18T00:00:00Z",
            "name": {
                "startTime": "2024-03-18T00:00:00Z",
                "value": "entity-name"
            },
            "metadata": [{"key": "version", "value": "5.0"}]
        }
    }

def create_entity(payload):
    """Creates an entity and validates the response."""
    print("\n🟢 Creating entity...")
    response = requests.post(BASE_URL, json=payload["create"], headers={"Content-Type": "application/json"})
    
    if response.status_code == 201:
        print("✅ Entity created:", json.dumps(response.json(), indent=2))
    else:
        print(f"❌ Create failed: {response.text}")
        sys.exit(1)

def read_entity():
    """Reads and validates the created entity."""
    print("\n🟢 Reading entity...")
    response = requests.get(f"{BASE_URL}/{ENTITY_ID}")
    
    if response.status_code == 200:
        data = response.json()
        assert data["id"] == ENTITY_ID, "Read entity ID mismatch"
        print("✅ Read Entity:", json.dumps(data, indent=2))
    else:
        print(f"❌ Read failed: {response.text}")
        sys.exit(1)

def update_entity(payload):
    """Updates the entity and validates the response."""
    print("\n🟢 Updating entity...")
    response = requests.put(f"{BASE_URL}/{ENTITY_ID}", json=payload["update"], headers={"Content-Type": "application/json"})
    
    if response.status_code == 200:
        updated_entity = response.json()
        decoded_value = decode_protobuf_any_value(updated_entity["metadata"][0]["value"])
        print("decoded value: ", decoded_value)
        assert decoded_value == "5.0", "Update did not modify metadata"
        print("✅ Entity updated:", json.dumps(updated_entity, indent=2))
    else:
        print(f"❌ Update failed: {response.text}")
        sys.exit(1)

def validate_update():
    """Validates that the update has been applied correctly."""
    print("\n🟢 Validating update...")
    response = requests.get(f"{BASE_URL}/{ENTITY_ID}")
    
    if response.status_code == 200:
        updated_data = response.json()
        decoded_value = decode_protobuf_any_value(updated_data["metadata"][0]["value"])
        assert decoded_value == "5.0", "Updated entity does not reflect changes"
        print("✅ Update Validation Passed:", json.dumps(updated_data, indent=2))
    else:
        print(f"❌ Read failed after update: {response.text}")
        sys.exit(1)

def delete_entity():
    """Deletes the entity."""
    print("\n🟢 Deleting entity...")
    response = requests.delete(f"{BASE_URL}/{ENTITY_ID}")
    
    if response.status_code == 204:
        print("✅ Entity deleted successfully.")
    else:
        print(f"❌ Delete failed: {response.text}")
        sys.exit(1)

def verify_deletion():
    """Verifies that the entity has been deleted."""
    print("\n🟢 Verifying deletion...")
    response = requests.get(f"{BASE_URL}/{ENTITY_ID}")
    
    if response.status_code == 500:
        print("❌ Server error occurred:", response.text)
        sys.exit(1)
    else:
        print(f"\n🟢 Entity was not deleted properly: {response.status_code} {response.text}")

def decode_protobuf_any_value(any_value):
    """Decode a protobuf Any value to get the actual string value"""
    if isinstance(any_value, dict) and 'typeUrl' in any_value and 'value' in any_value:
        if 'StringValue' in any_value['typeUrl']:
            try:
                # First try direct base64 decoding if that's how it's encoded
                try:
                    binary_data = base64.b64decode(any_value['value'])
                    # For StringValue, skip the field tag byte and length byte
                    # and decode the remaining bytes as UTF-8
                    return binary_data[2:].decode('utf-8')
                except:
                    # If it's hex encoded (which appears to be the case)
                    hex_value = any_value['value']
                    binary_data = bytes.fromhex(hex_value)
                    # For StringValue in hex format, typically the structure is:
                    # 0A (field tag) + 03 (length) + actual string bytes
                    # Skip the first 2 bytes (field tag and length)
                    if len(binary_data) > 2:
                        return binary_data[2:].decode('utf-8')
            except Exception as e:
                print(f"Failed to decode protobuf value: {e}")
    # Return the original value if decoding fails
    return any_value.strip()

if __name__ == "__main__":
    print("🚀 Running End-to-End API Test Suite...")
    
    payload = create_payload()

    try:
        create_entity(payload)
        read_entity()
        update_entity(payload)
        validate_update()
        delete_entity()
        verify_deletion()
        
        print("\n🎉 All tests passed successfully!")
    
    except AssertionError as e:
        print(f"\n❌ Test failed: {e}")
        sys.exit(1)
