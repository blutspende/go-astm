package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/blutspende/go-astm/v2/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type SimpleRecord struct {
	First  string `astm:"3"`
	Second string `astm:"4"`
	Third  string `astm:"5"`
}

func TestParseLine_SimpleRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second|third"
	target := SimpleRecord{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "second", target.Second)
	assert.Equal(t, "third", target.Third)
}

type MultitypeRecord struct {
	String    string    `astm:"3"`
	Int       int       `astm:"4"`
	Float32   float32   `astm:"5"`
	Float64   float64   `astm:"6"`
	ShortTime time.Time `astm:"7"`
	LongTime  time.Time `astm:"8"`
}

func TestParseLine_MultitypeRecord(t *testing.T) {
	// Arrange
	input := "T|1|string|4|3.14|3.1415926|20060102|20060102150405"
	target := MultitypeRecord{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, string("string"), target.String)
	assert.Equal(t, int(4), target.Int)
	assert.Equal(t, float32(3.14), target.Float32)
	assert.Equal(t, float64(3.1415926), target.Float64)
	expectedShortTime := time.Date(2006, 1, 2, 0, 0, 0, 0, config.TimeLocation)
	assert.Equal(t, expectedShortTime, target.ShortTime)
	expectedLongTime := time.Date(2006, 1, 2, 15, 04, 05, 0, config.TimeLocation)
	assert.Equal(t, expectedLongTime, target.LongTime)
}

type ComponentedRecord struct {
	First       string `astm:"3"`
	SecondComp1 string `astm:"4.1"`
	SecondComp2 string `astm:"4.2"`
	ThirdComp1  string `astm:"5.1"`
	ThirdComp2  string `astm:"5.2"`
	ThirdComp3  string `astm:"5.3"`
}

func TestParseLine_ComponentedRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second1^second2|third1^third2^third3"
	target := ComponentedRecord{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "second1", target.SecondComp1)
	assert.Equal(t, "second2", target.SecondComp2)
	assert.Equal(t, "third1", target.ThirdComp1)
	assert.Equal(t, "third2", target.ThirdComp2)
	assert.Equal(t, "third3", target.ThirdComp3)
}

type ArrayRecord struct {
	First string   `astm:"3"`
	Array []string `astm:"4"`
}

func TestParseLine_ArrayRecord(t *testing.T) {
	// Arrange
	input := "T|1|first|second1\\second2\\second3"
	target := ArrayRecord{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
	assert.Len(t, target.Array, 3)
	assert.Equal(t, "second1", target.Array[0])
	assert.Equal(t, "second2", target.Array[1])
	assert.Equal(t, "second3", target.Array[2])
}

type HeaderRecord struct {
	First string `astm:"3"`
}

func TestParseLine_HeaderRecord(t *testing.T) {
	// Arrange
	input := "H|\\^&|first"
	target := HeaderRecord{}
	// Act
	err := ParseLine(input, &target, "H", 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
}

type HeaderDelimiterChange struct {
	First string   `astm:"3"`
	Array []string `astm:"4"`
	Comp1 string   `astm:"5.1"`
	Comp2 string   `astm:"5.2"`
}

func TestParseLine_HeaderDelimiterChange(t *testing.T) {
	// Arrange
	input := "H/!*%/first/second1!second2/third1*third2"
	target := HeaderDelimiterChange{}
	// Act
	err := ParseLine(input, &target, "H", 0, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
	assert.Len(t, target.Array, 2)
	assert.Equal(t, "second1", target.Array[0])
	assert.Equal(t, "second2", target.Array[1])
	assert.Equal(t, "third1", target.Comp1)
	assert.Equal(t, "third2", target.Comp2)
	// Tear down
	config.Delimiters = models.DefaultDelimiters
}

func TestParseLine_MissingData(t *testing.T) {
	// Arrange
	input := "T|1|first||third"
	target := SimpleRecord{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "", target.Second)
	assert.Equal(t, "third", target.Third)
}

func TestParseLine_MissingDataAtTheEnd(t *testing.T) {
	// Arrange
	input := "T|1|first"
	target := SimpleRecord{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "first", target.First)
	assert.Equal(t, "", target.Second)
	assert.Equal(t, "", target.Third)
}

func TestParseLine_EmptyInput(t *testing.T) {
	// Arrange
	input := ""
	target := SimpleRecord{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Error(t, err, errmsg.LineParsing_ErrEmptyInput)
}

func TestParseLine_NotEnoughFields(t *testing.T) {
	// Arrange
	input := "T"
	target := SimpleRecord{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Error(t, err, errmsg.LineParsing_ErrMandatoryInputFieldsMissing)
}

type MissingRequiredField struct {
	First  string `astm:"3"`
	Second string `astm:"4,required"`
	Third  string `astm:"5"`
}

func TestParseLine_MissingRequiredField(t *testing.T) {
	// Arrange
	input := "T|1|first||third"
	target := MissingRequiredField{}
	// Act
	err := ParseLine(input, &target, "T", 1, config)
	// Assert
	assert.Error(t, err, errmsg.LineParsing_ErrRequiredFieldIsEmpty)
}

// Setup mock data for every test
var config *models.Configuration

func TestMain(m *testing.M) {
	Delimiters := models.DefaultDelimiters
	TimeLocation, _ := time.LoadLocation(string(constants.TIMEZONE_EUROPE_BERLIN))

	config = &models.Configuration{
		Delimiters:   Delimiters,
		TimeLocation: TimeLocation,
	}

	// Run all tests
	m.Run()
}
