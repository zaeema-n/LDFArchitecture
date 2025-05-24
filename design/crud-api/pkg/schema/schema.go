// Package schema provides functionality for generating and managing schema information
// for different types of data structures. It combines storage type inference and data
// type inference to create a complete schema representation that can be used for
// database schema generation, data validation, and API documentation.
//
// The package supports five main storage types:
//   - Scalar: Single values (numbers, strings, booleans)
//   - List: Arrays of values
//   - Map: Key-value pairs
//   - Tabular: Structured data with rows and columns
//   - Graph: Data with relationships between entities
//
// Each storage type can contain various data types (int, float, string, bool, etc.)
// and can be nested to represent complex data structures.
package schema

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	"lk/datafoundation/crud-api/pkg/storageinference"
	"lk/datafoundation/crud-api/pkg/typeinference"
)

// SchemaInfo represents the complete schema information for a data structure.
// It combines storage type information with detailed type information and
// maintains relationships between different parts of the schema.
//
// Fields:
//   - StorageType: Indicates how the data is organized (scalar, list, map, etc.)
//   - TypeInfo: Contains detailed information about the data type
//   - Fields: For tabular/graph data, contains schemas for each field
//   - Items: For list data, contains the schema for list items
//   - Properties: For map data, contains schemas for each property
type SchemaInfo struct {
	StorageType storageinference.StorageType // The storage type (tabular, graph, list, map, scalar)
	TypeInfo    *typeinference.TypeInfo      // The type information
	Fields      map[string]*SchemaInfo       // For tabular/graph data, contains field schemas
	Items       *SchemaInfo                  // For list data, contains item schema
	Properties  map[string]*SchemaInfo       // For map data, contains property schemas
}

// SchemaGenerator combines storage and type inference to generate complete schema information.
// It analyzes protobuf Any values to determine both how the data is stored and what
// types of data it contains.
type SchemaGenerator struct {
	storageInferrer *storageinference.StorageInferrer // Infers how data is organized
	typeInferrer    *typeinference.TypeInferrer       // Infers data types
}

// NewSchemaGenerator creates a new SchemaGenerator instance.
// It initializes both the storage inferrer and type inferrer with their default settings.
//
// Returns:
//   - *SchemaGenerator: A new instance ready to generate schemas
func NewSchemaGenerator() *SchemaGenerator {
	return &SchemaGenerator{
		storageInferrer: &storageinference.StorageInferrer{},
		typeInferrer:    &typeinference.TypeInferrer{},
	}
}

// GenerateSchema analyzes a protobuf Any value and generates complete schema information.
// The process involves two main steps:
//  1. Determine the storage type (how the data is organized)
//  2. Determine the data type (what kind of data it contains)
//
// The function handles different storage types differently:
//   - Scalar: Simple value types (int, float, string, bool)
//   - List: Arrays of values with consistent types
//   - Map: Key-value pairs with potentially different value types
//   - Tabular: Structured data with defined columns
//   - Graph: Data with relationships between entities
//
// Parameters:
//   - anyValue: A protobuf Any value containing the data to analyze
//
// Returns:
//   - *SchemaInfo: A complete schema representation of the data
//   - error: Any error that occurred during schema generation
func (sg *SchemaGenerator) GenerateSchema(anyValue *anypb.Any) (*SchemaInfo, error) {
	// Unpack the Any value to get the underlying message
	message, err := anyValue.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	// Get the struct value from the message
	structValue, ok := message.(*structpb.Struct)
	if !ok {
		// If not a struct, check if it's a scalar value
		switch m := message.(type) {
		case *structpb.Value:
			// Create a schema directly for the scalar value
			schema := &SchemaInfo{
				StorageType: storageinference.ScalarData,
				TypeInfo:    &typeinference.TypeInfo{},
			}
			switch m.GetKind().(type) {
			case *structpb.Value_StringValue:
				schema.TypeInfo.Type = typeinference.StringType
			case *structpb.Value_NumberValue:
				num := m.GetNumberValue()
				if num == float64(int64(num)) {
					schema.TypeInfo.Type = typeinference.IntType
				} else {
					schema.TypeInfo.Type = typeinference.FloatType
				}
			case *structpb.Value_BoolValue:
				schema.TypeInfo.Type = typeinference.BoolType
			case *structpb.Value_NullValue:
				schema.TypeInfo.Type = typeinference.NullType
				schema.TypeInfo.IsNullable = true
			default:
				return nil, fmt.Errorf("unsupported scalar type")
			}
			return schema, nil
		default:
			return nil, fmt.Errorf("expected struct value")
		}
	}

	// Check if this is a scalar value wrapped in a struct
	if len(structValue.Fields) == 1 {
		if value, ok := structValue.Fields["value"]; ok {
			// Create a schema directly for the scalar value
			schema := &SchemaInfo{
				StorageType: storageinference.ScalarData,
				TypeInfo:    &typeinference.TypeInfo{},
			}
			switch value.GetKind().(type) {
			case *structpb.Value_StringValue:
				schema.TypeInfo.Type = typeinference.StringType
			case *structpb.Value_NumberValue:
				num := value.GetNumberValue()
				if num == float64(int64(num)) {
					schema.TypeInfo.Type = typeinference.IntType
				} else {
					schema.TypeInfo.Type = typeinference.FloatType
				}
			case *structpb.Value_BoolValue:
				schema.TypeInfo.Type = typeinference.BoolType
			case *structpb.Value_NullValue:
				schema.TypeInfo.Type = typeinference.NullType
				schema.TypeInfo.IsNullable = true
			default:
				return nil, fmt.Errorf("unsupported scalar type")
			}
			return schema, nil
		}
	}

	// Determine storage type based on the structure of the data
	// The order of checks is important because:
	// 1. Graph data can contain both nodes and edges, which could be mistaken for other types
	// 2. Tabular data has a specific structure with columns and rows that should be identified before map/list
	// 3. Map data can contain nested structures that might look like lists
	// 4. List data is checked before scalar to handle arrays of values
	// 5. Scalar is the fallback for simple values
	var storageType storageinference.StorageType
	switch {
	case hasGraphStructure(structValue):
		storageType = storageinference.GraphData
	case hasTabularStructure(structValue):
		storageType = storageinference.TabularData
	case hasMapStructure(structValue):
		storageType = storageinference.MapData
	case hasListStructure(structValue):
		storageType = storageinference.ListData
	default:
		storageType = storageinference.ScalarData
	}

	// Determine the data type using type inference
	typeInfo, err := sg.typeInferrer.InferType(anyValue)
	if err != nil {
		return nil, fmt.Errorf("failed to infer data type: %v", err)
	}

	// Create the base schema info with storage type and type information
	schema := &SchemaInfo{
		StorageType: storageType,
		TypeInfo:    typeInfo,
	}

	// Handle different storage types with their specific processing functions
	switch storageType {
	case storageinference.TabularData:
		return sg.handleTabularData(structValue, schema)
	case storageinference.GraphData:
		return sg.handleGraphData(structValue, schema)
	case storageinference.ListData:
		return sg.handleListData(structValue, schema)
	case storageinference.MapData:
		return sg.handleMapData(structValue, schema)
	case storageinference.ScalarData:
		return sg.handleScalarData(structValue, schema)
	default:
		return nil, fmt.Errorf("unknown storage type: %v", storageType)
	}
}

// Helper functions to determine data structure type
func hasGraphStructure(structValue *structpb.Struct) bool {
	// Check for nodes or edges fields
	if _, hasNodes := structValue.Fields["nodes"]; hasNodes {
		return true
	}
	if _, hasEdges := structValue.Fields["edges"]; hasEdges {
		return true
	}
	return false
}

func hasListStructure(structValue *structpb.Struct) bool {
	// Check if any field is a list
	for _, field := range structValue.Fields {
		if _, ok := field.GetKind().(*structpb.Value_ListValue); ok {
			return true
		}
	}
	return false
}

func hasMapStructure(structValue *structpb.Struct) bool {
	// Check if any field is a struct
	for _, field := range structValue.Fields {
		if _, ok := field.GetKind().(*structpb.Value_StructValue); ok {
			return true
		}
	}
	return false
}

func hasTabularStructure(structValue *structpb.Struct) bool {
	// Check for both columns and rows fields directly in the struct
	columnsField, hasColumns := structValue.Fields["columns"]
	rowsField, hasRows := structValue.Fields["rows"]
	if !hasColumns || !hasRows {
		return false
	}

	// Verify columns is a list
	_, isColumnsList := columnsField.GetKind().(*structpb.Value_ListValue)
	if !isColumnsList {
		return false
	}

	// Verify rows is a list
	_, isRowsList := rowsField.GetKind().(*structpb.Value_ListValue)
	if !isRowsList {
		return false
	}

	return true
}

// Helper function to detect if a string is a date or datetime
func isDateOrDateTime(str string) (bool, bool) {
	// Try parsing as date (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", str); err == nil {
		return true, false
	}
	// Try parsing as datetime (YYYY-MM-DDTHH:MM:SSZ)
	if _, err := time.Parse(time.RFC3339, str); err == nil {
		return true, true
	}
	return false, false
}

// handleTabularData processes tabular data and generates field schemas.
func (sg *SchemaGenerator) handleTabularData(structValue *structpb.Struct, schema *SchemaInfo) (*SchemaInfo, error) {
	// Initialize the Fields map
	schema.Fields = make(map[string]*SchemaInfo)

	// Get columns and rows directly from the struct
	columnsField, hasColumns := structValue.Fields["columns"]
	rowsField, hasRows := structValue.Fields["rows"]
	if !hasColumns || !hasRows {
		return nil, fmt.Errorf("table must have both columns and rows")
	}

	// Verify columns is a list
	columnsList, ok := columnsField.GetKind().(*structpb.Value_ListValue)
	if !ok {
		return nil, fmt.Errorf("columns must be a list")
	}

	// Verify rows is a list
	rowsList, ok := rowsField.GetKind().(*structpb.Value_ListValue)
	if !ok {
		return nil, fmt.Errorf("rows must be a list")
	}

	// Get column names
	columnNames := make([]string, len(columnsList.ListValue.Values))
	for i, col := range columnsList.ListValue.Values {
		if strVal, ok := col.GetKind().(*structpb.Value_StringValue); ok {
			columnNames[i] = strVal.StringValue
		} else {
			return nil, fmt.Errorf("column names must be strings")
		}
	}

	// Process first row to determine types
	if len(rowsList.ListValue.Values) == 0 {
		return nil, fmt.Errorf("table must have at least one row")
	}

	firstRow := rowsList.ListValue.Values[0]
	rowValues, ok := firstRow.GetKind().(*structpb.Value_ListValue)
	if !ok {
		return nil, fmt.Errorf("row must be a list")
	}

	if len(rowValues.ListValue.Values) != len(columnNames) {
		return nil, fmt.Errorf("row length does not match number of columns")
	}

	// Create field schemas based on first row values
	for i, value := range rowValues.ListValue.Values {
		columnName := columnNames[i]
		fieldSchema := &SchemaInfo{
			StorageType: storageinference.ScalarData,
			TypeInfo:    &typeinference.TypeInfo{},
		}

		switch value.GetKind().(type) {
		case *structpb.Value_StringValue:
			str := value.GetStringValue()
			if isDate, isDateTime := isDateOrDateTime(str); isDate {
				if isDateTime {
					fieldSchema.TypeInfo.Type = typeinference.DateTimeType
				} else {
					fieldSchema.TypeInfo.Type = typeinference.DateType
				}
			} else {
				fieldSchema.TypeInfo.Type = typeinference.StringType
			}
		case *structpb.Value_NumberValue:
			num := value.GetNumberValue()
			if num == float64(int64(num)) {
				fieldSchema.TypeInfo.Type = typeinference.IntType
			} else {
				fieldSchema.TypeInfo.Type = typeinference.FloatType
			}
		case *structpb.Value_BoolValue:
			fieldSchema.TypeInfo.Type = typeinference.BoolType
		case *structpb.Value_NullValue:
			fieldSchema.TypeInfo.Type = typeinference.NullType
			fieldSchema.TypeInfo.IsNullable = true
		default:
			return nil, fmt.Errorf("unsupported value type in row")
		}

		schema.Fields[columnName] = fieldSchema
	}

	return schema, nil
}

// handleGraphData processes graph data and generates schemas for nodes and edges.
// Graph data is expected to be a struct with "nodes" and "edges" fields, where each field can be either:
//   - An array of objects with type and properties
//   - A map of type to property objects
//
// The function:
//  1. Processes node schemas from the "nodes" field
//  2. Processes edge schemas from the "edges" field
//  3. Combines them into a complete graph schema
//
// Parameters:
//   - structValue: The protobuf struct value containing graph data
//   - schema: The base schema to populate with graph information
//
// Returns:
//   - *SchemaInfo: The complete schema with node and edge information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleGraphData(structValue *structpb.Struct, schema *SchemaInfo) (*SchemaInfo, error) {
	// Initialize the schema fields map
	schema.Fields = make(map[string]*SchemaInfo)

	// Process nodes if present
	if nodes, ok := structValue.Fields["nodes"]; ok {
		// Create a map schema for nodes
		nodeSchema := &SchemaInfo{
			StorageType: storageinference.MapData,
			TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
			Properties:  make(map[string]*SchemaInfo),
		}

		// Handle both array and map formats for nodes
		switch nodeValue := nodes.GetKind().(type) {
		case *structpb.Value_ListValue:
			// Process array of nodes
			for _, node := range nodeValue.ListValue.Values {
				if nodeStruct, ok := node.GetKind().(*structpb.Value_StructValue); ok {
					// Get node type
					nodeType := "default"
					if typeField, ok := nodeStruct.StructValue.Fields["type"]; ok {
						if typeStr, ok := typeField.GetKind().(*structpb.Value_StringValue); ok {
							nodeType = typeStr.StringValue
						}
					}

					// Create a map schema for node properties
					propSchema := &SchemaInfo{
						StorageType: storageinference.MapData,
						TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
						Properties:  make(map[string]*SchemaInfo),
					}

					// Process properties - either directly in the node or in a properties field
					var properties *structpb.Struct
					if props, ok := nodeStruct.StructValue.Fields["properties"]; ok {
						if propStruct, ok := props.GetKind().(*structpb.Value_StructValue); ok {
							properties = propStruct.StructValue
						}
					} else {
						// Use the node struct itself as properties, excluding type and id fields
						properties = &structpb.Struct{
							Fields: make(map[string]*structpb.Value),
						}
						for k, v := range nodeStruct.StructValue.Fields {
							if k != "type" && k != "id" {
								properties.Fields[k] = v
							}
						}
					}

					if properties != nil {
						// Process each property
						for propName, propValue := range properties.Fields {
							// Create a schema for the property based on its type
							var propTypeSchema *SchemaInfo
							switch propValue.GetKind().(type) {
							case *structpb.Value_StringValue:
								str := propValue.GetStringValue()
								if isDate, isDateTime := isDateOrDateTime(str); isDate {
									if isDateTime {
										propTypeSchema = &SchemaInfo{
											StorageType: storageinference.ScalarData,
											TypeInfo:    &typeinference.TypeInfo{Type: typeinference.DateTimeType},
										}
									} else {
										propTypeSchema = &SchemaInfo{
											StorageType: storageinference.ScalarData,
											TypeInfo:    &typeinference.TypeInfo{Type: typeinference.DateType},
										}
									}
								} else {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
									}
								}
							case *structpb.Value_NumberValue:
								num := propValue.GetNumberValue()
								if num == float64(int64(num)) {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.IntType},
									}
								} else {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.FloatType},
									}
								}
							case *structpb.Value_BoolValue:
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.BoolType},
								}
							case *structpb.Value_NullValue:
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.NullType, IsNullable: true},
								}
							case *structpb.Value_StructValue:
								// For struct values, create a new Any value and generate schema
								propAny, err := anypb.New(&structpb.Struct{
									Fields: propValue.GetStructValue().Fields,
								})
								if err != nil {
									return nil, fmt.Errorf("failed to create property Any value: %v", err)
								}
								propTypeSchema, err = sg.GenerateSchema(propAny)
								if err != nil {
									return nil, fmt.Errorf("failed to generate property schema: %v", err)
								}
							default:
								return nil, fmt.Errorf("unsupported property type")
							}

							propSchema.Properties[propName] = propTypeSchema
						}

						// If we already have a schema for this node type, merge the properties
						if existingSchema, exists := nodeSchema.Properties[nodeType]; exists {
							for propName, propSchema := range propSchema.Properties {
								existingSchema.Properties[propName] = propSchema
							}
						} else {
							nodeSchema.Properties[nodeType] = propSchema
						}
					}
				}
			}
		case *structpb.Value_StructValue:
			// Process map of nodes
			for nodeType, nodeValue := range nodeValue.StructValue.Fields {
				if nodeStruct, ok := nodeValue.GetKind().(*structpb.Value_StructValue); ok {
					// Create a map schema for node properties
					propSchema := &SchemaInfo{
						StorageType: storageinference.MapData,
						TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
						Properties:  make(map[string]*SchemaInfo),
					}

					// Process each property
					for propName, propValue := range nodeStruct.StructValue.Fields {
						// Create a schema for the property based on its type
						var propTypeSchema *SchemaInfo
						switch propValue.GetKind().(type) {
						case *structpb.Value_StringValue:
							str := propValue.GetStringValue()
							if isDate, isDateTime := isDateOrDateTime(str); isDate {
								if isDateTime {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.DateTimeType},
									}
								} else {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.DateType},
									}
								}
							} else {
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
								}
							}
						case *structpb.Value_NumberValue:
							num := propValue.GetNumberValue()
							if num == float64(int64(num)) {
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.IntType},
								}
							} else {
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.FloatType},
								}
							}
						case *structpb.Value_BoolValue:
							propTypeSchema = &SchemaInfo{
								StorageType: storageinference.ScalarData,
								TypeInfo:    &typeinference.TypeInfo{Type: typeinference.BoolType},
							}
						case *structpb.Value_NullValue:
							propTypeSchema = &SchemaInfo{
								StorageType: storageinference.ScalarData,
								TypeInfo:    &typeinference.TypeInfo{Type: typeinference.NullType, IsNullable: true},
							}
						case *structpb.Value_StructValue:
							// For struct values, create a new Any value and generate schema
							propAny, err := anypb.New(&structpb.Struct{
								Fields: propValue.GetStructValue().Fields,
							})
							if err != nil {
								return nil, fmt.Errorf("failed to create property Any value: %v", err)
							}
							propTypeSchema, err = sg.GenerateSchema(propAny)
							if err != nil {
								return nil, fmt.Errorf("failed to generate property schema: %v", err)
							}
						default:
							return nil, fmt.Errorf("unsupported property type")
						}

						propSchema.Properties[propName] = propTypeSchema
					}

					nodeSchema.Properties[nodeType] = propSchema
				}
			}
		}

		schema.Fields["nodes"] = nodeSchema
	}

	// Process edges if present
	if edges, ok := structValue.Fields["edges"]; ok {
		// Create a map schema for edges
		edgeSchema := &SchemaInfo{
			StorageType: storageinference.MapData,
			TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
			Properties:  make(map[string]*SchemaInfo),
		}

		// Handle both array and map formats for edges
		switch edgeValue := edges.GetKind().(type) {
		case *structpb.Value_ListValue:
			// Process array of edges
			for _, edge := range edgeValue.ListValue.Values {
				if edgeStruct, ok := edge.GetKind().(*structpb.Value_StructValue); ok {
					// Get edge type
					edgeType := "default"
					if typeField, ok := edgeStruct.StructValue.Fields["type"]; ok {
						if typeStr, ok := typeField.GetKind().(*structpb.Value_StringValue); ok {
							edgeType = typeStr.StringValue
						}
					}

					// Create a map schema for edge properties
					propSchema := &SchemaInfo{
						StorageType: storageinference.MapData,
						TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
						Properties:  make(map[string]*SchemaInfo),
					}

					// Process properties - either directly in the edge or in a properties field
					var properties *structpb.Struct
					if props, ok := edgeStruct.StructValue.Fields["properties"]; ok {
						if propStruct, ok := props.GetKind().(*structpb.Value_StructValue); ok {
							properties = propStruct.StructValue
						}
					} else {
						// Use the edge struct itself as properties, excluding type, source, and target fields
						properties = &structpb.Struct{
							Fields: make(map[string]*structpb.Value),
						}
						for k, v := range edgeStruct.StructValue.Fields {
							if k != "type" && k != "source" && k != "target" {
								properties.Fields[k] = v
							}
						}
					}

					if properties != nil {
						// Process each property
						for propName, propValue := range properties.Fields {
							// Create a schema for the property based on its type
							var propTypeSchema *SchemaInfo
							switch propValue.GetKind().(type) {
							case *structpb.Value_StringValue:
								str := propValue.GetStringValue()
								if isDate, isDateTime := isDateOrDateTime(str); isDate {
									if isDateTime {
										propTypeSchema = &SchemaInfo{
											StorageType: storageinference.ScalarData,
											TypeInfo:    &typeinference.TypeInfo{Type: typeinference.DateTimeType},
										}
									} else {
										propTypeSchema = &SchemaInfo{
											StorageType: storageinference.ScalarData,
											TypeInfo:    &typeinference.TypeInfo{Type: typeinference.DateType},
										}
									}
								} else {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
									}
								}
							case *structpb.Value_NumberValue:
								num := propValue.GetNumberValue()
								if num == float64(int64(num)) {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.IntType},
									}
								} else {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.FloatType},
									}
								}
							case *structpb.Value_BoolValue:
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.BoolType},
								}
							case *structpb.Value_NullValue:
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.NullType, IsNullable: true},
								}
							case *structpb.Value_StructValue:
								// For struct values, create a new Any value and generate schema
								propAny, err := anypb.New(&structpb.Struct{
									Fields: propValue.GetStructValue().Fields,
								})
								if err != nil {
									return nil, fmt.Errorf("failed to create property Any value: %v", err)
								}
								propTypeSchema, err = sg.GenerateSchema(propAny)
								if err != nil {
									return nil, fmt.Errorf("failed to generate property schema: %v", err)
								}
							default:
								return nil, fmt.Errorf("unsupported property type")
							}

							propSchema.Properties[propName] = propTypeSchema
						}

						// If we already have a schema for this edge type, merge the properties
						if existingSchema, exists := edgeSchema.Properties[edgeType]; exists {
							for propName, propSchema := range propSchema.Properties {
								existingSchema.Properties[propName] = propSchema
							}
						} else {
							edgeSchema.Properties[edgeType] = propSchema
						}
					}
				}
			}
		case *structpb.Value_StructValue:
			// Process map of edges
			for edgeType, edgeValue := range edgeValue.StructValue.Fields {
				if edgeStruct, ok := edgeValue.GetKind().(*structpb.Value_StructValue); ok {
					// Create a map schema for edge properties
					propSchema := &SchemaInfo{
						StorageType: storageinference.MapData,
						TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
						Properties:  make(map[string]*SchemaInfo),
					}

					// Process each property
					for propName, propValue := range edgeStruct.StructValue.Fields {
						// Create a schema for the property based on its type
						var propTypeSchema *SchemaInfo
						switch propValue.GetKind().(type) {
						case *structpb.Value_StringValue:
							str := propValue.GetStringValue()
							if isDate, isDateTime := isDateOrDateTime(str); isDate {
								if isDateTime {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.DateTimeType},
									}
								} else {
									propTypeSchema = &SchemaInfo{
										StorageType: storageinference.ScalarData,
										TypeInfo:    &typeinference.TypeInfo{Type: typeinference.DateType},
									}
								}
							} else {
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
								}
							}
						case *structpb.Value_NumberValue:
							num := propValue.GetNumberValue()
							if num == float64(int64(num)) {
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.IntType},
								}
							} else {
								propTypeSchema = &SchemaInfo{
									StorageType: storageinference.ScalarData,
									TypeInfo:    &typeinference.TypeInfo{Type: typeinference.FloatType},
								}
							}
						case *structpb.Value_BoolValue:
							propTypeSchema = &SchemaInfo{
								StorageType: storageinference.ScalarData,
								TypeInfo:    &typeinference.TypeInfo{Type: typeinference.BoolType},
							}
						case *structpb.Value_NullValue:
							propTypeSchema = &SchemaInfo{
								StorageType: storageinference.ScalarData,
								TypeInfo:    &typeinference.TypeInfo{Type: typeinference.NullType, IsNullable: true},
							}
						case *structpb.Value_StructValue:
							// For struct values, create a new Any value and generate schema
							propAny, err := anypb.New(&structpb.Struct{
								Fields: propValue.GetStructValue().Fields,
							})
							if err != nil {
								return nil, fmt.Errorf("failed to create property Any value: %v", err)
							}
							propTypeSchema, err = sg.GenerateSchema(propAny)
							if err != nil {
								return nil, fmt.Errorf("failed to generate property schema: %v", err)
							}
						default:
							return nil, fmt.Errorf("unsupported property type")
						}

						propSchema.Properties[propName] = propTypeSchema
					}

					edgeSchema.Properties[edgeType] = propSchema
				}
			}
		}

		schema.Fields["edges"] = edgeSchema
	}

	// If neither nodes nor edges are present, return error
	if len(schema.Fields) == 0 {
		return nil, fmt.Errorf("graph data must contain either nodes or edges")
	}

	return schema, nil
}

// handleListData processes list data and generates item schema.
// List data can be represented in two ways:
//  1. Direct list: attributes is a list value
//  2. Wrapped list: attributes is a struct containing a list field
//
// The function:
//  1. Identifies the list value (either direct or in a struct field)
//  2. Creates a schema for the first item in the list
//  3. Sets the IsArray and ArrayType fields in the TypeInfo
//
// Parameters:
//   - structValue: The protobuf struct value containing list data
//   - schema: The base schema to populate with item information
//
// Returns:
//   - *SchemaInfo: The complete schema with item information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleListData(structValue *structpb.Struct, schema *SchemaInfo) (*SchemaInfo, error) {
	// Find the first list field
	var listField *structpb.Value
	for _, field := range structValue.Fields {
		if _, ok := field.GetKind().(*structpb.Value_ListValue); ok {
			listField = field
			break
		}
	}

	if listField == nil {
		return nil, fmt.Errorf("no list field found")
	}

	// Get the list value
	listVal, ok := listField.GetKind().(*structpb.Value_ListValue)
	if !ok {
		return nil, fmt.Errorf("field is not a list")
	}

	// If the list is empty, return the schema as is
	if len(listVal.ListValue.Values) == 0 {
		schema.TypeInfo.IsArray = true
		schema.TypeInfo.ArrayType = &typeinference.TypeInfo{Type: typeinference.StringType}
		return schema, nil
	}

	// Get the first item's value
	firstItem := listVal.ListValue.Values[0]

	// Create a schema for the item based on its type
	var itemSchema *SchemaInfo
	switch itemValue := firstItem.GetKind().(type) {
	case *structpb.Value_StructValue:
		// For struct values, create a new Any value and generate schema
		itemAny, err := anypb.New(&structpb.Struct{
			Fields: itemValue.StructValue.Fields,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create item Any value: %v", err)
		}
		itemSchema, err = sg.GenerateSchema(itemAny)
		if err != nil {
			return nil, fmt.Errorf("failed to generate item schema: %v", err)
		}
	case *structpb.Value_StringValue:
		// For string values
		itemSchema = &SchemaInfo{
			StorageType: storageinference.ScalarData,
			TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
		}
	case *structpb.Value_NumberValue:
		// For number values, check if it's an integer
		num := itemValue.NumberValue
		if num == float64(int64(num)) {
			itemSchema = &SchemaInfo{
				StorageType: storageinference.ScalarData,
				TypeInfo:    &typeinference.TypeInfo{Type: typeinference.IntType},
			}
		} else {
			itemSchema = &SchemaInfo{
				StorageType: storageinference.ScalarData,
				TypeInfo:    &typeinference.TypeInfo{Type: typeinference.FloatType},
			}
		}
	case *structpb.Value_BoolValue:
		// For boolean values
		itemSchema = &SchemaInfo{
			StorageType: storageinference.ScalarData,
			TypeInfo:    &typeinference.TypeInfo{Type: typeinference.BoolType},
		}
	case *structpb.Value_NullValue:
		// For null values
		itemSchema = &SchemaInfo{
			StorageType: storageinference.ScalarData,
			TypeInfo:    &typeinference.TypeInfo{Type: typeinference.NullType, IsNullable: true},
		}
	default:
		// For other types, create a new Any value and generate schema
		itemAny, err := anypb.New(firstItem)
		if err != nil {
			return nil, fmt.Errorf("failed to create item Any value: %v", err)
		}
		itemSchema, err = sg.GenerateSchema(itemAny)
		if err != nil {
			return nil, fmt.Errorf("failed to generate item schema: %v", err)
		}
	}

	schema.Items = itemSchema
	schema.TypeInfo.IsArray = true
	schema.TypeInfo.ArrayType = itemSchema.TypeInfo
	return schema, nil
}

// handleMapData processes map data and generates property schemas.
// Map data is expected to be a struct with an "attributes" field containing
// a struct with a field that contains key-value pairs.
//
// The function:
//  1. Extracts the attributes struct
//  2. Finds the first struct field that contains key-value pairs
//  3. For each property in the struct:
//     - Creates a new Any value for the property
//     - Recursively generates a schema for the property
//     - Adds the property schema to the Properties map
//
// Parameters:
//   - structValue: The protobuf struct value containing map data
//   - schema: The base schema to populate with property information
//
// Returns:
//   - *SchemaInfo: The complete schema with property information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleMapData(structValue *structpb.Struct, schema *SchemaInfo) (*SchemaInfo, error) {
	// Initialize the Properties map
	schema.Properties = make(map[string]*SchemaInfo)

	// Stack to keep track of nested structures to process
	type stackItem struct {
		structValue *structpb.Struct
		schema      *SchemaInfo
		fieldName   string
	}
	stack := []stackItem{{structValue: structValue, schema: schema}}

	// Process stack until empty
	for len(stack) > 0 {
		// Pop the top item
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Check if we have a "properties" field
		if props, ok := item.structValue.Fields["properties"]; ok {
			if propStruct, ok := props.GetKind().(*structpb.Value_StructValue); ok {
				// Process properties from the properties field
				for fieldName, fieldValue := range propStruct.StructValue.Fields {
					// Handle scalar values directly
					if _, ok := fieldValue.GetKind().(*structpb.Value_StringValue); ok {
						item.schema.Properties[fieldName] = &SchemaInfo{
							StorageType: storageinference.ScalarData,
							TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
						}
					} else if numValue, ok := fieldValue.GetKind().(*structpb.Value_NumberValue); ok {
						num := numValue.NumberValue
						if num == float64(int64(num)) {
							item.schema.Properties[fieldName] = &SchemaInfo{
								StorageType: storageinference.ScalarData,
								TypeInfo:    &typeinference.TypeInfo{Type: typeinference.IntType},
							}
						} else {
							item.schema.Properties[fieldName] = &SchemaInfo{
								StorageType: storageinference.ScalarData,
								TypeInfo:    &typeinference.TypeInfo{Type: typeinference.FloatType},
							}
						}
					} else if _, ok := fieldValue.GetKind().(*structpb.Value_BoolValue); ok {
						item.schema.Properties[fieldName] = &SchemaInfo{
							StorageType: storageinference.ScalarData,
							TypeInfo:    &typeinference.TypeInfo{Type: typeinference.BoolType},
						}
					} else if _, ok := fieldValue.GetKind().(*structpb.Value_NullValue); ok {
						item.schema.Properties[fieldName] = &SchemaInfo{
							StorageType: storageinference.ScalarData,
							TypeInfo:    &typeinference.TypeInfo{Type: typeinference.NullType, IsNullable: true},
						}
					} else if structValue, ok := fieldValue.GetKind().(*structpb.Value_StructValue); ok {
						// For nested structures, create a map schema and add to stack
						nestedSchema := &SchemaInfo{
							StorageType: storageinference.MapData,
							TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
							Properties:  make(map[string]*SchemaInfo),
						}
						item.schema.Properties[fieldName] = nestedSchema
						stack = append(stack, stackItem{
							structValue: structValue.StructValue,
							schema:      nestedSchema,
							fieldName:   fieldName,
						})
					}
				}
				continue
			}
		}

		// If no properties field, process fields directly
		for fieldName, fieldValue := range item.structValue.Fields {
			// Handle scalar values directly
			if _, ok := fieldValue.GetKind().(*structpb.Value_StringValue); ok {
				item.schema.Properties[fieldName] = &SchemaInfo{
					StorageType: storageinference.ScalarData,
					TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
				}
			} else if numValue, ok := fieldValue.GetKind().(*structpb.Value_NumberValue); ok {
				num := numValue.NumberValue
				if num == float64(int64(num)) {
					item.schema.Properties[fieldName] = &SchemaInfo{
						StorageType: storageinference.ScalarData,
						TypeInfo:    &typeinference.TypeInfo{Type: typeinference.IntType},
					}
				} else {
					item.schema.Properties[fieldName] = &SchemaInfo{
						StorageType: storageinference.ScalarData,
						TypeInfo:    &typeinference.TypeInfo{Type: typeinference.FloatType},
					}
				}
			} else if _, ok := fieldValue.GetKind().(*structpb.Value_BoolValue); ok {
				item.schema.Properties[fieldName] = &SchemaInfo{
					StorageType: storageinference.ScalarData,
					TypeInfo:    &typeinference.TypeInfo{Type: typeinference.BoolType},
				}
			} else if _, ok := fieldValue.GetKind().(*structpb.Value_NullValue); ok {
				item.schema.Properties[fieldName] = &SchemaInfo{
					StorageType: storageinference.ScalarData,
					TypeInfo:    &typeinference.TypeInfo{Type: typeinference.NullType, IsNullable: true},
				}
			} else if structValue, ok := fieldValue.GetKind().(*structpb.Value_StructValue); ok {
				// For nested structures, create a map schema and add to stack
				nestedSchema := &SchemaInfo{
					StorageType: storageinference.MapData,
					TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
					Properties:  make(map[string]*SchemaInfo),
				}
				item.schema.Properties[fieldName] = nestedSchema
				stack = append(stack, stackItem{
					structValue: structValue.StructValue,
					schema:      nestedSchema,
					fieldName:   fieldName,
				})
			}
		}
	}

	return schema, nil
}

// handleScalarData processes scalar data and generates the appropriate schema.
// Scalar data can be represented in two ways:
//  1. Direct value: attributes is a scalar value
//  2. Wrapped value: attributes is a struct containing a scalar value in any field
//
// The function:
//  1. Identifies the scalar value (either direct or in a struct field)
//  2. Uses type inference to determine the data type
//  3. Updates the schema with the inferred type information
//
// Parameters:
//   - structValue: The protobuf struct value containing scalar data
//   - schema: The base schema to populate with type information
//
// Returns:
//   - *SchemaInfo: The complete schema with type information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleScalarData(structValue *structpb.Struct, schema *SchemaInfo) (*SchemaInfo, error) {
	// Find the first scalar field
	var scalarField *structpb.Value
	for _, field := range structValue.Fields {
		switch field.GetKind().(type) {
		case *structpb.Value_StringValue, *structpb.Value_NumberValue, *structpb.Value_BoolValue, *structpb.Value_NullValue:
			scalarField = field
			break
		}
	}

	if scalarField == nil {
		return nil, fmt.Errorf("no scalar field found")
	}

	// Create a schema based on the scalar type
	switch scalarField.GetKind().(type) {
	case *structpb.Value_StringValue:
		schema.TypeInfo.Type = typeinference.StringType
	case *structpb.Value_NumberValue:
		// Check if the number is an integer
		num := scalarField.GetNumberValue()
		if num == float64(int64(num)) {
			schema.TypeInfo.Type = typeinference.IntType
		} else {
			schema.TypeInfo.Type = typeinference.FloatType
		}
	case *structpb.Value_BoolValue:
		schema.TypeInfo.Type = typeinference.BoolType
	case *structpb.Value_NullValue:
		schema.TypeInfo.Type = typeinference.NullType
		schema.TypeInfo.IsNullable = true
	}

	return schema, nil
}
