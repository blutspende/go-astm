package functions

import (
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Record struct {
	First  string `astm:"3"`
	Second string `astm:"4"`
}

type SingleRecordStruct struct {
	FirstRecord Record `astm:"R"`
}

func TestParseStruct_SingleLineStruct(t *testing.T) {
	// Arrange
	input := []string{
		"R|1|first|second",
	}
	target := SingleRecordStruct{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.FirstRecord.First)
	assert.Equal(t, "second", target.FirstRecord.Second)
}

type RecordArrayStruct struct {
	RecordArray []Record `astm:"R"`
}

func TestParseStruct_RecordArrayStruct(t *testing.T) {
	// Arrange
	input := []string{
		"R|1|first1|second1",
		"R|2|first2|second2",
	}
	target := RecordArrayStruct{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, target.RecordArray, 2)
	assert.Equal(t, "first1", target.RecordArray[0].First)
	assert.Equal(t, "second1", target.RecordArray[0].Second)
	assert.Equal(t, "first2", target.RecordArray[1].First)
	assert.Equal(t, "second2", target.RecordArray[1].Second)
}

type RecordType1 struct {
	First  string `astm:"3"`
	Second string `astm:"4"`
}
type RecordType2 struct {
	First  string `astm:"3"`
	Second string `astm:"4"`
}
type CompositeRecordStruct struct {
	Record1 RecordType1 `astm:"F"`
	Record2 RecordType2 `astm:"S"`
}
type CompositeMessage struct {
	CompositeRecordStruct CompositeRecordStruct
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
	err := ParseStruct(input, &target, &lineIndex, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "r1 first", target.CompositeRecordStruct.Record1.First)
	assert.Equal(t, "r1 second", target.CompositeRecordStruct.Record1.Second)
	assert.Equal(t, "r2 first", target.CompositeRecordStruct.Record2.First)
	assert.Equal(t, "r2 second", target.CompositeRecordStruct.Record2.Second)
}

type CompositeArrayMessage struct {
	CompositeRecordArray []CompositeRecordStruct
}

func TestParseStruct_CompositeArrayMessage(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|a1 r1 first|a1 r1 second",
		"S|1|a1 r2 first|a1 r2 second",
		"F|1|a2 r1 first|a2 r1 second",
		"S|1|a2 r2 first|a2 r2 second",
	}
	target := CompositeArrayMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 0, config)
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

func TestParseStruct_UnexpectedLineTypeError(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|r1 first|r1 second",
		"U|1|r2 first|r2 second",
	}
	target := CompositeMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 0, config)
	// Assert
	assert.Error(t, err, errmsg.LineParsing_ErrLineTypeNameMismatch)
}

func TestParseStruct_LinesDepletedError(t *testing.T) {
	// Arrange
	input := []string{
		"F|1|r1 first|r1 second",
	}
	target := CompositeMessage{}
	lineIndex := 0
	// Act
	err := ParseStruct(input, &target, &lineIndex, 0, config)
	// Assert
	assert.Error(t, err, errmsg.StructureParsing_ErrInputLinesDepleted)
}

// Note: config is setup in line parser test
// TODO: separate test packages better
