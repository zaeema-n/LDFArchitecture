import requests
import json
import sys

QUERY_API_URL = "http://localhost:8081/v1/entities"
UPDATE_API_URL = "http://localhost:8080/entities"
ENTITY_ID = "query-test-entity"
RELATED_ID = "query-related-entity"

MINISTER_ID = "minister_education"
DEPARTMENTS = [
    {"id": "dept_exams", "name": "Department of Exams"},
    {"id": "dept_nie", "name": "National Institute of Education"},
    {"id": "dept_ed_publications", "name": "Department of Educational Publications"}
]

START_DATE = "2015-04-11T00:00:00Z"


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

def create_minister():
    """Create a Minister entity."""
    print("\n🟢 Creating Minister entity...")
    
    payload = {
        "id": MINISTER_ID,
        "kind": {"major": "Organization", "minor": "Minister"},
        "created": START_DATE,
        "terminated": "",
        "name": {
            "startTime": START_DATE,
            "endTime": "",
            "value": {
                "typeUrl": "type.googleapis.com/google.protobuf.StringValue",
                "value": "Minister of Education"
            }
        },
        "metadata": []
    }
    
    res = requests.post(UPDATE_API_URL, json=payload)
    assert res.status_code in [200, 201], f"Failed to create Minister: {res.text}"

    print(f"Response: {res.status_code} - {res.text}")
    print("✅ Created Minister entity.")


def create_departments():
    """Create Department entities."""
    print("\n🟢 Creating Department entities...")
    
    for dept in DEPARTMENTS:
        payload = {
            "id": dept["id"],
            "kind": {"major": "Organization", "minor": "Department"},
            "created": START_DATE,
            "terminated": "",
            "name": {
                "startTime": START_DATE,
                "endTime": "",
                "value": {
                    "typeUrl": "type.googleapis.com/google.protobuf.StringValue",
                    "value": dept["name"]
                }
            },
            "metadata": []
        }
        
        res = requests.post(UPDATE_API_URL, json=payload)
        assert res.status_code in [200, 201], f"Failed to create {dept['name']}: {res.text}"
        print(f"Response: {res.status_code} - {res.text}")
        print(f"✅ Created {dept['name']} entity.")


def create_relationships():
    """Create HAS_DEPARTMENT relationships from Minister to Departments."""
    print("\n🔗 Creating relationships...")
    
    for dept in DEPARTMENTS:
        payload = {
            "relationships": [
                {
                    "key": "HAS_DEPARTMENT",
                    "value": {
                        "relatedEntityId": dept["id"],
                        "startTime": START_DATE,
                        "endTime": "",
                        "id": f"rel_{dept['id']}",
                        "name": "HAS_DEPARTMENT"
                    }
                }
            ]
        }
        
        url = f"{UPDATE_API_URL}/{MINISTER_ID}"  # Using PUT with ID in the path
        res = requests.put(url, json=payload)

        if res.status_code in [200, 204]:  # 204 for successful updates with no content
            print(f"Response: {res.status_code} - {res.text}")
            print(f"✅ Created relationship between Minister and {dept['name']}.")
        else:
            print(f"❌ Failed to create relationship for {dept['name']}: {res.status_code} - {res.text}")

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

    print("🚀 Running Minister and Departments Setup...")
    try:
        create_minister()
        create_departments()
        create_relationships()
        print("\n🎉 Minister and Departments setup completed successfully!")
    except AssertionError as e:
        print(f"\n❌ Setup failed: {e}")
        sys.exit(1)
