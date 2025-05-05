package e2e

import (
	"bytes"
	"github.com/blutspende/go-astm/v3/enums/encoding"
	"github.com/blutspende/go-astm/v3/models/astmmodels"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"testing"
	"time"
)

// Configuration struct for tests
var config astmmodels.Configuration

// Reset config to default values
func teardown() {
	config = astmmodels.DefaultConfiguration
	config.Encoding = encoding.UTF8
	config.Delimiters = astmmodels.DefaultDelimiters
	config.TimeLocation, _ = time.LoadLocation(config.TimeZone)
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
