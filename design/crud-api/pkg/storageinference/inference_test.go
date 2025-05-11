package storageinference

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// JSONToAny converts a JSON string to a protobuf Any value
func JSONToAny(jsonStr string) (*anypb.Any, error) {
	// First parse the JSON into a generic interface
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Convert the interface to a protobuf Struct
	structValue, err := structpb.NewStruct(data.(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("failed to create struct: %v", err)
	}

	// Convert the Struct to Any
	anyValue, err := anypb.New(structValue)
	if err != nil {
		return nil, fmt.Errorf("failed to create Any: %v", err)
	}

	return anyValue, nil
}

// TestStorageTypeInferenceFromEntity tests storage type inference using entity structure
func TestStorageTypeInferenceFromEntity(t *testing.T) {
	// Sample test cases for different storage types in attributes
	testCases := map[StorageType]string{
		TabularData: `{
			"attributes": {
				"columns": ["id", "name", "age"],
				"rows": [
					[1, "John", 30],
					[2, "Jane", 25]
				]
			}
		}`,
		ScalarData: `{
			"attributes": {
				"value": 42
			}
		}`,
		ListData: `{
			"attributes": {
				"items": [1, 2, 3, 4, 5]
			}
		}`,
		MapData: `{
			"attributes": {
				"key1": "value1",
				"key2": "value2",
				"key3": 42
			}
		}`,
	}

	inferrer := &StorageInferrer{}
	for expectedType, jsonStr := range testCases {
		t.Run(string(expectedType), func(t *testing.T) {
			// Convert JSON to Any
			anyValue, err := JSONToAny(jsonStr)
			assert.NoError(t, err)

			// Infer the storage type
			detectedType, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, expectedType, detectedType)
		})
	}
}

// TestComplexTabularEntity tests a more complex tabular structure in attributes
func TestComplexTabularEntity(t *testing.T) {
	complexTableJSON := `{
		"attributes": {
			"columns": ["id", "name", "age", "address"],
			"rows": [
				[1, "John", 30, {"city": "New York", "zip": "10001"}],
				[2, "Jane", 25, {"city": "Boston", "zip": "02108"}]
			],
			"metadata": {
				"total_rows": 2,
				"last_updated": "2024-03-20"
			}
		}
	}`

	anyValue, err := JSONToAny(complexTableJSON)
	assert.NoError(t, err)

	inferrer := &StorageInferrer{}
	detectedType, err := inferrer.InferType(anyValue)
	assert.NoError(t, err)
	assert.Equal(t, TabularData, detectedType)
}

// TestNestedListEntity tests a nested list structure in attributes
func TestNestedListEntity(t *testing.T) {
	nestedListJSON := `{
		"attributes": {
			"items": [
				[1, 2, 3],
				[4, 5, 6],
				[7, 8, 9]
			]
		}
	}`

	anyValue, err := JSONToAny(nestedListJSON)
	assert.NoError(t, err)

	inferrer := &StorageInferrer{}
	detectedType, err := inferrer.InferType(anyValue)
	assert.NoError(t, err)
	assert.Equal(t, ListData, detectedType)
}

// TestComplexMapEntity tests a more complex map structure in attributes
func TestComplexMapEntity(t *testing.T) {
	complexMapJSON := `{
		"attributes": {
			"user": {
				"name": "John",
				"age": 30,
				"address": {
					"city": "New York",
					"zip": "10001"
				}
			},
			"settings": {
				"theme": "dark",
				"notifications": true
			}
		}
	}`

	anyValue, err := JSONToAny(complexMapJSON)
	assert.NoError(t, err)

	inferrer := &StorageInferrer{}
	detectedType, err := inferrer.InferType(anyValue)
	assert.NoError(t, err)
	assert.Equal(t, MapData, detectedType)
}

// TestMixedEntity tests a structure that could be interpreted as multiple types
func TestMixedEntity(t *testing.T) {
	mixedDataJSON := `{
		"attributes": {
			"items": [1, 2, 3],
			"metadata": {
				"count": 3,
				"type": "numbers"
			}
		}
	}`

	anyValue, err := JSONToAny(mixedDataJSON)
	assert.NoError(t, err)

	inferrer := &StorageInferrer{}
	detectedType, err := inferrer.InferType(anyValue)
	assert.NoError(t, err)
	// This should be detected as ListData because it has a slice field
	assert.Equal(t, ListData, detectedType)
}

// TestInvalidJSON tests handling of invalid JSON input
func TestInvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`
	_, err := JSONToAny(invalidJSON)
	assert.Error(t, err)
}

// TestDirectScalarEntity tests scalar data with direct key-value pairs in attributes
func TestDirectScalarEntity(t *testing.T) {
	testCases := map[string]string{
		"integer": `{
			"attributes": {
				"x": 37
			}
		}`,
		"float": `{
			"attributes": {
				"pi": 3.14159
			}
		}`,
		"string": `{
			"attributes": {
				"name": "test"
			}
		}`,
		"boolean": `{
			"attributes": {
				"active": true
			}
		}`,
	}

	inferrer := &StorageInferrer{}
	for testName, jsonStr := range testCases {
		t.Run(testName, func(t *testing.T) {
			// Convert JSON to Any
			anyValue, err := JSONToAny(jsonStr)
			assert.NoError(t, err)

			// Infer the storage type
			detectedType, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, ScalarData, detectedType, "Expected scalar data for test case: %s", testName)
		})
	}
}

// TestComplexListEntity tests various complex list structures in attributes
func TestComplexListEntity(t *testing.T) {
	testCases := map[string]string{
		"nested_lists": `{
			"attributes": {
				"items": [
					[1, 2, 3],
					[4, 5, 6],
					[7, 8, 9]
				]
			}
		}`,
		"mixed_types": `{
			"attributes": {
				"items": [
					42,
					"string value",
					true,
					3.14,
					{"nested": "object"}
				]
			}
		}`,
		"list_of_objects": `{
			"attributes": {
				"items": [
					{"id": 1, "name": "item1"},
					{"id": 2, "name": "item2"},
					{"id": 3, "name": "item3"}
				]
			}
		}`,
		"list_with_metadata": `{
			"attributes": {
				"items": [1, 2, 3, 4, 5],
				"metadata": {
					"count": 5,
					"type": "numbers"
				}
			}
		}`,
		"heterogeneous_nested": `{
			"attributes": {
				"items": [
					[1, "two", 3],
					{"x": 1, "y": 2},
					[true, false],
					"string item"
				]
			}
		}`,
	}

	inferrer := &StorageInferrer{}
	for testName, jsonStr := range testCases {
		t.Run(testName, func(t *testing.T) {
			// Convert JSON to Any
			anyValue, err := JSONToAny(jsonStr)
			assert.NoError(t, err)

			// Infer the storage type
			detectedType, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, ListData, detectedType, "Expected list data for test case: %s", testName)
		})
	}
}

// TestAdvancedMapEntity tests advanced map structures with complex nested data
func TestAdvancedMapEntity(t *testing.T) {
	testCases := map[string]string{
		"deeply_nested": `{
			"attributes": {
				"system": {
					"config": {
						"network": {
							"interfaces": {
								"eth0": {
									"ip": "192.168.1.1",
									"mask": "255.255.255.0"
								},
								"eth1": {
									"ip": "10.0.0.1",
									"mask": "255.0.0.0"
								}
							},
							"dns": ["8.8.8.8", "8.8.4.4"]
						}
					}
				}
			}
		}`,
		"mixed_arrays": `{
			"attributes": {
				"data": {
					"numbers": [1, 2, 3, 4, 5],
					"strings": ["a", "b", "c"],
					"mixed": [1, "two", 3.0, true],
					"nested": [
						{"id": 1, "value": "one"},
						{"id": 2, "value": "two"}
					]
				}
			}
		}`,
		"complex_metrics": `{
			"attributes": {
				"metrics": {
					"cpu": {
						"usage": 45.2,
						"cores": 4,
						"temperature": 65.5
					},
					"memory": {
						"total": 16384,
						"used": 8192,
						"free": 8192
					},
					"disk": {
						"total": 1024000,
						"used": 512000,
						"free": 512000
					}
				}
			}
		}`,
		"multi_level": `{
			"attributes": {
				"organization": {
					"name": "Acme Corp",
					"departments": {
						"engineering": {
							"head": "John Doe",
							"teams": {
								"frontend": {
									"size": 5,
									"tech": ["React", "TypeScript"]
								},
								"backend": {
									"size": 8,
									"tech": ["Go", "PostgreSQL"]
								}
							}
						},
						"marketing": {
							"head": "Jane Smith",
							"budget": 100000
						}
					}
				}
			}
		}`,
		"heterogeneous_data": `{
			"attributes": {
				"mixed": {
					"primitive": 42,
					"text": "hello",
					"boolean": true,
					"null": null,
					"array": [1, "two", 3.0],
					"object": {
						"nested": "value",
						"numbers": [1, 2, 3],
						"flags": {
							"active": true,
							"verified": false
						}
					}
				}
			}
		}`,
	}

	inferrer := &StorageInferrer{}
	for testName, jsonStr := range testCases {
		t.Run(testName, func(t *testing.T) {
			// Convert JSON to Any
			anyValue, err := JSONToAny(jsonStr)
			assert.NoError(t, err)

			// Infer the storage type
			detectedType, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, MapData, detectedType, "Expected map data for test case: %s", testName)
		})
	}
}

// TestConsistentTabularEntity tests tabular data with consistent column types
func TestConsistentTabularEntity(t *testing.T) {
	testCases := map[string]string{
		"numeric_table": `{
			"attributes": {
				"columns": ["id", "age", "score", "rating"],
				"rows": [
					[1, 25, 95.5, 4.5],
					[2, 30, 88.0, 4.0],
					[3, 35, 92.5, 4.8]
				]
			}
		}`,
		"string_table": `{
			"attributes": {
				"columns": ["id", "name", "email", "department"],
				"rows": [
					["001", "John Doe", "john@example.com", "Engineering"],
					["002", "Jane Smith", "jane@example.com", "Marketing"],
					["003", "Bob Wilson", "bob@example.com", "Sales"]
				]
			}
		}`,
		"boolean_table": `{
			"attributes": {
				"columns": ["id", "is_active", "has_access", "is_verified"],
				"rows": [
					[1, true, true, false],
					[2, false, true, true],
					[3, true, false, true]
				]
			}
		}`,
		"date_table": `{
			"attributes": {
				"columns": ["id", "created_at", "updated_at", "expires_at"],
				"rows": [
					[1, "2024-01-01", "2024-03-01", "2024-12-31"],
					[2, "2024-01-15", "2024-03-15", "2024-12-31"],
					[3, "2024-02-01", "2024-03-20", "2024-12-31"]
				]
			}
		}`,
		"timestamp_table": `{
			"attributes": {
				"columns": ["id", "start_time", "end_time", "last_modified"],
				"rows": [
					[1, "2024-03-20T10:00:00Z", "2024-03-20T11:00:00Z", "2024-03-20T09:00:00Z"],
					[2, "2024-03-20T14:00:00Z", "2024-03-20T15:00:00Z", "2024-03-20T13:00:00Z"],
					[3, "2024-03-21T09:00:00Z", "2024-03-21T10:00:00Z", "2024-03-21T08:00:00Z"]
				]
			}
		}`,
	}

	inferrer := &StorageInferrer{}
	for testName, jsonStr := range testCases {
		t.Run(testName, func(t *testing.T) {
			// Convert JSON to Any
			anyValue, err := JSONToAny(jsonStr)
			assert.NoError(t, err)

			// Infer the storage type
			detectedType, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, TabularData, detectedType, "Expected tabular data for test case: %s", testName)
		})
	}
}

// TestAdvancedGraphEntity tests complex graph data structures
func TestAdvancedGraphEntity(t *testing.T) {
	testCases := map[string]string{
		"social_network": `{
			"attributes": {
				"nodes": [
					{"id": "user1", "type": "user", "properties": {"name": "Alice", "age": 30, "location": "NY"}},
					{"id": "user2", "type": "user", "properties": {"name": "Bob", "age": 25, "location": "SF"}},
					{"id": "user3", "type": "user", "properties": {"name": "Charlie", "age": 35, "location": "LA"}},
					{"id": "post1", "type": "post", "properties": {"title": "Hello", "content": "World", "created": "2024-03-20"}},
					{"id": "post2", "type": "post", "properties": {"title": "Graph", "content": "DB", "created": "2024-03-21"}}
				],
				"edges": [
					{"source": "user1", "target": "user2", "type": "follows", "properties": {"since": "2024-01-01"}},
					{"source": "user2", "target": "user3", "type": "follows", "properties": {"since": "2024-02-01"}},
					{"source": "user1", "target": "post1", "type": "created", "properties": {"timestamp": "2024-03-20T10:00:00Z"}},
					{"source": "user2", "target": "post1", "type": "likes", "properties": {"timestamp": "2024-03-20T11:00:00Z"}},
					{"source": "user3", "target": "post2", "type": "created", "properties": {"timestamp": "2024-03-21T09:00:00Z"}}
				]
			}
		}`,
		"dependency_graph": `{
			"attributes": {
				"nodes": [
					{"id": "pkg1", "type": "package", "properties": {"name": "core", "version": "1.0.0"}},
					{"id": "pkg2", "type": "package", "properties": {"name": "utils", "version": "2.1.0"}},
					{"id": "pkg3", "type": "package", "properties": {"name": "network", "version": "1.5.0"}},
					{"id": "pkg4", "type": "package", "properties": {"name": "database", "version": "3.0.0"}}
				],
				"edges": [
					{"source": "pkg2", "target": "pkg1", "type": "depends_on", "properties": {"version": ">=1.0.0"}},
					{"source": "pkg3", "target": "pkg1", "type": "depends_on", "properties": {"version": ">=1.0.0"}},
					{"source": "pkg4", "target": "pkg2", "type": "depends_on", "properties": {"version": ">=2.0.0"}},
					{"source": "pkg4", "target": "pkg3", "type": "depends_on", "properties": {"version": ">=1.5.0"}}
				]
			}
		}`,
		"workflow_graph": `{
			"attributes": {
				"nodes": [
					{"id": "task1", "type": "task", "properties": {"name": "fetch_data", "status": "completed"}},
					{"id": "task2", "type": "task", "properties": {"name": "process_data", "status": "running"}},
					{"id": "task3", "type": "task", "properties": {"name": "validate", "status": "pending"}},
					{"id": "task4", "type": "task", "properties": {"name": "store_results", "status": "pending"}}
				],
				"edges": [
					{"source": "task1", "target": "task2", "type": "triggers", "properties": {"condition": "success"}},
					{"source": "task2", "target": "task3", "type": "triggers", "properties": {"condition": "success"}},
					{"source": "task3", "target": "task4", "type": "triggers", "properties": {"condition": "success"}},
					{"source": "task1", "target": "task4", "type": "triggers", "properties": {"condition": "failure"}}
				]
			}
		}`,
		"knowledge_graph": `{
			"attributes": {
				"nodes": [
					{"id": "concept1", "type": "concept", "properties": {"name": "Machine Learning", "category": "AI"}},
					{"id": "concept2", "type": "concept", "properties": {"name": "Neural Networks", "category": "AI"}},
					{"id": "concept3", "type": "concept", "properties": {"name": "Deep Learning", "category": "AI"}},
					{"id": "concept4", "type": "concept", "properties": {"name": "Supervised Learning", "category": "ML"}}
				],
				"edges": [
					{"source": "concept2", "target": "concept1", "type": "is_a", "properties": {"confidence": 0.95}},
					{"source": "concept3", "target": "concept2", "type": "is_a", "properties": {"confidence": 0.90}},
					{"source": "concept4", "target": "concept1", "type": "is_a", "properties": {"confidence": 0.85}},
					{"source": "concept3", "target": "concept4", "type": "uses", "properties": {"confidence": 0.80}}
				]
			}
		}`,
		"network_topology": `{
			"attributes": {
				"nodes": [
					{"id": "router1", "type": "router", "properties": {"ip": "10.0.0.1", "model": "Cisco"}},
					{"id": "switch1", "type": "switch", "properties": {"ip": "10.0.0.2", "ports": 24}},
					{"id": "server1", "type": "server", "properties": {"ip": "10.0.0.3", "os": "Linux"}},
					{"id": "server2", "type": "server", "properties": {"ip": "10.0.0.4", "os": "Windows"}}
				],
				"edges": [
					{"source": "router1", "target": "switch1", "type": "connected_to", "properties": {"bandwidth": "1Gbps"}},
					{"source": "switch1", "target": "server1", "type": "connected_to", "properties": {"bandwidth": "100Mbps"}},
					{"source": "switch1", "target": "server2", "type": "connected_to", "properties": {"bandwidth": "100Mbps"}},
					{"source": "server1", "target": "server2", "type": "communicates_with", "properties": {"protocol": "HTTP"}}
				]
			}
		}`,
	}

	inferrer := &StorageInferrer{}
	for testName, jsonStr := range testCases {
		t.Run(testName, func(t *testing.T) {
			// Convert JSON to Any
			anyValue, err := JSONToAny(jsonStr)
			assert.NoError(t, err)

			// Infer the storage type
			detectedType, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, GraphData, detectedType, "Expected graph data for test case: %s", testName)
		})
	}
}
