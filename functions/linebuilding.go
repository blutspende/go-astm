package functions

import (
	"errors"
	"github.com/blutspende/go-astm/v3/constants/astmconst"
	"github.com/blutspende/go-astm/v3/errmsg"
	"github.com/blutspende/go-astm/v3/models"
	"github.com/blutspende/go-astm/v3/models/astmmodels"
	"math"
	"reflect"
	"strconv"
	"time"
)

func BuildLine(sourceStruct interface{}, lineTypeName string, sequenceNumber int, config *astmmodels.Configuration) (result string, err error) {
	// Process the target structure
	sourceTypes, sourceValues, sourceTypesLength, err := ProcessStructReflection(sourceStruct)
	if err != nil {
		return "", err
	}

	// Create a map to store field values indexed by FieldPos
	fieldMap := make(map[int]string)

	// Create an array to store already processed component fields to allow any sparse placement
	processedComponentFields := make([]int, 0)

	// Add line name
	fieldMap[1] = lineTypeName
	// If it's a header, add the other delimiters
	if lineTypeName == "H" {
		fieldMap[2] = config.Delimiters.Repeat +
			config.Delimiters.Component +
			config.Delimiters.Escape
	} else {
		// If it's not a header add the sequence number
		fieldMap[2] = strconv.Itoa(sequenceNumber)
	}

	// Iterate over the inputFields of the targetStruct struct
	for i := 0; i < sourceTypesLength; i++ {
		// Parse the sourceStruct field sourceFieldAnnotation
		sourceFieldAnnotation, err := ParseAstmFieldAnnotation(sourceTypes[i])
		if err != nil {
			if errors.Is(err, errmsg.AnnotationParsing_ErrMissingAstmAnnotation) {
				// If the annotation is missing, skip this field
				continue
			} else {
				return "", err
			}
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
					fieldValueString += config.Delimiters.Repeat
				}
			}
		} else if sourceFieldAnnotation.IsComponent {
			if contains(processedComponentFields, sourceFieldAnnotation.FieldPos) {
				// If the field is already processed, skip it
				continue
			}
			// Create a map to store the component values indexed by ComponentPos
			componentMap := make(map[int]string)
			// Iterate over the whole inputFields of the targetStruct struct to find the components anywhere
			for j := 0; j < sourceTypesLength; j++ {
				// Parse the targetStruct field targetFieldAnnotation
				currentFieldAnnotation, err := ParseAstmFieldAnnotation(sourceTypes[j])
				if err != nil {
					if errors.Is(err, errmsg.AnnotationParsing_ErrMissingAstmAnnotation) {
						// If the annotation is missing, skip this field
						continue
					} else {
						return "", err
					}
				}
				// If the field number is the same as the sourceFieldAnnotation, process it
				if currentFieldAnnotation.FieldPos == sourceFieldAnnotation.FieldPos {
					// Convert current component
					componentValue, err := convertField(sourceValues[j], currentFieldAnnotation, config)
					if err != nil {
						return "", err
					}
					// Store the value in the component map
					componentMap[currentFieldAnnotation.ComponentPos] = componentValue
				}
			}
			// Construct the result into the fieldValueString
			fieldValueString = constructResult(componentMap, config.Delimiters.Component, config.Notation)
			// Mark the field as processed
			processedComponentFields = append(processedComponentFields, sourceFieldAnnotation.FieldPos)
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
	result = constructResult(fieldMap, config.Delimiters.Field, config.Notation)

	return result, nil
}

func buildSubstructure(sourceStruct interface{}, config *astmmodels.Configuration) (result string, err error) {
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
			if errors.Is(err, errmsg.AnnotationParsing_ErrMissingAstmAnnotation) {
				// If the annotation is missing, skip this field
				continue
			} else {
				return "", err
			}
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
	result = constructResult(componentMap, config.Delimiters.Component, config.Notation)

	// Return result with no error
	return result, nil
}

func constructResult(fieldMap map[int]string, delimiter string, notation string) (result string) {
	// Determine how many fields to include by finding the biggest index
	lastIndex := 0
	for key := range fieldMap {
		// In short notation only non-empty fields are included at the end
		if key > lastIndex && (!(notation == astmconst.NOTATION_SHORT) || fieldMap[key] != "") {
			lastIndex = key
		}
	}
	// Iterate from one to the last index, building the result string
	result = ""
	for i := 1; i <= lastIndex; i++ {
		// If the field exists in the map, append its value to the result
		if value, exists := fieldMap[i]; exists {
			result += value
		}
		// Add the delimiter after all but the last element
		if i != lastIndex {
			result += delimiter
		}
	}
	// Return the built string
	return result
}

func convertField(field reflect.Value, annotation models.AstmFieldAnnotation, config *astmmodels.Configuration) (result string, err error) {
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
		precision := config.DefaultDecimalPrecision
		if value, exists := annotation.Attributes[astmconst.ATTRIBUTE_LENGTH]; exists {
			precision, err = strconv.Atoi(value)
			if err != nil {
				return "", errmsg.LineBuilding_ErrInvalidLengthAttributeValue
			}
		}
		result = strconv.FormatFloat(field.Float(), 'f', precision, field.Type().Bits())
		if !config.RoundLastDecimal && precision >= 0 {
			factor := math.Pow(10, float64(precision))
			truncated := math.Trunc(field.Float()*factor) / factor
			result = strconv.FormatFloat(truncated, 'f', precision, field.Type().Bits())
		}
		return result, nil
	case reflect.Struct:
		// Check for time.Time type (it reflects as a Struct)
		if field.Type() == reflect.TypeOf(time.Time{}) {
			timeFormat := "20060102"
			if _, exists := annotation.Attributes[astmconst.ATTRIBUTE_LONGDATE]; exists {
				timeFormat = "20060102150405"
			}
			// Check if the field is a time.Time
			timeValue, ok := field.Interface().(time.Time)
			if !ok {
				return "", errmsg.LineBuilding_ErrInvalidDateFormat
			}
			// Return empty if the time is zero
			if timeValue.IsZero() {
				return "", nil
			}
			// Convert the time to the config's timezone
			timeInLocation := timeValue.In(config.TimeLocation)
			// Format the date as a string
			result = timeInLocation.Format(timeFormat)
			return result, nil
		} else {
			// Note: option to handle other struct types here
		}
	}
	// Return error if no type match was found (each successful conversion returns with nil)
	return "", errmsg.LineBuilding_ErrUsupportedDataType
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
