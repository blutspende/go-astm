package astm

import (
	"github.com/blutspende/go-astm/v2/functions"
	"github.com/blutspende/go-astm/v2/models"
	"regexp"
	"time"
)

func Unmarshal(messageData []byte, targetStruct interface{}, configuration ...*models.Configuration) (err error) {
	// Load configuration
	config, err := loadConfiguration(configuration...)
	if err != nil {
		return err
	}
	// Convert encoding to UTF8
	utf8Data, err := functions.ConvertFromEncodingToUtf8(messageData, config)
	if err != nil {
		return err
	}
	// Split the message data into lines
	lines, err := functions.SliceLines(utf8Data, config)
	if err != nil {
		return err
	}
	// Parse the lines into the target structure
	lineIndex := 0
	err = functions.ParseStruct(lines, targetStruct, &lineIndex, 1, 0, config)
	if err != nil {
		return err
	}
	// Return nil if everything went well
	return nil
}

func Marshal(sourceStruct interface{}, configuration ...*models.Configuration) (result [][]byte, err error) {
	// Load configuration
	config, err := loadConfiguration(configuration...)
	if err != nil {
		return nil, err
	}
	// Build the lines from the source structure
	lines, err := functions.BuildStruct(sourceStruct, 1, 0, config)
	if err != nil {
		return nil, err
	}
	// Convert UTF8 string array to encoding
	result, err = functions.ConvertArrayFromUtf8ToEncoding(lines, config)
	if err != nil {
		return nil, err
	}
	// Return the result and no error if everything went well
	return result, nil
}

func IdentifyMessage(messageData []byte, configuration ...*models.Configuration) (messageType MessageType, err error) {
	// Load configuration
	config, err := loadConfiguration(configuration...)
	if err != nil {
		return MESSAGETYPE_UNKOWN, err
	}
	// Convert encoding to UTF8
	utf8Data, err := functions.ConvertFromEncodingToUtf8(messageData, config)
	if err != nil {
		return MESSAGETYPE_UNKOWN, err
	}
	// Split the message data into lines
	lines, err := functions.SliceLines(utf8Data, config)
	if err != nil {
		return MESSAGETYPE_UNKOWN, err
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
		return MESSAGETYPE_QUERY, nil
	case regexp.MustCompile(expressionOrder).MatchString(firstChars):
		return MESSAGETYPE_ORDERS_ONLY, nil
	case regexp.MustCompile(expressionOrderAndResult).MatchString(firstChars):
		return MESSAGETYPE_ORDERS_AND_RESULTS, nil
	case regexp.MustCompile(expressionManyOrderAndResult).MatchString(firstChars):
		return MESSAGETYPE_ORDERS_AND_RESULTS, nil
	}
	// If no match was found return unknown
	return MESSAGETYPE_UNKOWN, err
}

func loadConfiguration(configuration ...*models.Configuration) (config *models.Configuration, err error) {
	if len(configuration) > 0 {
		config = configuration[0]
	} else {
		config = &models.DefaultConfiguration
		config.Internal.Delimiters = models.DefaultDelimiters
	}
	config.Internal.TimeLocation, err = time.LoadLocation(config.TimeZone)
	if err != nil {
		return nil, err
	}
	return config, nil
}
