package models

import "github.com/blutspende/go-astm/v2/constants"

// Configuration struct for the whole process
type Configuration struct {
	Encoding                   string
	LineSeparator              string
	AutoDetectLineSeparator    bool
	TimeZone                   string
	EnforceSequenceNumberCheck bool
	Notation                   string
	Internal                   InternalConfiguration
}

var DefaultConfiguration = Configuration{
	Encoding:                   constants.ENCODING_ISO8859_1,
	LineSeparator:              constants.LF,
	AutoDetectLineSeparator:    true,
	TimeZone:                   constants.TIMEZONE_EUROPE_BERLIN,
	EnforceSequenceNumberCheck: true,
	Notation:                   constants.NOTATION_STANDARD,
}
