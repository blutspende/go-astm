package astmmodels

import (
	"github.com/blutspende/bloodlab-common/encoding"
	"github.com/blutspende/bloodlab-common/timezone"
	"github.com/blutspende/go-astm/v3/enums/lineseparator"
	"github.com/blutspende/go-astm/v3/enums/notation"
	"time"
)

// Configuration struct for the whole process
type Configuration struct {
	Encoding                   encoding.Encoding
	LineSeparator              string
	AutoDetectLineSeparator    bool
	TimeZone                   timezone.TimeZone
	EnforceSequenceNumberCheck bool
	Notation                   string
	DefaultDecimalPrecision    int
	RoundLastDecimal           bool
	KeepShortDateTimeZone      bool
	Delimiters                 Delimiters
	TimeLocation               *time.Location
}

var DefaultConfiguration = Configuration{
	Encoding:                   encoding.ISO8859_1,
	LineSeparator:              lineseparator.LF,
	AutoDetectLineSeparator:    true,
	TimeZone:                   timezone.EuropeBerlin,
	EnforceSequenceNumberCheck: true,
	Notation:                   notation.Standard,
	DefaultDecimalPrecision:    3,
	RoundLastDecimal:           true,
	KeepShortDateTimeZone:      true,
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
