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
		storageType = storageinference.ListData
	case *structpb.Value_StructValue:
		// Check if it's a map (has properties)
		if _, ok := attr.StructValue.Fields["properties"]; ok {
			storageType = storageinference.MapData
		} else {
			// Check if it's a list (has items)
			for _, field := range attr.StructValue.Fields {
				if _, ok := field.GetKind().(*structpb.Value_ListValue); ok {
					storageType = storageinference.ListData
					break
				}
			}
			// If no list found, treat as scalar
			if storageType == "" {
				storageType = storageinference.ScalarData
			}
		}
	default:
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

	// Get the properties field
	properties, ok := attrStruct.StructValue.Fields["properties"]
	if !ok {
		return nil, fmt.Errorf("properties field not found in attributes")
	}

	// Get the struct value from properties
	propsStruct, ok := properties.GetKind().(*structpb.Value_StructValue)
	if !ok {
		return nil, fmt.Errorf("properties is not a struct")
	}

	// Generate schemas for each property
	schema.Properties = make(map[string]*SchemaInfo)
	for propName, propValue := range propsStruct.StructValue.Fields {
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
		// If attributes is a struct, look for a "value" field
		value, ok := attr.StructValue.Fields["value"]
		if !ok {
			return nil, fmt.Errorf("value field not found in attributes")
		}
		// Create a new Any value for the scalar value
		valueAny, err := anypb.New(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"attributes": value,
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
