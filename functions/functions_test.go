package functions

import (
	"github.com/blutspende/go-astm/v2/models"
	"testing"
	"time"
)

// Configuration struct for tests
var config *models.Configuration

// Reset config to default values
func teardown() {
	config = &models.Configuration{}
	*config = models.DefaultConfiguration
	config.Internal.Delimiters = models.DefaultDelimiters
	config.Internal.TimeLocation, _ = time.LoadLocation(config.TimeZone)
}

// Setup mock data for every test
func TestMain(m *testing.M) {
	// Set up configuration
	teardown()
	// Run all tests
	m.Run()
}
