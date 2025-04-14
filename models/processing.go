package models

import "github.com/blutspende/go-astm/v2/constants"

// Configuration struct for the whole process
type Configuration struct {
	Encoding                   string
	LineSeparator              string
	TimeZone                   string
	EnforceSequenceNumberCheck bool
	Internal                   InternalConfiguration
	Notation                   string
}

var DefaultConfiguration = Configuration{
	Encoding:                   "ISO8859-1",
	LineSeparator:              "",
	TimeZone:                   "Europe/Berlin",
	EnforceSequenceNumberCheck: true,
	Notation:                   constants.NOTATION_STANDARD,
}
