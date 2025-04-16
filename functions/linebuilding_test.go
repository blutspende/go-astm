package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
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
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|second|third", result)
}

func TestBuildLine_MultitypeRecord(t *testing.T) {
	// Arrange
	source := MultitypeRecord{
		String:  "string",
		Int:     3,
		Float32: 3.14,
		Float64: 3.14159265,
		Date:    time.Date(2006, 1, 2, 0, 0, 0, 0, config.Internal.TimeLocation),
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|string|3|3.14|3.14159265|20060102", result)
}

func TestBuildLine_MultitypeLengthRecord(t *testing.T) {
	// Arrange
	source := MultitypeLengthRecord{
		FloatFull: 3.14159265,
		FloatFix3: 3.14159265,
		FloatFix0: 3.14159265,
		LongDate:  time.Date(2006, 1, 2, 15, 04, 05, 0, config.Internal.TimeLocation),
		ShortDate: time.Date(2006, 1, 2, 15, 04, 05, 0, config.Internal.TimeLocation),
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|3.14159265|3.142|3|20060102150405|20060102", result)
}

func TestBuildLine_MultitypeLengthRecordRound(t *testing.T) {
	// Arrange
	source := MultitypeLengthRecord{
		FloatFull: 3.77777,
		FloatFix3: 3.77777,
		FloatFix0: 3.77777,
	}
	config.RoundFixedNumbers = true
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|3.77777|3.778|4||", result)
	// Tear down
	teardown()
}

func TestBuildLine_MultitypeLengthRecordTruncate(t *testing.T) {
	// Arrange
	source := MultitypeLengthRecord{
		FloatFull: 3.77777,
		FloatFix3: 3.77777,
		FloatFix0: 3.77777,
	}
	config.RoundFixedNumbers = false
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|3.77777|3.777|3||", result)
	// Tear down
	teardown()
}

func TestBuildLine_MultitypeRecordEmpty(t *testing.T) {
	// Arrange
	source := MultitypeRecord{}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1||0|0|0|", result)
}

func TestBuildLine_MultitypeLengthRecordEmpty(t *testing.T) {
	// Arrange
	source := MultitypeLengthRecord{}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|0|0.000|0||", result)
}

func TestBuildLine_MultitypePointerRecord(t *testing.T) {
	// Arrange
	String := "string"
	Int := 3
	Float32 := float32(3.14)
	Float64 := float64(3.14159265)
	Date := time.Date(2006, 1, 2, 0, 0, 0, 0, config.Internal.TimeLocation)
	source := MultitypePointerRecord{
		String:  &String,
		Int:     &Int,
		Float32: &Float32,
		Float64: &Float64,
		Date:    &Date,
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|string|3|3.14|3.14159265|20060102", result)
}

func TestBuildLine_MultitypePointerRecordEmpty(t *testing.T) {
	// Arrange
	source := MultitypePointerRecord{}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|||||", result)
}

func TestBuildLine_UnorderedRecord(t *testing.T) {
	// Arrange
	source := UnorderedRecord{
		First:  "first",
		Second: "second",
		Third:  "third",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
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
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first||third", result)
}

func TestBuildLine_MissingDataAtEndStandardNotation(t *testing.T) {
	// Arrange
	source := SimpleRecord{
		First:  "first",
		Second: "",
		Third:  "",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
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
	result, err := BuildLine(source, "T", 1, config)
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
	result, err := BuildLine(source, "D", 3, config)
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
	result, err := BuildLine(source, "H", 0, config)
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
	result, err := BuildLine(source, "H", 0, config)
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
	result, err := BuildLine(source, "H", 0, config)
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
	result, err := BuildLine(source, "T", 1, config)
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
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|second1^second2|third1^third2^third3", result)
}

func TestBuildLine_EnumRecord(t *testing.T) {
	// Arrange
	source := EnumRecord{
		Enum: "enum",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|enum", result)
}

func TestBuildLine_ReservedFieldRecord(t *testing.T) {
	// Arrange
	source := ReservedFieldRecord{}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Error(t, err, errmsg.LineBuilding_ErrReservedFieldPosReference)
	assert.Equal(t, "", result)
}
