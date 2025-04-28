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

func ParseAstmFieldAnnotation(input reflect.StructField) (result models.AstmFieldAnnotation, err error) {
	// Get the "astm" tag value and check if it is empty
	raw := input.Tag.Get("astm")
	if raw == "" {
		return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrMissingAstmAnnotation
	}

	// Parse the annotation string
	result, err = parseAstmFieldAnnotationString(raw)
	if err != nil {
		return models.AstmFieldAnnotation{}, err
	}

	// Determine if the field is an array or not
	result.IsArray = input.Type.Kind() == reflect.Slice || input.Type.Kind() == reflect.Array

	// Determine if the field is a substructure or not (excluding the time.Time type)
	var checkType reflect.Type
	if result.IsArray {
		checkType = input.Type.Elem()
	} else {
		checkType = input.Type
	}
	result.IsSubstructure = checkType.Kind() == reflect.Struct && checkType != reflect.TypeOf(time.Time{})

	// Check illegal combinations
	if result.IsComponent && result.IsArray {
		return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrIllegalComponentArray
	}
	if result.IsComponent && result.IsSubstructure {
		return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrIllegalComponentSubstructure
	}

	// All okay, return the result and no error
	return result, nil
}

func parseAstmFieldAnnotationString(input string) (result models.AstmFieldAnnotation, err error) {
	result.Raw = input

	// Separate attributes and the field definition
	fieldDef, attributes := splitByFirst(input, ",")

	// Parse and save attributes
	result.Attributes, err = parseAttributes(attributes, []string{
		astmconst.ATTRIBUTE_REQUIRED,
		astmconst.ATTRIBUTE_LONGDATE,
		astmconst.ATTRIBUTE_LENGTH,
	})
	if err != nil {
		return models.AstmFieldAnnotation{}, err
	}

	// Split field and component (if any) and parse them
	segments := strings.Split(fieldDef, ".")
	if len(segments) > 2 {
		return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation
	}
	if len(segments) == 2 {
		result.IsComponent = true
		result.ComponentPos, err = strconv.Atoi(segments[1])
		if err != nil {
			return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation
		}
	}
	result.FieldPos, err = strconv.Atoi(segments[0])
	if err != nil {
		return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation
	}

	return result, nil
}

func ParseAstmStructAnnotation(input reflect.StructField) (result models.AstmStructAnnotation, err error) {
	// Get the "astm" tag value
	raw := input.Tag.Get("astm")
	result.Raw = raw

	// Determine if the struct is composite (no tag) or not
	result.IsComposite = raw == ""

	// Determine if the field is an array or not
	result.IsArray = input.Type.Kind() == reflect.Slice || input.Type.Kind() == reflect.Array

	// Composite has no tag so further parsing is not needed
	if result.IsComposite {
		return result, nil
	}

	// Separate attributes and the struct name, and save the name
	attributes := ""
	result.StructName, attributes = splitByFirst(raw, ",")

	// Parse and save attributes
	result.Attributes, err = parseAttributes(attributes, []string{
		astmconst.ATTRIBUTE_OPTIONAL,
		astmconst.ATTRIBUTE_SUBNAME,
	})

	return result, err
}

func parseAttributes(input string, valids []string) (result map[string]string, err error) {
	// Initialize the result map
	result = make(map[string]string)
	// Check for empty input (if empty still return a usable empty map)
	if input == "" {
		return result, nil
	}
	// Split the input string by commas
	attributes := strings.Split(input, ",")
	// Iterate over the attributes and parse them
	for _, attribute := range attributes {
		// Split each attribute by the colon
		attributeParts := strings.Split(attribute, ":")
		if len(attributeParts) > 2 {
			return nil, errmsg.AnnotationParsing_ErrInvalidAstmAttributeFormat
		}
		// Check if the attribute is valid
		if !isInList(attributeParts[0], valids) {
			return nil, errmsg.AnnotationParsing_ErrInvalidAstmAttribute
		}
		// Save the attribute name and value (if present)
		if len(attributeParts) == 2 {
			result[attributeParts[0]] = attributeParts[1]
		} else {
			result[attributeParts[0]] = ""
		}
	}
	// Return the result map and no error
	return result, nil
}

func splitByFirst(input string, delimiter string) (before string, after string) {
	index := strings.Index(input, delimiter) // Find the first occurrence of the comma
	if index == -1 {
		return input, "" // No comma, return whole string and empty second part
	}
	return input[:index], input[index+1:] // Split at the first comma
}
func isInList(target string, list []string) bool {
	set := make(map[string]struct{})
	for _, item := range list {
		set[item] = struct{}{}
	}
	_, exists := set[target]
	return exists
}

func ProcessStructReflection(inputStruct interface{}) (outputTypes []reflect.StructField, outputValues []reflect.Value, length int, err error) {
	// Ensure the inputStruct is a pointer to a struct
	targetPtrValue := reflect.ValueOf(inputStruct)
	if targetPtrValue.Kind() != reflect.Ptr {
		// If inputStruct is not a pointer, take its address
		targetPtrValue = reflect.New(reflect.TypeOf(inputStruct))
		targetPtrValue.Elem().Set(reflect.ValueOf(inputStruct))
	}
	if targetPtrValue.Elem().Kind() != reflect.Struct {
		// inputStruct must be a pointer to a struct
		return nil, nil, 0, errmsg.AnnotationParsing_ErrInvalidInputStruct
	}

	// Get the underlying struct
	targetValue := targetPtrValue.Elem()
	targetType := targetPtrValue.Type().Elem()

	// Allocate the results
	outputTypes = make([]reflect.StructField, targetValue.NumField())
	outputValues = make([]reflect.Value, targetType.NumField())
	length = targetType.NumField()

	// Iterate and save outputTypes and outputValues
	for i := 0; i < targetType.NumField(); i++ {
		outputTypes[i] = targetType.Field(i)
		outputValues[i] = targetValue.Field(i)
	}

	// Return the results
	return outputTypes, outputValues, length, nil
}
