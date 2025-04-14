package functions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildStruct_SingleLineStruct(t *testing.T) {
	// Arrange
	source := SingleRecordStruct{
		FirstRecord: SimpleRecord{
			First:  "first",
			Second: "second",
			Third:  "third",
		},
	}
	// Act
	result, err := BuildStruct(&source, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	expected := "R|1|first|second|third"
	assert.Equal(t, expected, result)
}

func TestBuildStruct_RecordArrayStruct(t *testing.T) {
	// Arrange
	source := RecordArrayStruct{
		RecordArray: []SimpleRecord{
			{
				First:  "first1",
				Second: "second1",
				Third:  "third1",
			},
			{
				First:  "first2",
				Second: "second2",
				Third:  "third2",
			},
		},
	}
	// Act
	result, err := BuildStruct(&source, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	expected := "R|1|first1|second1|third1\n"
	expected += "R|2|first2|second2|third2"
	assert.Equal(t, expected, result)
}

func TestBuildStruct_CompositeMessage(t *testing.T) {
	// Arrange
	source := CompositeMessage{
		CompositeRecordStruct: CompositeRecordStruct{
			Record1: RecordType1{
				First:  "r1 first",
				Second: "r1 second",
			},
			Record2: RecordType2{
				First:  "r2 first",
				Second: "r2 second",
			},
		},
	}
	// Act
	result, err := BuildStruct(&source, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	expected := "F|1|r1 first|r1 second\n"
	expected += "S|1|r2 first|r2 second"
	assert.Equal(t, expected, result)
}

func TestBuildStruct_CompositeArrayMessage(t *testing.T) {
	// Arrange
	source := CompositeArrayMessage{
		CompositeRecordArray: []CompositeRecordStruct{
			{
				Record1: RecordType1{
					First:  "a1 r1 first",
					Second: "a1 r1 second",
				},
				Record2: RecordType2{
					First:  "a1 r2 first",
					Second: "a1 r2 second",
				},
			},
			{
				Record1: RecordType1{
					First:  "a2 r1 first",
					Second: "a2 r1 second",
				},
				Record2: RecordType2{
					First:  "a2 r2 first",
					Second: "a2 r2 second",
				},
			},
		},
	}
	// Act
	result, err := BuildStruct(&source, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	expected := "F|1|a1 r1 first|a1 r1 second\n"
	expected += "S|1|a1 r2 first|a1 r2 second\n"
	expected += "F|2|a2 r1 first|a2 r1 second\n"
	expected += "S|1|a2 r2 first|a2 r2 second"
	assert.Equal(t, expected, result)
}
