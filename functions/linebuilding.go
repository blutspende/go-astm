package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/blutspende/go-astm/v2/models"
	"math"
	"reflect"
	"sort"
	"strconv"
	"time"
)

func BuildLine(sourceStruct interface{}, lineTypeName string, sequenceNumber int, config *models.Configuration) (result string, err error) {
	// Process the target structure
	sourceTypes, sourceValues, sourceTypesLength, err := ProcessStructReflection(sourceStruct)
	if err != nil {
		return "", err
	}

	// Create a map to store field values indexed by FieldPos
	fieldMap := make(map[int]string)

	// Add line name
	fieldMap[1] = lineTypeName
	// If it's a header, add the other delimiters
	if lineTypeName == "H" {
		fieldMap[2] = config.Internal.Delimiters.Repeat +
			config.Internal.Delimiters.Component +
			config.Internal.Delimiters.Escape
	} else {
		// If it's not a header add the sequence number
		fieldMap[2] = strconv.Itoa(sequenceNumber)
	}

	// Iterate over the inputFields of the targetStruct struct
	for i := 0; i < sourceTypesLength; i++ {
		// Parse the sourceStruct field sourceFieldAnnotation
		sourceFieldAnnotation, err := ParseAstmFieldAnnotation(sourceTypes[i])
		if err != nil {
			return "", err
		}

		// Check for fieldPos not being lower than 3 (first 2 are reserved for line name and sequence number)
		if sourceFieldAnnotation.FieldPos < 3 {
			return "", errmsg.LineBuilding_ErrReservedFieldPosReference
		}

		fieldValueString := ""
		// If the field is an array, iterate over its elements and use the Repeat delimiter
		if sourceFieldAnnotation.IsArray {
			for j := 0; j < sourceValues[i].Len(); j++ {
				elementValue := sourceValues[i].Index(j)
				convertedValue := ""
				if sourceFieldAnnotation.IsSubstructure {
					// If the field is a substructure use buildSubstructure to process it
					convertedValue, err = buildSubstructure(elementValue.Interface(), config)
					if err != nil {
						return "", err
					}
				} else {
					// Simple field, convert it directly
					convertedValue, err = convertField(elementValue, sourceFieldAnnotation, config)
					if err != nil {
						return "", err
					}
				}
				fieldValueString += convertedValue
				if j < sourceValues[i].Len()-1 {
					fieldValueString += config.Internal.Delimiters.Repeat
				}
			}
		} else if sourceFieldAnnotation.IsComponent {
			// If the field is a component, iterate over sourceTypes until a field is not a component
			// Note: components for the same field have to come sequentially, or it will break
			componentFieldString := ""
			for ; i < len(sourceTypes); i++ {
				// Parse the targetStruct field targetFieldAnnotation
				currentFieldAnnotation, err := ParseAstmFieldAnnotation(sourceTypes[i])
				if err != nil {
					return "", err
				}
				// If the field is not the same field anymore, break the loop
				if currentFieldAnnotation.FieldPos != sourceFieldAnnotation.FieldPos {
					i--
					break
				}

				// Convert current component
				componentValue, err := convertField(sourceValues[i], currentFieldAnnotation, config)
				if err != nil {
					return "", err
				}

				// Add the component value and a component delimiter to the field string
				componentFieldString += componentValue + config.Internal.Delimiters.Component
			}
			// Remove the last component delimiter
			if len(componentFieldString) > 0 {
				componentFieldString = componentFieldString[:len(componentFieldString)-1]
			}
			// Set the field value string to the component field string
			fieldValueString = componentFieldString
		} else if sourceFieldAnnotation.IsSubstructure {
			// If the field is a substructure use buildSubstructure to process it
			fieldValueString, err = buildSubstructure(sourceValues[i].Interface(), config)
			if err != nil {
				return "", err
			}
		} else {
			// If the field is not an array, convert it directly
			fieldValueString, err = convertField(sourceValues[i], sourceFieldAnnotation, config)
			if err != nil {
				return "", err
			}
		}

		// Store the field value in the map using FieldPos as the key
		fieldMap[sourceFieldAnnotation.FieldPos] = fieldValueString
	}

	// Construct the result string based on the field map
	result = constructResult(fieldMap, config.Internal.Delimiters.Field, config.Notation)

	return result, nil
}

func buildSubstructure(sourceStruct interface{}, config *models.Configuration) (result string, err error) {
	// Process the target structure
	sourceTypes, sourceValues, sourceTypesLength, err := ProcessStructReflection(sourceStruct)
	if err != nil {
		return "", err
	}

	// Create a map to store component values indexed by FieldPos
	componentMap := make(map[int]string)

	// Iterate over the inputFields of the targetStruct struct
	for i := 0; i < sourceTypesLength; i++ {
		// Parse the sourceStruct field sourceFieldAnnotation
		sourceFieldAnnotation, err := ParseAstmFieldAnnotation(sourceTypes[i])
		if err != nil {
			return "", err
		}
		// Convert the component directly
		componentValueString, err := convertField(sourceValues[i], sourceFieldAnnotation, config)
		if err != nil {
			return "", err
		}
		// Store the component value in the map using FieldPos as the key
		componentMap[sourceFieldAnnotation.FieldPos] = componentValueString
	}

	// Construct the result string
	result = constructResult(componentMap, config.Internal.Delimiters.Component, constants.NOTATION_STANDARD)

	// Return result with no error
	return result, nil
}

func constructResult(fieldMap map[int]string, delimiter string, notation string) (result string) {
	// Sort the keys of the map
	keys := make([]int, 0, len(fieldMap))
	for k := range fieldMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Determine how many fields to include based on the notation
	lastElementIndex := len(fieldMap) - 1
	// In short notation the empty fields in the end are skipped
	if notation == constants.NOTATION_SHORT {
		for i, key := range keys {
			if fieldMap[key] != "" {
				lastElementIndex = i
			}
		}
	}

	// Construct the result string
	for i, key := range keys {
		result += fieldMap[key]
		// Add the field delimiter if not the last field
		if i < lastElementIndex {
			result += delimiter
		}
		// Break when we reach the targeted last element
		if i == lastElementIndex {
			break
		}
	}

	return result
}

func convertField(field reflect.Value, annotation models.AstmFieldAnnotation, config *models.Configuration) (result string, err error) {
	// Check if the field is a pointer, nil returns empty, otherwise dereference it
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return "", nil
		}
		field = field.Elem()
	}
	// Format the result as a string based on the field type
	switch field.Kind() {
	case reflect.String:
		if field.Type().ConvertibleTo(reflect.TypeOf("")) {
			result = field.String()
		} else {
			return "", errmsg.LineBuilding_ErrUsupportedDataType
		}
		return result, nil
	case reflect.Int:
		result = strconv.Itoa(int(field.Int()))
		return result, nil
	case reflect.Float32, reflect.Float64:
		precision := -1
		if annotation.Attribute == constants.ATTRIBUTE_LENGTH {
			precision = annotation.AttributeValue
		}
		result = strconv.FormatFloat(field.Float(), 'f', precision, field.Type().Bits())
		if !config.RoundFixedNumbers && precision >= 0 {
			factor := math.Pow(10, float64(precision))
			truncated := math.Trunc(field.Float()*factor) / factor
			result = strconv.FormatFloat(truncated, 'f', precision, field.Type().Bits())
		}
		return result, nil
	case reflect.Struct:
		// Check for time.Time type (it reflects as a Struct)
		if field.Type() == reflect.TypeOf(time.Time{}) {
			timeFormat := "20060102"
			if annotation.Attribute == constants.ATTRIBUTE_LONGDATE {
				timeFormat = "20060102150405"
			}
			timeValue, ok := field.Interface().(time.Time)
			if !ok {
				return "", errmsg.LineBuilding_ErrInvalidDateFormat
			}
			if timeValue.IsZero() {
				result = ""
			} else {
				result = timeValue.In(config.Internal.TimeLocation).Format(timeFormat) // Format the date as a string
			}
			return result, nil
		} else {
			// Note: option to handle other struct types here
		}
	}
	// Return error if no type match was found (each successful conversion returns with nil)
	return "", errmsg.LineBuilding_ErrUsupportedDataType
}
