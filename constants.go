package astm

const (
	ANNOTATION_DELIMITER = "delimiter" // annotation that triggers the delimiters in the scanner to be reset
	ANNOTATION_REQUIRED  = "require"   // field-annotation: by default all fields are optinal
	ANNOTATION_OPTIONAL  = "optional"  // record-annotation: by default all records are mandatory
	ANNOTATION_SEQUENCE  = "sequence"  // indicating that a sequence number should be generated (output only)
	ANNOTATION_LONGDATE  = "longdate"
	ANNOTATION_LENGTH    = "length" // used for specifying the decimal length of float fields - astm:"1,length:2" (output only)
)

type Encoding int

const EncodingUTF8 Encoding = 1
const EncodingASCII Encoding = 2
const EncodingWindows1250 Encoding = 3
const EncodingWindows1251 Encoding = 4
const EncodingWindows1252 Encoding = 5
const EncodingDOS852 Encoding = 6
const EncodingDOS855 Encoding = 7
const EncodingDOS866 Encoding = 8
const EncodingISO8859_1 Encoding = 9

type Timezone string

const TimezoneUTC Timezone = "UTC"
const TimezoneEuropeBerlin Timezone = "Europe/Berlin"
const TimezoneEuropeBudapest Timezone = "Europe/Budapest"
const TimezoneEuropeLondon Timezone = "Europe/London"

type LineBreak int

const CR LineBreak = 0x0D
const LF LineBreak = 0x0A
const CRLF LineBreak = 0x0D0A

/*
	Notation defines how the output format is build

ShortNotation will skip all delimiters to the right of the last value
StandardNotation will always produce as many delimiters as there are values in the export-format
*/
type Notation int

const StandardNotation = 1
const ShortNotation = 2
