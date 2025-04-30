package astm

import (
	"github.com/blutspende/go-astm/v3/enums/messagetype"
	"github.com/blutspende/go-astm/v3/functions"
	"github.com/blutspende/go-astm/v3/models/astmmodels"
	"regexp"
)

func IdentifyMessage(messageData []byte, configuration ...astmmodels.Configuration) (messageType string, err error) {
	// Load configuration
	config, err := loadConfiguration(configuration...)
	if err != nil {
		return "", err
	}
	// Convert encoding to UTF8
	utf8Data, err := functions.ConvertFromEncodingToUtf8(messageData, config)
	if err != nil {
		return "", err
	}
	// Split the message data into lines
	lines, err := functions.SliceLines(utf8Data, config)
	if err != nil {
		return "", err
	}
	// Extract the first characters from each line
	firstChars := ""
	for _, line := range lines {
		if len(line) > 0 {
			firstChars += string(line[0])
		}
	}
	// TODO: verify these regexes to be correct
	// Set up the possible message types regexes
	expressionQuery := "^(HQ+)+L?$"
	expressionOrder := "^(H(PM?C?M?OM?C?M?)+)+L?$"
	expressionOrderAndResult := "^(H(PM*C?M*OM*C?M*(RM*C?M*)+)+)+L?$"
	expressionManyOrderAndResult := "^(H(PM*C?M*(OM*C?M*(RM*C?M*)+)*)+)L?$"
	// Check the first characters against the regexes and return the message type
	switch {
	case regexp.MustCompile(expressionQuery).MatchString(firstChars):
		return messagetype.Query, nil
	case regexp.MustCompile(expressionOrder).MatchString(firstChars):
		return messagetype.Order, nil
	case regexp.MustCompile(expressionOrderAndResult).MatchString(firstChars):
		return messagetype.Result, nil
	case regexp.MustCompile(expressionManyOrderAndResult).MatchString(firstChars):
		return messagetype.Result, nil
	}
	// If no match was found return unknown
	return messagetype.Unidentified, err
}
