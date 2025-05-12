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
				"attributes": {
					"value": 42
				}
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
				"attributes": {
					"numbers": [1, 2, 3]
				}
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
		"list_data_with_different_name": {
			input: `{
				"attributes": {
					"values": [1, 2, 3]
				}
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
					"properties": {
						"name": "John",
						"age": 30,
						"active": true
					}
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
					"properties": {
						"empty_str": "",
						"zero": 0,
						"null_val": null
					}
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
		"graph_data_with_nodes": {
			input: `{
				"attributes": {
					"nodes": {
						"person": {
							"name": "John",
							"age": 30,
							"active": true
						}
					}
				}
			}`,
			expected: `{
				"storage_type": "graph",
				"type_info": {
					"type": "string"
				},
				"fields": {
					"nodes": {
						"storage_type": "map",
						"type_info": {
							"type": "string"
						},
						"properties": {
							"person": {
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
							}
						}
					}
				}
			}`,
		},
		"graph_data_with_edges": {
			input: `{
				"attributes": {
					"edges": {
						"knows": {
							"since": "2020-01-01",
							"strength": 0.8
						}
					}
				}
			}`,
			expected: `{
				"storage_type": "graph",
				"type_info": {
					"type": "string"
				},
				"fields": {
					"edges": {
						"storage_type": "map",
						"type_info": {
							"type": "string"
						},
						"properties": {
							"knows": {
								"storage_type": "map",
								"type_info": {
									"type": "string"
								},
								"properties": {
									"since": {
										"storage_type": "scalar",
										"type_info": {
											"type": "date"
										}
									},
									"strength": {
										"storage_type": "scalar",
										"type_info": {
											"type": "float"
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"graph_data_with_both": {
			input: `{
				"attributes": {
					"nodes": {
						"person": {
							"name": "John",
							"age": 30
						}
					},
					"edges": {
						"knows": {
							"since": "2020-01-01"
						}
					}
				}
			}`,
			expected: `{
				"storage_type": "graph",
				"type_info": {
					"type": "string"
				},
				"fields": {
					"nodes": {
						"storage_type": "map",
						"type_info": {
							"type": "string"
						},
						"properties": {
							"person": {
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
									}
								}
							}
						}
					},
					"edges": {
						"storage_type": "map",
						"type_info": {
							"type": "string"
						},
						"properties": {
							"knows": {
								"storage_type": "map",
								"type_info": {
									"type": "string"
								},
								"properties": {
									"since": {
										"storage_type": "scalar",
										"type_info": {
											"type": "date"
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"social_network": {
			input: `{
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
			expected: `{
				"storage_type": "graph",
				"type_info": {
					"type": "string"
				},
				"fields": {
					"nodes": {
						"storage_type": "map",
						"type_info": {
							"type": "string"
						},
						"properties": {
							"user": {
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
									"location": {
										"storage_type": "scalar",
										"type_info": {
											"type": "string"
										}
									}
								}
							},
							"post": {
								"storage_type": "map",
								"type_info": {
									"type": "string"
								},
								"properties": {
									"title": {
										"storage_type": "scalar",
										"type_info": {
											"type": "string"
										}
									},
									"content": {
										"storage_type": "scalar",
										"type_info": {
											"type": "string"
										}
									},
									"created": {
										"storage_type": "scalar",
										"type_info": {
											"type": "date"
										}
									}
								}
							}
						}
					},
					"edges": {
						"storage_type": "map",
						"type_info": {
							"type": "string"
						},
						"properties": {
							"follows": {
								"storage_type": "map",
								"type_info": {
									"type": "string"
								},
								"properties": {
									"since": {
										"storage_type": "scalar",
										"type_info": {
											"type": "date"
										}
									}
								}
							},
							"created": {
								"storage_type": "map",
								"type_info": {
									"type": "string"
								},
								"properties": {
									"timestamp": {
										"storage_type": "scalar",
										"type_info": {
											"type": "datetime"
										}
									}
								}
							},
							"likes": {
								"storage_type": "map",
								"type_info": {
									"type": "string"
								},
								"properties": {
									"timestamp": {
										"storage_type": "scalar",
										"type_info": {
											"type": "datetime"
										}
									}
								}
							}
						}
					}
				}
			}`,
		},
		"dependency_graph": {
			input: `{
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
			expected: `{
				"storage_type": "graph",
				"type_info": {
					"type": "string"
				},
				"fields": {
					"nodes": {
						"storage_type": "map",
						"type_info": {
							"type": "string"
						},
						"properties": {
							"package": {
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
									"version": {
										"storage_type": "scalar",
										"type_info": {
											"type": "string"
										}
									}
								}
							}
						}
					},
					"edges": {
						"storage_type": "map",
						"type_info": {
							"type": "string"
						},
						"properties": {
							"depends_on": {
								"storage_type": "map",
								"type_info": {
									"type": "string"
								},
								"properties": {
									"version": {
										"storage_type": "scalar",
										"type_info": {
											"type": "string"
										}
									}
								}
							}
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
