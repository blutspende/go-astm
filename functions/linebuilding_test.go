package functions

import (
	"github.com/blutspende/go-astm/v2/constants/astmconst"
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
		Date:    time.Date(2006, 1, 2, 0, 0, 0, 0, config.TimeLocation),
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
		LongDate:  time.Date(2006, 1, 2, 15, 04, 05, 0, config.TimeLocation),
		ShortDate: time.Date(2006, 1, 2, 15, 04, 05, 0, config.TimeLocation),
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
	// Teardown
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
	// Teardown
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
	Date := time.Date(2006, 1, 2, 0, 0, 0, 0, config.TimeLocation)
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
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first", result)
	// Teardown
	teardown()
}
func TestBuildLine_MissingDataAtEndWithComponentsShortNotation(t *testing.T) {
	// Arrange
	source := ComponentedRecord{
		First: "first",
	}
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first", result)
	// Teardown
	teardown()
}
func TestBuildLine_DifferentHeaderAndSequence(t *testing.T) {
	// Arrange
	source := SimpleRecord{
		First:  "first",
		Second: "second",
		Third:  "third",
	}
	// Act
	result, err := BuildLine(source, "D", 3, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "D|3|first|second|third", result)
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
	config.Delimiters.Field = "/"
	config.Delimiters.Repeat = "!"
	config.Delimiters.Component = "*"
	config.Delimiters.Escape = "%"
	// Act
	result, err := BuildLine(source, "H", 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "H/!*%/first", result)
	// Teardown
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
	config.Delimiters.Field = "/"
	config.Delimiters.Repeat = "!"
	config.Delimiters.Component = "*"
	config.Delimiters.Escape = "%"
	// Act
	result, err := BuildLine(source, "H", 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "H/!*%/first/second1!second2/third1*third2", result)
	// Teardown
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
func TestBuildLine_ComponentedRecordShortNotation(t *testing.T) {
	// Arrange
	source := ComponentedRecord{
		First:       "first",
		SecondComp1: "second1",
		SecondComp2: "second2",
		ThirdComp1:  "third1",
		ThirdComp2:  "third2",
		ThirdComp3:  "third3",
	}
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|second1^second2|third1^third2^third3", result)
	// Teardown
	teardown()
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
func TestBuildLine_SparseFieldRecord(t *testing.T) {
	// Arrange
	source := SparseFieldRecord{
		Field3: "field3",
		Field5: "field5",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|field3||field5", result)
}
func TestBuildLine_SubstructureRecord(t *testing.T) {
	// Arrange
	source := SubstructureRecord{
		First: "first",
		Second: SubstructureField{
			FirstComponent:  "firstComponent",
			SecondComponent: "secondComponent",
			ThirdComponent:  "thirdComponent",
		},
		Third: "third",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|firstComponent^secondComponent^thirdComponent|third", result)
}
func TestBuildLine_SubstructureRecordMissingData(t *testing.T) {
	// Arrange
	source := SubstructureRecord{
		Second: SubstructureField{
			FirstComponent: "firstComponent",
		},
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1||firstComponent^^|", result)
}
func TestBuildLine_SubstructureRecordMissingDataShortNotation(t *testing.T) {
	// Arrange
	source := SubstructureRecord{
		Second: SubstructureField{
			FirstComponent: "firstComponent",
		},
	}
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1||firstComponent", result)
	// Teardown
	teardown()
}
func TestBuildLine_SparseSubstructureRecord(t *testing.T) {
	// Arrange
	source := SparseSubstructureRecord{
		First: "first",
		Second: SparseSubstructureField{
			Component1: "component1",
			Component3: "component3",
			Component6: "component6",
		},
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|component1^^component3^^^component6", result)
}
func TestBuildLine_SubstructureArrayRecord(t *testing.T) {
	// Arrange
	source := SubstructureArrayRecord{
		First: "first",
		Second: []SubstructureField{
			SubstructureField{
				FirstComponent:  "r1c1",
				SecondComponent: "r1c2",
				ThirdComponent:  "r1c3",
			},
			SubstructureField{
				FirstComponent:  "r2c1",
				SecondComponent: "r2c2",
				ThirdComponent:  "r2c3",
			},
		},
		Third: "third",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|r1c1^r1c2^r1c3\\r2c1^r2c2^r2c3|third", result)
}

func TestBuildLine_TimeLineTimeZone(t *testing.T) {
	// Arrange
	source := TimeRecord{
		Time: time.Date(2006, 03, 06, 16, 44, 29, 0, config.TimeLocation).UTC(),
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|20060306164429", result)
}

func TestBuildLine_WrongComponentOrder(t *testing.T) {
	// Arrange
	source := WrongComponentOrderRecord{
		First: "first",
		Comp2: "comp2",
		Comp1: "comp1",
		Comp3: "comp3",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|first|comp1^comp2^comp3", result)
}

func TestBuildLine_WrongComponentPlacement(t *testing.T) {
	// Arrange
	source := WrongComponentPlacementRecord{
		Field1: "field1",
		Comp1:  "comp1",
		Field2: "field2",
		Comp2:  "comp2",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|field1|comp1^comp2|field2", result)
}

func TestBuildLine_MultipleWrongComponentPlacement(t *testing.T) {
	// Arrange
	source := MultipleWrongComponentPlacementRecord{
		Field3: "field3",
		Comp41: "comp41",
		Field5: "field5",
		Comp62: "comp62",
		Comp42: "comp42",
		Field7: "field7",
		Comp61: "comp61",
		Field8: "field8",
	}
	// Act
	result, err := BuildLine(source, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "T|1|field3|comp41^comp42|field5|comp61^comp62|field7|field8", result)
}
