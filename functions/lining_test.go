package functions

import (
	"github.com/blutspende/go-astm/v3/constants/astmconst"
	"github.com/blutspende/go-astm/v3/errmsg"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Default line separator config
func TestSliceLines_Lf(t *testing.T) {
	// Arrange
	input := "first\nsecond"
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLines_Cr(t *testing.T) {
	// Arrange
	input := "first\rsecond"
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLines_LfCr(t *testing.T) {
	// Arrange
	input := "first\n\rsecond"
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLines_CrLf(t *testing.T) {
	// Arrange
	input := "first\r\nsecond"
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}

func TestSliceLines_LfSpace(t *testing.T) {
	// Arrange
	input := "first\nsecond "
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLines_LfEmptyLine(t *testing.T) {
	// Arrange
	input := "first\nsecond\n"
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
}
func TestSliceLines_Complex(t *testing.T) {
	// Arrange
	input := "first \n\rsecond \n\r\n\r  third\n\r"
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 3)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "second", lines[1])
	assert.Equal(t, "third", lines[2])
}
func TestSliceLines_Empty(t *testing.T) {
	// Arrange
	input := ""
	// Act
	_, err := SliceLines(input, config)
	// Assert
	assert.EqualError(t, err, errmsg.Lining_ErrEmptyInput.Error())
}
func TestSliceLines_Invalid(t *testing.T) {
	// Arrange
	input := "first\r\r\nsecond"
	// Act
	_, err := SliceLines(input, config)
	// Assert
	assert.EqualError(t, err, errmsg.Lining_ErrInvalidLinebreak.Error())
}
func TestSliceLines_SingleLine(t *testing.T) {
	// Arrange
	input := "single line"
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 1)
	assert.Equal(t, "single line", lines[0])
}

// Explicit line separator config
func TestSliceLines_ExplicitCr(t *testing.T) {
	// Arrange
	input := "first\r\nsecond"
	config.LineSeparator = astmconst.CR
	config.AutoDetectLineSeparator = false
	// Act
	lines, err := SliceLines(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "first", lines[0])
	assert.Equal(t, "\nsecond", lines[1])
	// Teardown
	teardown()
}

// Lines building
func TestBuildLines_Default(t *testing.T) {
	// Arrange
	input := []string{"first", "second"}
	// Act
	output := BuildLines(input, config)
	// Assert
	assert.Equal(t, "first\nsecond", output)
}
func TestBuildLines_ExplicitLFCR(t *testing.T) {
	// Arrange
	input := []string{"first", "second"}
	config.LineSeparator = astmconst.LFCR
	config.AutoDetectLineSeparator = false
	// Act
	output := BuildLines(input, config)
	// Assert
	assert.Equal(t, "first\n\rsecond", output)
	// Teardown
	teardown()
}
