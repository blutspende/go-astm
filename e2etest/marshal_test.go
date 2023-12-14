package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/blutspende/go-astm/lib/standardlis2a2"
	"github.com/blutspende/go-astm/lis2a2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
)

type IllFormatedButLegal struct {
	GeneratedSequence int    `astm:"1,sequence"`
	ThirdfieldArray1  string `astm:"2.1.3"`
	FirstFieldArray1  string `astm:"2.1.1"`
	FirstFieldArray2  string `astm:"2.2.1"`
	SecondfieldArray2 string `astm:"2.2.2"`
	SomeEmptyField    string `astm:"3"`
}

type MinimalMessageMarshal struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Ill        IllFormatedButLegal       `astm:"?"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestSimpleMarshal(t *testing.T) {
	var msg MinimalMessageMarshal

	msg.Header.AccessPassword = "password"
	msg.Header.Version = "0.1.0"
	msg.Header.SenderNameOrID = "test"

	msg.Ill.ThirdfieldArray1 = "third-arr1"
	msg.Ill.FirstFieldArray1 = "first-arr1"
	msg.Ill.FirstFieldArray2 = "first-arr2"
	msg.Ill.SecondfieldArray2 = "second-arr2"

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
	}

	assert.Nil(t, err)

	assert.Equal(t, "H|\\^&||password|test||||||||0.1.0|", string(lines[0]))
	assert.Equal(t, "?|1|first-arr1^^third-arr1\\first-arr2^second-arr2|", string(lines[1]))
	assert.Equal(t, "L|1|", string(lines[2]))
}

type ArrayMessageMarshal struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Patient    []standardlis2a2.Patient  `astm:"P"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestGenverateSequence(t *testing.T) {

	var msg ArrayMessageMarshal

	msg.Patient = make([]standardlis2a2.Patient, 2)
	msg.Patient[0].LastName = "Firstus'"
	msg.Patient[0].FirstName = "Firstie"
	msg.Patient[1].LastName = "Secundus'"
	msg.Patient[1].FirstName = "Secundie"

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	assert.Nil(t, err)
	// output on screen
	for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
	}

	assert.Equal(t, "H|\\^&||||||||||||", string(lines[0]))
	assert.Equal(t, "P|1||||Firstus'^Firstie|||||||||||||||||||||||||||||", string(lines[1]))
	assert.Equal(t, "P|2||||Secundus'^Secundie|||||||||||||||||||||||||||||", string(lines[2]))
	assert.Equal(t, "L|1|", string(lines[3]))
}

type PatientResult struct {
	Patient standardlis2a2.Patient  `astm:"P"`
	Result  []standardlis2a2.Result `astm:"R"`
}

type ArrayNestedStructMessageMarshal struct {
	Header        standardlis2a2.Header `astm:"H"`
	PatientResult []PatientResult
	Terminator    standardlis2a2.Terminator `astm:"L"`
}

func TestNestedStruct(t *testing.T) {
	var msg ArrayNestedStructMessageMarshal

	msg.PatientResult = make([]PatientResult, 2)
	msg.PatientResult[0].Patient.FirstName = "Chuck"
	msg.PatientResult[0].Patient.LastName = "Norris"
	msg.PatientResult[0].Patient.Religion = "Binaries"
	msg.PatientResult[0].Result = make([]standardlis2a2.Result, 2)
	msg.PatientResult[0].Result[0].ManufacturersTestName = "SulfurBloodCount"
	msg.PatientResult[0].Result[0].MeasurementValueOfDevice = "100"
	msg.PatientResult[0].Result[0].Units = "%"
	msg.PatientResult[0].Result[1].ManufacturersTestName = "Catblood"
	msg.PatientResult[0].Result[1].MeasurementValueOfDevice = ">100000"
	msg.PatientResult[0].Result[1].Units = "U/l"
	msg.PatientResult[1].Patient.FirstName = "Eric"
	msg.PatientResult[1].Patient.LastName = "Cartman"
	msg.PatientResult[1].Patient.Religion = "None"
	msg.PatientResult[1].Result = make([]standardlis2a2.Result, 1)
	msg.PatientResult[1].Result[0].ManufacturersTestName = "Fungenes"
	msg.PatientResult[1].Result[0].MeasurementValueOfDevice = "present"
	msg.PatientResult[1].Result[0].Units = "none"

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	assert.Nil(t, err)
	// output on screen
	for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
	}

	// TODO: ? with this struct the Patient sequence will not advance - always P|1...
	assert.Equal(t, "H|\\^&||||||||||||", string(lines[0]))
	assert.Equal(t, "P|1||||Norris^Chuck|||||||||||||||||||||||Binaries||||||", string(lines[1]))
	assert.Equal(t, "R|1|^^^^SulfurBloodCount^^|^^100|%||||||^|||", string(lines[2]))
	assert.Equal(t, "R|2|^^^^Catblood^^|^^>100000|U/l||||||^|||", string(lines[3]))
	assert.Equal(t, "P|1||||Cartman^Eric|||||||||||||||||||||||None||||||", string(lines[4]))
	assert.Equal(t, "R|1|^^^^Fungenes^^|^^present|none||||||^|||", string(lines[5]))
	assert.Equal(t, "L|1|", string(lines[6]))
}

type TimeTestMessageMarshal struct {
	Header standardlis2a2.Header `astm:"H"`
}

/*
Test provides current time as UTC and expects the converter to stream as Belrin-Time
*/
func TestTimeLocalization(t *testing.T) {

	var msg TimeTestMessageMarshal

	europeBerlin, err := time.LoadLocation("Europe/Berlin")
	assert.Nil(t, err)

	testTime := time.Now()
	timeInBerlin := time.Now().In(europeBerlin)

	msg.Header.DateAndTime = testTime.UTC()

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)
	assert.Nil(t, err)

	assert.Equal(t, fmt.Sprintf("H|\\^&||||||||||||%s", timeInBerlin.Format("20060102150405")), string(lines[0]))
}

type TestMarshalEnum string

const (
	SomeTestMarshalEnum1 TestMarshalEnum = "Something"
	SomeTestMarshalEnum2 TestMarshalEnum = "SomethingElse"
)

type TestMarshalEnumRecord struct {
	Field TestMarshalEnum
}

type TestMarshalEnumMessage struct {
	Record TestMarshalEnumRecord `astm:"X"`
}

/*
Marshalling of enums
*/
func TestEnumMarshal(t *testing.T) {
	var msg TestMarshalEnumMessage

	msg.Record.Field = SomeTestMarshalEnum2

	lines, err := lis2a2.Marshal(msg, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	assert.Nil(t, err)
	// output on screen
	for _, line := range lines {
		linestr := string(line)
		fmt.Println(linestr)
	}

}

type TestCorrectFieldEnumeration struct {
	Request OrderRequestV5 `astm:"R"`
}

/*
Testing field assignments
*/
type OrderRequestV5 struct {
	// ID                  uuid.UUID `json:"id" db:"id"`
	SequenceNumber      int    `astm:"2,sequence" db:"sequence_number"` // 8.4.2 (see https://samson-rus.com/wp-content/files/LIS2-A2.pdf)
	SpecimenID          string `astm:"3" db:"specimen_id"`              // 8.4.3
	CodeOfSpecimen1     string `astm:"4.1.1" db:"code_of_specimen_1"`   // 8.4.4
	TypeOfSpecimen1     string `astm:"4.1.2" db:"type_of_specimen_1"`
	CodeOfDonor1        string `astm:"4.1.3" db:"code_of_donor_1"`
	TypeOfDonorSample1  string `astm:"4.1.4" db:"type_of_donor_sample_1"`
	CodeOfSpecimen2     string `astm:"4.2.1" db:"code_of_specimen_2"` // 8.4.4
	TypeOfSpecimen2     string `astm:"4.2.2" db:"type_of_specimen_2"`
	CodeOfDonor2        string `astm:"4.2.3" db:"code_of_donor_2"`
	TypeOfDonorSample2  string `astm:"4.2.4" db:"type_of_donor_sample_2"`
	UniversalTestID     string `astm:"5.1" db:"universal_test_id"`      // 8.4.5
	UniversalTestIDName string `astm:"5.2" db:"universal_test_id_name"` // 8.4.5
	UniversalTestIDType string `astm:"5.3" db:"universal_test_id_type"` // 8.4.5
	ManufacturesTestID  string `astm:"5.4" db:"manufactures_test_id"`
	// Priority                    OrderPriority `astm:"6" db:"priority"`                               // 8.4.6
	RequestedOrderDateTime      time.Time `astm:"7,longdate" db:"requested_order_date_time"`     // 8.4.7
	SpecimenCollectionDateTime  time.Time `astm:"8,longdate" db:"specimen_collection_date_time"` // 8.4.8
	CollectionEndTime           time.Time `astm:"9,longdate" db:"collection_end_time"`           // 8.4.9
	CollectionVolume            string    `astm:"10" db:"collection_volume"`                     // 8.4.10
	CollectorID                 string    `astm:"11" db:"collector_id"`                          // 8.4.11
	ActionCode                  string    `astm:"12" db:"action_code"`                           // 8.4.12
	DangerCode                  string    `astm:"13" db:"danger_code"`                           // 8.4.13
	RelevantClinicalInformation string    `astm:"14" db:"relevant_clinical_information"`         // 8.4.14
	DateTimeSpecimenReceived    string    `astm:"15" db:"date_time_specimen_received"`           // 8.4.15
	SpecimenTypeSource          string    `astm:"16" db:"specimen_type_source"`                  // 8.4.16
	OrderingPhysician           string    `astm:"17" db:"ordering_physician"`                    // 8.4.17
	PhysicianTelephone          string    `astm:"18" db:"physician_telephone"`                   // 8.4.18
	UserField1                  string    `astm:"19" db:"user_field_1"`                          // 8.4.19
	UserField2                  string    `astm:"20" db:"user_field_2"`                          // 8.4.20
	LaboratoryField1            string    `astm:"21" db:"laboratory_field_1"`
	LaboratoryField2            string    `astm:"22" db:"laboratory_field_2"`
	// ProtocolMessageHistoryID    uuid.UUID     `db:"message_history_id"`
	CreatedAt time.Time `db:"created_at"`
}

func TestFieldEnumeration(t *testing.T) {
	var orq TestCorrectFieldEnumeration

	orq.Request.ActionCode = "N"
	record, err := lis2a2.Marshal(orq, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.StandardNotation)

	assert.Nil(t, err)

	assert.Equal(t, "R|1||^^^\\^^^|^^^|||||||N||||||||||", string(record[0]))

}

/*
Testing bug: one delimiter too much
*/
type TestOneDlimiterTooMuchStruct struct {
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestOneDlimiterTooMuch(t *testing.T) {

	var record TestOneDlimiterTooMuchStruct

	record.Terminator.TerminatorCode = "N"
	filedata, err := lis2a2.Marshal(record, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.StandardNotation)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(filedata))

	assert.Equal(t, "L|1|N", string(filedata[0]))
}

// --------------------------------------------------------------
// Testing bug: German Language encoding
// --------------------------------------------------------------
type TestGermanLanguageDecoderRecord struct {
	Patient standardlis2a2.Patient `astm:"P"`
}

func TestGermanLanguageDecoder(t *testing.T) {

	var record TestGermanLanguageDecoderRecord

	record.Patient.FirstName = "Högendäg"
	record.Patient.LastName = "Nügendiß"
	filedata, err := lis2a2.Marshal(record, lis2a2.EncodingWindows1252, lis2a2.TimezoneEuropeBerlin, lis2a2.StandardNotation)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(filedata))

	expectedWindows1252 := helperEncode(charmap.Windows1252, []byte("P|1||||Nügendiß^Högendäg|||||||||||||||||||||||||||||"))

	assert.Equal(t, expectedWindows1252, filedata[0])

	// test for iso8859_1
	filedata, err = lis2a2.Marshal(record, lis2a2.EncodingISO8859_1, lis2a2.TimezoneEuropeBerlin, lis2a2.StandardNotation)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(filedata))

	expectedWindowsISO8859_1 := helperEncode(charmap.ISO8859_1, []byte("P|1||||Nügendiß^Högendäg|||||||||||||||||||||||||||||"))

	assert.Equal(t, expectedWindowsISO8859_1, filedata[0])
}

// Due to a bug that panics
func TestFailMarshalOnlyHeader(t *testing.T) {

	var header standardlis2a2.Header

	message, err := lis2a2.Marshal(header, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	assert.Nil(t, err)
	assert.NotNil(t, message)
}

// Pointers to structures dont work for marshalling
func TestFailOnHeaderIsA_Reference(t *testing.T) {

	var header standardlis2a2.Header

	_, err := lis2a2.Marshal(&header, lis2a2.EncodingASCII, lis2a2.TimezoneEuropeBerlin, lis2a2.ShortNotation)

	assert.NotNil(t, err)
}
