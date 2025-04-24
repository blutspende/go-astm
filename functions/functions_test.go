package functions

import (
	"github.com/blutspende/go-astm/v2/models/astmmodels"
	"testing"
	"time"
)

// Configuration struct for tests
var config *astmmodels.Configuration

// Reset config to default values
func teardown() {
	config = &astmmodels.Configuration{}
	*config = astmmodels.DefaultConfiguration
	config.Delimiters = astmmodels.DefaultDelimiters
	config.TimeLocation, _ = time.LoadLocation(config.TimeZone)
}

// Setup mock data for every test
func TestMain(m *testing.M) {
	// Set up configuration
	teardown()
	// Run all tests
	m.Run()
}

// Common test structures

// Annotation records
type AnnotatedLine struct {
	Field string `astm:"3.2,length:4"`
}
type AnnotatedArrayLine struct {
	Field []string `astm:"3,length:4"`
}
type Line struct {
	Field string `astm:"3"`
}
type SingleLineStruct struct {
	Lines Line `astm:"L"`
}
type AnnotatedArrayStruct struct {
	Lines []Line `astm:"L,required"`
}
type CompositeStruct struct {
	Composite AnnotatedArrayStruct
}
type CompositeArrayStruct struct {
	Composite []AnnotatedArrayStruct
}
type Substructure struct {
	FirstComponent  string `astm:"1"`
	SecondComponent string `astm:"2"`
}
type IllegalComponentArray struct {
	ComponentArray []string `astm:"3.1"`
}
type IllegalComponentSubstructure struct {
	ComponentSubstructure Substructure `astm:"3.1"`
}
type SubstructuredLine struct {
	Field Substructure   `astm:"3"`
	Array []Substructure `astm:"4"`
}
type TimeLine struct {
	Time time.Time `astm:"3"`
}

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
	String  string    `astm:"3"`
	Int     int       `astm:"4"`
	Float32 float32   `astm:"5"`
	Float64 float64   `astm:"6"`
	Date    time.Time `astm:"7"`
}
type MultitypeLengthRecord struct {
	FloatFull float64   `astm:"3"`
	FloatFix3 float64   `astm:"4,length:3"`
	FloatFix0 float64   `astm:"5,length:0"`
	LongDate  time.Time `astm:"6,longdate"`
	ShortDate time.Time `astm:"7"`
}
type MultitypePointerRecord struct {
	String  *string    `astm:"3"`
	Int     *int       `astm:"4"`
	Float32 *float32   `astm:"5"`
	Float64 *float64   `astm:"6"`
	Date    *time.Time `astm:"7"`
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
type EnumString string
type EnumRecord struct {
	Enum EnumString `astm:"3"`
}
type ReservedFieldRecord struct {
	TypeName  string `astm:"1"`
	SeqNumber string `astm:"2"`
}
type SparseFieldRecord struct {
	Field3 string `astm:"3"`
	Field5 string `astm:"5"`
}
type SubstructureField struct {
	FirstComponent  string `astm:"1"`
	SecondComponent string `astm:"2"`
	ThirdComponent  string `astm:"3"`
}
type SubstructureRecord struct {
	First  string            `astm:"3"`
	Second SubstructureField `astm:"4"`
	Third  string            `astm:"5"`
}
type SubstructureArrayRecord struct {
	First  string              `astm:"3"`
	Second []SubstructureField `astm:"4"`
	Third  string              `astm:"5"`
}
type SparseSubstructureField struct {
	Component1 string `astm:"1"`
	Component3 string `astm:"3"`
	Component6 string `astm:"6"`
}
type SparseSubstructureRecord struct {
	First  string                  `astm:"3"`
	Second SparseSubstructureField `astm:"4"`
}
type TimeRecord struct {
	Time time.Time `astm:"3,longdate"`
}
type WrongComponentOrderRecord struct {
	First string `astm:"3"`
	Comp2 string `astm:"4.2"`
	Comp1 string `astm:"4.1"`
	Comp3 string `astm:"4.3"`
}
type WrongComponentPlacementRecord struct {
	Field1 string `astm:"3"`
	Comp1  string `astm:"4.1"`
	Field2 string `astm:"5"`
	Comp2  string `astm:"4.2"`
}
type MultipleWrongComponentPlacementRecord struct {
	Field3 string `astm:"3"`
	Comp41 string `astm:"4.1"`
	Field5 string `astm:"5"`
	Comp62 string `astm:"6.2"`
	Comp42 string `astm:"4.2"`
	Field7 string `astm:"7"`
	Comp61 string `astm:"6.1"`
	Field8 string `astm:"8"`
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
type OptionalArrayAtTheEndMessage struct {
	First    RecordType1   `astm:"F"`
	Optional []RecordType2 `astm:"A,optional"`
}
type OptionalAtTheEndMessage struct {
	First    RecordType1 `astm:"F"`
	Optional RecordType2 `astm:"O,optional"`
}
