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
