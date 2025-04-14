package models

import "github.com/blutspende/go-astm/v2/constants"

// Configuration struct for the whole process
type Configuration struct {
	Encoding                   string
	LineSeparator              string
	AutoDetectLineSeparator    bool
	TimeZone                   string
	EnforceSequenceNumberCheck bool
	Internal                   InternalConfiguration
	Notation                   string
}

var DefaultConfiguration = Configuration{
	Encoding:                   "ISO8859-1",
	LineSeparator:              constants.LF,
	AutoDetectLineSeparator:    true,
	TimeZone:                   "Europe/Berlin",
	EnforceSequenceNumberCheck: true,
	Notation:                   constants.NOTATION_STANDARD,
}
