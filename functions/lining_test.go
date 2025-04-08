package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceLinesLf(t *testing.T) {
	// Arrange
	input := "first\nsecond"
	// Act
	lines, err := SliceLines(input)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLinesCr(t *testing.T) {
	// Arrange
	input := "first\rsecond"
	// Act
	lines, err := SliceLines(input)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLinesLfCr(t *testing.T) {
	// Arrange
	input := "first\n\rsecond"
	// Act
	lines, err := SliceLines(input)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLinesCrLf(t *testing.T) {
	// Arrange
	input := "first\r\nsecond"
	// Act
	lines, err := SliceLines(input)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}

func TestSliceLinesLfSpace(t *testing.T) {
	// Arrange
	input := "first\nsecond "
	// Act
	lines, err := SliceLines(input)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLinesLfEmptyLine(t *testing.T) {
	// Arrange
	input := "first\nsecond\n"
	// Act
	lines, err := SliceLines(input)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLinesComplex(t *testing.T) {
	// Arrange
	input := "first \n\rsecond \n\r\n\r  third\n\r"
	// Act
	lines, err := SliceLines(input)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 3)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
	assert.Equal(t, "third", lines[2])
}
func TestSliceLinesEmpty(t *testing.T) {
	// Arrange
	input := ""
	// Act
	_, err := SliceLines(input)
	// Assert
	assert.Error(t, err, errmsg.Parsing_ErrNotEnoughLines)
}
func TestSliceLinesOneLine(t *testing.T) {
	// Arrange
	input := "first"
	// Act
	_, err := SliceLines(input)
	// Assert
	assert.Error(t, err, errmsg.Parsing_ErrNotEnoughLines)
}
func TestSliceLinesInvalid(t *testing.T) {
	// Arrange
	input := "first\r\r\nsecond"
	// Act
	_, err := SliceLines(input)
	// Assert
	assert.Error(t, err, errmsg.Parsing_ErrInvalidLinebreak)
}

func TestBuildLinesLF(t *testing.T) {
	// Arrange
	input := []string{"first", "second"}
	// Act
	output := BuildLines(input, constants.LF)
	// Assert
	assert.Equal(t, "first\nsecond", output)
}
func TestBuildLinesLFCR(t *testing.T) {
	// Arrange
	input := []string{"first", "second"}
	// Act
	output := BuildLines(input, constants.LFCR)
	// Assert
	assert.Equal(t, "first\n\rsecond", output)
}
