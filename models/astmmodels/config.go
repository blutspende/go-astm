package astmmodels

import (
	"github.com/blutspende/go-astm/v2/constants/astmconst"
	"time"
)

// Configuration struct for the whole process
type Configuration struct {
	Encoding                   string
	LineSeparator              string
	AutoDetectLineSeparator    bool
	TimeZone                   string
	EnforceSequenceNumberCheck bool
	Notation                   string
	DefaultDecimalPrecision    int
	RoundLastDecimal           bool
	Delimiters                 Delimiters
	TimeLocation               *time.Location
}

var DefaultConfiguration = Configuration{
	Encoding:                   astmconst.ENCODING_ISO8859_1,
	LineSeparator:              astmconst.LF,
	AutoDetectLineSeparator:    true,
	TimeZone:                   astmconst.TIMEZONE_EUROPE_BERLIN,
	EnforceSequenceNumberCheck: true,
	Notation:                   astmconst.NOTATION_STANDARD,
	DefaultDecimalPrecision:    3,
	RoundLastDecimal:           true,
	Delimiters:                 DefaultDelimiters,
	TimeLocation:               nil,
}

// Delimiters used in ASTM parsing
type Delimiters struct {
	Field     string
	Repeat    string
	Component string
	Escape    string
}

var DefaultDelimiters = Delimiters{
	Field:     `|`,
	Repeat:    `\`,
	Component: `^`,
	Escape:    `&`,
}
