package astm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifyOrderMessage(t *testing.T) {

	astm := "H|\\^&|||LIS|||||NEO|||LIS2A2|20220928182311\n"
	astm = astm + "P|1||||^|||||||||||||||||||||||||||||\n"
	astm = astm + "O|1|idk1||^^^Pool_Cell||||R|||N||||Blood^Product|||||||||||||||\n"
	astm = astm + "L|1|N\n"

	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeOrdersOnly, messageType)
}

func TestIdentifyOrderMessageWithMultiHeader(t *testing.T) {

	astm := "H|\\^&|||LIS|||||NEO|||LIS2A2|20220928182311\n"
	astm = astm + "P|1||||^|||||||||||||||||||||||||||||\n"
	astm = astm + "O|1|idk1||^^^Pool_Cell||||R|||N||||Blood^Product|||||||||||||||\n"
	astm = astm + "H|\\^&|||LIS|||||NEO|||LIS2A2|20220928182311\n"
	astm = astm + "P|1||||^|||||||||||||||||||||||||||||\n"
	astm = astm + "O|1|idk1||^^^Pool_Cell||||R|||N||||Blood^Product|||||||||||||||\n"
	astm = astm + "H|\\^&|||LIS|||||NEO|||LIS2A2|20220928182311\n"
	astm = astm + "P|1||||^|||||||||||||||||||||||||||||\n"
	astm = astm + "O|1|idk1||^^^Pool_Cell||||R|||N||||Blood^Product|||||||||||||||\n"
	astm = astm + "H|\\^&|||LIS|||||NEO|||LIS2A2|20220928182311\n"
	astm = astm + "P|1||||^|||||||||||||||||||||||||||||\n"
	astm = astm + "O|1|idk1||^^^Pool_Cell||||R|||N||||Blood^Product|||||||||||||||\n"
	astm = astm + "L|1|N\n"

	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeOrdersOnly, messageType)
}

func TestIdentifyQuery(t *testing.T) {

	astm := `H|\^&|||RVT|||||LIS|||LIS2-A2|20200302132021
Q|1|VALI200301||ALL
Q|2|VALI200302||ALL
Q|3|VALI200303||ALL
Q|4|VALI200304||ALL
Q|5|VALI200305||ALL
L|1|N`

	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeQuery, messageType)
}

func TestIdentifyQueryWithMultiHeader(t *testing.T) {
	astm := `H|\^&|||RVT|||||LIS|||LIS2-A2|20200302132021
Q|1|VALI200301||ALL
Q|2|VALI200302||ALL
Q|3|VALI200303||ALL
Q|4|VALI200304||ALL
Q|5|VALI200305||ALL
H|\^&|||RVT|||||LIS|||LIS2-A2|20200302132021
Q|1|VALI200301||ALL
Q|2|VALI200302||ALL
Q|3|VALI200303||ALL
Q|4|VALI200304||ALL
H|\^&|||RVT|||||LIS|||LIS2-A2|20200302132021
Q|1|VALI200301||ALL
Q|5|VALI200305||ALL
L|1|N`

	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeQuery, messageType)
}

func TestIdentifyOrderAndResult(t *testing.T) {

	astm := `H|\^&|||RVT|||||LIS|||LIS2-A2|20200302131145
P|1||||^^^^|||U|||||||||||||||||Main||||||||||
O|1|VAL99999903||^^^Pool_Cell|R||||||||||^||||||||||F||||||
R|1|^^^Pool_Cell 1|0^0^8.8|||||F||Immucor||20200226153444|5030100389|
R|2|^^^Pool_Cell|Negative|||||F||immucor||20200226153444|5030100389|
L|1|N`
	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)

	assert.Nil(t, err)
	assert.Equal(t, MessageTypeOrdersAndResults, messageType)
}

func TestIdentifyOrderAndResultWithMultiHeader(t *testing.T) {
	astm := `H|\^&|||RVT|||||LIS|||LIS2-A2|20200302131145
P|1||||^^^^|||U|||||||||||||||||Main||||||||||
O|1|VAL99999903||^^^Pool_Cell|R||||||||||^||||||||||F||||||
R|1|^^^Pool_Cell 1|0^0^8.8|||||F||Immucor||20200226153444|5030100389|
R|2|^^^Pool_Cell|Negative|||||F||immucor||20200226153444|5030100389|
H|\^&|||RVT|||||LIS|||LIS2-A2|20200302131145
P|1||||^^^^|||U|||||||||||||||||Main||||||||||
O|1|VAL99999903||^^^Pool_Cell|R||||||||||^||||||||||F||||||
R|1|^^^Pool_Cell 1|0^0^8.8|||||F||Immucor||20200226153444|5030100389|
R|2|^^^Pool_Cell|Negative|||||F||immucor||20200226153444|5030100389|
H|\^&|||RVT|||||LIS|||LIS2-A2|20200302131145
P|1||||^^^^|||U|||||||||||||||||Main||||||||||
O|1|VAL99999903||^^^Pool_Cell|R||||||||||^||||||||||F||||||
R|1|^^^Pool_Cell 1|0^0^8.8|||||F||Immucor||20200226153444|5030100389|
R|2|^^^Pool_Cell|Negative|||||F||immucor||20200226153444|5030100389|
L|1|N`
	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)

	assert.Nil(t, err)
	assert.Equal(t, MessageTypeOrdersAndResults, messageType)
}

func TestIdentifyWithEmptyLines(t *testing.T) {

	astm := `H|\^&|||RVT|||||LIS|||LIS2-A2|20200302132021
Q|1|VALI200301||ALL
Q|2|VALI200302||ALL

Q|4|VALI200304||ALL
Q|5|VALI200305||ALL
L|1|N

`

	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeQuery, messageType)
}

// -----------------------------------------------------------------------------------
// The bug was that this Transmission contains one "P" and then mutlitple orders
// Default Multi Message Was not processing those corerctly
// -----------------------------------------------------------------------------------
func TestIdentifyHPORCOROCOROC(t *testing.T) {
	data := ""

	data = data + "H|\\^&|||Bio-Rad|IH v5.1||||||||20230805142035\n"
	data = data + "P|1||AA5E2ACC29||^|||||||||||||||||||||||||||^\n"
	data = data + "O|1||AA5E2ACC29^^^\\^^^|F\n"
	data = data + "R|1|^^^AntiA^MO01A^Blutgruppe: ABO/D  (5048)^|\n"
	data = data + "C|1|ID-Diluent 2^^05761.04.41^20250228\\^^^|\n"
	data = data + "R|2|^^^AntiB^MO01A^Blutgruppe: ABO/D  (5048)^|\n"
	data = data + "O|2||AA5E2ACC29^^^\\^^^|^^^MO10^^33619^|\n"
	data = data + "R|1|^^^AntiA^MO10^Blutgruppe Best�tigung: A,B,D (5005)^|\n"
	data = data + "C|1|ID-Diluent 2^^05761.04.41^20250228\\^^^|\n"
	data = data + "R|2|^^^AntiB^MO10^Blutgruppe Best�tigung: A,B,D (5005)^|\n"
	data = data + "C|1|ID-Diluent 2^^05761.04.41^20250228\\^^^|\n"
	data = data + "O|3||AA5E2ACC29^^^\\^^^|^^^PR07C^^33619^|\n"
	data = data + "R|1|^^^cellA1^PR07C^Serumgegenprobe: A1,B,O (5052)^|\n"
	data = data + "C|1|ID-DiaCell A1^^06012.49.1^20230821\\^^^|\n"
	data = data + "L|1|N\n"

	messageType, err := IdentifyMessage([]byte(data), EncodingUTF8)
	assert.Nil(t, err)

	assert.Equal(t, MessageTypeOrdersAndResults, messageType)

}

func TestIdentifyOrderAndMultipleResultsWithManufacturerDefinedField(t *testing.T) {
	astm := `H|\^&|||H550^909YAXH02732^1.2.1.4|||||||P|LIS2-A2|20240906090907
P|1|||||||||||||||||||||||||||||||||||
O|1|70424906396^^013148^6||^^^DIF|R|20240906090745|||||||||BLOOD||||||||||F|||||
C|1|I|NON_COMPLIANT_DATA^WBC^ABNORMAL_DIFFERENTIA\NON_COMPLIANT_DATA^RBC^RBC_DBL\NON_COMPLIANT_DATA^RBC^ABNORMAL_MCH\NON_COMPLIANT_DATA^PLT^PC_MODE\NON_COMPLIANT_DATA^PLT^SEP_RBC_PLT\SUSPECTED_PATHOLOGY^^ANEMIA\SUSPECTED_PATHOLOGY^^LEUKOPENB9IA\SUSPECTED_PATHOLOGY^^LYMPHOPENIA\SUSPECTED_PATHOLOGY^^LARGE_IMMATURE_CELLS\SUSPECTED_PATHOLOGY^^EXTREM_NEUTROPENIA|I
M|1|REAGENT|CLEANER\DILUENT\LYSE|240415I1(^20240902000000^20241202\240423H1(^20240905000000^20250305\240411M11^20240828000000^20241028
R|1|^^^MCV^787-2|56.1|um3|80.0 - 96.0^REFERENCE_RANGE|L||W||LABOR^^USER|20240906090745||
R|2|^^^NEU#^751-8|0.00|10E3/uL|1.60 - 7.00^REFERENCE_RANGE|LL||W||LABOR^^USER|20240906090745||
R|3|^^^NEU%^770-8|12.5|%|40.0 - 73.0^REFERENCE_RANGE|L||W||LABOR^^USER|20240906090745||
R|4|^^^RDW-CV^788-0|23.0|%|11.0 - 17.0^REFERENCE_RANGE|H||F||LABOR^^USER|20240906090745||
L|1|N`
	messageType, err := IdentifyMessage([]byte(astm), EncodingUTF8)

	assert.Nil(t, err)
	assert.Equal(t, MessageTypeOrdersAndResults, messageType)
}
