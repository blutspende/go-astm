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

	// Separate attributes and the field definitions
	mainParts := strings.Split(result.Raw, ",")
	if len(mainParts) > 2 {
		return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrTooManyAttributes
	}

	// If there is an attribute parse it
	if len(mainParts) == 2 {
		result.HasAttribute = true
		// Parse attribute value if there is any
		attributeParts := strings.Split(mainParts[1], ":")
		if len(attributeParts) > 2 {
			return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrInvalidAstmAttribute
		} else if len(attributeParts) == 2 {
			result.HasAttributeValue = true
			result.AttributeValue, err = strconv.Atoi(attributeParts[1])
			if err != nil {
				return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrInvalidAstmAttribute
			}
		}
		// Check attribute value to be any of the allowed values
		if !isValidAttribute(attributeParts[0]) {
			return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrInvalidAstmAttribute
		}
		result.Attribute = attributeParts[0]
	}

	// Split field and component (if any) and parse them
	segments := strings.Split(mainParts[0], ".")
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

	// Separate attribute (if any) and the struct name
	mainParts := strings.Split(result.Raw, ",")
	if len(mainParts) > 2 {
		return models.AstmStructAnnotation{}, errmsg.AnnotationParsing_ErrTooManyAttributes
	}
	// If there is an attribute parse it
	if len(mainParts) == 2 {
		result.HasAttribute = true
		// Check attribute value to be any of the allowed values
		if !isValidAttribute(mainParts[1]) {
			return models.AstmStructAnnotation{}, errmsg.AnnotationParsing_ErrInvalidAstmAttribute
		}
		result.Attribute = mainParts[1]
	}

	// Validate and save the struct name
	if len(mainParts[0]) != 1 {
		return models.AstmStructAnnotation{}, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation
	}
	result.StructName = mainParts[0]

	return result, err
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

func isValidAttribute(attribute string) bool {
	validAttributes := []string{
		astmconst.ATTRIBUTE_REQUIRED,
		astmconst.ATTRIBUTE_OPTIONAL,
		astmconst.ATTRIBUTE_LONGDATE,
		astmconst.ATTRIBUTE_LENGTH,
	}
	return isInList(attribute, validAttributes)
}
func isInList(target string, list []string) bool {
	set := make(map[string]struct{})
	for _, item := range list {
		set[item] = struct{}{}
	}
	_, exists := set[target]
	return exists
}
