package models

import (
	"github.com/blutspende/go-astm/v2/constants/astmconst"
)

// Configuration struct for the whole process
type Configuration struct {
	Encoding                   string
	LineSeparator              string
	AutoDetectLineSeparator    bool
	TimeZone                   string
	EnforceSequenceNumberCheck bool
	Notation                   string
	RoundFixedNumbers          bool
	Internal                   InternalConfiguration
}

var DefaultConfiguration = Configuration{
	Encoding:                   astmconst.ENCODING_ISO8859_1,
	LineSeparator:              astmconst.LF,
	AutoDetectLineSeparator:    true,
	TimeZone:                   astmconst.TIMEZONE_EUROPE_BERLIN,
	EnforceSequenceNumberCheck: true,
	Notation:                   astmconst.NOTATION_STANDARD,
	RoundFixedNumbers:          true,
}
