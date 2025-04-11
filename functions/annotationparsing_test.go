package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseAstmFieldAnnotationString_SingleValue(t *testing.T) {
	// Arrange
	input := "4"
	// Act
	result, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "4", result.Raw)
	assert.Equal(t, 4, result.FieldPos)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, 0, result.ComponentPos)
	assert.Equal(t, false, result.HasAttribute)
	assert.Equal(t, "", result.Attribute)
	assert.Equal(t, false, result.HasAttributeValue)
	assert.Equal(t, 0, result.AttributeValue)
}
func TestParseAstmFieldAnnotationString_Componented(t *testing.T) {
	// Arrange
	input := "4.1"
	// Act
	result, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "4.1", result.Raw)
	assert.Equal(t, 4, result.FieldPos)
	assert.Equal(t, true, result.IsComponent)
	assert.Equal(t, 1, result.ComponentPos)
	assert.Equal(t, false, result.HasAttribute)
	assert.Equal(t, "", result.Attribute)
	assert.Equal(t, false, result.HasAttributeValue)
	assert.Equal(t, 0, result.AttributeValue)
}
func TestParseAstmFieldAnnotationString_Attributed(t *testing.T) {
	// Arrange
	input := "4,required"
	// Act
	result, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "4,required", result.Raw)
	assert.Equal(t, 4, result.FieldPos)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, 0, result.ComponentPos)
	assert.Equal(t, true, result.HasAttribute)
	assert.Equal(t, constants.ATTRIBUTE_REQUIRED, result.Attribute)
	assert.Equal(t, false, result.HasAttributeValue)
	assert.Equal(t, 0, result.AttributeValue)
}
func TestParseAstmFieldAnnotationString_AttributedValue(t *testing.T) {
	// Arrange
	input := "4,length:2"
	// Act
	result, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "4,length:2", result.Raw)
	assert.Equal(t, 4, result.FieldPos)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, 0, result.ComponentPos)
	assert.Equal(t, true, result.HasAttribute)
	assert.Equal(t, constants.ATTRIBUTE_LENGTH, result.Attribute)
	assert.Equal(t, true, result.HasAttributeValue)
	assert.Equal(t, 2, result.AttributeValue)
}
func TestParseAstmFieldAnnotationString_Complex(t *testing.T) {
	// Arrange
	input := "3.2,length:4"
	// Act
	result, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "3.2,length:4", result.Raw)
	assert.Equal(t, 3, result.FieldPos)
	assert.Equal(t, true, result.IsComponent)
	assert.Equal(t, 2, result.ComponentPos)
	assert.Equal(t, true, result.HasAttribute)
	assert.Equal(t, constants.ATTRIBUTE_LENGTH, result.Attribute)
	assert.Equal(t, true, result.HasAttributeValue)
	assert.Equal(t, 4, result.AttributeValue)
}
func TestParseAstmFieldAnnotationString_InvalidAttribute(t *testing.T) {
	// Arrange
	input := "4.1,something"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Error(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAttribute)
}
func TestParseAstmFieldAnnotationString_InvalidAnnotationTooManyParts(t *testing.T) {
	// Arrange
	input := "2.1.2"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Error(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation)
}
func TestParseAstmFieldAnnotationString_InvalidAnnotationTooManyPartsWithAttribute(t *testing.T) {
	// Arrange
	input := "4.1.3,something"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Error(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation)
}
func TestParseAstmFieldAnnotationString_InvalidNumber(t *testing.T) {
	// Arrange
	input := "4.f"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Error(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation)
}
func TestParseAstmFieldAnnotationString_TooManyAttributes(t *testing.T) {
	// Arrange
	input := "4,something,otherthing"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Error(t, err, errmsg.AnnotationParsing_ErrTooManyAttributes)
}

type AnnotatedLine struct {
	Field string `astm:"3.2,length:4"`
}

func TestParseAstmFieldAnnotation_AnnotatedStruct(t *testing.T) {
	// Arrange
	var input AnnotatedLine
	field, _ := reflect.TypeOf(input).FieldByName("Field")
	// Act
	result, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "3.2,length:4", result.Raw)
	assert.Equal(t, false, result.IsArray)
	assert.Equal(t, 3, result.FieldPos)
	assert.Equal(t, true, result.IsComponent)
	assert.Equal(t, 2, result.ComponentPos)
	assert.Equal(t, true, result.HasAttribute)
	assert.Equal(t, constants.ATTRIBUTE_LENGTH, result.Attribute)
	assert.Equal(t, true, result.HasAttributeValue)
	assert.Equal(t, 4, result.AttributeValue)
}

type AnnotatedArrayLine struct {
	Field []string `astm:"3.2,length:4"`
}

func TestParseAstmFieldAnnotation_AnnotatedArrayStruct(t *testing.T) {
	// Arrange
	var input AnnotatedArrayLine
	field, _ := reflect.TypeOf(input).FieldByName("Field")
	// Act
	result, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "3.2,length:4", result.Raw)
	assert.Equal(t, true, result.IsArray)
	assert.Equal(t, 3, result.FieldPos)
	assert.Equal(t, true, result.IsComponent)
	assert.Equal(t, 2, result.ComponentPos)
	assert.Equal(t, true, result.HasAttribute)
	assert.Equal(t, constants.ATTRIBUTE_LENGTH, result.Attribute)
	assert.Equal(t, true, result.HasAttributeValue)
	assert.Equal(t, 4, result.AttributeValue)
}

type Line struct {
	Field string `astm:"2"`
}
type SingleLineStruct struct {
	Lines Line `astm:"L"`
}

func TestParseAstmStructAnnotation_SingleLineStruct(t *testing.T) {
	// Arrange
	var input SingleLineStruct
	field, _ := reflect.TypeOf(input).FieldByName("Lines")
	// Act
	result, err := ParseAstmStructAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "L", result.Raw)
	assert.Equal(t, false, result.IsComposite)
	assert.Equal(t, false, result.IsArray)
	assert.Equal(t, "L", result.StructName)
	assert.Equal(t, false, result.HasAttribute)
	assert.Equal(t, "", result.Attribute)
}

type AnnotatedArrayStruct struct {
	Lines []Line `astm:"L,required"`
}

func TestParseAstmStructAnnotation_AnnotatedArrayStruct(t *testing.T) {
	// Arrange
	var input AnnotatedArrayStruct
	field, _ := reflect.TypeOf(input).FieldByName("Lines")
	// Act
	result, err := ParseAstmStructAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "L,required", result.Raw)
	assert.Equal(t, false, result.IsComposite)
	assert.Equal(t, true, result.IsArray)
	assert.Equal(t, "L", result.StructName)
	assert.Equal(t, true, result.HasAttribute)
	assert.Equal(t, constants.ATTRIBUTE_REQUIRED, result.Attribute)
}

type CompositeStruct struct {
	Composite AnnotatedArrayStruct
}

func TestParseAstmStructAnnotation_CompositeStruct(t *testing.T) {
	// Arrange
	var input CompositeStruct
	field, _ := reflect.TypeOf(input).FieldByName("Composite")
	// Act
	result, err := ParseAstmStructAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "", result.Raw)
	assert.Equal(t, true, result.IsComposite)
	assert.Equal(t, false, result.IsArray)
	assert.Equal(t, "", result.StructName)
	assert.Equal(t, false, result.HasAttribute)
	assert.Equal(t, "", result.Attribute)
}

type CompositeArrayStruct struct {
	Composite []AnnotatedArrayStruct
}

func TestParseAstmStructAnnotation_CompositeArrayStruct(t *testing.T) {
	// Arrange
	var input CompositeArrayStruct
	field, _ := reflect.TypeOf(input).FieldByName("Composite")
	// Act
	result, err := ParseAstmStructAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "", result.Raw)
	assert.Equal(t, true, result.IsComposite)
	assert.Equal(t, true, result.IsArray)
	assert.Equal(t, "", result.StructName)
	assert.Equal(t, false, result.HasAttribute)
	assert.Equal(t, "", result.Attribute)
}
