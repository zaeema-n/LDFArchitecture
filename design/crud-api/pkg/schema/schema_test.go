package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"lk/datafoundation/crud-api/pkg/typeinference"
)

// SchemaInfoJSON represents the JSON structure of SchemaInfo
type SchemaInfoJSON struct {
	StorageType string                     `json:"storage_type"`
	TypeInfo    *TypeInfoJSON              `json:"type_info"`
	Fields      map[string]*SchemaInfoJSON `json:"fields,omitempty"`
	Items       *SchemaInfoJSON            `json:"items,omitempty"`
	Properties  map[string]*SchemaInfoJSON `json:"properties,omitempty"`
}

// TypeInfoJSON represents the JSON structure of TypeInfo
type TypeInfoJSON struct {
	Type       string        `json:"type"`
	IsNullable bool          `json:"is_nullable,omitempty"`
	IsArray    bool          `json:"is_array,omitempty"`
	ArrayType  *TypeInfoJSON `json:"array_type,omitempty"`
}

// SchemaInfoToJSON converts a SchemaInfo to its JSON representation
func SchemaInfoToJSON(schema *SchemaInfo) (*SchemaInfoJSON, error) {
	if schema == nil {
		return nil, nil
	}

	jsonSchema := &SchemaInfoJSON{
		StorageType: string(schema.StorageType),
		TypeInfo:    TypeInfoToJSON(schema.TypeInfo),
	}

	if schema.Fields != nil {
		jsonSchema.Fields = make(map[string]*SchemaInfoJSON)
		for k, v := range schema.Fields {
			fieldJSON, err := SchemaInfoToJSON(v)
			if err != nil {
				return nil, err
			}
			jsonSchema.Fields[k] = fieldJSON
		}
	}

	if schema.Items != nil {
		itemsJSON, err := SchemaInfoToJSON(schema.Items)
		if err != nil {
			return nil, err
		}
		jsonSchema.Items = itemsJSON
	}

	if schema.Properties != nil {
		jsonSchema.Properties = make(map[string]*SchemaInfoJSON)
		for k, v := range schema.Properties {
			propJSON, err := SchemaInfoToJSON(v)
			if err != nil {
				return nil, err
			}
			jsonSchema.Properties[k] = propJSON
		}
	}

	return jsonSchema, nil
}

// TypeInfoToJSON converts a TypeInfo to its JSON representation
func TypeInfoToJSON(typeInfo *typeinference.TypeInfo) *TypeInfoJSON {
	if typeInfo == nil {
		return nil
	}

	jsonTypeInfo := &TypeInfoJSON{
		Type:       string(typeInfo.Type),
		IsNullable: typeInfo.IsNullable,
		IsArray:    typeInfo.IsArray,
	}

	if typeInfo.ArrayType != nil {
		jsonTypeInfo.ArrayType = TypeInfoToJSON(typeInfo.ArrayType)
	}

	return jsonTypeInfo
}

// JSONToAny converts a JSON string to a protobuf Any value
func JSONToAny(jsonStr string) (*anypb.Any, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}

	structValue, err := structpb.NewStruct(data.(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	anyValue, err := anypb.New(structValue)
	if err != nil {
		return nil, err
	}

	return anyValue, nil
}

// TestSchemaGeneration tests the schema generation for different data structures
func TestSchemaGeneration(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected string
	}{
		"scalar_data": {
			input: `{
				"attributes": 42
			}`,
			expected: `{
				"storage_type": "scalar",
				"type_info": {
					"type": "int"
				}
			}`,
		},
		"list_data": {
			input: `{
				"attributes": [1, 2, 3]
			}`,
			expected: `{
				"storage_type": "list",
				"type_info": {
					"type": "string",
					"is_array": true,
					"array_type": {
						"type": "int"
					}
				},
				"items": {
					"storage_type": "scalar",
					"type_info": {
						"type": "int"
					}
				}
			}`,
		},
		"map_data": {
			input: `{
				"attributes": {
					"name": "John",
					"age": 30,
					"active": true
				}
			}`,
			expected: `{
				"storage_type": "map",
				"type_info": {
					"type": "string"
				},
				"properties": {
					"name": {
						"storage_type": "scalar",
						"type_info": {
							"type": "string"
						}
					},
					"age": {
						"storage_type": "scalar",
						"type_info": {
							"type": "int"
						}
					},
					"active": {
						"storage_type": "scalar",
						"type_info": {
							"type": "bool"
						}
					}
				}
			}`,
		},
		"empty_values": {
			input: `{
				"attributes": {
					"empty_str": "",
					"zero": 0,
					"null_val": null
				}
			}`,
			expected: `{
				"storage_type": "map",
				"type_info": {
					"type": "string"
				},
				"properties": {
					"empty_str": {
						"storage_type": "scalar",
						"type_info": {
							"type": "string"
						}
					},
					"zero": {
						"storage_type": "scalar",
						"type_info": {
							"type": "int"
						}
					},
					"null_val": {
						"storage_type": "scalar",
						"type_info": {
							"type": "null",
							"is_nullable": true
						}
					}
				}
			}`,
		},
	}

	generator := NewSchemaGenerator()
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			// Parse input
			anyValue, err := JSONToAny(tc.input)
			assert.NoError(t, err)

			// Generate schema
			schema, err := generator.GenerateSchema(anyValue)
			assert.NoError(t, err)

			// Convert schema to JSON
			schemaJSON, err := SchemaInfoToJSON(schema)
			assert.NoError(t, err)

			// Parse expected JSON
			var expectedJSON SchemaInfoJSON
			err = json.Unmarshal([]byte(tc.expected), &expectedJSON)
			assert.NoError(t, err)

			// Compare schemas
			assert.Equal(t, expectedJSON, *schemaJSON)
		})
	}
}
