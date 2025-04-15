package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Note: structures come from functions_test.go

func TestBuildLine_SimpleRecord(t *testing.T) {
	// Arrange
	source := SimpleRecord{
		First:  "first",
		Second: "second",
		Third:  "third",
	}
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|second|third", result)
}

func TestBuildLine_MultitypeRecord(t *testing.T) {
	// Arrange
	source := MultitypeRecord{
		String:     "string",
		Int:        4,
		Float32:    3.14,
		Float64:    3.1415926,
		Float64Cut: 3.1415926,
		ShortTime:  time.Date(2006, 1, 2, 0, 0, 0, 0, config.Internal.TimeLocation),
		LongTime:   time.Date(2006, 1, 2, 15, 04, 05, 0, config.Internal.TimeLocation),
	}
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	//TODO: question should the digit cutting round or truncate?
	assert.Equal(t, "T|1|string|4|3.14|3.1415926|3.142|20060102|20060102150405", result)
}

func TestBuildLine_MultitypeEmptyRecord(t *testing.T) {
	// Arrange
	source := MultitypeRecord{}
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1||0|0|0|0.000||", result)
}

func TestBuildLine_UnorderedRecord(t *testing.T) {
	// Arrange
	source := UnorderedRecord{
		First:  "first",
		Second: "second",
		Third:  "third",
	}
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|second|third", result)
}

func TestBuildLine_MissingData(t *testing.T) {
	// Arrange
	source := SimpleRecord{
		First:  "first",
		Second: "",
		Third:  "third",
	}
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first||third", result)
}

func TestBuildLine_MissingDataAtEndLongNotation(t *testing.T) {
	// Arrange
	source := SimpleRecord{
		First:  "first",
		Second: "",
		Third:  "",
	}
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first||", result)
}

func TestBuildLine_MissingDataAtEndShortNotation(t *testing.T) {
	// Arrange
	source := SimpleRecord{
		First:  "first",
		Second: "",
		Third:  "",
	}
	config.Notation = constants.NOTATION_SHORT
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first", result)
	// Tear down
	teardown()
}

func TestBuildLine_DifferentHeaderAndSequence(t *testing.T) {
	// Arrange
	source := SimpleRecord{
		First:  "first",
		Second: "second",
		Third:  "third",
	}
	config.Notation = constants.NOTATION_SHORT
	// Act
	result, err := BuildLine(&source, "D", 3, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "D|3|first|second|third", result)
	// Tear down
	teardown()
}

func TestBuildLine_HeaderRecord(t *testing.T) {
	// Arrange
	source := HeaderRecord{
		First: "first",
	}
	// Act
	result, err := BuildLine(&source, "H", 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "H|\\^&|first", result)
}

func TestBuildLine_HeaderRecordCustomDelimiters(t *testing.T) {
	// Arrange
	source := HeaderRecord{
		First: "first",
	}
	config.Internal.Delimiters.Field = "/"
	config.Internal.Delimiters.Repeat = "!"
	config.Internal.Delimiters.Component = "*"
	config.Internal.Delimiters.Escape = "%"
	// Act
	result, err := BuildLine(&source, "H", 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "H/!*%/first", result)
	// Tear down
	teardown()
}

func TestBuildLine_HeaderDelimiterChange(t *testing.T) {
	// Arrange
	source := HeaderDelimiterChange{
		First: "first",
		Array: []string{"second1", "second2"},
		Comp1: "third1",
		Comp2: "third2",
	}
	config.Internal.Delimiters.Field = "/"
	config.Internal.Delimiters.Repeat = "!"
	config.Internal.Delimiters.Component = "*"
	config.Internal.Delimiters.Escape = "%"
	// Act
	result, err := BuildLine(&source, "H", 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "H/!*%/first/second1!second2/third1*third2", result)
	// Tear down
	teardown()
}

func TestBuildLine_ArrayRecord(t *testing.T) {
	// Arrange
	source := ArrayRecord{
		First: "first",
		Array: []string{"second1", "second2", "second3"},
	}
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|second1\\second2\\second3", result)
}

func TestBuildLine_ComponentedRecord(t *testing.T) {
	// Arrange
	source := ComponentedRecord{
		First:       "first",
		SecondComp1: "second1",
		SecondComp2: "second2",
		ThirdComp1:  "third1",
		ThirdComp2:  "third2",
		ThirdComp3:  "third3",
	}
	// Act
	result, err := BuildLine(&source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|second1^second2|third1^third2^third3", result)
}
