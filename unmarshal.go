package astm

import (
	"github.com/blutspende/go-astm/v3/functions"
	"github.com/blutspende/go-astm/v3/models/astmmodels"
)

func Unmarshal(messageData []byte, targetStruct interface{}, configuration ...astmmodels.Configuration) (err error) {
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
