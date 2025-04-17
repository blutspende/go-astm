package e2e

import (
	"bytes"
	"github.com/blutspende/go-astm/v2/constants/astmconst"
	"github.com/blutspende/go-astm/v2/models"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"testing"
	"time"
)

// Configuration struct for tests
var config *models.Configuration

// Reset config to default values
func teardown() {
	config = &models.Configuration{}
	*config = models.DefaultConfiguration
	config.Encoding = astmconst.ENCODING_UTF8
	config.Internal.Delimiters = models.DefaultDelimiters
	config.Internal.TimeLocation, _ = time.LoadLocation(config.TimeZone)
}

// Setup default config and run all tests
func TestMain(m *testing.M) {
	// Set up configuration
	teardown()
	// Run all tests
	m.Run()
}

// Encoding helper function
func helperEncode(charmap *charmap.Charmap, data []byte) []byte {
	e := charmap.NewEncoder()
	var b bytes.Buffer
	writer := transform.NewWriter(&b, e)
	writer.Write(data)
	resultdata := b.Bytes()
	writer.Close()
	return resultdata
}
