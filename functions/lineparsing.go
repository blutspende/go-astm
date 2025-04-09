package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/blutspende/go-astm/v2/models"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// TODO: Add error handling
// TODO: figue out how to better pass parent data
func ParseLine(inputLine string, targetStruct interface{}, lineTypeName string, sequenceNumber int, config models.Configuration) (err error) {
	// Check for input line length
	if len(inputLine) == 0 {
		return
	}
	// Handle header special case
	if inputLine[0] == 'H' {
		// Check if the inputLine is long enough to contain delimiters
		if len(inputLine) < 5 {
			return
		}
		config.Delimiters.Field = string(inputLine[1])
		config.Delimiters.Repeat = string(inputLine[2])
		config.Delimiters.Component = string(inputLine[3])
		config.Delimiters.Escape = string(inputLine[4])
	}

	// Split the input with the field delimiter
	inputFields := strings.Split(inputLine, config.Delimiters.Field)

	// Check for validity with parent data
	if len(inputFields) < 2 {
		return
	}
	if inputFields[0] != lineTypeName {
		return errmsg.LineParsing_ErrLineTypeNameMismatch
	}
	if inputFields[1] != strconv.Itoa(sequenceNumber) {
		return errmsg.LineParsing_ErrSequenceNumberMismatch
	}

	// Process the target structure
	targetTypes, targetValues, targetFieldCount, _ := ProcessStructReflection(targetStruct)

	// Iterate over the inputFields of the targetStruct struct
	for i, targetType := range targetTypes {
		// Parse the targetStruct field targetFieldAnnotation
		targetFieldAnnotation, _ := ParseAstmFieldAnnotation(targetType)

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
			setField(targetValues[i], components[targetFieldAnnotation.ComponentPos-1], config)
			// Field is not an array or component (normal singular field)
		} else {
			setField(targetValues[i], inputField, config)
		}
		//TODO: handle componented array case
		// |comp1^comp2^comp3\comp1^comp2^comp3\comp1^comp2^comp3|

		// Check if there are more inputFields in the input not mapped to the struct
		if i == targetFieldCount-1 && len(inputFields) > targetFieldAnnotation.FieldPos {
			// TODO: this could be a warning about lost data
			//return
		}
	}
	// Return nil if everything went well
	return nil
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
