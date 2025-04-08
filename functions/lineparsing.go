package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/models"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// TODO: Add error handling
// TODO: figue out how to better pass parent data
func ParseLine(inputLine string, targetStruct interface{}, lineName string, squenceNum int, config models.Configuration) {
	// Split the input with the field delimiter
	inputFields := strings.Split(inputLine, config.Delimiters.Field)

	// Check for validity with parent data
	if len(inputFields) < 2 {
		return
	}
	if inputFields[0] != lineName || inputFields[1] != strconv.Itoa(squenceNum) {
		return
	}

	// Ensure the targetStruct is a pointer to a struct
	targetPtrValue := reflect.ValueOf(targetStruct)
	if targetPtrValue.Kind() != reflect.Ptr || targetPtrValue.Elem().Kind() != reflect.Struct {
		// targetStruct must be a pointer to a struct
		return
	}
	// Get the underlying struct
	targetValue := targetPtrValue.Elem()
	targetType := targetPtrValue.Type()

	// Iterate over the inputFields of the targetStruct struct
	for i := 0; i < targetType.NumField(); i++ {
		// Parse the targetStruct field targetFieldAnnotation
		targetFieldAnnotation, _ := ParseAstmFieldAnnotation(targetType.Field(i))

		// Not enough inputFields in the input inputLine
		if len(inputFields) < targetFieldAnnotation.FieldPos {
			return
		}
		// Save the current inputField
		inputField := inputFields[targetFieldAnnotation.FieldPos-1]

		//Check if there is any data
		if inputField == "" {
			if targetFieldAnnotation.Attribute == constants.ATTRIBUTE_REQUIRE {
				// Error: required field is empty
				return
			} else {
				// Non required field can be skipped
				continue
			}
		}

		// |rep1\rep2\rep3|
		// Field is an array
		if targetFieldAnnotation.IsArray {
			repeats := strings.Split(inputField, config.Delimiters.Repeat)
			for repeat := range repeats {
				//TODO: handle array setting
				_ = repeat
			}
			// |comp1^comp2^comp3|
			// Field is a component
		} else if targetFieldAnnotation.IsComponent {
			components := strings.Split(inputField, config.Delimiters.Component)
			// Not enough components in the inputField
			if len(components) < targetFieldAnnotation.ComponentPos {
				return
			}
			setField(targetValue.Field(i), components[targetFieldAnnotation.ComponentPos-1], config)
			// Field is not an array or component (normal singular field)
		} else {
			setField(targetValue.Field(i), inputField, config)
		}
		//TODO: handle componented array case
		// |comp1^comp2^comp3\comp1^comp2^comp3\comp1^comp2^comp3|

		// Check if there are more inputFields in the input not mapped to the struct
		if i == targetType.NumField()-1 && len(inputFields) > targetFieldAnnotation.FieldPos {
			// TODO: this could be a warning about lost data
			//return
		}
	}
}

func setField(field reflect.Value, value string, config models.Configuration) {
	// Ensure the field is settable
	if !field.CanSet() {
		// Field is not settable
		return
	}

	// Set the field value
	switch field.Kind() {
	case reflect.String:
		field.Set(reflect.ValueOf(value))
	case reflect.Int:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			//TODO: handle error
			return
		}
		field.Set(reflect.ValueOf(num))
	case reflect.Float32:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			//TODO: handle error
			return
		}
		field.Set(reflect.ValueOf(num))
	case reflect.Float64:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			//TODO: handle error
			return
		}
		field.Set(reflect.ValueOf(num))
	// Check for time.Time type (it reflects as a Struct)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			timeFormat := ""
			switch len(value) {
			case 8:
				timeFormat = "20060102" // YYYYMMDD
			case 14:
				timeFormat = "20060102150405" // YYYYMMDDHHMMSS
			default:
				//TODO: handle error
				return
			}
			timeInLocation, err := time.ParseInLocation(timeFormat, value, config.TimeLocation)
			if err != nil {
				// TODO: handle parsing error
				return
			}
			field.Set(reflect.ValueOf(timeInLocation))
		} else {
			// Option to handle other struct types here
		}
	default:
		//TODO: handle other types
	}
}
