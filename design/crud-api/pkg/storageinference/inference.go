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

// InferType attempts to determine the data type from a protobuf Any value
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
				// Check if it's a scalar structure (has "value" field)
				if value, ok := attrStruct.Fields["value"]; ok {
					switch value.GetKind().(type) {
					case *structpb.Value_NumberValue, *structpb.Value_StringValue, *structpb.Value_BoolValue:
						return ScalarData, nil
					}
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
