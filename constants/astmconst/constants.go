package astmconst

// Message format attributes

const ATTRIBUTE_REQUIRED string = "required" // field-annotation: by default all fields are optinal
const ATTRIBUTE_OPTIONAL string = "optional" // record-annotation: by default all records are mandatory
const ATTRIBUTE_LONGDATE string = "longdate" // Indicating that the date should be formatted as date and time (output only)
const ATTRIBUTE_LENGTH string = "length"     // used for specifying the decimal length of float fields - astm:"1,length:2" (output only)
const ATTRIBUTE_SUBNAME string = "subname"   // used for specifying a subname for a record - astm:"M,subname:MATRIX"

// Public functions parameters

const ENCODING_UTF8 string = "UTF8"
const ENCODING_ASCII string = "ASCII"
const ENCODING_WINDOWS1250 string = "Windows1250"
const ENCODING_WINDOWS1251 string = "Windows1251"
const ENCODING_WINDOWS1252 string = "Windows1252"
const ENCODING_DOS852 string = "DOS852"
const ENCODING_DOS855 string = "DOS855"
const ENCODING_DOS866 string = "DOS866"
const ENCODING_ISO8859_1 string = "ISO8859-1"

const TIMEZONE_UTC string = "UTC"
const TIMEZONE_EUROPE_BERLIN string = "Europe/Berlin"
const TIMEZONE_EUROPE_BUDAPEST string = "Europe/Budapest"
const TIMEZONE_EUROPE_LONDON string = "Europe/London"

var LF string = string(byte(0x0A))
var CR string = string(byte(0x0D))
var LFCR string = string([]byte{byte(0x0A), byte(0x0D)})
var CRLF string = string([]byte{byte(0x0D), byte(0x0A)})

// NOTATION_STANDARD will always produce as many delimiters as there are values in the export-format
// NOTATION_SHORT will skip all delimiters to the right of the last value
const NOTATION_STANDARD string = "STANDARD"
const NOTATION_SHORT string = "SHORT"

// Possible results from IdentifyMessage

const MESSAGETYPE_UNIDENTIFIED string = "UNIDENTIFIED"
const MESSAGETYPE_QUERY string = "QUERY"
const MESSAGETYPE_ORDER string = "ORDER"
const MESSAGETYPE_RESULT string = "RESULT"
