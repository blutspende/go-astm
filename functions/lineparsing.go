package functions

import (
	"github.com/blutspende/go-astm/v2/constants/astmconst"
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
	// Name checking is always enforced, but instead of error it is returned in the nameOk variable
	if !nameOk {
		return nameOk, nil
	}
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

		// Check for fieldPos not being lower than 3 (first 2 are reserved for line name and sequence number)
		if targetFieldAnnotation.FieldPos < 3 {
			return nameOk, errmsg.LineParsing_ErrReservedFieldPosReference
		}

		// Not enough inputFields or empty inputField
		if len(inputFields) < targetFieldAnnotation.FieldPos || inputFields[targetFieldAnnotation.FieldPos-1] == "" {
			// If the field is required it's an error, otherwise skip it
			if targetFieldAnnotation.Attribute == astmconst.ATTRIBUTE_REQUIRED {
				return nameOk, errmsg.LineParsing_ErrRequiredInputFieldMissing
			} else {
				continue
			}
		}
		// Save the current inputField
		inputField := inputFields[targetFieldAnnotation.FieldPos-1]

		if targetFieldAnnotation.IsArray {
			// |rep1\rep2\rep3|
			// Field is an array
			repeats := strings.Split(inputField, config.Internal.Delimiters.Repeat)
			arrayType := reflect.SliceOf(targetValues[i].Type().Elem())
			arrayValue := reflect.MakeSlice(arrayType, len(repeats), len(repeats))
			for j, repeat := range repeats {
				if targetFieldAnnotation.IsSubstructure {
					// |comp1^comp2^comp3\comp1^comp2^comp3\comp1^comp2^comp3|
					// Substructures (with components) in the array: use parseSubstructure
					err = parseSubstructure(repeat, arrayValue.Index(j).Addr().Interface(), config)
					if err != nil {
						return nameOk, err
					}
				} else {
					// |value1\value2\value3|
					// Simple values in the array
					err = setField(repeat, arrayValue.Index(j), config)
					if err != nil {
						return nameOk, err
					}
				}

			}
			targetValues[i].Set(arrayValue)
		} else if targetFieldAnnotation.IsComponent {
			// |comp1^comp2^comp3|
			// Field is a component
			components := strings.Split(inputField, config.Internal.Delimiters.Component)
			// Not enough components in the inputField
			if len(components) < targetFieldAnnotation.ComponentPos {
				// Error if the component is required, skip otherwise
				if targetFieldAnnotation.Attribute == astmconst.ATTRIBUTE_REQUIRED {
					return nameOk, errmsg.LineParsing_ErrInputComponentsMissing
				} else {
					continue
				}
			}
			err = setField(components[targetFieldAnnotation.ComponentPos-1], targetValues[i], config)
			if err != nil {
				return nameOk, err
			}
		} else if targetFieldAnnotation.IsSubstructure {
			// |comp1^comp2^comp3|
			// If the field is a substructure use parseSubstructure to process it
			err = parseSubstructure(inputField, targetValues[i].Addr().Interface(), config)
			if err != nil {
				return nameOk, err
			}
		} else {
			// |field|
			// Field is not an array or component (normal singular field)
			err = setField(inputField, targetValues[i], config)
			if err != nil {
				return nameOk, err
			}
		}
		// Note: this could be a place to produce warnings about lost data
		// if i == targetFieldCount-1 && len(inputFields) > targetFieldAnnotation.FieldPos
	}
	// Return no error if everything went well
	return nameOk, nil
}

func parseSubstructure(inputString string, targetStruct interface{}, config *models.Configuration) (err error) {
	// Split the input with the field delimiter
	inputFields := strings.Split(inputString, config.Internal.Delimiters.Component)

	// Process the target structure
	targetTypes, targetValues, _, err := ProcessStructReflection(targetStruct)
	if err != nil {
		return err
	}

	// Iterate over the inputFields of the targetStruct struct
	for i, targetType := range targetTypes {
		// Parse the targetStruct field targetFieldAnnotation
		targetFieldAnnotation, err := ParseAstmFieldAnnotation(targetType)
		if err != nil {
			return err
		}

		// Not enough inputFields or empty inputField
		if len(inputFields) < targetFieldAnnotation.FieldPos || inputFields[targetFieldAnnotation.FieldPos-1] == "" {
			// If the field is required it's an error, otherwise skip it
			if targetFieldAnnotation.Attribute == astmconst.ATTRIBUTE_REQUIRED {
				return errmsg.LineParsing_ErrRequiredInputFieldMissing
			} else {
				continue
			}
		}
		// Save the current inputField
		inputField := inputFields[targetFieldAnnotation.FieldPos-1]

		// Set field is value
		err = setField(inputField, targetValues[i], config)
		if err != nil {
			return err
		}
	}

	// Return no error if everything went well
	return nil
}

func setField(value string, field reflect.Value, config *models.Configuration) (err error) {
	// Ensure the field is settable
	if !field.CanSet() {
		// Field is not settable
		return errmsg.LineParsing_ErrNonSettableField
	}
	// Set the field value
	switch field.Kind() {
	case reflect.String:
		if field.Type().ConvertibleTo(reflect.TypeOf("")) {
			field.Set(reflect.ValueOf(value).Convert(field.Type()))
		} else {
			field.Set(reflect.ValueOf(value))
		}
		return nil
	case reflect.Int:
		num, err := strconv.Atoi(value)
		if err != nil {
			return errmsg.LineParsing_ErrDataParsingError
		}
		field.Set(reflect.ValueOf(num))
		return nil
	case reflect.Float32:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return errmsg.LineParsing_ErrDataParsingError
		}
		field.Set(reflect.ValueOf(float32(num)))
		return nil
	case reflect.Float64:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errmsg.LineParsing_ErrDataParsingError
		}
		field.Set(reflect.ValueOf(num))
		return nil
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
			return nil
		} else {
			// Note: option to handle other struct types here
		}
	}
	// Return error if no type match was found (each successful parsing returns nil)
	return errmsg.LineParsing_ErrUsupportedDataType
}
