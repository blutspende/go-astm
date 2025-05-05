package functions

import (
	"github.com/blutspende/go-astm/enums/encoding"
	"github.com/blutspende/go-astm/errmsg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoding_ConvertFromEncodingToUtf8(t *testing.T) {
	// Arrange
	input := []byte("űúőéóüöáßäüöë")
	config.Encoding = encoding.UTF8
	// Act
	result, err := ConvertFromEncodingToUtf8(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "űúőéóüöáßäüöë", result)
	// Teardown
	teardown()
}

func TestEncoding_ConvertFromEncodingToUtf8Win1252(t *testing.T) {
	// Arrange
	input := []byte{0xDF, 0xE4, 0xFC, 0xF6, 0xEB}
	config.Encoding = encoding.Windows1252
	// Act
	result, err := ConvertFromEncodingToUtf8(input, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "ßäüöë", result)
	// Teardown
	teardown()
}

func TestEncoding_ConvertFromUtf8ToEncodingWin1252(t *testing.T) {
	// Arrange
	input := "ßäüöë"
	config.Encoding = encoding.Windows1252
	// Act
	result, err := ConvertFromUtf8ToEncoding(input, config)
	// Assert
	assert.Nil(t, err)
	expected := []byte{0xDF, 0xE4, 0xFC, 0xF6, 0xEB}
	assert.Equal(t, expected, result)
	// Teardown
	teardown()
}

func TestEncoding_ConvertFromEncodingToUtf8InvalidEncoding(t *testing.T) {
	// Arrange
	input := []byte("invalid encoding")
	config.Encoding = "invalid_encoding"
	// Act
	_, err := ConvertFromEncodingToUtf8(input, config)
	// Assert
	assert.Error(t, err, errmsg.ErrEncodingInvalidEncoding.Error())
	assert.Equal(t, "invalid_encoding: "+errmsg.ErrEncodingInvalidEncoding.Error(), err.Error())
	// Teardown
	teardown()
}
