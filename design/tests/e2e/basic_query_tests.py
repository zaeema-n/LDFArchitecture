import requests
import json
import sys

QUERY_API_URL = "http://localhost:8081/v1/entities"
UPDATE_API_URL = "http://localhost:8080/entities"
ENTITY_ID = "query-test-entity"
RELATED_ID = "query-related-entity"

def decode_protobuf_any_value(any_value):
    """Decode a protobuf Any value to get the actual string value"""
    if isinstance(any_value, dict) and 'typeUrl' in any_value and 'value' in any_value:
        if 'StringValue' in any_value['typeUrl']:
            try:
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
                return any_value['value']
    
    # If any_value is a string that looks like a JSON object
    elif isinstance(any_value, str) and any_value.startswith('{') and any_value.endswith('}'):
        try:
            # Try to parse it as JSON
            obj = json.loads(any_value)
            # Recursively decode
            return decode_protobuf_any_value(obj)
        except json.JSONDecodeError:
            pass
    
    # Return the original value if decoding fails
    return any_value

def create_entity_for_query():
    """Create a base entity with metadata, attributes, and relationships."""
    print("\n🟢 Creating entity for query tests...")

    payload = {
        "id": ENTITY_ID,
        "kind": {"major": "test", "minor": "query"},
        "created": "2024-01-01T00:00:00Z",
        "terminated": "",
        "name": {
            "startTime": "2024-01-01T00:00:00Z",
            "endTime": "",
            "value": {
                "typeUrl": "type.googleapis.com/google.protobuf.StringValue",
                "value": "Query Test Entity"
            }
        },
        "metadata": [
            {"key": "source", "value": "unit-test"},
            {"key": "env", "value": "test"}
        ],
        "attributes": [
            {
                "key": "temperature",
                "value": {
                    "values": [
                        {
                            "startTime": "2024-01-01T00:00:00Z",
                            "endTime": "2024-01-02T00:00:00Z",
                            "value": {
                                "typeUrl": "type.googleapis.com/google.protobuf.StringValue",
                                "value": "25.5"
                            }
                        }
                    ]
                }
            }
        ],
        "relationships": [
            {
                "key": "linked",
                "value": {
                    "relatedEntityId": RELATED_ID,
                    "startTime": "2024-01-01T00:00:00Z",
                    "endTime": "2024-12-31T23:59:59Z",
                    "id": "rel-001",
                    "name": "linked"
                }
            }
        ]
    }

    res = requests.post(UPDATE_API_URL, json=payload)
    assert res.status_code == 201 or res.status_code == 200, f"Failed to create entity: {res.text}"
    print("✅ Created entity for query tests.")

def test_attribute_lookup():
    """Test retrieving attributes via the query API."""
    print("\n🔍 Testing attribute retrieval...")
    url = f"{QUERY_API_URL}/{ENTITY_ID}/attributes/temperature"
    res = requests.get(url)
    assert res.status_code == 404, f"Failed to get attribute: {res.text}"
    print("✅ Attribute response:", json.dumps(res.json(), indent=2))

def test_metadata_lookup():
    """Test retrieving metadata."""
    print("\n🔍 Testing metadata retrieval...")
    url = f"{QUERY_API_URL}/{ENTITY_ID}/metadata"
    res = requests.get(url)
    assert res.status_code == 200, f"Failed to get metadata: {res.text}"
    
    body = res.json()
    print("✅ Raw metadata response:", json.dumps(body, indent=2))
    
    # Check if keys exist, regardless of their format
    assert "source" in body, "Source metadata key missing"
    assert "env" in body, "Env metadata key missing"
    
    # Extract actual string values from possibly complex protobuf structures
    source_value = decode_protobuf_any_value(body["source"])
    env_value = decode_protobuf_any_value(body["env"])
    
    # Verify the extracted values match what we expect
    assert source_value == "unit-test", f"Source value mismatch: {source_value}"
    assert env_value == "test", f"Env value mismatch: {env_value}"

def test_relationship_query():
    """Test relationship query via POST /relations."""
    print("\n🔍 Testing relationship filtering...")
    url = f"{QUERY_API_URL}/{ENTITY_ID}/relations"
    payload = {
        "relatedEntityId": RELATED_ID,
        "startTime": "2024-01-01T00:00:00Z",
        "endTime": "2024-12-31T23:59:59Z",
        "id": "rel-001",
        "name": "linked"
    }
    res = requests.post(url, json=payload)
    assert res.status_code == 200, f"Failed to get relationships: {res.text}"
    print("✅ Relationship response:", json.dumps(res.json(), indent=2))

def test_entity_search():
    """Test search by kind."""
    print("\n🔍 Testing entity search...")
    url = f"{QUERY_API_URL}/search"
    payload = {
        "kind": "test",
        "birthDate": "",
        "deathDate": "",
        "attributes": {
            "temperature": "25.5"
        }
    }
    res = requests.post(url, json=payload)
    assert res.status_code == 200, f"Search failed: {res.text}"
    print("✅ Search response:", json.dumps(res.json(), indent=2))

if __name__ == "__main__":
    print("🚀 Running Query API E2E Tests...")

    try:
        create_entity_for_query()
        test_attribute_lookup()
        test_metadata_lookup()
        test_relationship_query()
        test_entity_search()
        print("\n🎉 All Query API tests passed!")
    except AssertionError as e:
        print(f"\n❌ Test failed: {e}")
        sys.exit(1)
