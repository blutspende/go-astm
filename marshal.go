package astm

import (
	"github.com/blutspende/go-astm/v3/functions"
	"github.com/blutspende/go-astm/v3/models/astmmodels"
)

func Marshal(sourceStruct interface{}, configuration ...astmmodels.Configuration) (result [][]byte, err error) {
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
