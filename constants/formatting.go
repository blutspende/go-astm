package constants

type Encoding string

const ENCODING_UTF8 Encoding = "UTF8"
const ENCODING_ASCII Encoding = "ASCII"
const ENCODING_WINDOWS1250 Encoding = "Windows1250"
const ENCODING_WINDOWS1251 Encoding = "Windows1251"
const ENCODING_WINDOWS1252 Encoding = "Windows1252"
const ENCODING_DOS852 Encoding = "DOS852"
const ENCODING_DOS855 Encoding = "DOS855"
const ENCODING_DOS866 Encoding = "DOS866"
const ENCODING_ISO8859_1 Encoding = "ISO8859-1"

type Timezone string

const TIMEZONE_UTC Timezone = "UTC"
const TIMEZONE_EUROPE_BERLIN Timezone = "Europe/Berlin"
const TIMEZONE_EUROPE_BUDAPEST Timezone = "Europe/Budapest"
const TIMEZONE_EUROPE_LONDON Timezone = "Europe/London"

var LF string = string(byte(0x0A))
var CR string = string(byte(0x0D))
var LFCR string = string([]byte{byte(0x0A), byte(0x0D)})
var CRLF string = string([]byte{byte(0x0D), byte(0x0A)})

// Notation defines how the output format is build
// StandardNotation will always produce as many delimiters as there are values in the export-format
// ShortNotation will skip all delimiters to the right of the last value
type Notation int

const StandardNotation = 1
const ShortNotation = 2
