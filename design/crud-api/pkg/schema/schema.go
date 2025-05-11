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
	// First, determine the storage type
	storageType, err := sg.storageInferrer.InferType(anyValue)
	if err != nil {
		return nil, fmt.Errorf("failed to infer storage type: %v", err)
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
		return schema, nil
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

// handleGraphData processes graph data and generates field schemas.
// Currently, graph data is handled similarly to tabular data, as both
// are structured data types. In the future, this could be extended to
// handle relationships between entities.
//
// Parameters:
//   - anyValue: The protobuf Any value containing graph data
//   - schema: The base schema to populate with field information
//
// Returns:
//   - *SchemaInfo: The complete schema with field information
//   - error: Any error that occurred during processing
func (sg *SchemaGenerator) handleGraphData(anyValue *anypb.Any, schema *SchemaInfo) (*SchemaInfo, error) {
	// For now, handle graph data similar to tabular data
	return sg.handleTabularData(anyValue, schema)
}

// handleListData processes list data and generates item schema.
// List data is expected to be a struct with an "attributes" field containing
// an array of values.
//
// The function:
//  1. Extracts the list from attributes
//  2. If the list is empty, returns the base schema
//  3. Otherwise, generates a schema for the first item and sets it as the Items schema
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

	// Get the list value from attributes
	attrList, ok := attributes.GetKind().(*structpb.Value_ListValue)
	if !ok {
		return nil, fmt.Errorf("attributes is not a list")
	}

	// If the list is empty, return the schema as is
	if len(attrList.ListValue.Values) == 0 {
		return schema, nil
	}

	// Create a new Any value for the first item
	itemAny, err := anypb.New(&structpb.Struct{
		Fields: map[string]*structpb.Value{
			"attributes": attrList.ListValue.Values[0],
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
	return schema, nil
}

// handleMapData processes map data and generates property schemas.
// Map data is expected to be a struct with an "attributes" field containing
// a struct with property definitions.
//
// The function:
//  1. Extracts the attributes struct
//  2. For each property in the struct:
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

	// Generate schemas for each property
	schema.Properties = make(map[string]*SchemaInfo)
	for propName, propValue := range attrStruct.StructValue.Fields {
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
