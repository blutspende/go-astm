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

// Common test structures

// Single line records
type SimpleRecord struct {
	First  string `astm:"3"`
	Second string `astm:"4"`
	Third  string `astm:"5"`
}

type UnorderedRecord struct {
	First  string `astm:"3"`
	Third  string `astm:"5"`
	Second string `astm:"4"`
}

type MultitypeRecord struct {
	String     string    `astm:"3"`
	Int        int       `astm:"4"`
	Float32    float32   `astm:"5"`
	Float64    float64   `astm:"6"`
	Float64Cut float64   `astm:"7,length:3"`
	ShortTime  time.Time `astm:"8"`
	LongTime   time.Time `astm:"9,longdate"`
}

type ComponentedRecord struct {
	First       string `astm:"3"`
	SecondComp1 string `astm:"4.1"`
	SecondComp2 string `astm:"4.2"`
	ThirdComp1  string `astm:"5.1"`
	ThirdComp2  string `astm:"5.2"`
	ThirdComp3  string `astm:"5.3"`
}

type ArrayRecord struct {
	First string   `astm:"3"`
	Array []string `astm:"4"`
}

type HeaderRecord struct {
	First string `astm:"3"`
}

type HeaderDelimiterChange struct {
	First string   `astm:"3"`
	Array []string `astm:"4"`
	Comp1 string   `astm:"5.1"`
	Comp2 string   `astm:"5.2"`
}

type MissingRequiredField struct {
	First  string `astm:"3"`
	Second string `astm:"4,required"`
	Third  string `astm:"5"`
}

type RecordType1 struct {
	First  string `astm:"3"`
	Second string `astm:"4"`
}
type RecordType2 struct {
	First  string `astm:"3"`
	Second string `astm:"4"`
}

// Structures

type SingleRecordStruct struct {
	FirstRecord SimpleRecord `astm:"R"`
}

type RecordArrayStruct struct {
	RecordArray []SimpleRecord `astm:"R"`
}

type CompositeRecordStruct struct {
	Record1 RecordType1 `astm:"F"`
	Record2 RecordType2 `astm:"S"`
}
type CompositeMessage struct {
	CompositeRecordStruct CompositeRecordStruct
}

type CompositeArrayMessage struct {
	CompositeRecordArray []CompositeRecordStruct
}

type OptionalMessage struct {
	First    RecordType1 `astm:"F"`
	Optional RecordType2 `astm:"S,optional"`
	Third    RecordType1 `astm:"T"`
}

type OptionalArrayMessage struct {
	First    RecordType1   `astm:"F"`
	Optional []RecordType2 `astm:"A,optional"`
	Last     RecordType1   `astm:"L"`
}
