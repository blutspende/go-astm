package functions

import (
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Note: structures come from functions_test.go

func TestParseLine_SimpleRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second|third"
	target := SimpleRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
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
	nameOk, err := ParseLine(input, &target, "T", 1, config)
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
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "string", target.String)
	assert.Equal(t, 3, target.Int)
	assert.Equal(t, float32(3.14), target.Float32)
	assert.Equal(t, float64(3.14159265), target.Float64)
	expectedShortTime := time.Date(2006, 1, 1, 23, 0, 0, 0, time.UTC)
	assert.Equal(t, expectedShortTime, target.Date)
}

func TestParseLine_ComponentedRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second1^second2|third1^third2^third3"
	target := ComponentedRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
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
	nameOk, err := ParseLine(input, &target, "T", 1, config)
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
	nameOk, err := ParseLine(input, &target, "H", 0, config)
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
	nameOk, err := ParseLine(input, &target, "H", 0, config)
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
	target := SimpleRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "", target.Second)
	assert.Equal(t, "third", target.Third)
}

func TestParseLine_MissingDataAtTheEnd(t *testing.T) {
	// Arrange
	input := "T|1|first"
	target := SimpleRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
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
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	assert.Equal(t, EnumString("enum"), target.Enum)
}

func TestParseLine_RecordTypeNameMismatch(t *testing.T) {
	// Arrange
	input := "W|1|first|second|third"
	target := SimpleRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.False(t, nameOk)
}

func TestParseLine_EmptyInput(t *testing.T) {
	// Arrange
	input := ""
	target := SimpleRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.False(t, nameOk)
	assert.Error(t, err, errmsg.LineParsing_ErrEmptyInput)
}

func TestParseLine_MandatoryFieldsMissing(t *testing.T) {
	// Arrange
	input := "T"
	target := SimpleRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.False(t, nameOk)
	assert.Error(t, err, errmsg.LineParsing_ErrMandatoryInputFieldsMissing)
}

func TestParseLine_MissingRequiredField(t *testing.T) {
	// Arrange
	input := "T|1|first||third"
	target := MissingRequiredField{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.True(t, nameOk)
	assert.Error(t, err, errmsg.LineParsing_ErrRequiredInputFieldMissing)
}

func TestParseLine_NotEnoughInputFields(t *testing.T) {
	// Arrange
	input := "T|1|first"
	target := MissingRequiredField{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.True(t, nameOk)
	assert.Error(t, err, errmsg.LineParsing_ErrRequiredInputFieldMissing)
}

func TestParseLine_SequenceNumberMismatch(t *testing.T) {
	// Arrange
	input := "T|2|first|second|third"
	target := SimpleRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.True(t, nameOk)
	assert.Error(t, err, errmsg.LineParsing_ErrSequenceNumberMismatch)
}

func TestParseLine_SequenceNumberMismatchWithoutEnforcing(t *testing.T) {
	// Arrange
	input := "T|2|first|second|third"
	target := SimpleRecord{}
	config.EnforceSequenceNumberCheck = false
	// Act
	_, err := ParseLine(input, &target, "T", 1, config)
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
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.True(t, nameOk)
	assert.Error(t, err, errmsg.LineParsing_ErrReservedFieldPosReference)
}

func TestParseLine_SubstructureRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|firstComponent^secondComponent^thirdComponent|third"
	target := SubstructureRecord{}
	// Act
	nameOk, err := ParseLine(input, &target, "T", 1, config)
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
	nameOk, err := ParseLine(input, &target, "T", 1, config)
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
	nameOk, err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.True(t, nameOk)
	expectedTime := time.Date(2006, 03, 06, 16, 44, 29, 0, config.Internal.TimeLocation).UTC()
	assert.Equal(t, expectedTime, target.Time)
}
