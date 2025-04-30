package e2e

import (
	"fmt"
	"github.com/blutspende/go-astm/v3/constants/astmconst"
	"github.com/blutspende/go-astm/v3/models/messageformat/lis02a2"
	"testing"
	"time"

	"github.com/blutspende/go-astm/v3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/charmap"
)

type MissingComponentMessage struct {
	MissingComponent MissingComponentRecord `astm:"M"`
}
type MissingComponentRecord struct {
	Combined   string `astm:"3"`
	Component1 string `astm:"4.1"`
	Component2 string `astm:"4.2"`
}

func TestMissingComponent(t *testing.T) {
	// Arrange
	testMessage := MissingComponentMessage{
		MissingComponent: MissingComponentRecord{
			Combined:   "First^Second",
			Component2: "Second",
		},
	}
	// Act
	lines, err := astm.Marshal(testMessage, config)
	// Assert
	assert.Nil(t, err)
	expectedResult := "M|1|First^Second|^Second"
	assert.Equal(t, expectedResult, string(lines[0]))
}

type IllFormattedSubstructure struct {
	ThirdComp string `astm:"3"`
	FirstComp string `astm:"1"`
}
type WellFormattedSubstructure struct {
	FirstComp  string `astm:"1"`
	SecondComp string `astm:"2"`
}
type IllFormatedButLegal struct {
	Ill        IllFormattedSubstructure  `astm:"3"`
	Well       WellFormattedSubstructure `astm:"4"`
	EmptyField string                    `astm:"5"`
}
type IllFormattedMinimalMessage struct {
	Header     lis02a2.Header      `astm:"H"`
	Substruct  IllFormatedButLegal `astm:"?"`
	Terminator lis02a2.Terminator  `astm:"L"`
}

func TestIllFormattedSubstructured(t *testing.T) {
	// Arrange
	var message IllFormattedMinimalMessage
	message.Header.AccessPassword = "password"
	message.Header.Version = "0.1.0"
	message.Header.SenderNameOrID = "test"
	message.Substruct.Ill.FirstComp = "struct1-comp1"
	message.Substruct.Ill.ThirdComp = "struct1-comp3"
	message.Substruct.Well.FirstComp = "struct2-comp1"
	message.Substruct.Well.SecondComp = "struct2-comp2"
	// Act
	lines, err := astm.Marshal(message, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 3)
	assert.Equal(t, "H|\\^&||password|test||||||||0.1.0|", string(lines[0]))
	assert.Equal(t, "?|1|struct1-comp1^^struct1-comp3|struct2-comp1^struct2-comp2|", string(lines[1]))
	assert.Equal(t, "L|1|", string(lines[2]))
}

type ArrayMessage struct {
	Header     lis02a2.Header     `astm:"H"`
	Patient    []lis02a2.Patient  `astm:"P"`
	Terminator lis02a2.Terminator `astm:"L"`
}

func TestGenerateSequence(t *testing.T) {
	// Arrange
	var msg ArrayMessage
	msg.Patient = make([]lis02a2.Patient, 2)
	msg.Patient[0].LastName = "Firstus'"
	msg.Patient[0].FirstName = "Firstie"
	msg.Patient[1].LastName = "Secundus'"
	msg.Patient[1].FirstName = "Secundie"
	// Act
	lines, err := astm.Marshal(msg, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 4)
	assert.Equal(t, "H|\\^&||||||||||||", string(lines[0]))
	assert.Equal(t, "P|1||||Firstus'^Firstie|||||||||||||||||||||||||||||", string(lines[1]))
	assert.Equal(t, "P|2||||Secundus'^Secundie|||||||||||||||||||||||||||||", string(lines[2]))
	assert.Equal(t, "L|1|", string(lines[3]))
}

type PatientResult struct {
	Patient lis02a2.Patient  `astm:"P"`
	Result  []lis02a2.Result `astm:"R"`
}
type ArrayNestedStructMessage struct {
	Header        lis02a2.Header `astm:"H"`
	PatientResult []PatientResult
	Terminator    lis02a2.Terminator `astm:"L"`
}

func TestNestedStruct(t *testing.T) {
	// Arrange
	var msg ArrayNestedStructMessage
	msg.PatientResult = make([]PatientResult, 2)
	msg.PatientResult[0].Patient.FirstName = "Chuck"
	msg.PatientResult[0].Patient.LastName = "Norris"
	msg.PatientResult[0].Patient.Religion = "Binaries"
	msg.PatientResult[0].Result = make([]lis02a2.Result, 2)
	msg.PatientResult[0].Result[0].UniversalTestID.ManufacturersTestName = "SulfurBloodCount"
	msg.PatientResult[0].Result[0].MeasurementValueOfDevice = "100"
	msg.PatientResult[0].Result[0].Units = "%"
	msg.PatientResult[0].Result[1].UniversalTestID.ManufacturersTestName = "Catblood"
	msg.PatientResult[0].Result[1].MeasurementValueOfDevice = ">100000"
	msg.PatientResult[0].Result[1].Units = "U/l"
	msg.PatientResult[1].Patient.FirstName = "Eric"
	msg.PatientResult[1].Patient.LastName = "Cartman"
	msg.PatientResult[1].Patient.Religion = "None"
	msg.PatientResult[1].Result = make([]lis02a2.Result, 1)
	msg.PatientResult[1].Result[0].UniversalTestID.ManufacturersTestName = "Fungenes"
	msg.PatientResult[1].Result[0].MeasurementValueOfDevice = "present"
	msg.PatientResult[1].Result[0].Units = "none"
	// Act
	lines, err := astm.Marshal(msg, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 7)
	assert.Equal(t, "H|\\^&||||||||||||", string(lines[0]))
	assert.Equal(t, "P|1||||Norris^Chuck|||||||||||||||||||||||Binaries||||||", string(lines[1]))
	assert.Equal(t, "R|1|^^^^SulfurBloodCount^^|^^100|%||||||^|||", string(lines[2]))
	assert.Equal(t, "R|2|^^^^Catblood^^|^^>100000|U/l||||||^|||", string(lines[3]))
	assert.Equal(t, "P|2||||Cartman^Eric|||||||||||||||||||||||None||||||", string(lines[4]))
	assert.Equal(t, "R|1|^^^^Fungenes^^|^^present|none||||||^|||", string(lines[5]))
	assert.Equal(t, "L|1|", string(lines[6]))
}

type HeaderMessage struct {
	Header lis02a2.Header `astm:"H"`
}

func TestTimeLocalization(t *testing.T) {
	// Note: Test provides current time as UTC and expects the converter to stream as Berlin-Time
	// Arrange
	var msg HeaderMessage
	europeBerlin, err := time.LoadLocation("Europe/Berlin")
	testTime := time.Now()
	timeInBerlin := time.Now().In(europeBerlin)
	msg.Header.DateAndTime = testTime.UTC()
	// Act
	lines, err := astm.Marshal(msg, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("H|\\^&||||||||||||%s", timeInBerlin.Format("20060102150405")), string(lines[0]))
}

type MarshalEnum string

const SomeTestMarshalEnum1 MarshalEnum = "Something"
const SomeTestMarshalEnum2 MarshalEnum = "SomethingElse"

type MarshalEnumRecord struct {
	Field MarshalEnum `astm:"3"`
}
type MarshalEnumMessage struct {
	Record MarshalEnumRecord `astm:"X"`
}

func TestEnumMarshal(t *testing.T) {
	// Arrange
	var msg MarshalEnumMessage
	msg.Record.Field = SomeTestMarshalEnum2
	// Act
	lines, err := astm.Marshal(msg, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 1)
	assert.Equal(t, "X|1|SomethingElse", string(lines[0]))
}

type SpecimenDonorSubstructure struct {
	CodeOfSpecimen    string `astm:"1"` // 8.4.4
	TypeOfSpecimen    string `astm:"2"`
	CodeOfDonor       string `astm:"3"`
	TypeOfDonorSample string `astm:"4"`
}
type OrderRequestV5 struct {
	SpecimenID                  string                      `astm:"3" db:"specimen_id"`              // 8.4.3
	SpecimenDonors              []SpecimenDonorSubstructure `astm:"4"`                               // 8.4.4
	UniversalTestID             string                      `astm:"5.1" db:"universal_test_id"`      // 8.4.5
	UniversalTestIDName         string                      `astm:"5.2" db:"universal_test_id_name"` // 8.4.5
	UniversalTestIDType         string                      `astm:"5.3" db:"universal_test_id_type"` // 8.4.5
	ManufacturesTestID          string                      `astm:"5.4" db:"manufactures_test_id"`
	RequestedOrderDateTime      time.Time                   `astm:"7,longdate" db:"requested_order_date_time"`     // 8.4.7
	SpecimenCollectionDateTime  time.Time                   `astm:"8,longdate" db:"specimen_collection_date_time"` // 8.4.8
	CollectionEndTime           time.Time                   `astm:"9,longdate" db:"collection_end_time"`           // 8.4.9
	CollectionVolume            string                      `astm:"10" db:"collection_volume"`                     // 8.4.10
	CollectorID                 string                      `astm:"11" db:"collector_id"`                          // 8.4.11
	ActionCode                  string                      `astm:"12" db:"action_code"`                           // 8.4.12
	DangerCode                  string                      `astm:"13" db:"danger_code"`                           // 8.4.13
	RelevantClinicalInformation string                      `astm:"14" db:"relevant_clinical_information"`         // 8.4.14
	DateTimeSpecimenReceived    string                      `astm:"15" db:"date_time_specimen_received"`           // 8.4.15
	SpecimenTypeSource          string                      `astm:"16" db:"specimen_type_source"`                  // 8.4.16
	OrderingPhysician           string                      `astm:"17" db:"ordering_physician"`                    // 8.4.17
	PhysicianTelephone          string                      `astm:"18" db:"physician_telephone"`                   // 8.4.18
	UserField1                  string                      `astm:"19" db:"user_field_1"`                          // 8.4.19
	UserField2                  string                      `astm:"20" db:"user_field_2"`                          // 8.4.20
	LaboratoryField1            string                      `astm:"21" db:"laboratory_field_1"`
	LaboratoryField2            string                      `astm:"22" db:"laboratory_field_2"`
	CreatedAt                   time.Time                   `db:"created_at"`
}
type FieldEnumerationMessage struct {
	Request OrderRequestV5 `astm:"R"`
}

func TestFieldEnumeration(t *testing.T) {
	// Arrange
	var orq FieldEnumerationMessage
	orq.Request.ActionCode = "N"
	// Act
	lines, err := astm.Marshal(orq, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 1)
	assert.Equal(t, "R|1|||^^^|||||||N||||||||||", string(lines[0]))
}

type GermanLanguageDecoderMessage struct {
	Patient lis02a2.Patient `astm:"P"`
}

func TestGermanLanguageDecoder_Windows1252(t *testing.T) {
	// Arrange
	var record GermanLanguageDecoderMessage
	record.Patient.FirstName = "Högendäg"
	record.Patient.LastName = "Nügendiß"
	config.Encoding = astmconst.ENCODING_WINDOWS1252
	// Act
	lines, err := astm.Marshal(record, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 1)
	expected := helperEncode(charmap.Windows1252, []byte("P|1||||Nügendiß^Högendäg|||||||||||||||||||||||||||||"))
	assert.Equal(t, expected, lines[0])
	// Teardown
	teardown()
}
func TestGermanLanguageDecoder_ISO8859_1(t *testing.T) {
	// Arrange
	var record GermanLanguageDecoderMessage
	record.Patient.FirstName = "Högendäg"
	record.Patient.LastName = "Nügendiß"
	config.Encoding = astmconst.ENCODING_ISO8859_1
	// Act
	lines, err := astm.Marshal(record, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 1)
	expected := helperEncode(charmap.ISO8859_1, []byte("P|1||||Nügendiß^Högendäg|||||||||||||||||||||||||||||"))
	assert.Equal(t, expected, lines[0])
	// Teardown
	teardown()
}

func TestMarshalOnlyEmptyHeader(t *testing.T) {
	// Arrange
	var message HeaderMessage
	config.Encoding = astmconst.ENCODING_ASCII
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	lines, err := astm.Marshal(message, config)
	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, lines)
	// Teardown
	teardown()
}

func TestPointerInput(t *testing.T) {
	// Arrange
	var message HeaderMessage
	// Act
	_, err := astm.Marshal(&message, config)
	// Assert
	assert.Nil(t, err)
}

func TestQueryMessageNoQueryData(t *testing.T) {
	// Arrange
	var query lis02a2.QueryMessage
	query.Terminator.TerminatorCode = "N"
	// Act
	lines, err := astm.Marshal(query, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "H|\\^&||||||||||||", string(lines[0]))
	assert.Equal(t, "L|1|N", string(lines[1]))
}

func TestQueryMessage(t *testing.T) {
	// Arrange
	var query lis02a2.QueryMessage
	query.Terminator.TerminatorCode = "N"
	query.Queries = []lis02a2.Query{
		{
			StartingRangeIDNumber: "SampleCode1",
			UniversalTestID:       "ALL",
		},
	}
	// Act
	lines, err := astm.Marshal(query, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "H|\\^&||||||||||||", string(lines[0]))
	assert.Equal(t, "Q|1|SampleCode1||ALL||||||||", string(lines[1]))
	assert.Equal(t, "L|1|N", string(lines[2]))
}

func TestMarshalMultipleOrder(t *testing.T) {
	// Arrange
	msg := lis02a2.OrderMessage{
		PatientOrders: []lis02a2.PatientOrder{
			{
				Patient: lis02a2.Patient{
					LabAssignedPatientID: "Mate",
				},
				Orders: []lis02a2.Order{
					{
						SpecimenID: "Samplecode1",
						UniversalTestID: lis02a2.StandardUniversalTestID{
							UniversalTestID: "Brains",
						},
					},
					{
						SpecimenID: "Samplecode1",
						UniversalTestID: lis02a2.StandardUniversalTestID{
							UniversalTestID: "Gutts",
						},
					},
				},
			},
			{
				Patient: lis02a2.Patient{
					LabAssignedPatientID: "Stephan",
				},
				Orders: []lis02a2.Order{
					{
						SpecimenID: "Samplecode2",
						UniversalTestID: lis02a2.StandardUniversalTestID{
							UniversalTestID: "Looks",
						},
					},
					{
						SpecimenID: "Samplecode2",
						UniversalTestID: lis02a2.StandardUniversalTestID{
							UniversalTestID: "Money",
						},
					},
				},
			},
		},
		Terminator: lis02a2.Terminator{
			TerminatorCode: "N",
		},
	}
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	lines, err := astm.Marshal(msg, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 8)
	assert.Equal(t, "H|\\^&", string(lines[0]))
	assert.Equal(t, "P|1||Mate", string(lines[1]))
	assert.Equal(t, "O|1|Samplecode1||Brains", string(lines[2]))
	assert.Equal(t, "O|2|Samplecode1||Gutts", string(lines[3]))
	assert.Equal(t, "P|2||Stephan", string(lines[4]))
	assert.Equal(t, "O|1|Samplecode2||Looks", string(lines[5]))
	assert.Equal(t, "O|2|Samplecode2||Money", string(lines[6]))
	assert.Equal(t, "L|1|N", string(lines[7]))
	// Teardown
	teardown()
}

func TestShorthandOnStandardMessage(t *testing.T) {
	// Arrange
	msg := lis02a2.StandardPOCRMessage{
		Header: lis02a2.Header{
			SenderNameOrID: "LIS",
			ReceiverID:     "NonExistentTestSystem",
			DateAndTime:    time.Now(),
		},
		PatientOrderCommentedResults: []lis02a2.PatientOrderCommentedResult{
			{
				Patient: lis02a2.Patient{},
				OrderCommentedResults: []lis02a2.OrderCommentedResult{
					{
						Order: lis02a2.Order{
							SpecimenID: "VAL24981209",
							UniversalTestID: lis02a2.StandardUniversalTestID{
								UniversalTestID: "Pool_Cell",
							},
							Priority:     "R",
							ActionCode:   "N",
							SpecimenType: "TestData",
						},
					},
					{
						Order: lis02a2.Order{
							SpecimenID: "VAL24981210",
							UniversalTestID: lis02a2.StandardUniversalTestID{
								UniversalTestID: "Pool_Cell",
							},
							Priority:     "R",
							ActionCode:   "N",
							SpecimenType: "TestData",
						},
					},
				},
			},
		},
		Terminator: lis02a2.Terminator{
			TerminatorCode: "N",
		},
	}
	config.Encoding = astmconst.ENCODING_ASCII
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	lines, err := astm.Marshal(msg, config)
	// Assert
	assert.Nil(t, err)
	assert.Len(t, lines, 6)
	assert.Equal(t, []byte("M|1"), lines[1])
	assert.Equal(t, []byte("P|1"), lines[2])
	assert.Equal(t, []byte("O|1|VAL24981209||Pool_Cell|R||||||N||||TestData"), lines[3])
	assert.Equal(t, []byte("O|2|VAL24981210||Pool_Cell|R||||||N||||TestData"), lines[4])
	assert.Equal(t, []byte("L|1|N"), lines[5])
	// Teardown
	teardown()
}

func TestEmbeddedStructsAndArrays(t *testing.T) {
	// Arrange
	message := HoribaYumizenMessage{
		ExtraTests: struct {
			ArrayOfInt     []int     `astm:"3"`
			ArrayOfFloat32 []float32 `astm:"4"`
			ArrayOfFloat64 []float64 `astm:"5"`
		}(struct {
			ArrayOfInt     []int
			ArrayOfFloat32 []float32
			ArrayOfFloat64 []float64
		}{
			ArrayOfInt:     []int{1, 2, 3},
			ArrayOfFloat32: []float32{4.1, 4.2, 4.3},
			ArrayOfFloat64: []float64{5.111, 5.222},
		}),
		Manufacturer: ManufacturerInfo{
			F2:       "REAGENT",
			Reagents: []string{"CLEANER", "DILUENT", "LYSE"},
			ReagentInfo: []TraceabilityReagentInfo{
				{
					SerialNumber:   "240415I1(",
					UsedAtDateTime: "20240902000000",
					UseByDate:      "20241202",
				},
				{
					SerialNumber:   "240423H1(",
					UsedAtDateTime: "20240905000000",
					UseByDate:      "20250305",
				},
				{
					SerialNumber:   "240411M11",
					UsedAtDateTime: "20240828000000",
					UseByDate:      "20241028",
				},
			},
		},
		Terminator: lis02a2.Terminator{
			TerminatorCode: "N",
		},
	}
	config.Encoding = astmconst.ENCODING_UTF8
	config.TimeZone = astmconst.TIMEZONE_UTC
	// Act
	lines, err := astm.Marshal(message, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "M|1|REAGENT|CLEANER\\DILUENT\\LYSE|240415I1(^20240902000000^20241202\\240423H1(^20240905000000^20250305\\240411M11^20240828000000^20241028", string(lines[0]))
	assert.Equal(t, "E|1|1\\2\\3|4.100\\4.200\\4.300|5.111\\5.222", string(lines[1]))
	assert.Equal(t, "L|1|N", string(lines[2]))
	// Teardown
	teardown()
}

func TestEmbeddedStructsAndArraysShortNotation(t *testing.T) {
	// Arrange
	message := HoribaYumizenMessage{
		ExtraTests: struct {
			ArrayOfInt     []int     `astm:"3"`
			ArrayOfFloat32 []float32 `astm:"4"`
			ArrayOfFloat64 []float64 `astm:"5"`
		}{
			ArrayOfInt: []int{1, 2, 3},
		},
		Manufacturer: ManufacturerInfo{
			F2:       "REAGENT",
			Reagents: []string{"DILUENT", "LYSE"},
		},
		Terminator: lis02a2.Terminator{
			TerminatorCode: "N",
		},
	}
	config.Encoding = astmconst.ENCODING_UTF8
	config.TimeZone = astmconst.TIMEZONE_UTC
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	lines, err := astm.Marshal(message, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "M|1|REAGENT|DILUENT\\LYSE", string(lines[0]))
	assert.Equal(t, "E|1|1\\2\\3", string(lines[1]))
	assert.Equal(t, "L|1|N", string(lines[2]))
	// Teardown
	teardown()
}

type CustomDecimalLength struct {
	Float1 float32 `astm:"3,length:4"`
	Float2 float64 `astm:"4,length:1"`
	Floats []struct {
		EmbeddedFloat1 float32 `astm:"1,length:7"`
		EmbeddedFloat2 float64 `astm:"2,length:2"`
		EmbeddedFloat3 float64 `astm:"3"`
		EmbeddedFloat4 float32 `astm:"4,length:3"`
	} `astm:"5"`
}

func TestCustomDecimalLengthAnnotation(t *testing.T) {
	// Arrange
	message := struct {
		DecimalLength CustomDecimalLength `astm:"D"`
	}{
		DecimalLength: CustomDecimalLength{
			Float1: 0.34567,
			Float2: 0.40001,
			Floats: []struct {
				EmbeddedFloat1 float32 `astm:"1,length:7"`
				EmbeddedFloat2 float64 `astm:"2,length:2"`
				EmbeddedFloat3 float64 `astm:"3"`
				EmbeddedFloat4 float32 `astm:"4,length:3"`
			}{
				{
					EmbeddedFloat1: 0.123456711,
					EmbeddedFloat2: 0.984654321,
					EmbeddedFloat3: 0.234444,
					EmbeddedFloat4: 0.345444,
				},
				{
					EmbeddedFloat1: 0.99,
					EmbeddedFloat2: 0.1122334455,
					EmbeddedFloat3: 0.2233445566,
					EmbeddedFloat4: 0.3344556677,
				},
			},
		},
	}
	config.Encoding = astmconst.ENCODING_UTF8
	config.TimeZone = astmconst.TIMEZONE_UTC
	config.Notation = astmconst.NOTATION_SHORT
	// Act
	lines, err := astm.Marshal(message, config)
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "D|1|0.3457|0.4|0.1234567^0.98^0.234^0.345\\0.9900000^0.11^0.223^0.334", string(lines[0]))
	// Teardown
	teardown()
}
