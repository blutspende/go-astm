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

func ParseLine(inputLine string, targetStruct interface{}, lineTypeName string, sequenceNumber int, config *models.Configuration) (nameOk bool, err error) {
	// Check for input line length
	if len(inputLine) == 0 {
		return false, errmsg.LineParsing_ErrEmptyInput
	}
	// Handle header special case
	if inputLine[0] == 'H' {
		// Check if the inputLine is long enough to contain delimiters
		if len(inputLine) < 5 {
			return false, errmsg.LineParsing_ErrHeaderTooShort
		}
		// Override delimiters
		config.Internal.Delimiters.Field = string(inputLine[1])
		config.Internal.Delimiters.Repeat = string(inputLine[2])
		config.Internal.Delimiters.Component = string(inputLine[3])
		config.Internal.Delimiters.Escape = string(inputLine[4])
	}

	// Split the input with the field delimiter
	inputFields := strings.Split(inputLine, config.Internal.Delimiters.Field)

	// Check for validity with parent data
	if len(inputFields) < 2 {
		return false, errmsg.LineParsing_ErrMandatoryInputFieldsMissing
	}
	nameOk = inputFields[0] == lineTypeName
	// Name checking is never enforced
	//if !nameOk && config.EnforceRecordNameCheck {
	//	return nameOk, errmsg.LineParsing_ErrLineTypeNameMismatch
	//}
	if inputFields[1] != strconv.Itoa(sequenceNumber) && inputLine[0] != 'H' && config.EnforceSequenceNumberCheck {
		return nameOk, errmsg.LineParsing_ErrSequenceNumberMismatch
	}

	// Process the target structure
	targetTypes, targetValues, _, err := ProcessStructReflection(targetStruct)
	if err != nil {
		return nameOk, err
	}

	// Iterate over the inputFields of the targetStruct struct
	for i, targetType := range targetTypes {
		// Parse the targetStruct field targetFieldAnnotation
		targetFieldAnnotation, err := ParseAstmFieldAnnotation(targetType)
		if err != nil {
			return nameOk, err
		}

		// Not enough inputFields in the input inputLine
		if len(inputFields) < targetFieldAnnotation.FieldPos {
			// If the field is required it's an error, otherwise skip it
			if targetFieldAnnotation.Attribute == constants.ATTRIBUTE_REQUIRED {
				return nameOk, errmsg.LineParsing_ErrInputFieldsMissing
			} else {
				continue
			}
		}
		// Save the current inputField
		inputField := inputFields[targetFieldAnnotation.FieldPos-1]

		//Check if there is any data
		if inputField == "" {
			if targetFieldAnnotation.Attribute == constants.ATTRIBUTE_REQUIRED {
				return nameOk, errmsg.LineParsing_ErrRequiredFieldIsEmpty
			} else {
				// Non required field can be skipped
				continue
			}
		}

		// |rep1\rep2\rep3|
		// Field is an array
		if targetFieldAnnotation.IsArray {
			repeats := strings.Split(inputField, config.Internal.Delimiters.Repeat)
			arrayType := reflect.SliceOf(targetValues[i].Type().Elem())
			arrayValue := reflect.MakeSlice(arrayType, len(repeats), len(repeats))
			for j, repeat := range repeats {
				err = setField(arrayValue.Index(j), repeat, config)
				if err != nil {
					return nameOk, err
				}
			}
			targetValues[i].Set(arrayValue)
			// |comp1^comp2^comp3|
			// Field is a component
		} else if targetFieldAnnotation.IsComponent {
			components := strings.Split(inputField, config.Internal.Delimiters.Component)
			// Not enough components in the inputField
			if len(components) < targetFieldAnnotation.ComponentPos {
				return nameOk, errmsg.LineParsing_ErrInputComponentsMissing
			}
			err = setField(targetValues[i], components[targetFieldAnnotation.ComponentPos-1], config)
			if err != nil {
				return nameOk, err
			}
			// Field is not an array or component (normal singular field)
		} else {
			err = setField(targetValues[i], inputField, config)
			if err != nil {
				return nameOk, err
			}
		}
		//TODO: handle componented array case
		// |comp1^comp2^comp3\comp1^comp2^comp3\comp1^comp2^comp3|

		// Check if there are more inputFields in the input not mapped to the struct
		//if i == targetFieldCount-1 && len(inputFields) > targetFieldAnnotation.FieldPos {
		// Note: this could be a warning about lost data
		//}
	}
	// Return nil if everything went well
	return nameOk, nil
}

func setField(field reflect.Value, value string, config *models.Configuration) (err error) {
	// Ensure the field is settable
	if !field.CanSet() {
		// Field is not settable
		return errmsg.LineParsing_ErrNonSettableField
	}

	// Set the field value
	switch field.Kind() {
	case reflect.String:
		field.Set(reflect.ValueOf(value))
	case reflect.Int:
		num, err := strconv.Atoi(value)
		if err != nil {
			return errmsg.LineParsing_ErrDataParsingError
		}
		field.Set(reflect.ValueOf(num))
	case reflect.Float32:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return errmsg.LineParsing_ErrDataParsingError
		}
		field.Set(reflect.ValueOf(float32(num)))
	case reflect.Float64:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errmsg.LineParsing_ErrDataParsingError
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
				return errmsg.LineParsing_ErrInvalidDateFormat
			}
			timeInLocation, err := time.ParseInLocation(timeFormat, value, config.Internal.TimeLocation)
			if err != nil {
				return errmsg.LineParsing_ErrDataParsingError
			}
			field.Set(reflect.ValueOf(timeInLocation))
		} else {
			// Note: option to handle other struct types here
		}
	default:
		return errmsg.LineParsing_ErrUsupportedDataType
	}
	// Return nil if everything went well
	return nil
}
