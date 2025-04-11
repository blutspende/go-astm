package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/blutspende/go-astm/v2/models"
	"reflect"
	"strconv"
	"strings"
)

func ParseAstmFieldAnnotation(input reflect.StructField) (result models.AstmFieldAnnotation, err error) {
	// Get the "astm" tag value and check if it is empty
	raw := input.Tag.Get("astm")
	if raw == "" {
		return models.AstmFieldAnnotation{}, errmsg.AnnotationParsing_ErrMissingAstmAnnotation
	}

	// Parse the annotation string
	result, err = parseAstmFieldAnnotationString(raw)

	// Determine if the field is an array or not
	result.IsArray = input.Type.Kind() == reflect.Slice || input.Type.Kind() == reflect.Array

	return result, err
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

// TODO: find a better place for this
func ProcessStructReflection(targetStruct interface{}) (targetTypes []reflect.StructField, targetValues []reflect.Value, length int, err error) {
	// Ensure the targetStruct is a pointer to a struct
	targetPtrValue := reflect.ValueOf(targetStruct)
	if targetPtrValue.Kind() != reflect.Ptr || targetPtrValue.Elem().Kind() != reflect.Struct {
		// targetStruct must be a pointer to a struct
		return nil, nil, 0, errmsg.AnnotationParsing_ErrInvalidTargetStruct
	}

	// Get the underlying struct
	targetValue := targetPtrValue.Elem()
	targetType := targetPtrValue.Type().Elem()

	// Allocate the results
	targetTypes = make([]reflect.StructField, targetValue.NumField())
	targetValues = make([]reflect.Value, targetType.NumField())
	length = targetType.NumField()

	// Iterate and save targetTypes and targetValues
	for i := 0; i < targetType.NumField(); i++ {
		targetTypes[i] = targetType.Field(i)
		targetValues[i] = targetValue.Field(i)
	}

	// Return the results
	return targetTypes, targetValues, length, nil
}

func isValidAttribute(attribute string) bool {
	validAttributes := []string{
		constants.ATTRIBUTE_DELIMITER,
		constants.ATTRIBUTE_REQUIRED,
		constants.ATTRIBUTE_OPTIONAL,
		constants.ATTRIBUTE_SEQUENCE,
		constants.ATTRIBUTE_LONGDATE,
		constants.ATTRIBUTE_LENGTH,
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
