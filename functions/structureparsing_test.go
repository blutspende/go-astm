package functions

import (
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseStruct_SingleLineStruct(t *testing.T) {
	// Arrange
	input := []string{
		"R|1|first|second|third",
	}
	target := SingleRecordStruct{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.FirstRecord.First)
	assert.Equal(t, "second", target.FirstRecord.Second)
	assert.Equal(t, "third", target.FirstRecord.Third)
}

func TestParseStruct_RecordArrayStruct(t *testing.T) {
	// Arrange
	input := []string{
		"R|1|first1|second1|third1",
		"R|2|first2|second2|third2",
	}
	target := RecordArrayStruct{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, target.RecordArray, 2)
	assert.Equal(t, "first1", target.RecordArray[0].First)
	assert.Equal(t, "second1", target.RecordArray[0].Second)
	assert.Equal(t, "third1", target.RecordArray[0].Third)
	assert.Equal(t, "first2", target.RecordArray[1].First)
	assert.Equal(t, "second2", target.RecordArray[1].Second)
	assert.Equal(t, "third2", target.RecordArray[1].Third)
}

func TestParseStruct_CompositeMessage(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|r1 first|r1 second",
		"S|1|r2 first|r2 second",
	}
	target := CompositeMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "r1 first", target.CompositeRecordStruct.Record1.First)
	assert.Equal(t, "r1 second", target.CompositeRecordStruct.Record1.Second)
	assert.Equal(t, "r2 first", target.CompositeRecordStruct.Record2.First)
	assert.Equal(t, "r2 second", target.CompositeRecordStruct.Record2.Second)
}

func TestParseStruct_CompositeArrayMessage(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|a1 r1 first|a1 r1 second",
		"S|1|a1 r2 first|a1 r2 second",
		"F|2|a2 r1 first|a2 r1 second",
		"S|1|a2 r2 first|a2 r2 second",
	}
	target := CompositeArrayMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, target.CompositeRecordArray, 2)
	assert.Equal(t, "a1 r1 first", target.CompositeRecordArray[0].Record1.First)
	assert.Equal(t, "a1 r1 second", target.CompositeRecordArray[0].Record1.Second)
	assert.Equal(t, "a1 r2 first", target.CompositeRecordArray[0].Record2.First)
	assert.Equal(t, "a1 r2 second", target.CompositeRecordArray[0].Record2.Second)
	assert.Equal(t, "a2 r1 first", target.CompositeRecordArray[1].Record1.First)
	assert.Equal(t, "a2 r1 second", target.CompositeRecordArray[1].Record1.Second)
	assert.Equal(t, "a2 r2 first", target.CompositeRecordArray[1].Record2.First)
	assert.Equal(t, "a2 r2 second", target.CompositeRecordArray[1].Record2.Second)
}

func TestParseStruct_OptionalMessage(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first|second",
		"T|1|first|second",
	}
	target := OptionalMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Equal(t, "second", target.First.Second)
	assert.Equal(t, "", target.Optional.First)
	assert.Equal(t, "", target.Optional.Second)
	assert.Equal(t, "first", target.Third.First)
	assert.Equal(t, "second", target.Third.Second)
}

func TestParseStruct_OptionalArrayMessage(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first|second",
		"L|1|first|second",
	}
	target := OptionalArrayMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Equal(t, "second", target.First.Second)
	assert.Len(t, target.Optional, 0)
	assert.Equal(t, "first", target.Last.First)
	assert.Equal(t, "second", target.Last.Second)
}
func TestParseStruct_OptionalArrayMessageWithData(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first|second",
		"A|1|first|second",
		"A|2|first|second",
		"L|1|first|second",
	}
	target := OptionalArrayMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Equal(t, "second", target.First.Second)
	assert.Len(t, target.Optional, 2)
	assert.Equal(t, "first", target.Last.First)
	assert.Equal(t, "second", target.Last.Second)
}
func TestParseStruct_OptionalArrayAtTheEndMessageWithMissingOptionalData(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first|second",
	}
	target := OptionalArrayAtTheEndMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Equal(t, "second", target.First.Second)
	assert.Len(t, target.Optional, 0)
}
func TestParseStruct_OptionalAtTheEndMessageWithMissingOptionalData(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first|second",
	}
	target := OptionalAtTheEndMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Equal(t, "second", target.First.Second)
	assert.Equal(t, "", target.Optional.First)
	assert.Equal(t, "", target.Optional.Second)
}
func TestParseStruct_UnexpectedLineTypeError(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|r1 first|r1 second",
		"U|1|r2 first|r2 second",
	}
	target := CompositeMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.EqualError(t, err, errmsg.LineParsing_ErrLineTypeNameMismatch.Error())
}
func TestParseStruct_LinesDepletedError(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|r1 first|r1 second",
	}
	target := CompositeMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.EqualError(t, err, errmsg.StructureParsing_ErrInputLinesDepleted.Error())
}
