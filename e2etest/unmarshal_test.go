package e2e

import (
	"bytes"
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/models"
	"testing"
	"time"

	"github.com/blutspende/go-astm/v2"
	"github.com/blutspende/go-astm/v2/lib/standardlis2a2"

	"github.com/stretchr/testify/assert"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Configuration struct for tests
var config *models.Configuration

// Reset config to default values
func teardown() {
	config = &models.Configuration{}
	*config = models.DefaultConfiguration
	config.Encoding = constants.ENCODING_UTF8
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

type ComponentedTestMessage struct {
	Componented Componented `astm:"C"`
}
type Componented struct {
	Combined   string `astm:"3"`
	Component1 string `astm:"4.1"`
	Component2 string `astm:"4.2"`
}

func TestComponentMessage(t *testing.T) {
	// Arrange
	fileData := ""
	fileData = fileData + "C|1|First^Second|First^Second\n"
	var message ComponentedTestMessage

	// Act
	err := astm.Unmarshal([]byte(fileData), &message, config)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "First^Second", message.Componented.Combined)
	assert.Equal(t, "First", message.Componented.Component1)
	assert.Equal(t, "Second", message.Componented.Component2)
}

type NoSequenceMessage struct {
	NoSeq NoSequenceRecord `astm:"N"`
}
type NoSequenceRecord struct {
	Second string `astm:"2"`
	Third  string `astm:"3"`
}

func TestNoSequenceMessage(t *testing.T) {
	// Arrange
	fileData := ""
	fileData = fileData + "N|Second|Third\n"
	var message NoSequenceMessage
	config.EnforceSequenceNumberCheck = false

	// Act
	err := astm.Unmarshal([]byte(fileData), &message, config)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "Second", message.NoSeq.Second)
	assert.Equal(t, "Third", message.NoSeq.Third)

	// Tear down
	teardown()
}

type MinimalMessage struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestReadMinimalMessage(t *testing.T) {
	fileData := ""
	fileData = fileData + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\n"
	fileData = fileData + "L|1|N\n"

	var message MinimalMessage
	err := astm.Unmarshal([]byte(fileData), &message, config)

	assert.Nil(t, err)

	locale, err := time.LoadLocation("Europe/Berlin")
	assert.Nil(t, err)

	assert.Equal(t, "Bio-Rad", message.Header.SenderNameOrID)
	assert.Equal(t, "IH v5.2", message.Header.SenderStreetAddress)
	assert.Equal(t, "", message.Header.Comment)

	localtime := message.Header.DateAndTime.In(locale)
	assert.Equal(t, "20220315194227", localtime.Format("20060102150405"))
}

type FullSingleASTMMessage struct {
	Header       standardlis2a2.Header       `astm:"H"`
	Manufacturer standardlis2a2.Manufacturer `astm:"M,optional"`
	Patient      standardlis2a2.Patient      `astm:"P"`
	Order        standardlis2a2.Order        `astm:"O"`
	Result       standardlis2a2.Result       `astm:"R"`
	Terminator   standardlis2a2.Terminator   `astm:"L"`
}

func TestFullSingleASTMMessage(t *testing.T) {
	fileData := ""
	fileData = fileData + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\n"
	fileData = fileData + "P|1||1010868845||Testus^Test||19400607|M||||||||||||||||||||||||^\n"
	fileData = fileData + "O|1|1122206642|specimen1^^^\\specimen2^^^|^^^MO10^^28343^|R|20220311103217|20220311103217|||||||||||11||||20220311114103|||P\n"
	fileData = fileData + "R|1|^^^AntiA^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|40^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\n"
	fileData = fileData + "L|1|N\n"

	var message FullSingleASTMMessage
	err := astm.Unmarshal([]byte(fileData), &message, config)

	assert.Nil(t, err)

	assert.Equal(t, "Testus", message.Patient.LastName)
	assert.Equal(t, "Test", message.Patient.FirstName)
	assert.Equal(t, "19400607", message.Patient.DOB.Format("20060102")) // dates are read okay
	assert.Equal(t, "specimen1^^^\\specimen2^^^", message.Order.InstrumentSpecimenID)
	assert.Equal(t, "lalina", message.Result.OperatorIDPerformed)
}

// -----------------------------------------------------------------------------------
// TEST
// -----------------------------------------------------------------------------------
// Testing a rather more complex structure with optional and arrays on the
// structure as well as on the Records
// -----------------------------------------------------------------------------------
type MessagePORC struct {
	Header       standardlis2a2.Header       `astm:"H"`
	Manufacturer standardlis2a2.Manufacturer `astm:"M,optional"`
	OrderResults []struct {
		Patient         standardlis2a2.Patient `astm:"P"`
		Order           standardlis2a2.Order   `astm:"O"`
		CommentedResult []struct {
			Result  standardlis2a2.Result    `astm:"R"`
			Comment []standardlis2a2.Comment `astm:"C,optional"`
		}
	}
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestNestedStructure(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||1010868845||Testus^Test||19400607|M||||||||||||||||||||||||^\r"
	data = data + "O|1|1122206642|1122206642^^^\\1122206642^^^|^^^MO10^^28343^|R|20220311103217|20220311103217|||||||||||11||||20220311114103|||P\r"
	data = data + "R|1|^^^AntiA^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|40^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\r"
	data = data + "C|1|FirstComment^^05761.03.12^20240131\\^^^|CAS^5005352062212117030^50053.52.06^20221231^4||\r"
	data = data + "C|2|SecondComment^^05761.03.12^20240131\\^^^|CAS^5005352062212117030^50053.52.06^20221231^4||\r"
	data = data + "R|2|^^^AntiB^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|0^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\r"
	data = data + "C|1|ID-Diluent 2^^05761.03.12^20240131\\^^^|CAS^5005352062212117030^50053.52.06^20221231^5||\r"
	data = data + "R|3|^^^AntiD^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|0^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\r"
	data = data + "P|2||1010868845||Testis^Tost||19400607|M||||||||||||||||||||||||^\r"
	data = data + "O|1|1122206642|1122206642^^^\\1122206642^^^|^^^MO10^^28343^|R|20220311103217|20220311103217|||||||||||11||||20220311114103|||P\r"
	data = data + "R|1|^^^AntiA^MO10^Bloodgroup: A,B,D Confirmation for Patients (DiaClon) (5005)^|40^^|C||||R||lalina^|20220311114103||11|IH-1000|0300768|lalina\r"
	data = data + "L|1|N\r"

	var message MessagePORC
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)

	// patients have been read in an array
	assert.Equal(t, "Testus", message.OrderResults[0].Patient.LastName)
	assert.Equal(t, "Testis", message.OrderResults[1].Patient.LastName)

	// the array of comments was produced in two entries into the array
	assert.Equal(t, 2, len(message.OrderResults[0].CommentedResult[0].Comment))
	assert.Equal(t, "FirstComment^^05761.03.12^20240131\\^^^", message.OrderResults[0].CommentedResult[0].Comment[0].CommentSource)
	assert.Equal(t, "SecondComment^^05761.03.12^20240131\\^^^", message.OrderResults[0].CommentedResult[0].Comment[1].CommentSource)
}

// -----------------------------------------------------------------------------------
// TEST
// -----------------------------------------------------------------------------------
// Custom Delimiters : In the header there is a delimiter-field that can change
// the default delimiters
// -----------------------------------------------------------------------------------
type MessageCustomDelimiterTest struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Patient    standardlis2a2.Patient    `astm:"P"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestCustomDelimiters(t *testing.T) {
	data := ""
	data = data + "H|\\#&|||Bio-Rad|IH v5.2||||||||20220315194227\n"
	data = data + "P|1||1010868845||Testus#Test||19400607|M||||||||||||||||||||||||^\n"
	data = data + "L|1|N\n"

	var message MessageCustomDelimiterTest
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)

	// the delimiter is now "#" instead of "^"; this should have been recognizes
	assert.Equal(t, "Testus", message.Patient.LastName)
	assert.Equal(t, "Test", message.Patient.FirstName)

}

// -----------------------------------------------------------------------------------
// TEST
// -----------------------------------------------------------------------------------
// Custom Structures can be defined and mapped to records
//
//	aside:    also testing float values
//
// -----------------------------------------------------------------------------------
type CompleteOutOfStandardCustomRecord struct {
	SequenceNumber string  `astm:"2"`
	F2             string  `astm:"3"`
	F3             string  `astm:"4"`
	Float32Value   float32 `astm:"5"`
	Float64Value   float64 `astm:"6"`
}

type MessageWithOutOfStandardCustomRecord struct {
	Header       standardlis2a2.Header             `astm:"H"`
	CustomRecord CompleteOutOfStandardCustomRecord `astm:"X"`
	Terminator   standardlis2a2.Terminator         `astm:"L"`
}

func TestCustomRecord(t *testing.T) {
	data := ""
	data = data + "H|\\#&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "X|1|A|B|4.14159|2.172\r"
	data = data + "L|1|N\r" // ! mixed line-endings (should not matter)

	var message MessageWithOutOfStandardCustomRecord
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)
	assert.Equal(t, float32(4.14159), message.CustomRecord.Float32Value)
	assert.Equal(t, float64(2.172), message.CustomRecord.Float64Value)
}

// test timezone ok
// test delimiters ok
// test optionals ok
// test encoding
// line ending 0a or 0d or 0d0a all okay ? ok

//type valami struct {
//	Field11 string `astm:"1"`
//	Field12 string `astm:"2"`
//	Field13 string `astm:"3"`
//}
//type arraysor struct {
//	Sor []valami `astm:"2"`
//}
//
//type arraysor struct {
//	Sor string `astm:"2.1"`
//	Sor string `astm:"2.2"`
//	Sor string `astm:"2.3"`
//}

//type SubMessageRecord struct {
//	Field11 string `astm:"3.1.1"`
//	Field12 string `astm:"3.1.2"`
//	Field13 string `astm:"3.1.3"`
//	Field21 string `astm:"3.2.1"`
//	Field22 string `astm:"3.2.2"`
//	Field23 string `astm:"3.2.3"`
//}
//
//type MessageWithSubArrayRecord struct {
//	Anonymous struct { // Testing wether this annoynmous structure is recused into
//		Rec SubMessageRecord `astm:"?"`
//	}
//
//	AnonymousArray []struct { // This anynymous array of structures needs to be created and scanned
//		Rec SubMessageRecord `astm:"!"`
//	}
//}
//
//func TestArrayMapping(t *testing.T) {
//
//	data := "?|1|a^^c\\d^e^f|\r"
//	data = data + "!|1|x^y\\z^^|\r"
//	data = data + "!|1|1^2^3\\4^5^6|\r"
//
//	var message MessageWithSubArrayRecord
//	err := astm.Unmarshal([]byte(data), &message,
//		constants.ENCODING_UTF8, constants.TIMEZONE_EUROPE_BERLIN)
//
//	assert.Nil(t, err)
//
//	assert.Equal(t, "a", message.Anonymous.Rec.Field11)
//	assert.Equal(t, "c", message.Anonymous.Rec.Field13)
//
//	assert.Equal(t, "d", message.Anonymous.Rec.Field21)
//	assert.Equal(t, "e", message.Anonymous.Rec.Field22)
//	assert.Equal(t, "f", message.Anonymous.Rec.Field23)
//
//	// now test that the subarray values have been read
//	assert.Equal(t, 2, len(message.AnonymousArray))
//	assert.Equal(t, "x", message.AnonymousArray[0].Rec.Field11)
//	assert.Equal(t, "y", message.AnonymousArray[0].Rec.Field12)
//	assert.Equal(t, "z", message.AnonymousArray[0].Rec.Field21)
//	assert.Equal(t, "", message.AnonymousArray[0].Rec.Field22)
//
//	assert.Equal(t, "1", message.AnonymousArray[1].Rec.Field11)
//	assert.Equal(t, "2", message.AnonymousArray[1].Rec.Field12)
//	assert.Equal(t, "4", message.AnonymousArray[1].Rec.Field21)
//	assert.Equal(t, "5", message.AnonymousArray[1].Rec.Field22)
//}

type SomeEnum string

const (
	EnumValue1 SomeEnum = "EnumValue1"
	EnumValue2 SomeEnum = "EnumValue2"
)

type SomeEnumRecord struct {
	Value SomeEnum `astm:"3"`
}

type TestUnmarshalEnumMessage struct {
	Record SomeEnumRecord `astm:"E"`
}

// TODO enum value
func TestEnumEncoding(t *testing.T) {
	data := "E|1|EnumValue1|\r"

	var message TestUnmarshalEnumMessage
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)

	assert.Equal(t, EnumValue1, message.Record.Value)
}

// -----------------------------------------------------------------------------------
// TEST
// -----------------------------------------------------------------------------------
// Access time.Time in a struct with multiple components
//
//	aside:    also testing time.time values
//
// -----------------------------------------------------------------------------------
type AccessTimeComment struct {
	SequenceNumber              int       `astm:"2,sequence" db:"sequence_number"`            // 3.2.8 - ih_com_host_connection_manual_astm_1394_en_h009164_v1_8.pdf
	DescriptionOfReagent        string    `astm:"3.1.1"  db:"description_of_reagent"`         //
	BarcodeOfReagent            string    `astm:"3.1.2" db:"barcode_of_reagent"`              //
	LotNumberOfReagent          string    `astm:"3.1.3" db:"lot_number_of_reagent"`           //
	ExpirationDateOfReagent     time.Time `astm:"3.1.4" db:"expiration_date_of_reagent"`      //
	DescriptionOfReagent2       string    `astm:"3.2.1" db:"description_of_reagent_2"`        //
	BarcodeOfReagent2           string    `astm:"3.2.2" db:"barcode_of_reagent_2"`            //
	LotNumberOfReagent2         string    `astm:"3.2.3" db:"lot_number_of_reagent_2"`         //
	ExpirationDateOfReagent2    time.Time `astm:"3.2.4" db:"expiration_date_of_reagent_2"`    //
	TypeOfTestMedia             string    `astm:"4.1" db:"type_of_test_media"`                //
	PlateOrIDCardBarcode        string    `astm:"4.2" db:"plate_or_id_card_barcode"`          //
	LotNumberOfCassetteOrPlate  string    `astm:"4.3" db:"lot_number_of_cassette_or_plate"`   //
	ExpDateForIDCardOrPlate     time.Time `astm:"4.4" db:"exp_date_for_id_card_or_plate"`     //
	IDCardOrPlateRealWellNumber int       `astm:"4.5" db:"id_card_or_plate_real_well_number"` //
	Comment                     string    `astm:"5" db:"comment"`                             //
	FileName                    string    `astm:"6" db:"file_name"`                           //
}
type MessageTimeAccess struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Comment    AccessTimeComment         `astm:"C"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestComponentAccessOnTime(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "C|1|FirstComment^^05761.03.12^20240131\\^^^|CAS^5005352062212117030^50053.52.06^20221231^4||\r"
	data = data + "L|1|N\r"

	var message MessageTimeAccess
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)

	location, err := time.LoadLocation(string(constants.TIMEZONE_EUROPE_BERLIN))
	assert.Nil(t, err, "Can not parse timezone")

	expDate, err := time.ParseInLocation("20060102", "20240131", location)
	assert.Nil(t, err, "Can not parse date")
	assert.Equal(t, expDate, message.Comment.ExpirationDateOfReagent)
}

type TestCommentNoneBugComment struct {
	SequenceNumber int       `astm:"2,sequence"`
	Field1         time.Time `astm:"3.1.4"` // out of bounds with component index
	Field2         time.Time `astm:"3.2.4"` // out of bounds with repeat index
	Field3         time.Time `astm:"4.4"`
}
type TestCommentNoneBugMessage struct {
	Field TestCommentNoneBugComment `astm:"C"`
}

type TestCommentNoneBugCommentCrash struct {
	SequenceNumber int       `astm:"2,sequence"`
	Field1         time.Time `astm:"3.1.4,required"` // out of bounds with component index
	Field2         time.Time `astm:"3.2.4"`          // out of bounds with repeat index
	Field3         time.Time `astm:"4.4,required"`
}

type TestCommentNoneBugMessageCrash struct {
	Field TestCommentNoneBugComment `astm:"C"`
}

func TestCommentNoneBug(t *testing.T) {
	data := ""
	data = data + "C|1|^^^||\r"

	var message TestCommentNoneBugMessage
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)

	assert.Equal(t, time.Time{}, message.Field.Field1)
	assert.Equal(t, time.Time{}, message.Field.Field2)
	assert.Equal(t, time.Time{}, message.Field.Field3)

	/* var crash TestCommentNoneBugMessageCrash
	err := lis2a2.Unmarshal([]byte(data), &crash,
		lis2a2.ENCODING_UTF8, lis2a2.TimezoneEuropeBerlin)
	assert.NotNil(t, err) */
}

// -----------------------------------------------------------------------------------
// TEST a german message Win1252 Encoded
// -----------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------
type MessageGermanLanguageTest struct {
	Header     standardlis2a2.Header     `astm:"H"`
	Patient    standardlis2a2.Patient    `astm:"P"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}

func TestGermanLanguage(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\n"
	data = data + "P|1||1010868845||König^#$§?/+öäüß||19400607|M||||||||||||||||||||||||^\r"
	data = data + "L|1|N\n" // ! mixed line-endings (should not matter)

	var message MessageGermanLanguageTest

	config.Encoding = constants.ENCODING_WINDOWS1252
	encdata := helperEncode(charmap.Windows1252, []byte(data))
	err := astm.Unmarshal([]byte(encdata), &message, config)
	assert.Nil(t, err)
	assert.Equal(t, "König", message.Patient.LastName)
	assert.Equal(t, "#$§?/+öäüß", message.Patient.FirstName)
	teardown()

	config.Encoding = constants.ENCODING_ISO8859_1
	encdata = helperEncode(charmap.ISO8859_1, []byte(data))
	err = astm.Unmarshal([]byte(encdata), &message, config)
	assert.Nil(t, err)
	assert.Equal(t, "König", message.Patient.LastName)
	assert.Equal(t, "#$§?/+öäüß", message.Patient.FirstName)
	teardown()
}

func TestTransmissionWithoutLTerminator(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||\r"
	data = data + "P|1||DIA-27-079-5-1\r"

	var message standardlis2a2.DefaultMessage
	config.Encoding = constants.ENCODING_WINDOWS1252
	err := astm.Unmarshal([]byte(data), &message, config)
	assert.NotNil(t, err)
	teardown()
}

func TestFullMultipleASTMMessage(t *testing.T) {
	var data string

	// Message 1
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSQVIGG3||20220715071219\r"
	data = data + "R|1|^^^SARSQVIGG3|2598,88|BAU/ml|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSQVIGG3||20220715071219\r"
	data = data + "R|1|^^^SARSQVIGG3|3636,64|BAU/ml|\r"
	data = data + "L|1|N\r"

	// Message 2
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSNCPIGG||20220715071219\r"
	data = data + "R|1|^^^SARSNCPIGG|0,08|Ratio|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSNCPIGG||20220715071219\r"
	data = data + "R|1|^^^SARSNCPIGG|0,20|Ratio|\r"
	data = data + "L|1|N\r"

	// Message 3
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSNEUTRA||20220715071219\r"
	data = data + "R|1|^^^SARSNEUTRA|99,39|% IH|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSNEUTRA||20220715071219\r"
	data = data + "R|1|^^^SARSNEUTRA|99,23|% IH|\r"
	data = data + "L|1|N\r"

	// Message 4
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|>10|Ratio|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|>10|Ratio|\r"
	data = data + "P|3||DIA-01-061-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|4,87|Ratio|\r"
	data = data + "P|4||DIA-01-055-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|4,14|Ratio|\r"
	data = data + "L|1|N"

	var message standardlis2a2.DefaultMultiMessage
	err := astm.Unmarshal(
		[]byte(data), &message, config)

	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, 4, len(message.Messages))
}

func TestFullMultipleASTMMessageWithWrongInput(t *testing.T) {
	var data string

	// Message 1
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSQVIGG3||20220715071219\r"
	data = data + "R|1|^^^SARSQVIGG3|2598,88|BAU/ml|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSQVIGG3||20220715071219\r"
	data = data + "R|1|^^^SARSQVIGG3|3636,64|BAU/ml|\r"
	data = data + "L|1|N\r"

	// Message 2
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSNCPIGG||20220715071219\r"
	data = data + "R|1|^^^SARSNCPIGG|0,08|Ratio|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSNCPIGG||20220715071219\r"
	data = data + "R|1|^^^SARSNCPIGG|0,20|Ratio|\r"
	data = data + "L|1|N\r"

	// Message 3
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSNEUTRA||20220715071219\r"
	data = data + "R|1|^^^SARSNEUTRA|99,39|% IH|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSNEUTRA||20220715071219\r"
	data = data + "R|1|^^^SARSNEUTRA|99,23|% IH|\r"
	data = data + "L|1|N\r"

	// Message 4
	data = data + "H|\\^&|||Bio-Rad|IH v5.2||||||||20220315194227\r"
	data = data + "P|1||DIA-01-085-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|>10|Ratio|\r"
	data = data + "P|2||DIA-01-056-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|>10|Ratio|\r"
	data = data + "P|3||DIA-01-061-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|4,87|Ratio|\r"
	data = data + "P|4||DIA-01-055-7-1\r"
	data = data + "O|1|||^^^SARSCOV2IGA||20220715071219\r"
	data = data + "R|1|^^^SARSCOV2IGA|4,14|Ratio|\r"
	data = data + "L|1|N"

	var message standardlis2a2.DefaultMessage
	err := astm.Unmarshal(
		[]byte(data),
		&message, config)

	assert.NotNil(t, err)
}

func helperEncode(charmap *charmap.Charmap, data []byte) []byte {
	e := charmap.NewEncoder()
	var b bytes.Buffer
	writer := transform.NewWriter(&b, e)
	writer.Write([]byte(data))
	resultdata := b.Bytes()
	writer.Close()
	return resultdata
}

func TestFailOnUndisciplinedMultipleCRCRatEndOfLine(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||\u000d\u000d"
	data = data + "P|1||DIA-04-066-7-1\u000d\u000d"
	data = data + "O|1|||^^^SARS-CoV-2 NeutraLISA||20220715071342\u000d\u000d"
	data = data + "R|1|^^^SARS-CoV-2 NeutraLISA|99,66|% IH|\u000d\u000d"
	data = data + "L|1|N\u000d\u000d"

	var message standardlis2a2.DefaultMessage
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)
}

func TestMultipleMessagesInOne(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||\u000d\u000d"
	data = data + "P|1||DIA-04-066-7-1\u000d\u000d"
	data = data + "O|1|||^^^SARS-CoV-2 NeutraLISA||20220715071342\u000d\u000d"
	data = data + "R|1|^^^SARS-CoV-2 NeutraLISA|12,5|% IH|\u000d\u000d"
	data = data + "L|1|N\u000d\u000d"
	data = data + "H|\\^&|||\u000d\u000d"
	data = data + "P|1||DIA-04-066-7-2\u000d\u000d"
	data = data + "O|1|||^^^SARS-CoV-2 NeutraLISA||20220715071343\u000d\u000d"
	data = data + "R|1|^^^SARS-CoV-2 NeutraLISA|99,66|% IH|\u000d\u000d"
	data = data + "L|1|N\u000d\u000d"

	var message standardlis2a2.DefaultMultiMessage
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(message.Messages))
	assert.Equal(t, "DIA-04-066-7-2", message.Messages[1].PatientOrderCommentedResults[0].Patient.LabAssignedPatientID)

	assert.Equal(t, "12,5", message.Messages[0].PatientOrderCommentedResults[0].OrderCommentedResults[0].CommentedResults[0].Result.DataMeasurementValue)
	assert.Equal(t, "99,66", message.Messages[1].PatientOrderCommentedResults[0].OrderCommentedResults[0].CommentedResults[0].Result.DataMeasurementValue)
}

func TestNullValuesShouldGiveQualifiedError(t *testing.T) {

	var message standardlis2a2.DefaultMultiMessage
	err := astm.Unmarshal(nil /*giving null as input*/, &message, config)

	assert.NotNil(t, err)
	assert.Equal(t, "message has nil value - aborting", err.Error())
}

/* A Result transmission needs to process multiple orders/results per patient*/
func TestUnmarshalMultipleOrdersAndResultsForOnePatient(t *testing.T) {
	data := ""
	data = data + "H|\\^&|||RVT|||||LIS|||LIS2-A2|20240709103536\r"
	data = data + "P|1||||^^^^|||U|||||||||||||||||Main||||||||||\r"
	data = data + "O|1|CL5G2A118S||^^^Pool_Cell|R||||||||||^||||||||||F||||||\r"
	data = data + "R|1|^^^Pool_Cell 1|0^0^3.0|||||F||saidam||20240625092245|5030100461|\r"
	data = data + "R|2|^^^Pool_Cell|Negative|||||F||peilja||20240625092245|5030100461|\r"
	data = data + "O|2|CL5G2A118S||^^^Weak_D1|R||||||||||^||||||||||F||||||\r"
	data = data + "R|1|^^^Weak_D1|0^0^0.0|||||F||saidam||20240626193019|5030100461|\r"
	data = data + "R|2|^^^Weak_D1|Negative|||||F||SCHMMI||20240626193019|5030100461|\r"
	data = data + "L|1|N\r"

	var message standardlis2a2.DefaultMessage
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)
}

type MessageMadeForTheNextTest struct {
	//Header       standardlis2a2.Header `astm:"H"`
	Manufacturer ManufacturerInfo `astm:"M,optional"`
	ExtraTests   struct {
		SequenceNumber int       `astm:"2,sequence"`
		ArrayOfInt     []int     `astm:"3"`
		ArrayOfFloat32 []float32 `astm:"4"`
		ArrayOfFloat64 []float64 `astm:"5"`
	} `astm:"E,optional"`
	Terminator standardlis2a2.Terminator `astm:"L"`
}
type ManufacturerInfo struct {
	SequenceNumber int                       `astm:"2,sequence"`
	F2             string                    `astm:"3"`
	Reagents       []string                  `astm:"4"`
	ReagentInfo    []TraceabilityReagentInfo `astm:"5"`
}

type TraceabilityReagentInfo struct {
	SerialNumber   string `astm:"1.1"`
	UsedAtDateTime string `astm:"1.2"`
	UseByDate      string `astm:"1.3"`
}

/* Test a funny format with an array found with the yumizen Horiba instrument */
func TestHoribaYumizenManufacturerRecordWithArray(t *testing.T) {
	data := ""
	data = data + "M|1|REAGENT|CLEANER\\DILUENT\\LYSE|240415I1(^20240902000000^20241202\\240423H1(^20240905000000^20250305\\240411M11^20240828000000^20241028\n"
	data = data + "E|1|1\\2\\3|6.0\\7.8|5.887\\88.1045|"
	data = data + "L|1|N\n"

	var message MessageMadeForTheNextTest
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)

	assert.Equal(t, []string{"CLEANER", "DILUENT", "LYSE"}, message.Manufacturer.Reagents)
	assert.Equal(t, []int{1, 2, 3}, message.ExtraTests.ArrayOfInt)
	assert.Equal(t, []float32{6.0, 7.8}, message.ExtraTests.ArrayOfFloat32)
	assert.Equal(t, []float64{5.887, 88.1045}, message.ExtraTests.ArrayOfFloat64)

	assert.Equal(t, "240415I1(", message.Manufacturer.ReagentInfo[0].SerialNumber)
	assert.Equal(t, "20240902000000", message.Manufacturer.ReagentInfo[0].UsedAtDateTime)
	assert.Equal(t, "20241202", message.Manufacturer.ReagentInfo[0].UseByDate)
	assert.Equal(t, "240423H1(", message.Manufacturer.ReagentInfo[1].SerialNumber)
	assert.Equal(t, "20240905000000", message.Manufacturer.ReagentInfo[1].UsedAtDateTime)
	assert.Equal(t, "20250305", message.Manufacturer.ReagentInfo[1].UseByDate)
	assert.Equal(t, "240411M11", message.Manufacturer.ReagentInfo[2].SerialNumber)
	assert.Equal(t, "20240828000000", message.Manufacturer.ReagentInfo[2].UsedAtDateTime)
	assert.Equal(t, "20241028", message.Manufacturer.ReagentInfo[2].UseByDate)

}

func TestUnmarshalSliceOfOneRecordType(t *testing.T) {
	data := "M|1|HISTOGRAM|RBC/PLT|RbcAlongRes|FLOATLE-stream/deflate:base64^ZV+UPhH/dfBm3vuue/z+n6fHzyG4XyaZfbL8Cw3DNM0jAplGGny6P0G+7ZK2NtcHva6TE8yaLJtpib72B4x7cgUv11WW1OelZr0W9mqcosxlnjWWnak1qq2GyzGW5FMO2MSVm/JDsb1WZmuPYwdsJiL8SOWP3yAYw5ZoXSW46asmO8sx16sQunLFL6yI1epSOYaFfPlKuZSidF81VtSqJKtixTzqrRarDJdy9TyjnXFyhSsV6lNe1UvnDIRVIrVGsTVXb67A3YNvYtdibsetwovzXgNXI/zG8OGOaMd9jXAtuG2PbsT9gfAd+J+tJMMdHrKmbeT5mXTuYq5e1fcp8O1lfH3Puwk5i78beg92P/QV2CnsA+0vsr7CHZq9/zDeMvRf7G+wR7H3Y+7G/wz6AfRD7B+wx7EPYE9iHsY9gZ7GPYh/DPo49hT2N/Qv2DPZZ7F85xhBX6CLcoRT/JS4ZpLYnSOjJ2cK17XPO69S3pLLpds1RXiD18p3AtJtrplume+BFJXC/dFUpMLxHlmxHOtcI8kra7Dvh77Buw87Buxb8K+GTsf24u9ELsAuxC7CPsW7FuxF2Hfhn07tg/bj30HdjH2ndiLse/CXoK9FHsZdgl2APtu7FLse7DLsO/FLsc2sS1shS3Yy7Hvw74fuwL7AewHsR/CrsSuwq7Gfhg7iP0I9qPTj2KuwQ9hPYT2Kvxl6D/RT209jPYq/Ffg47jP089jrsF7BrsNdjb8B+EXsj9kvYNvYm7JexX8GuxX4V+zXs17E3Y2/BjmC/gV2H/Sb2Vuy32I5K9usov+vFP1HPvgaxZxqY821JuhuZt1Gmi5qYu0kCZTHmj0lkxTbWsE1S6+OsIy5G9B3W8u75d8B/8ImW/znnveJk6qLUTYPOft4LbiqmIEWpY/ZdYRg/0QK89omI7SPNadTwEKUg3VUwt9yLntpH4aon10kA7z3498n+D7FM2wfYZO0890jI6wf5z2015+D9Dn1EOd7Hufb2cem1aR82wtoULKlXPrmmF9WfqWBmk3fUIJaqU41VENrdbn6LyTA1Ssz7uA8nS56tz1cOvmUg6d0tfKaVz3vb6G/bpd1K3r+EOturi+9ufbRM/oe+Jk6pz7lK9z65x7OG7+3qCu+2+K/kPmBTI+vzGw==\n"
	data += "M|2|HISTOGRAM|RBC/PLT|PltAlongRes|FLOATLE-stream/deflate:base64^Y2AAAQ4nMMXQA6IdgMghPS3IoXrKWsdVH00cIXIN9iA5AA==|FLOATLE-stream/deflate:base64^7dR/TJR1HAfwR7gZ1Khj0tjlaOAUzRoRSzIpnu/7uRlj0MjRGDpsmCwZ69yYIhkDTzQlIK8gFYWChBDHUQxA+DOkNDqbN6g2whUeHWH/3Rs3vvbvc83/u89/o+90jSxOGkTL5JmXzXC35gPCFJBcwl3Dg38d0Orc66NNDudjDQsqozUGu6U/bW+spO7Wvk4dgt8mByuhwZlCs3qI1yQGeVbCxskj10Vtngb5Mn1lfEOIn4BrVY7qMRfVleIvfqfeLpKD/hYlklzi7RipT0EPGILVzYwqNESXmM2OChE5rUBNE6kMzzu8Shl4XbbI9srhtXkidE8Bry8SquESrikV1WEnuO6U2Gau4Npq4aM5I9JDzonB5EahNTWL/N5WYXe7ICKDuoR5e7dwNfYJnXVQNKhHhLf2kkhNuCKsR6+JgE47f1vCi/kOKClXob1lPlTDTvBT3YENHi4w+N+F6jBXDMcugCb1bgTlumObWYOipoVoHfDAhOHlBYvQv3QxOgK8cT5sGWo2Locp8QEUZvjgQL4v0k7ke0bSXCHQPwhPvjePR+GXzh3nAtXDethmNSEMb2BePikRB0nX4SzY1hsFjXwDwajuL5ETh0TyQyfdYhRYlCfMQziImLRkTKswjOisFjxc/hwcpYeDXHwa33edw2thnjzvEY8diC7ocS0LI6EXVrt6NMl4TjO5ORt38H6tv0/P2dqBxK5YxdKLXv5pw9KFDv5aw0ZC1+mfPS8dLKDM7MxAuhr3DuPsRFGzj7VazGx2eB1K3n72OIAVpoPskoNldYfY5zAWduSyUx5cRt5grzcn77lx57fodoT9CmlXxI5v06+YPY/SsIRdj9HxOPu+Q0sjO5fS8132fo+mJhiOnaDr+0itPklbM7a2nKLvaWzqL6NxOdaNV9D5A4S6VNG6GrJXDb0/hN+KWppbsCT4DN0/gvv6Otqfxe3x5+hfj+u7P+YeNMCW08h9OI9e4yfciyZ8UdvM/fiUZp/GZtNGun2QWaddDsS5p10qyLZl/R7GuaWWnWTbMemn1Ds16a9dHsW5r102yAZoM0G6LZRZp9R7Nhmo3Q7Hua/TD5Hx13/pFmozS7TLMrNPuJZmM0u0qzazQbp9nPNPuFZnaaXafZrzT7jeclpb5NUurWzlMs1nlKzUYHpXLIQSnTOSrmUUfFlKhSSu2qqWfEv3Ho5T9HmnruiGnRT0sBUzstPVOR8EfU0+I5I7557njD5q3Nx1Ij3i5tTOkoIZ0f9NxBwi3UJm26f/89/JrezlX+Wf7pcb99bv\n"
	data += "M|3|MATRIX|LMNE|LMNEResAbs|FLOATLE-stream/deflate:base64^Y2AAggf/XRjgtIMDiAkA|FLOATLE-stream/deflate:base64^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\n"
	message := struct {
		ImageData []struct {
			SequenceNumber  int    `astm:"2"`
			Field1          string `astm:"3"`
			Field2          string `astm:"4"`
			Field3          string `astm:"5"`
			Stream1Encoding string `astm:"6.1"`
			Stream1Data     string `astm:"6.2"`
			Stream2Encoding string `astm:"7.1"`
			Stream2Data     string `astm:"7.2"`
		} `astm:"M"`
	}{}
	err := astm.Unmarshal([]byte(data), &message, config)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(message.ImageData))
	assert.Equal(t, "HISTOGRAM", message.ImageData[0].Field1)
	assert.Equal(t, "MATRIX", message.ImageData[2].Field1)
}

// Fix a bug in which the decodes message couldn't exceed 4096 bytes
func TestEncodingVeryLongCharsets(t *testing.T) {

	veryLongMessage := []byte{}
	for i := 0; i < 100000; i++ {
		veryLongMessage = append(veryLongMessage, []byte("Ich bin ein sehr langer Text")...)
	}

	encoded, err := astm.EncodeCharsetToUTF8From(charmap.Windows1252, veryLongMessage)
	assert.Nil(t, err)

	// in this case the encoding should equal the coded part (no special characters)
	assert.Equal(t, len(veryLongMessage), len(encoded))

	for i := 0; i < len(veryLongMessage); i++ {
		assert.Equal(t, veryLongMessage[i], encoded[i])
	}
}
