package functions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildStruct_SingleLineStruct(t *testing.T) {
	// Arrange
	source := SingleRecordStruct{
		FirstRecord: ThreeFieldRecord{
			First:  "first",
			Second: "second",
			Third:  "third",
		},
	}
	// Act
	result, err := BuildStruct(source, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "R|1|first|second|third", result[0])
}

func TestBuildStruct_RecordArrayStruct(t *testing.T) {
	// Arrange
	source := RecordArrayStruct{
		RecordArray: []ThreeFieldRecord{
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
	result, err := BuildStruct(source, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "R|1|first1|second1|third1", result[0])
	assert.Equal(t, "R|2|first2|second2|third2", result[1])
}

func TestBuildStruct_CompositeMessage(t *testing.T) {
	// Arrange
	source := CompositeMessage{
		CompositeRecordStruct: CompositeRecordStruct{
			Record1: RecordType1{
				First:  "r1 first",
				Second: 2,
			},
			Record2: RecordType2{
				First:  1,
				Second: "r2 second",
			},
		},
	}
	// Act
	result, err := BuildStruct(source, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "F|1|r1 first|2", result[0])
	assert.Equal(t, "S|1|1|r2 second", result[1])
}

func TestBuildStruct_CompositeArrayMessage(t *testing.T) {
	// Arrange
	source := CompositeArrayMessage{
		CompositeRecordArray: []CompositeRecordStruct{
			{
				Record1: RecordType1{
					First:  "a1 r1 first",
					Second: 112,
				},
				Record2: RecordType2{
					First:  121,
					Second: "a1 r2 second",
				},
			},
			{
				Record1: RecordType1{
					First:  "a2 r1 first",
					Second: 212,
				},
				Record2: RecordType2{
					First:  221,
					Second: "a2 r2 second",
				},
			},
		},
	}
	// Act
	result, err := BuildStruct(source, 1, 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, result, 4)
	assert.Equal(t, "F|1|a1 r1 first|112", result[0])
	assert.Equal(t, "S|1|121|a1 r2 second", result[1])
	assert.Equal(t, "F|2|a2 r1 first|212", result[2])
	assert.Equal(t, "S|1|221|a2 r2 second", result[3])
}
