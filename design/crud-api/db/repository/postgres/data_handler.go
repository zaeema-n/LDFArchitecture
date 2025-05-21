package postgres

import (
	"fmt"
	"log"

	pb "lk/datafoundation/crud-api/lk/datafoundation/crud-api"
	"lk/datafoundation/crud-api/pkg/schema"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// UnmarshalAnyToString attempts to unmarshal an Any protobuf message to a string value
func UnmarshalAnyToString(anyValue *anypb.Any) (string, error) {
	if anyValue == nil {
		return "", nil
	}

	var stringValue wrapperspb.StringValue
	if err := anyValue.UnmarshalTo(&stringValue); err != nil {
		return "", fmt.Errorf("error unmarshaling to string value: %v", err)
	}
	return stringValue.Value, nil
}

// UnmarshalTimeBasedValueList attempts to unmarshal a TimeBasedValueList from an Any protobuf message
func UnmarshalTimeBasedValueList(anyValue *anypb.Any) ([]interface{}, error) {
	if anyValue == nil {
		return nil, nil
	}

	var timeBasedValueList pb.TimeBasedValueList
	if err := anyValue.UnmarshalTo(&timeBasedValueList); err != nil {
		return nil, fmt.Errorf("error unmarshaling to TimeBasedValueList: %v", err)
	}

	// Convert TimeBasedValueList to []interface{}
	result := make([]interface{}, len(timeBasedValueList.Values))
	for i, v := range timeBasedValueList.Values {
		result[i] = v
	}
	return result, nil
}

// UnmarshalEntityAttributes unmarshals the attributes map from a protobuf Entity
func UnmarshalEntityAttributes(attributes map[string]*anypb.Any) (map[string]interface{}, error) {
	if attributes == nil {
		return nil, nil
	}

	result := make(map[string]interface{})
	for key, value := range attributes {
		if value == nil {
			continue
		}

		// Try to unmarshal as string first
		if strValue, err := UnmarshalAnyToString(value); err == nil {
			result[key] = strValue
			continue
		}

		// Try to unmarshal as TimeBasedValueList
		if timeBasedValues, err := UnmarshalTimeBasedValueList(value); err == nil {
			result[key] = timeBasedValues
			continue
		}

		log.Printf("Warning: Could not unmarshal attribute %s with type %s", key, value.TypeUrl)
	}

	return result, nil
}

func HandleAttributes(attributes map[string]*pb.TimeBasedValueList) (map[string]interface{}, error) {
	log.Printf("--------------Handling Attributes------------------")
	log.Printf("Handling attributes: %v", attributes)

	if attributes == nil {
		return nil, nil
	}

	// Print each attribute's key and values
	for key, value := range attributes {
		if value != nil {
			log.Printf("Attribute - Key: %s, Values: %v", key, value.Values)
		}
		log.Printf("Attribute - Key: %s", key)
		log.Printf("Attribute - Values: %v", value.Values)
	}

	result := make(map[string]interface{})
	for key, value := range attributes {
		if value != nil {
			values := value.Values
			for i, val := range values {
				log.Printf("Processing value %d: %v", i, val)
				if val != nil {
					log.Printf("Value details - StartTime: %s, EndTime: %s, Value: %v",
						val.GetStartTime(),
						val.GetEndTime(),
						val.GetValue())

					// Generate schema directly from the Any value
					schemaGenerator := schema.NewSchemaGenerator()
					schemaInfo, err := schemaGenerator.GenerateSchema(val.GetValue())
					if err != nil {
						log.Printf("Failed to generate schema: %v", err)
						continue
					}

					// Log the schema information
					log.Printf("Schema for value %d: StorageType=%v, TypeInfo=%v",
						i,
						schemaInfo.StorageType,
						schemaInfo.TypeInfo)
				}
			}
			result[key] = values
		}
	}
	return result, nil
}
