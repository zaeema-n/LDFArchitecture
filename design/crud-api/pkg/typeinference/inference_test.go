package typeinference

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
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	structValue, err := structpb.NewStruct(data.(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("failed to create struct: %v", err)
	}

	anyValue, err := anypb.New(structValue)
	if err != nil {
		return nil, fmt.Errorf("failed to create Any: %v", err)
	}

	return anyValue, nil
}

// TestScalarTypes tests type inference for scalar data types
func TestScalarTypes(t *testing.T) {
	testCases := map[string]struct {
		json     string
		expected DataType
	}{
		"empty_string": {
			json:     `{"attributes": ""}`,
			expected: StringType,
		},
		"zero_int": {
			json:     `{"attributes": 0}`,
			expected: IntType,
		},
		"zero_float": {
			json:     `{"attributes": 0.1}`,
			expected: FloatType,
		},
		"integer": {
			json:     `{"attributes": 42}`,
			expected: IntType,
		},
		"float": {
			json:     `{"attributes": 3.14}`,
			expected: FloatType,
		},
		"string": {
			json:     `{"attributes": "hello"}`,
			expected: StringType,
		},
		"boolean": {
			json:     `{"attributes": true}`,
			expected: BoolType,
		},
		"null": {
			json:     `{"attributes": null}`,
			expected: NullType,
		},
		"missing_attributes": {
			json:     `{}`,
			expected: NullType,
		},
	}

	inferrer := &TypeInferrer{}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			anyValue, err := JSONToAny(tc.json)
			assert.NoError(t, err)

			typeInfo, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, typeInfo.Type)
		})
	}
}

// TestListTypes tests type inference for list data types
func TestListTypes(t *testing.T) {
	testCases := map[string]struct {
		json     string
		expected *TypeInfo
	}{
		"empty_list": {
			json: `{"attributes": []}`,
			expected: &TypeInfo{
				Type:    StringType,
				IsArray: true,
			},
		},
		"list_with_empty_values": {
			json: `{"attributes": ["", 0, 0.0, null]}`,
			expected: &TypeInfo{
				Type:    StringType,
				IsArray: true,
				ArrayType: &TypeInfo{
					Type: StringType,
				},
			},
		},
		"list_of_integers": {
			json: `{"attributes": [1, 2, 3]}`,
			expected: &TypeInfo{
				Type:    StringType,
				IsArray: true,
				ArrayType: &TypeInfo{
					Type: IntType,
				},
			},
		},
		"list_of_strings": {
			json: `{"attributes": ["a", "b", "c"]}`,
			expected: &TypeInfo{
				Type:    StringType,
				IsArray: true,
				ArrayType: &TypeInfo{
					Type: StringType,
				},
			},
		},
		"list_with_null": {
			json: `{"attributes": null}`,
			expected: &TypeInfo{
				Type:       NullType,
				IsNullable: true,
			},
		},
	}

	inferrer := &TypeInferrer{}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			anyValue, err := JSONToAny(tc.json)
			assert.NoError(t, err)

			typeInfo, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected.Type, typeInfo.Type)
			assert.Equal(t, tc.expected.IsArray, typeInfo.IsArray)
			if tc.expected.ArrayType != nil {
				assert.NotNil(t, typeInfo.ArrayType)
				assert.Equal(t, tc.expected.ArrayType.Type, typeInfo.ArrayType.Type)
			}
		})
	}
}

// TestMapTypes tests type inference for map data types
func TestMapTypes(t *testing.T) {
	testCases := map[string]struct {
		json     string
		expected *TypeInfo
	}{
		"empty_map": {
			json: `{"attributes": {}}`,
			expected: &TypeInfo{
				Type: StringType,
			},
		},
		"map_with_empty_values": {
			json: `{"attributes": {"empty_str": "", "zero": 0, "null_val": null}}`,
			expected: &TypeInfo{
				Type: StringType,
			},
		},
		"simple_map": {
			json: `{"attributes": {"key": "value"}}`,
			expected: &TypeInfo{
				Type: StringType,
			},
		},
		"map_with_null": {
			json: `{"attributes": null}`,
			expected: &TypeInfo{
				Type:       NullType,
				IsNullable: true,
			},
		},
	}

	inferrer := &TypeInferrer{}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			anyValue, err := JSONToAny(tc.json)
			assert.NoError(t, err)

			typeInfo, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected.Type, typeInfo.Type)
		})
	}
}

// TestSpecialTypes tests type inference for special data types
func TestSpecialTypes(t *testing.T) {
	testCases := map[string]struct {
		json     string
		expected DataType
	}{
		"date": {
			json:     `{"attributes": "2024-03-20"}`,
			expected: DateType,
		},
		"time": {
			json:     `{"attributes": "14:30:00"}`,
			expected: TimeType,
		},
		"datetime": {
			json:     `{"attributes": "2024-03-20T14:30:00Z"}`,
			expected: DateTimeType,
		},
	}

	inferrer := &TypeInferrer{}
	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			anyValue, err := JSONToAny(tc.json)
			assert.NoError(t, err)

			typeInfo, err := inferrer.InferType(anyValue)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, typeInfo.Type)
		})
	}
}
