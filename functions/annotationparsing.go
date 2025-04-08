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
		result.HasComponent = true
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
	// Get the "astm" tag value and check if it is empty
	raw := input.Tag.Get("astm")
	if raw == "" {
		return models.AstmStructAnnotation{}, errmsg.AnnotationParsing_ErrMissingAstmAnnotation
	}
	result.Raw = raw

	// Determine if the field is an array or not
	result.IsArray = input.Type.Kind() == reflect.Slice || input.Type.Kind() == reflect.Array

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

func isValidAttribute(attribute string) bool {
	validAttributes := []string{
		constants.ATTRIBUTE_DELIMITER,
		constants.ATTRIBUTE_REQUIRE,
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
