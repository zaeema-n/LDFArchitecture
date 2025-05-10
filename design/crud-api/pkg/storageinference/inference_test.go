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
