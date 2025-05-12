package storageinference

import (
	"reflect"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// StorageType represents the type of data structure
type StorageType string

const (
	TabularData StorageType = "tabular"
	ScalarData  StorageType = "scalar"
	ListData    StorageType = "list"
	MapData     StorageType = "map"
	GraphData   StorageType = "graph"
)

// TypeInferrer provides functionality to infer data types from protobuf Any values
type StorageInferrer struct{}

// InferType attempts to determine the storage type from a protobuf Any value.
// The function follows a hierarchical approach to identify the storage type:
//
// 1. First, it unpacks the Any value to get the underlying message:
//
//   - Uses UnmarshalNew() to convert the Any value to its concrete type
//
//   - Handles any unmarshaling errors that might occur
//
//     2. For structpb.Struct messages (most common case):
//     a. Checks for the presence of an "attributes" field
//     b. Based on the kind of the attributes field:
//
//   - If it's a StructValue:
//
//   - Checks for tabular structure (has both "columns" and "rows" fields)
//
//   - Checks for graph structure (has both "nodes" and "edges" fields)
//
//   - Checks for list structure (has "items" field with ListValue)
//
//   - Checks for scalar structure (single field with scalar value)
//
//   - If none of the above, defaults to MapData
//
//   - If it's a ListValue: returns ListData
//
//   - If it's a NumberValue, StringValue, or BoolValue: returns ScalarData
//
//   - For other cases: defaults to ScalarData
//
// 3. For non-structpb.Struct messages:
//   - Uses reflection to determine the type:
//   - Slice/Array types return ListData
//   - Map types return MapData
//   - All other types return ScalarData
//
// The function returns one of the following StorageType values:
// - TabularData: For data with columns and rows structure
// - GraphData: For data with nodes and edges structure
// - ListData: For array-like data structures
// - MapData: For key-value pair structures
// - ScalarData: For single value data
//
// Example JSON structures for each type:
//
// TabularData:
//
//	{
//	  "attributes": {
//	    "columns": ["id", "name"],
//	    "rows": [[1, "John"], [2, "Jane"]]
//	  }
//	}
//
// GraphData:
//
//	{
//	  "attributes": {
//	    "nodes": [{"id": "1", "type": "user"}],
//	    "edges": [{"source": "1", "target": "2"}]
//	  }
//	}
//
// ListData:
//
//	{
//	  "attributes": {
//	    "items": [1, 2, 3]
//	  }
//	}
//
// MapData:
//
//	{
//	  "attributes": {
//	    "key1": "value1",
//	    "key2": "value2"
//	  }
//	}
//
// ScalarData:
//
//	{
//	  "attributes": {
//	    "value": 42
//	  }
//	}
//
// Error handling:
// - Returns error if Any value cannot be unmarshaled
// - Returns error if the message type is not supported
// - Returns error if the structure is invalid
//
// Note: The function prioritizes specific structures over generic ones.
// For example, if a structure has both "items" and "nodes"/"edges",
// it will be classified based on the more specific structure first.
func (ti *StorageInferrer) InferType(anyValue *anypb.Any) (StorageType, error) {
	// Unpack the Any value to get the underlying message
	message, err := anyValue.UnmarshalNew()
	if err != nil {
		return "", err
	}

	// Get the reflection value of the message
	rv := reflect.ValueOf(message)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// Handle structpb.Struct
	if structValue, ok := message.(*structpb.Struct); ok {
		// Get the attributes field
		if attributes, ok := structValue.Fields["attributes"]; ok {
			switch attributes.GetKind().(type) {
			case *structpb.Value_StructValue:
				attrStruct := attributes.GetStructValue()
				// Check if it's a tabular structure
				if isTabular(attrStruct) {
					return TabularData, nil
				}
				// Check if it's a graph structure
				if isGraph(attrStruct) {
					return GraphData, nil
				}
				// Check if it's a list structure (has "items" field)
				if items, ok := attrStruct.Fields["items"]; ok {
					if _, ok := items.GetKind().(*structpb.Value_ListValue); ok {
						return ListData, nil
					}
				}
				// Check if it's a scalar structure
				if isScalar(attrStruct) {
					return ScalarData, nil
				}
				return MapData, nil
			case *structpb.Value_ListValue:
				return ListData, nil
			case *structpb.Value_NumberValue, *structpb.Value_StringValue, *structpb.Value_BoolValue:
				return ScalarData, nil
			default:
				return ScalarData, nil
			}
		}
	}

	// If not a structpb.Struct, check the direct type
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return ListData, nil
	case reflect.Map:
		return MapData, nil
	default:
		return ScalarData, nil
	}
}

// isTabular checks if a struct represents tabular data
func isTabular(structValue *structpb.Struct) bool {
	// A struct is considered tabular if it has both columns and rows fields
	_, hasColumns := structValue.Fields["columns"]
	_, hasRows := structValue.Fields["rows"]
	return hasColumns && hasRows
}

// isGraph checks if a struct represents graph data
func isGraph(structValue *structpb.Struct) bool {
	// A struct is considered a graph if it has both nodes and edges fields
	_, hasNodes := structValue.Fields["nodes"]
	_, hasEdges := structValue.Fields["edges"]
	return hasNodes && hasEdges
}

// isScalar checks if a struct represents scalar data
func isScalar(structValue *structpb.Struct) bool {
	// A struct is considered scalar if it has exactly one field with a scalar value
	if len(structValue.Fields) != 1 {
		return false
	}

	// Get the single value
	var value *structpb.Value
	for _, v := range structValue.Fields {
		value = v
		break
	}

	// Check if the value is scalar
	switch value.GetKind().(type) {
	case *structpb.Value_NumberValue, *structpb.Value_StringValue, *structpb.Value_BoolValue:
		return true
	default:
		return false
	}
}
