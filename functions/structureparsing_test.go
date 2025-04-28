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
		"F|1|r1 first|12",
		"S|1|21|r2 second",
	}
	target := CompositeMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "r1 first", target.CompositeRecordStruct.Record1.First)
	assert.Equal(t, 12, target.CompositeRecordStruct.Record1.Second)
	assert.Equal(t, 21, target.CompositeRecordStruct.Record2.First)
	assert.Equal(t, "r2 second", target.CompositeRecordStruct.Record2.Second)
}

func TestParseStruct_CompositeArrayMessage(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|a1 r1 first|112",
		"S|1|121|a1 r2 second",
		"F|2|a2 r1 first|212",
		"S|1|221|a2 r2 second",
	}
	target := CompositeArrayMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, target.CompositeRecordArray, 2)
	assert.Equal(t, "a1 r1 first", target.CompositeRecordArray[0].Record1.First)
	assert.Equal(t, 112, target.CompositeRecordArray[0].Record1.Second)
	assert.Equal(t, 121, target.CompositeRecordArray[0].Record2.First)
	assert.Equal(t, "a1 r2 second", target.CompositeRecordArray[0].Record2.Second)
	assert.Equal(t, "a2 r1 first", target.CompositeRecordArray[1].Record1.First)
	assert.Equal(t, 212, target.CompositeRecordArray[1].Record1.Second)
	assert.Equal(t, 221, target.CompositeRecordArray[1].Record2.First)
	assert.Equal(t, "a2 r2 second", target.CompositeRecordArray[1].Record2.Second)
}

func TestParseStruct_OptionalMessage(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first",
		"T|1|first",
	}
	target := OptionalMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Equal(t, "", target.Optional.First)
	assert.Equal(t, "first", target.Third.First)
}

func TestParseStruct_OptionalArrayMessage(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first",
		"L|1|first",
	}
	target := OptionalArrayMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Len(t, target.Optional, 0)
	assert.Equal(t, "first", target.Last.First)
}
func TestParseStruct_OptionalArrayMessageWithData(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first",
		"A|1|first",
		"A|2|first",
		"L|1|first",
	}
	target := OptionalArrayMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Len(t, target.Optional, 2)
	assert.Equal(t, "first", target.Last.First)
}
func TestParseStruct_OptionalArrayAtTheEndMessageWithMissingOptionalData(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first",
	}
	target := OptionalArrayAtTheEndMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Len(t, target.Optional, 0)
}
func TestParseStruct_OptionalAtTheEndMessageWithMissingOptionalData(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|first",
	}
	target := OptionalAtTheEndMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First.First)
	assert.Equal(t, "", target.Optional.First)
}
func TestParseStruct_UnexpectedLineTypeError(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|r1 first|12",
		"U|1|21|r2 second",
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
		"F|1|r1 first|12",
	}
	target := CompositeMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.EqualError(t, err, errmsg.StructureParsing_ErrInputLinesDepleted.Error())
}
func TestParseStruct_SubnameMessage(t *testing.T) {
	// Arrange
	input := []string{
		"R|1|FIRST|r1 first|12",
		"R|1|SECOND|21|r2 second",
	}
	target := SubnameMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "FIRST", target.Record1.Subname)
	assert.Equal(t, "r1 first", target.Record1.First)
	assert.Equal(t, 12, target.Record1.Second)
	assert.Equal(t, "SECOND", target.Record2.Subname)
	assert.Equal(t, 21, target.Record2.First)
	assert.Equal(t, "r2 second", target.Record2.Second)
}
func TestParseStruct_SubnameArrayMessage(t *testing.T) {
	// Arrange
	input := []string{
		"R|1|FIRST|a1 first|12",
		"R|2|FIRST|a2 first|22",
		"R|1|SECOND|1|second",
	}
	target := SubnameArrayMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, target.Array, 2)
	assert.Equal(t, "FIRST", target.Array[0].Subname)
	assert.Equal(t, "a1 first", target.Array[0].First)
	assert.Equal(t, 12, target.Array[0].Second)
	assert.Equal(t, "FIRST", target.Array[1].Subname)
	assert.Equal(t, "a2 first", target.Array[1].First)
	assert.Equal(t, 22, target.Array[1].Second)
	assert.Equal(t, "SECOND", target.Record2.Subname)
	assert.Equal(t, 1, target.Record2.First)
	assert.Equal(t, "second", target.Record2.Second)
}
func TestParseStruct_SubnameMultiArrayMessage(t *testing.T) {
	// Arrange
	input := []string{
		"R|1|FIRST|1 a1 first|112",
		"R|2|FIRST|1 a2 first|122",
		"R|3|SECOND|221|2 a1 second",
		"R|4|SECOND|222|2 a2 second",
	}
	target := SubnameMultiArrayMessage{}
	lineIndex := 0
	config.EnforceSequenceNumberCheck = false
	// Act
	err := ParseStruct(input, &target, &lineIndex, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, target.Array1, 2)
	assert.Equal(t, "FIRST", target.Array1[0].Subname)
	assert.Equal(t, "1 a1 first", target.Array1[0].First)
	assert.Equal(t, 112, target.Array1[0].Second)
	assert.Equal(t, "FIRST", target.Array1[1].Subname)
	assert.Equal(t, "1 a2 first", target.Array1[1].First)
	assert.Equal(t, 122, target.Array1[1].Second)
	assert.Len(t, target.Array2, 2)
	assert.Equal(t, "SECOND", target.Array2[0].Subname)
	assert.Equal(t, 221, target.Array2[0].First)
	assert.Equal(t, "2 a1 second", target.Array2[0].Second)
	assert.Equal(t, "SECOND", target.Array2[1].Subname)
	assert.Equal(t, 222, target.Array2[1].First)
	assert.Equal(t, "2 a2 second", target.Array2[1].Second)
	// Teardown
	teardown()
}
