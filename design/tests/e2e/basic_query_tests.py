import requests
import json
import sys

QUERY_API_URL = "http://localhost:8081/v1/entities"
UPDATE_API_URL = "http://localhost:8080/entities"
ENTITY_ID = "query-test-entity"
RELATED_ID = "query-related-entity"


"""
The current tests only contain metadata validation.
"""

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

    payload_child = {
        "id": RELATED_ID,
        "kind": {"major": "test", "minor": "child"},
        "created": "2024-01-01T00:00:00Z",
        "terminated": "",
        "name": {
            "startTime": "2024-01-01T00:00:00Z",
            "endTime": "",
            "value": {
                "typeUrl": "type.googleapis.com/google.protobuf.StringValue",
                "value": "Query Test Entity Child"
            }
        },
        "metadata": [
            {"key": "source", "value": "unit-test-1"},
            {"key": "env", "value": "test-1"}
        ],
        "attributes": [
            {
                "key": "humidity",
                "value": {
                    "values": [
                        {
                            "startTime": "2024-01-01T00:00:00Z",
                            "endTime": "2024-01-02T00:00:00Z",
                            "value": {
                                "typeUrl": "type.googleapis.com/google.protobuf.StringValue",
                                "value": "10.5"
                            }
                        }
                    ]
                }
            }
        ],
        "relationships": [
        ]
    }

    payload_source = {
        "id": ENTITY_ID,
        "kind": {"major": "test", "minor": "parent"},
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
                "key": "rel-001",
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

    res = requests.post(UPDATE_API_URL, json=payload_child)
    assert res.status_code == 201 or res.status_code == 200, f"Failed to create entity: {res.text}"
    print("✅ Created entity for query tests.")

    res = requests.post(UPDATE_API_URL, json=payload_source)
    assert res.status_code == 201 or res.status_code == 200, f"Failed to create entity: {res.text}"
    print("✅ Created entity for query tests.")

def test_attribute_lookup():
    """Test retrieving attributes via the query API."""
    print("\n🔍 Testing attribute retrieval...")
    url = f"{QUERY_API_URL}/{ENTITY_ID}/attributes/temperature"
    res = requests.get(url)
    assert res.status_code == 404, f"Failed to get attribute: {res.text}"
    
    # Add response body validation
    body = res.json()
    assert isinstance(body, dict), "Response should be a dictionary"
    assert "error" in body, "Error message should be present in 404 response"
    print("✅ Attribute response:", json.dumps(res.json(), indent=2))

def test_metadata_lookup():
    """Test retrieving metadata."""
    print("\n🔍 Testing metadata retrieval...")
    url = f"{QUERY_API_URL}/{ENTITY_ID}/metadata"
    res = requests.get(url)
    assert res.status_code == 200, f"Failed to get metadata: {res.text}"
    
    body = res.json()
    print("✅ Raw metadata response:", json.dumps(body, indent=2))
    
    # Enhanced metadata validation
    assert isinstance(body, dict), "Metadata response should be a dictionary"
    assert len(body) == 2, f"Expected 2 metadata entries, got {len(body)}"
    assert "source" in body, "Source metadata key missing"
    assert "env" in body, "Env metadata key missing"
    
    source_value = decode_protobuf_any_value(body["source"])
    env_value = decode_protobuf_any_value(body["env"])
    
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
    
    body = res.json()
    # Add relationship response validation
    assert isinstance(body, list), "Relationship response should be a list"
    assert len(body) > 0, "Expected at least one relationship"
    
    relationship = body[0]
    assert "relatedEntityId" in relationship, "Relationship should have relatedEntityId"
    assert relationship["relatedEntityId"] == RELATED_ID, "Related entity ID mismatch"
    assert relationship["name"] == "linked", "Relationship name mismatch"
    assert relationship["id"] == "rel-001", "Relationship ID mismatch"
    print("✅ Relationship response:", json.dumps(res.json(), indent=2))

def test_entity_search():
    """Test search by entity ID."""
    print("\n🔍 Testing entity search...")
    url = f"{QUERY_API_URL}/search"
    payload = {
        "id": ENTITY_ID,
        "created": "",
        "terminated": ""
    }
    res = requests.post(url, json=payload)
    assert res.status_code == 200, f"Search failed: {res.text}"
    
    body = res.json()
    # Add search response validation
    ## FIXME: Make sure to implement the entities/search and update this test case
    assert isinstance(body, dict), "Search response should be a dictionary"
    assert "body" in body, "Search response should have a 'body' field"
    assert isinstance(body["body"], list), "Search response body should be a list"
    assert len(body["body"]) == 0, "Expected an empty list in search response"


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
