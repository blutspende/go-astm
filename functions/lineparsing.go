package functions

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/blutspende/go-astm/v3/constants"
	"github.com/blutspende/go-astm/v3/errmsg"
	"github.com/blutspende/go-astm/v3/models"
	"github.com/blutspende/go-astm/v3/models/astmmodels"
)

func ParseLine(inputLine string, targetStruct interface{}, recordAnnotation models.AstmStructAnnotation, sequenceNumber int, config *astmmodels.Configuration) (nameOk bool, err error) {
	// Check for input line length
	if len(inputLine) == 0 {
		return false, errmsg.ErrLineParsingEmptyInput
	}

	// Setup inputFields variable to split the inputLine into
	var inputFields []string

	// Handle header special case
	if inputLine[0] == 'H' {
		// Check if the inputLine is long enough to contain delimiters
		if len(inputLine) < 5 {
			return false, errmsg.ErrLineParsingHeaderTooShort
		}
		// Override delimiters
		config.Delimiters.Field = string(inputLine[1])
		config.Delimiters.Repeat = string(inputLine[2])
		config.Delimiters.Component = string(inputLine[3])
		config.Delimiters.Escape = string(inputLine[4])
		// Place the fix segment into the inputFields
		inputFields = []string{inputLine[0:1], inputLine[1:5]}
		// Add the rest of the inputLine split by the field delimiter
		inputFields = append(inputFields, splitStringWithEscape(inputLine[6:], config.Delimiters.Field, config.Delimiters.Escape)...)
	} else {
		// Split the input with the field delimiter
		inputFields = splitStringWithEscape(inputLine, config.Delimiters.Field, config.Delimiters.Escape)
	}

	// Check for minimum number of input fields (first two fields are mandatory)
	if len(inputFields) < 2 {
		return false, errmsg.ErrLineParsingMandatoryInputFieldsMissing
	}

	// Check for mach of name and subname
	// Note: name checking is always enforced, but instead of error it is returned in the nameOk variable
	if inputFields[0] != recordAnnotation.StructName {
		return false, nil
	}
	if subname, exists := recordAnnotation.Attributes[constants.AttributeSubname]; exists {
		// If subname is given at least 3 fields are required
		if len(inputFields) < 3 {
			return false, errmsg.ErrLineParsingMandatoryInputFieldsMissing
		}
		// Check for subname match
		if inputFields[2] != subname {
			return false, nil
		}
	}

	// Check for validity of the sequence number (error only if enforced)
	if inputFields[1] != strconv.Itoa(sequenceNumber) && inputLine[0] != 'H' && config.EnforceSequenceNumberCheck {
		return true, errmsg.ErrLineParsingSequenceNumberMismatch
	}

	// Process the target structure
	targetTypes, targetValues, _, err := ProcessStructReflection(targetStruct)
	if err != nil {
		return true, err
	}

	// Iterate over the inputFields of the targetStruct struct
	for i, targetType := range targetTypes {
		// Parse the targetStruct field targetFieldAnnotation
		targetFieldAnnotation, err := ParseAstmFieldAnnotation(targetType)
		if err != nil {
			if errors.Is(err, errmsg.ErrAnnotationParsingMissingAstmAnnotation) {
				// If the annotation is missing, skip this field
				continue
			} else {
				return true, err
			}
		}

		// Check for fieldPos not being lower than 3 (first 2 are reserved for line name and sequence number)
		if targetFieldAnnotation.FieldPos < 3 {
			return true, errmsg.ErrLineParsingReservedFieldPosReference
		}

		// Not enough inputFields or empty inputField
		if len(inputFields) < targetFieldAnnotation.FieldPos || inputFields[targetFieldAnnotation.FieldPos-1] == "" {
			// If the field is required it's an error, otherwise skip it
			if _, exists := targetFieldAnnotation.Attributes[constants.AttributeRequired]; exists {
				return true, errmsg.ErrLineParsingRequiredInputFieldMissing
			} else {
				continue
			}
		}
		// Save the current inputField
		inputField := inputFields[targetFieldAnnotation.FieldPos-1]

		if targetFieldAnnotation.IsArray {
			// |rep1\rep2\rep3|
			// Field is an array
			repeats := splitStringWithEscape(inputField, config.Delimiters.Repeat, config.Delimiters.Escape)
			arrayType := reflect.SliceOf(targetValues[i].Type().Elem())
			arrayValue := reflect.MakeSlice(arrayType, len(repeats), len(repeats))
			for j, repeat := range repeats {
				if targetFieldAnnotation.IsSubstructure {
					// |comp1^comp2^comp3\comp1^comp2^comp3\comp1^comp2^comp3|
					// Substructures (with components) in the array: use parseSubstructure
					err = parseSubstructure(repeat, arrayValue.Index(j).Addr().Interface(), config)
					if err != nil {
						return true, err
					}
				} else {
					// |value1\value2\value3|
					// Simple values in the array
					err = setField(repeat, arrayValue.Index(j), targetFieldAnnotation, config)
					if err != nil {
						return true, err
					}
				}

			}
			targetValues[i].Set(arrayValue)
		} else if targetFieldAnnotation.IsComponent {
			// |comp1^comp2^comp3|
			// Field is a component
			components := splitStringWithEscape(inputField, config.Delimiters.Component, config.Delimiters.Escape)
			// Not enough components in the inputField
			if len(components) < targetFieldAnnotation.ComponentPos {
				// Error if the component is required, skip otherwise
				if _, exists := targetFieldAnnotation.Attributes[constants.AttributeRequired]; exists {
					return true, errmsg.ErrLineParsingInputComponentsMissing
				} else {
					continue
				}
			}
			err = setField(components[targetFieldAnnotation.ComponentPos-1], targetValues[i], targetFieldAnnotation, config)
			if err != nil {
				return true, err
			}
		} else if targetFieldAnnotation.IsSubstructure {
			// |comp1^comp2^comp3|
			// If the field is a substructure use parseSubstructure to process it
			err = parseSubstructure(inputField, targetValues[i].Addr().Interface(), config)
			if err != nil {
				return true, err
			}
		} else {
			// |field|
			// Field is not an array or component (normal singular field)
			err = setField(inputField, targetValues[i], targetFieldAnnotation, config)
			if err != nil {
				return true, err
			}
		}
		// Note: this could be a place to produce warnings about lost data
		// if i == targetFieldCount-1 && len(inputFields) > targetFieldAnnotation.FieldPos
	}
	// Return no error if everything went well
	return true, nil
}

func parseSubstructure(inputString string, targetStruct interface{}, config *astmmodels.Configuration) (err error) {
	// Split the input with the field delimiter
	inputFields := splitStringWithEscape(inputString, config.Delimiters.Component, config.Delimiters.Escape)

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
			if errors.Is(err, errmsg.ErrAnnotationParsingMissingAstmAnnotation) {
				// If the annotation is missing, skip this field
				continue
			} else {
				return err
			}
		}

		// Not enough inputFields or empty inputField
		if len(inputFields) < targetFieldAnnotation.FieldPos || inputFields[targetFieldAnnotation.FieldPos-1] == "" {
			// If the field is required it's an error, otherwise skip it
			if _, exists := targetFieldAnnotation.Attributes[constants.AttributeRequired]; exists {
				return errmsg.ErrLineParsingRequiredInputFieldMissing
			} else {
				continue
			}
		}
		// Save the current inputField
		inputField := inputFields[targetFieldAnnotation.FieldPos-1]

		// Set field is value
		err = setField(inputField, targetValues[i], targetFieldAnnotation, config)
		if err != nil {
			return err
		}
	}

	// Return no error if everything went well
	return nil
}

func setField(value string, field reflect.Value, annotation models.AstmFieldAnnotation, config *astmmodels.Configuration) (err error) {
	// Ensure the field is settable
	if !field.CanSet() {
		// Field is not settable
		return errmsg.ErrLineParsingNonSettableField
	}
	// Set the field value
	switch field.Kind() {
	case reflect.String:
		escaped := filterStringEscapeChars(value, config.Delimiters.Escape)
		if field.Type().ConvertibleTo(reflect.TypeOf("")) {
			field.Set(reflect.ValueOf(escaped).Convert(field.Type()))
		} else {
			field.Set(reflect.ValueOf(escaped))
		}
		return nil
	case reflect.Int:
		num, err := strconv.Atoi(value)
		if err != nil {
			return errmsg.ErrLineParsingDataParsingError
		}
		field.Set(reflect.ValueOf(num))
		return nil
	case reflect.Float32:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return errmsg.ErrLineParsingDataParsingError
		}
		field.Set(reflect.ValueOf(float32(num)))
		return nil
	case reflect.Float64:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errmsg.ErrLineParsingDataParsingError
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
				return errmsg.ErrLineParsingInvalidDateFormat
			}
			timeInLocation, err := time.ParseInLocation(timeFormat, value, config.TimeLocation)
			if err != nil {
				return errmsg.ErrLineParsingDataParsingError
			}
			if _, exists := annotation.Attributes[constants.AttributeLongdate]; !exists && config.KeepShortDateTimeZone {
				// Keep the short date time zone
				timeInLocation = timeInLocation.In(config.TimeLocation)
			} else {
				// Set the time to UTC
				timeInLocation = timeInLocation.UTC()
			}
			field.Set(reflect.ValueOf(timeInLocation))
			return nil
		} else {
			// Note: option to handle other struct types here
		}
	}
	// Return error if no type match was found (each successful parsing returns nil)
	return errmsg.ErrLineParsingUnsupportedDataType
}

func splitStringWithEscape(input, delimiter, escape string) []string {
	var result []string
	delimiterRune := rune(delimiter[0])
	escapeRune := rune(escape[0])
	inputRunes := []rune(input)
	start := 0
	for i := 0; i < len(inputRunes); i++ {
		if inputRunes[i] == delimiterRune {
			result = append(result, string(inputRunes[start:i]))
			start = i + 1
		}
		if inputRunes[i] == escapeRune {
			if i+1 < len(inputRunes) && string(inputRunes[i+1]) == "Z" {
				for j := i + 2; j < len(inputRunes); j++ {
					if string(inputRunes[j]) == string(escapeRune) {
						i = j
						break
					}
				}
			} else {
				i++
			}
			continue
		}
	}

	if start <= len(inputRunes)-1 {
		result = append(result, string(inputRunes[start:]))
	}

	return result
}

func filterStringEscapeChars(input string, escape string) string {
	var builder strings.Builder
	escapeRune := rune(escape[0])
	inputRunes := []rune(input)
	for i := 0; i < len(inputRunes); i++ {
		if inputRunes[i] == escapeRune {
			i++
			if i < len(inputRunes) {
				builder.WriteRune(inputRunes[i])
			}
		} else {
			builder.WriteRune(inputRunes[i])
		}
	}
	return builder.String()
}
