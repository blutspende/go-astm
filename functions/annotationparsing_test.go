package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// Field annotation tests
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
	assert.Equal(t, false, result.IsSubstructure)
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
	assert.Equal(t, false, result.IsSubstructure)
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
	assert.Equal(t, false, result.IsSubstructure)
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
	assert.Equal(t, false, result.IsSubstructure)
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
	assert.Equal(t, false, result.IsSubstructure)
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
func TestParseAstmFieldAnnotation_AnnotatedStruct(t *testing.T) {
	// Arrange
	var input AnnotatedLine
	field, _ := reflect.TypeOf(input).FieldByName("Field")
	// Act
	result, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "3.2,length:4", result.Raw)
	assert.Equal(t, 3, result.FieldPos)
	assert.Equal(t, false, result.IsArray)
	assert.Equal(t, true, result.IsComponent)
	assert.Equal(t, 2, result.ComponentPos)
	assert.Equal(t, false, result.IsSubstructure)
	assert.Equal(t, true, result.HasAttribute)
	assert.Equal(t, constants.ATTRIBUTE_LENGTH, result.Attribute)
	assert.Equal(t, true, result.HasAttributeValue)
	assert.Equal(t, 4, result.AttributeValue)
}
func TestParseAstmFieldAnnotation_AnnotatedArrayStruct(t *testing.T) {
	// Arrange
	var input AnnotatedArrayLine
	field, _ := reflect.TypeOf(input).FieldByName("Field")
	// Act
	result, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "3,length:4", result.Raw)
	assert.Equal(t, 3, result.FieldPos)
	assert.Equal(t, true, result.IsArray)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, false, result.IsSubstructure)
	assert.Equal(t, true, result.HasAttribute)
	assert.Equal(t, constants.ATTRIBUTE_LENGTH, result.Attribute)
	assert.Equal(t, true, result.HasAttributeValue)
	assert.Equal(t, 4, result.AttributeValue)
}
func TestParseAstmFieldAnnotation_Substructure(t *testing.T) {
	// Arrange
	var input SubstructuredLine
	field, _ := reflect.TypeOf(input).FieldByName("Field")
	// Act
	result, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "3", result.Raw)
	assert.Equal(t, 3, result.FieldPos)
	assert.Equal(t, false, result.IsArray)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, true, result.IsSubstructure)
	assert.Equal(t, false, result.HasAttribute)
	assert.Equal(t, false, result.HasAttributeValue)
}
func TestParseAstmFieldAnnotation_SubstructureArray(t *testing.T) {
	// Arrange
	var input SubstructuredLine
	field, _ := reflect.TypeOf(input).FieldByName("Array")
	// Act
	result, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "4", result.Raw)
	assert.Equal(t, 4, result.FieldPos)
	assert.Equal(t, true, result.IsArray)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, true, result.IsSubstructure)
	assert.Equal(t, false, result.HasAttribute)
	assert.Equal(t, false, result.HasAttributeValue)
}
func TestParseAstmFieldAnnotation_IllegalComponentArray(t *testing.T) {
	// Arrange
	var input IllegalComponentArray
	field, _ := reflect.TypeOf(input).FieldByName("ComponentArray")
	// Act
	_, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Error(t, err, errmsg.AnnotationParsing_ErrIllegalComponentArray)
}
func TestParseAstmFieldAnnotation_IllegalComponentSubstructure(t *testing.T) {
	// Arrange
	var input IllegalComponentSubstructure
	field, _ := reflect.TypeOf(input).FieldByName("ComponentSubstructure")
	// Act
	_, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Error(t, err, errmsg.AnnotationParsing_ErrIllegalComponentSubstructure)
}
func TestParseAstmFieldAnnotation_TimeLine(t *testing.T) {
	// Arrange
	var input TimeLine
	field, _ := reflect.TypeOf(input).FieldByName("Time")
	// Act
	result, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "3", result.Raw)
	assert.Equal(t, 3, result.FieldPos)
	assert.Equal(t, false, result.IsArray)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, false, result.IsSubstructure)
	assert.Equal(t, false, result.HasAttribute)
	assert.Equal(t, false, result.HasAttributeValue)
}

// Struct annotation tests
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

// ProcessStructReflection tests
func TestProcessStructReflection_SimpleRecord(t *testing.T) {
	// Arrange
	input := SimpleRecord{}
	// Act
	types, values, length, err := ProcessStructReflection(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 3, length)
	assert.Len(t, types, 3)
	assert.Len(t, values, 3)
	assert.Equal(t, "First", types[0].Name)
}
func TestProcessStructReflection_CompositeRecordStruct(t *testing.T) {
	// Arrange
	input := CompositeRecordStruct{}
	// Act
	types, values, length, err := ProcessStructReflection(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 2, length)
	assert.Len(t, types, 2)
	assert.Len(t, values, 2)
	assert.Equal(t, "Record1", types[0].Name)
}
func TestProcessStructReflection_SimpleRecordPointer(t *testing.T) {
	// Arrange
	input := SimpleRecord{}
	// Act
	_, _, length, err := ProcessStructReflection(&input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 3, length)
}
