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
	expectedShortTime := time.Date(2006, 1, 2, 0, 0, 0, 0, config.Internal.TimeLocation)
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
	// Tear down
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

func TestParseLine_NotEnoughFields(t *testing.T) {
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
	assert.Error(t, err, errmsg.LineParsing_ErrRequiredFieldIsEmpty)
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
	// Tear down
	teardown()
}
