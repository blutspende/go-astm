package astm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
)

// Fix a bug in which the decodes message couldnt exceed 4096 bytes
func TestEncodingVeryLongCharsets(t *testing.T) {

	veryLongMessage := []byte{}
	for i := 0; i < 100000; i++ {
		veryLongMessage = append(veryLongMessage, []byte("Ich bin ein sehr langer Text")...)
	}

	encoded, err := EncodeCharsetToUTF8From(charmap.Windows1252, veryLongMessage)
	assert.Nil(t, err)

	// in this case the encoding should equal the coded part (no special characters)
	assert.Equal(t, len(veryLongMessage), len(encoded))

	for i := 0; i < len(veryLongMessage); i++ {
		assert.Equal(t, veryLongMessage[i], encoded[i])
	}
}
