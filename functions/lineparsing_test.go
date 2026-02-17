package functions

import (
	"testing"
	"time"

	"github.com/blutspende/go-astm/v3/errmsg"
	"github.com/blutspende/go-astm/v3/models/astmmodels"
	"github.com/blutspende/go-astm/v3/models/messageformat/lis02a2"
	"github.com/stretchr/testify/assert"
)

// Note: structures come from functions_test.go

func TestParseLine_SimpleRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second|third"
	target := ThreeFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "second", target.Second)
	assert.Equal(t, "third", target.Third)
}

func TestParseLine_UnorderedRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second|third"
	target := UnorderedRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "second", target.Second)
	assert.Equal(t, "third", target.Third)
}

func TestParseLine_MultitypeRecord(t *testing.T) {
	// Arrange
	input := "T|1|string|3|3.14|3.14159265|20060102"
	target := MultitypeRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "string", target.String)
	assert.Equal(t, 3, target.Int)
	assert.Equal(t, float32(3.14), target.Float32)
	assert.Equal(t, 3.14159265, target.Float64)
	expectedShortTime := time.Date(2006, 1, 2, 0, 0, 0, 0, config.TimeLocation)
	assert.Equal(t, expectedShortTime, target.Date)
}

func TestParseLine_ComponentedRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second1^second2|third1^third2^third3"
	target := ComponentedRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "second1", target.SecondComp1)
	assert.Equal(t, "second2", target.SecondComp2)
	assert.Equal(t, "third1", target.ThirdComp1)
	assert.Equal(t, "third2", target.ThirdComp2)
	assert.Equal(t, "third3", target.ThirdComp3)
}

func TestParseLine_ArrayRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second1\\second2\\second3"
	target := ArrayRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Len(t, target.Array, 3)
	assert.Equal(t, "second1", target.Array[0])
	assert.Equal(t, "second2", target.Array[1])
	assert.Equal(t, "second3", target.Array[2])
}

func TestParseLine_HeaderRecord(t *testing.T) {
	// Arrange
	input := "H|\\^&|first"
	target := HeaderRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("H"), 0, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
}

func TestParseLine_HeaderDelimiterChange(t *testing.T) {
	// Arrange
	input := "H/!*%/first/second1!second2/third1*third2"
	target := HeaderDelimiterChange{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("H"), 0, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Len(t, target.Array, 2)
	assert.Equal(t, "second1", target.Array[0])
	assert.Equal(t, "second2", target.Array[1])
	assert.Equal(t, "third1", target.Comp1)
	assert.Equal(t, "third2", target.Comp2)
	// Teardown
	teardown()
}

func TestParseLine_MissingData(t *testing.T) {
	// Arrange
	input := "T|1|first||third"
	target := ThreeFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "", target.Second)
	assert.Equal(t, "third", target.Third)
}

func TestParseLine_MissingComponent(t *testing.T) {
	// Arrange
	input := "T|1|first^second"
	target := RequiredComponentRecord{}
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "second", target.Second)
	assert.Equal(t, "", target.Third)
}

func TestParseLine_MissingRequiredComponent(t *testing.T) {
	// Arrange
	input := "T|1|first"
	target := RequiredComponentRecord{}
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.EqualError(t, err, errmsg.ErrLineParsingInputComponentsMissing.Error())
}

func TestParseLine_MissingDataAtTheEnd(t *testing.T) {
	// Arrange
	input := "T|1|first"
	target := ThreeFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "", target.Second)
	assert.Equal(t, "", target.Third)
}

func TestParseLine_EnumRecord(t *testing.T) {
	// Arrange
	input := "T|1|enum"
	target := EnumRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, EnumString("enum"), target.Enum)
}

func TestParseLine_RecordTypeNameMismatch(t *testing.T) {
	// Arrange
	input := "W|1|first|second|third"
	target := ThreeFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.False(t, nameOk)
}

func TestParseLine_PatientEscapeSequence(t *testing.T) {
	// Arrange
	input := `P|1||||\ZTest1\ \ZTest2\^\ZTest3\|||U|||||\ZTest4\`
	target := lis02a2.Patient{}
	cfg := &astmmodels.DefaultConfiguration
	cfg.Delimiters = astmmodels.Delimiters{
		Field:     "|",
		Repeat:    "@",
		Component: "^",
		Escape:    "\\",
	}

	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("P"), 1, cfg)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "ZTest1 ZTest2", target.LastName)
	assert.Equal(t, "ZTest3", target.FirstName)
	assert.Equal(t, "U", target.Gender)
	assert.Equal(t, "ZTest4", target.AttendingPhysicianID)
}

func TestParseLine_PatientEscapeSequenceEmptyFirstName(t *testing.T) {
	// Arrange
	input := `P|1||||\ZTest1\ \ZTest2\^|||U|||||\ZTest4\`
	target := lis02a2.Patient{}
	cfg := &astmmodels.DefaultConfiguration
	cfg.Delimiters = astmmodels.Delimiters{
		Field:     "|",
		Repeat:    "@",
		Component: "^",
		Escape:    "\\",
	}

	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("P"), 1, cfg)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "ZTest1 ZTest2", target.LastName)
	assert.Equal(t, "", target.FirstName)
	assert.Equal(t, "U", target.Gender)
	assert.Equal(t, "ZTest4", target.AttendingPhysicianID)
}

func TestParseLine_EmptyInput(t *testing.T) {
	// Arrange
	input := ""
	target := ThreeFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.False(t, nameOk)
	assert.EqualError(t, err, errmsg.ErrLineParsingEmptyInput.Error())
}

func TestParseLine_MandatoryFieldsMissing(t *testing.T) {
	// Arrange
	input := "T"
	target := ThreeFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.False(t, nameOk)
	assert.EqualError(t, err, errmsg.ErrLineParsingMandatoryInputFieldsMissing.Error())
}

func TestParseLine_MissingRequiredField(t *testing.T) {
	// Arrange
	input := "T|1|first||third"
	target := RequiredFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.True(t, nameOk)
	assert.EqualError(t, err, errmsg.ErrLineParsingRequiredInputFieldMissing.Error())
}

func TestParseLine_NotEnoughInputFields(t *testing.T) {
	// Arrange
	input := "T|1|first"
	target := RequiredFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.True(t, nameOk)
	assert.EqualError(t, err, errmsg.ErrLineParsingRequiredInputFieldMissing.Error())
}

func TestParseLine_SequenceNumberMismatch(t *testing.T) {
	// Arrange
	input := "T|2|first|second|third"
	target := ThreeFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.True(t, nameOk)
	assert.EqualError(t, err, errmsg.ErrLineParsingSequenceNumberMismatch.Error())
}

func TestParseLine_SequenceNumberMismatchWithoutEnforcing(t *testing.T) {
	// Arrange
	input := "T|2|first|second|third"
	target := ThreeFieldRecord{}
	config.EnforceSequenceNumberCheck = false
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	// Teardown
	teardown()
}

func TestParseLine_ReservedFieldRecord(t *testing.T) {
	// Arrange
	input := "T|1"
	target := ReservedFieldRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.True(t, nameOk)
	assert.EqualError(t, err, errmsg.ErrLineParsingReservedFieldPosReference.Error())
}

func TestParseLine_SubstructureRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|firstComponent^secondComponent^thirdComponent|third"
	target := SubstructureRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "firstComponent", target.Second.FirstComponent)
	assert.Equal(t, "secondComponent", target.Second.SecondComponent)
	assert.Equal(t, "thirdComponent", target.Second.ThirdComponent)
	assert.Equal(t, "third", target.Third)
}

func TestParseLine_SubstructureArrayRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|r1c1^r1c2^r1c3\\r2c1^r2c2^r2c3|third"
	target := SubstructureArrayRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Len(t, target.Second, 2)
	assert.Equal(t, "r1c1", target.Second[0].FirstComponent)
	assert.Equal(t, "r1c2", target.Second[0].SecondComponent)
	assert.Equal(t, "r1c3", target.Second[0].ThirdComponent)
	assert.Equal(t, "r2c1", target.Second[1].FirstComponent)
	assert.Equal(t, "r2c2", target.Second[1].SecondComponent)
	assert.Equal(t, "r2c3", target.Second[1].ThirdComponent)
	assert.Equal(t, "third", target.Third)
}

func TestParseLine_TimeLineTimeZone(t *testing.T) {
	// Arrange
	input := "T|1|20060306164429"
	target := TimeRecord{}
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	expectedTime := time.Date(2006, 03, 06, 16, 44, 29, 0, config.TimeLocation).UTC()
	assert.Equal(t, expectedTime, target.Time)
}

func TestParseLine_WrongComponentOrder(t *testing.T) {
	// Arrange
	input := "T|1|first|comp1^comp2^comp3"
	target := WrongComponentOrderRecord{}
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "comp1", target.Comp1)
	assert.Equal(t, "comp2", target.Comp2)
	assert.Equal(t, "comp3", target.Comp3)
}

func TestParseLine_WrongComponentPlacement(t *testing.T) {
	// Arrange
	input := "T|1|field1|comp1^comp2|field2"
	target := WrongComponentPlacementRecord{}
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "field1", target.Field1)
	assert.Equal(t, "comp1", target.Comp1)
	assert.Equal(t, "field2", target.Field2)
	assert.Equal(t, "comp2", target.Comp2)
}

func TestParseLine_MissingAnnotation(t *testing.T) {
	// Arrange
	input := "T|1|field3|field4"
	target := MissingAnnotationRecord{}
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "field3", target.Field3)
	assert.Equal(t, "field4", target.Field4)
}

func TestParseLine_ShortDateKeepTimeZone(t *testing.T) {
	// Arrange
	input := "T|1|20060304"
	target := ShortDateRecord{}
	config.KeepShortDateTimeZone = true
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	expectedTime := time.Date(2006, 03, 04, 0, 0, 0, 0, config.TimeLocation)
	assert.Equal(t, expectedTime, target.Time)
	// Teardown
	teardown()
}

func TestParseLine_ShortDateDontKeepTimeZone(t *testing.T) {
	// Arrange
	input := "T|1|20060304"
	target := ShortDateRecord{}
	config.KeepShortDateTimeZone = false
	// Act
	_, err := ParseLine(input, &target, createStructAnnotation("T"), 1, config)
	// Assert
	assert.Nil(t, err)
	expectedTime := time.Date(2006, 03, 04, 0, 0, 0, 0, config.TimeLocation).UTC()
	assert.Equal(t, expectedTime, target.Time)
	// Teardown
	teardown()
}

func TestSplitStringWithEscape_NoEscape(t *testing.T) {
	// Arrange
	input := "no&|split"
	// Act
	result := splitStringWithEscape(input, config.Delimiters.Field, config.Delimiters.Escape)
	// Assert
	assert.Len(t, result, 1)
	assert.Equal(t, "no&|split", result[0])
}

func TestSplitStringWithEscape_Mixed(t *testing.T) {
	// Arrange
	input := "no&^split^second&&^third"
	// Act
	result := splitStringWithEscape(input, config.Delimiters.Component, config.Delimiters.Escape)
	// Assert
	assert.Len(t, result, 3)
	assert.Equal(t, "no&^split", result[0])
	assert.Equal(t, "second&&", result[1])
	assert.Equal(t, "third", result[2])
}

func TestSplitStringWithEscape_EmptyFields(t *testing.T) {
	// Arrange
	input := "first||third"
	// Act
	result := splitStringWithEscape(input, config.Delimiters.Field, config.Delimiters.Escape)
	// Assert
	assert.Len(t, result, 3)
	assert.Equal(t, "first", result[0])
	assert.Equal(t, "", result[1])
	assert.Equal(t, "third", result[2])
}

func TestSplitStringWithEscape_EmptyInput(t *testing.T) {
	// Arrange
	input := ""
	// Act
	result := splitStringWithEscape(input, config.Delimiters.Field, config.Delimiters.Escape)
	// Assert
	assert.Len(t, result, 0)
}

func TestSplitStringWithEscape_Unicode(t *testing.T) {
	// Arrange
	input := "first|őáúäö&||third"
	// Act
	result := splitStringWithEscape(input, config.Delimiters.Field, config.Delimiters.Escape)
	// Assert
	assert.Len(t, result, 3)
	assert.Equal(t, "first", result[0])
	assert.Equal(t, "őáúäö&|", result[1])
	assert.Equal(t, "third", result[2])
}

func TestFilterEscapeChars_Delimiters(t *testing.T) {
	// Arrange
	input := "escaped&| and&^ and&&"
	// Act
	result := filterStringEscapeChars(input, config.Delimiters.Escape)
	// Assert
	assert.Equal(t, "escaped| and^ and&", result)
}

func TestFilterEscapeChars_Multiple(t *testing.T) {
	// Arrange
	input := "esc&&&|ape"
	// Act
	result := filterStringEscapeChars(input, config.Delimiters.Escape)
	// Assert
	assert.Equal(t, "esc&|ape", result)
}

// Note: this should be invalid, but for simplicity we allow it by escaping the nothing
func TestFilterEscapeChars_AtTheEnd(t *testing.T) {
	// Arrange
	input := "escape&"
	// Act
	result := filterStringEscapeChars(input, config.Delimiters.Escape)
	// Assert
	assert.Equal(t, "escape", result)
}

func TestFilterEscapeChars_Empty(t *testing.T) {
	// Arrange
	input := ""
	// Act
	result := filterStringEscapeChars(input, config.Delimiters.Escape)
	// Assert
	assert.Equal(t, "", result)
}

func TestFilterEscapeChars_Unicode(t *testing.T) {
	// Arrange
	input := "őáúäö&|"
	// Act
	result := filterStringEscapeChars(input, config.Delimiters.Escape)
	// Assert
	assert.Equal(t, "őáúäö|", result)
}
