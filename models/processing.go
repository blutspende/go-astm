package models

// Configuration struct for the whole process
type Configuration struct {
	Encoding                   string
	LineSeparator              string
	TimeZone                   string
	EnforceSequenceNumberCheck bool
	Internal                   InternalConfiguration
}

var DefaultConfiguration = Configuration{
	Encoding:                   "ISO8859-1",
	LineSeparator:              "",
	TimeZone:                   "Europe/Berlin",
	EnforceSequenceNumberCheck: true,
}
