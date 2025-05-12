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
	storageInferrer *storageinference.StorageInferrer
	typeInferrer    *typeinference.TypeInferrer
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
	// Unpack the Any value
	message, err := anyValue.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	// Get the struct value
	structValue, ok := message.(*structpb.Struct)
	if !ok {
		return nil, fmt.Errorf("expected struct value")
	}

	// Get the attributes field
	attributes, ok := structValue.Fields["attributes"]
	if !ok {
		return nil, fmt.Errorf("attributes field not found")
	}

	// Determine storage type based on the structure of attributes
	var storageType storageinference.StorageType
	switch attr := attributes.GetKind().(type) {
	case *structpb.Value_ListValue:
		// Direct list values are always ListData
		storageType = storageinference.ListData
	case *structpb.Value_StructValue:
		// Check if it's a graph (has nodes or edges field)
		if _, hasNodes := attr.StructValue.Fields["nodes"]; hasNodes {
			storageType = storageinference.GraphData
		} else if _, hasEdges := attr.StructValue.Fields["edges"]; hasEdges {
			storageType = storageinference.GraphData
		} else {
			// Check if it's a list (has any list field)
			for _, field := range attr.StructValue.Fields {
				if _, ok := field.GetKind().(*structpb.Value_ListValue); ok {
					storageType = storageinference.ListData
					break
				}
			}
			// If not a list, check if it's a map (has any struct field)
			if storageType == "" {
				for _, field := range attr.StructValue.Fields {
					if _, ok := field.GetKind().(*structpb.Value_StructValue); ok {
						storageType = storageinference.MapData
						break
					}
				}
			}
			// If no list or map found, treat as scalar
			if storageType == "" {
				storageType = storageinference.ScalarData
			}
		}
	default:
		// Any other value type is treated as scalar
		storageType = storageinference.ScalarData
	}

	// Then, determine the data type
	typeInfo, err := sg.typeInferrer.InferType(anyValue)
	if err != nil {
		return nil, fmt.Errorf("failed to infer data type: %v", err)
	}

	// Create the base schema info
	schema := &SchemaInfo{
		StorageType: storageType,
		TypeInfo:    typeInfo,
	}

	// Handle different storage types
	switch storageType {
	case storageinference.TabularData:
		return sg.handleTabularData(anyValue, schema)
	case storageinference.GraphData:
		return sg.handleGraphData(anyValue, schema)
	case storageinference.ListData:
		return sg.handleListData(anyValue, schema)
	case storageinference.MapData:
		return sg.handleMapData(anyValue, schema)
	case storageinference.ScalarData:
		return sg.handleScalarData(anyValue, schema)
	default:
		return nil, fmt.Errorf("unknown storage type: %v", storageType)
	}
}

// handleTabularData processes tabular data and generates field schemas.
// Tabular data is expected to be a struct with an "attributes" field containing
// a struct with field definitions.
//
// The function:
//  1. Extracts the attributes struct
//  2. For each field in the struct:
//     - Creates a new Any value for the field
//     - Recursively generates a schema for the field
//     - Adds the field schema to the Fields map
//
// Parameters:
//   - anyValue: The protobuf Any value containing tabular data
//   - schema: The base schema to populate with field information
//
// Returns:
//   - *SchemaInfo: The complete schema with field information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleTabularData(anyValue *anypb.Any, schema *SchemaInfo) (*SchemaInfo, error) {
	// Unpack the Any value
	message, err := anyValue.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	// Get the struct value
	structValue, ok := message.(*structpb.Struct)
	if !ok {
		return nil, fmt.Errorf("expected struct value for tabular data")
	}

	// Get the attributes field
	attributes, ok := structValue.Fields["attributes"]
	if !ok {
		return nil, fmt.Errorf("attributes field not found")
	}

	// Get the struct value from attributes
	attrStruct, ok := attributes.GetKind().(*structpb.Value_StructValue)
	if !ok {
		return nil, fmt.Errorf("attributes is not a struct")
	}

	// Generate schemas for each field
	schema.Fields = make(map[string]*SchemaInfo)
	for fieldName, fieldValue := range attrStruct.StructValue.Fields {
		// Create a new Any value for the field
		fieldAny, err := anypb.New(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"attributes": fieldValue,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create field Any value: %v", err)
		}

		// Generate schema for the field
		fieldSchema, err := sg.GenerateSchema(fieldAny)
		if err != nil {
			return nil, fmt.Errorf("failed to generate field schema: %v", err)
		}

		schema.Fields[fieldName] = fieldSchema
	}

	return schema, nil
}

// handleGraphData processes graph data and generates schemas for nodes and edges.
// Graph data is expected to be a struct with an "attributes" field containing
// a struct with "nodes" and "edges" fields, where each field is an array of objects
// with properties.
//
// The function:
//  1. Extracts the attributes struct
//  2. Processes node schemas from the "nodes" field
//  3. Processes edge schemas from the "edges" field
//  4. Combines them into a complete graph schema
//
// Parameters:
//   - anyValue: The protobuf Any value containing graph data
//   - schema: The base schema to populate with graph information
//
// Returns:
//   - *SchemaInfo: The complete schema with node and edge information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleGraphData(anyValue *anypb.Any, schema *SchemaInfo) (*SchemaInfo, error) {
	// Unpack the Any value
	message, err := anyValue.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	// Get the struct value
	structValue, ok := message.(*structpb.Struct)
	if !ok {
		return nil, fmt.Errorf("expected struct value for graph data")
	}

	// Get the attributes field
	attributes, ok := structValue.Fields["attributes"]
	if !ok {
		return nil, fmt.Errorf("attributes field not found")
	}

	// Get the struct value from attributes
	attrStruct, ok := attributes.GetKind().(*structpb.Value_StructValue)
	if !ok {
		return nil, fmt.Errorf("attributes is not a struct")
	}

	// Initialize the schema fields
	schema.Fields = make(map[string]*SchemaInfo)

	// Process nodes if present
	if nodes, ok := attrStruct.StructValue.Fields["nodes"]; ok {
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

					// Get node properties
					if props, ok := nodeStruct.StructValue.Fields["properties"]; ok {
						if propStruct, ok := props.GetKind().(*structpb.Value_StructValue); ok {
							// Create a map schema for node properties
							propSchema := &SchemaInfo{
								StorageType: storageinference.MapData,
								TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
								Properties:  make(map[string]*SchemaInfo),
							}

							// Process each property
							for propName, propValue := range propStruct.StructValue.Fields {
								// Create a new Any value for the property
								propAny, err := anypb.New(&structpb.Struct{
									Fields: map[string]*structpb.Value{
										"attributes": propValue,
									},
								})
								if err != nil {
									return nil, fmt.Errorf("failed to create property Any value: %v", err)
								}

								// Generate schema for the property
								propTypeSchema, err := sg.GenerateSchema(propAny)
								if err != nil {
									return nil, fmt.Errorf("failed to generate property schema: %v", err)
								}

								propSchema.Properties[propName] = propTypeSchema
							}

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
						// Create a new Any value for the property
						propAny, err := anypb.New(&structpb.Struct{
							Fields: map[string]*structpb.Value{
								"attributes": propValue,
							},
						})
						if err != nil {
							return nil, fmt.Errorf("failed to create property Any value: %v", err)
						}

						// Generate schema for the property
						propTypeSchema, err := sg.GenerateSchema(propAny)
						if err != nil {
							return nil, fmt.Errorf("failed to generate property schema: %v", err)
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
	if edges, ok := attrStruct.StructValue.Fields["edges"]; ok {
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

					// Get edge properties
					if props, ok := edgeStruct.StructValue.Fields["properties"]; ok {
						if propStruct, ok := props.GetKind().(*structpb.Value_StructValue); ok {
							// Create a map schema for edge properties
							propSchema := &SchemaInfo{
								StorageType: storageinference.MapData,
								TypeInfo:    &typeinference.TypeInfo{Type: typeinference.StringType},
								Properties:  make(map[string]*SchemaInfo),
							}

							// Process each property
							for propName, propValue := range propStruct.StructValue.Fields {
								// Create a new Any value for the property
								propAny, err := anypb.New(&structpb.Struct{
									Fields: map[string]*structpb.Value{
										"attributes": propValue,
									},
								})
								if err != nil {
									return nil, fmt.Errorf("failed to create property Any value: %v", err)
								}

								// Generate schema for the property
								propTypeSchema, err := sg.GenerateSchema(propAny)
								if err != nil {
									return nil, fmt.Errorf("failed to generate property schema: %v", err)
								}

								propSchema.Properties[propName] = propTypeSchema
							}

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
						// Create a new Any value for the property
						propAny, err := anypb.New(&structpb.Struct{
							Fields: map[string]*structpb.Value{
								"attributes": propValue,
							},
						})
						if err != nil {
							return nil, fmt.Errorf("failed to create property Any value: %v", err)
						}

						// Generate schema for the property
						propTypeSchema, err := sg.GenerateSchema(propAny)
						if err != nil {
							return nil, fmt.Errorf("failed to generate property schema: %v", err)
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
//   - anyValue: The protobuf Any value containing list data
//   - schema: The base schema to populate with item information
//
// Returns:
//   - *SchemaInfo: The complete schema with item information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleListData(anyValue *anypb.Any, schema *SchemaInfo) (*SchemaInfo, error) {
	// Unpack the Any value
	message, err := anyValue.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	// Get the struct value
	structValue, ok := message.(*structpb.Struct)
	if !ok {
		return nil, fmt.Errorf("expected struct value for list data")
	}

	// Get the attributes field
	attributes, ok := structValue.Fields["attributes"]
	if !ok {
		return nil, fmt.Errorf("attributes field not found")
	}

	// Handle both direct list values and lists wrapped in a struct
	switch attr := attributes.GetKind().(type) {
	case *structpb.Value_ListValue:
		// If attributes is a direct list, use it
		if len(attr.ListValue.Values) == 0 {
			schema.TypeInfo.IsArray = true
			schema.TypeInfo.ArrayType = &typeinference.TypeInfo{Type: typeinference.StringType}
			return schema, nil
		}

		// Create a new Any value for the first item
		itemAny, err := anypb.New(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"attributes": attr.ListValue.Values[0],
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create item Any value: %v", err)
		}

		// Generate schema for the item
		itemSchema, err := sg.GenerateSchema(itemAny)
		if err != nil {
			return nil, fmt.Errorf("failed to generate item schema: %v", err)
		}

		schema.Items = itemSchema
		schema.TypeInfo.IsArray = true
		schema.TypeInfo.ArrayType = itemSchema.TypeInfo
		return schema, nil

	case *structpb.Value_StructValue:
		// If attributes is a struct, find the first list field
		var listField *structpb.Value
		var listFieldName string
		for name, field := range attr.StructValue.Fields {
			if _, ok := field.GetKind().(*structpb.Value_ListValue); ok {
				listField = field
				listFieldName = name
				break
			}
		}

		if listField == nil {
			return nil, fmt.Errorf("no list field found in attributes")
		}

		// Get the list value
		listVal, ok := listField.GetKind().(*structpb.Value_ListValue)
		if !ok {
			return nil, fmt.Errorf("field %s is not a list", listFieldName)
		}

		// If the list is empty, return the schema as is
		if len(listVal.ListValue.Values) == 0 {
			schema.TypeInfo.IsArray = true
			schema.TypeInfo.ArrayType = &typeinference.TypeInfo{Type: typeinference.StringType}
			return schema, nil
		}

		// Create a new Any value for the first item
		itemAny, err := anypb.New(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"attributes": listVal.ListValue.Values[0],
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create item Any value: %v", err)
		}

		// Generate schema for the item
		itemSchema, err := sg.GenerateSchema(itemAny)
		if err != nil {
			return nil, fmt.Errorf("failed to generate item schema: %v", err)
		}

		schema.Items = itemSchema
		schema.TypeInfo.IsArray = true
		schema.TypeInfo.ArrayType = itemSchema.TypeInfo
		return schema, nil

	default:
		return nil, fmt.Errorf("attributes is not a list or struct containing a list")
	}
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
//   - anyValue: The protobuf Any value containing map data
//   - schema: The base schema to populate with property information
//
// Returns:
//   - *SchemaInfo: The complete schema with property information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleMapData(anyValue *anypb.Any, schema *SchemaInfo) (*SchemaInfo, error) {
	// Unpack the Any value
	message, err := anyValue.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	// Get the struct value
	structValue, ok := message.(*structpb.Struct)
	if !ok {
		return nil, fmt.Errorf("expected struct value for map data")
	}

	// Get the attributes field
	attributes, ok := structValue.Fields["attributes"]
	if !ok {
		return nil, fmt.Errorf("attributes field not found")
	}

	// Get the struct value from attributes
	attrStruct, ok := attributes.GetKind().(*structpb.Value_StructValue)
	if !ok {
		return nil, fmt.Errorf("attributes is not a struct")
	}

	// Find the first struct field that contains key-value pairs
	var mapField *structpb.Value
	var mapFieldName string
	for name, field := range attrStruct.StructValue.Fields {
		if _, ok := field.GetKind().(*structpb.Value_StructValue); ok {
			mapField = field
			mapFieldName = name
			break
		}
	}

	if mapField == nil {
		return nil, fmt.Errorf("no map field found in attributes")
	}

	// Get the struct value from the map field
	mapStruct, ok := mapField.GetKind().(*structpb.Value_StructValue)
	if !ok {
		return nil, fmt.Errorf("field %s is not a struct", mapFieldName)
	}

	// Generate schemas for each property
	schema.Properties = make(map[string]*SchemaInfo)
	for propName, propValue := range mapStruct.StructValue.Fields {
		// Create a new Any value for the property
		propAny, err := anypb.New(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"attributes": propValue,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create property Any value: %v", err)
		}

		// Generate schema for the property
		propSchema, err := sg.GenerateSchema(propAny)
		if err != nil {
			return nil, fmt.Errorf("failed to generate property schema: %v", err)
		}

		schema.Properties[propName] = propSchema
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
//   - anyValue: The protobuf Any value containing scalar data
//   - schema: The base schema to populate with type information
//
// Returns:
//   - *SchemaInfo: The complete schema with type information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleScalarData(anyValue *anypb.Any, schema *SchemaInfo) (*SchemaInfo, error) {
	// Unpack the Any value
	message, err := anyValue.UnmarshalNew()
	if err != nil {
		return nil, err
	}

	// Get the struct value
	structValue, ok := message.(*structpb.Struct)
	if !ok {
		return nil, fmt.Errorf("expected struct value for scalar data")
	}

	// Get the attributes field
	attributes, ok := structValue.Fields["attributes"]
	if !ok {
		return nil, fmt.Errorf("attributes field not found")
	}

	// Handle both direct values and values wrapped in a struct
	switch attr := attributes.GetKind().(type) {
	case *structpb.Value_StructValue:
		// If attributes is a struct, find the first scalar field
		var scalarField *structpb.Value
		for _, field := range attr.StructValue.Fields {
			// Check if the field is a scalar value (not a struct or list)
			switch field.GetKind().(type) {
			case *structpb.Value_NumberValue, *structpb.Value_StringValue,
				*structpb.Value_BoolValue, *structpb.Value_NullValue:
				scalarField = field
			}
		}

		if scalarField == nil {
			return nil, fmt.Errorf("no scalar field found in attributes")
		}

		// Create a new Any value for the scalar value
		valueAny, err := anypb.New(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"attributes": scalarField,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create value Any: %v", err)
		}
		// Use type inference directly
		typeInfo, err := sg.typeInferrer.InferType(valueAny)
		if err != nil {
			return nil, fmt.Errorf("failed to infer type: %v", err)
		}
		schema.TypeInfo = typeInfo
		return schema, nil

	default:
		// If attributes is a direct value, use it directly
		valueAny, err := anypb.New(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"attributes": attributes,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create value Any: %v", err)
		}
		// Use type inference directly
		typeInfo, err := sg.typeInferrer.InferType(valueAny)
		if err != nil {
			return nil, fmt.Errorf("failed to infer type: %v", err)
		}
		schema.TypeInfo = typeInfo
		return schema, nil
	}
}
