package functions

import (
	"github.com/blutspende/go-astm/v2/constants/astmconst"
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
	assert.Empty(t, result.Attributes)
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
	assert.Empty(t, result.Attributes)
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
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_REQUIRED)
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
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_LENGTH)
	assert.Equal(t, "2", result.Attributes[astmconst.ATTRIBUTE_LENGTH])
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
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_LENGTH)
	assert.Equal(t, "4", result.Attributes[astmconst.ATTRIBUTE_LENGTH])
}
func TestParseAstmFieldAnnotationString_InvalidAttribute(t *testing.T) {
	// Arrange
	input := "4.1,something"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAttribute.Error())
}
func TestParseAstmFieldAnnotationString_InvalidAttributeFormat(t *testing.T) {
	// Arrange
	input := "4.1,length:2:3"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAttributeFormat.Error())
}
func TestParseAstmFieldAnnotationString_InvalidAnnotationTooManyParts(t *testing.T) {
	// Arrange
	input := "2.1.2"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation.Error())
}
func TestParseAstmFieldAnnotationString_InvalidAnnotationTooManyPartsWithAttribute(t *testing.T) {
	// Arrange
	input := "4.1.3,required"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation.Error())
}
func TestParseAstmFieldAnnotationString_InvalidNumber(t *testing.T) {
	// Arrange
	input := "4.f"
	// Act
	_, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAnnotation.Error())
}
func TestParseAstmFieldAnnotationString_MultipleAttributes(t *testing.T) {
	// Arrange
	input := "4,required,longdate"
	// Act
	result, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "4,required,longdate", result.Raw)
	assert.Equal(t, 4, result.FieldPos)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, 0, result.ComponentPos)
	assert.Equal(t, false, result.IsSubstructure)
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_REQUIRED)
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_LONGDATE)
}
func TestParseAstmFieldAnnotationString_ValuedMultipleAttributes(t *testing.T) {
	// Arrange
	input := "4,length:3,required"
	// Act
	result, err := parseAstmFieldAnnotationString(input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "4,length:3,required", result.Raw)
	assert.Equal(t, 4, result.FieldPos)
	assert.Equal(t, false, result.IsComponent)
	assert.Equal(t, 0, result.ComponentPos)
	assert.Equal(t, false, result.IsSubstructure)
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_LENGTH)
	assert.Equal(t, "3", result.Attributes[astmconst.ATTRIBUTE_LENGTH])
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_REQUIRED)
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
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_LENGTH)
	assert.Equal(t, "4", result.Attributes[astmconst.ATTRIBUTE_LENGTH])
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
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_LENGTH)
	assert.Equal(t, "4", result.Attributes[astmconst.ATTRIBUTE_LENGTH])
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
	assert.Empty(t, result.Attributes)
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
	assert.Empty(t, result.Attributes)
}
func TestParseAstmFieldAnnotation_IllegalComponentArray(t *testing.T) {
	// Arrange
	var input IllegalComponentArray
	field, _ := reflect.TypeOf(input).FieldByName("ComponentArray")
	// Act
	_, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrIllegalComponentArray.Error())
}
func TestParseAstmFieldAnnotation_IllegalComponentSubstructure(t *testing.T) {
	// Arrange
	var input IllegalComponentSubstructure
	field, _ := reflect.TypeOf(input).FieldByName("ComponentSubstructure")
	// Act
	_, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrIllegalComponentSubstructure.Error())
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
	assert.Empty(t, result.Attributes)
}
func TestParseAstmFieldAnnotation_InvalidFieldAttribute(t *testing.T) {
	// Arrange
	var input InvalidFieldAttribute
	field, _ := reflect.TypeOf(input).FieldByName("First")
	// Act
	_, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAttribute.Error())
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
	assert.Empty(t, result.Attributes)
}
func TestParseAstmStructAnnotation_AnnotatedArrayStruct(t *testing.T) {
	// Arrange
	var input AnnotatedArrayStruct
	field, _ := reflect.TypeOf(input).FieldByName("Lines")
	// Act
	result, err := ParseAstmStructAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "L,optional", result.Raw)
	assert.Equal(t, false, result.IsComposite)
	assert.Equal(t, true, result.IsArray)
	assert.Equal(t, "L", result.StructName)
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_OPTIONAL)
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
	assert.Empty(t, result.Attributes)
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
	assert.Empty(t, result.Attributes)
}
func TestParseAstmStructAnnotation_InvalidStructAttribute(t *testing.T) {
	// Arrange
	var input InvalidStructAttribute
	field, _ := reflect.TypeOf(input).FieldByName("Record")
	// Act
	_, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAttribute.Error())
}
func TestParseAstmStructAnnotation_TooManyStructNameAttributeValues(t *testing.T) {
	// Arrange
	var input TooManyStructNameAttributeValues
	field, _ := reflect.TypeOf(input).FieldByName("Record")
	// Act
	_, err := ParseAstmFieldAnnotation(field)
	// Assert
	assert.EqualError(t, err, errmsg.AnnotationParsing_ErrInvalidAstmAttributeFormat.Error())
}
func TestParseAstmStructAnnotation_SubnameAttribute(t *testing.T) {
	// Arrange
	var input SubnameAttribute
	field, _ := reflect.TypeOf(input).FieldByName("Record")
	// Act
	result, err := ParseAstmStructAnnotation(field)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "R,subname:SUBNAME", result.Raw)
	assert.Equal(t, false, result.IsComposite)
	assert.Equal(t, false, result.IsArray)
	assert.Equal(t, "R", result.StructName)
	assert.Contains(t, result.Attributes, astmconst.ATTRIBUTE_SUBNAME)
	assert.Equal(t, "SUBNAME", result.Attributes[astmconst.ATTRIBUTE_SUBNAME])
}

// ProcessStructReflection tests
func TestProcessStructReflection_SimpleRecord(t *testing.T) {
	// Arrange
	input := ThreeFieldRecord{}
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
	input := ThreeFieldRecord{}
	// Act
	_, _, length, err := ProcessStructReflection(&input)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 3, length)
}
