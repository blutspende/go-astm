package astm

import (
	"github.com/blutspende/go-astm/models/astmmodels"
	"time"
)

func loadConfiguration(configuration ...astmmodels.Configuration) (config *astmmodels.Configuration, err error) {
	if len(configuration) > 0 {
		config = &configuration[0]
	} else {
		config = &astmmodels.DefaultConfiguration
	}
	if config.Delimiters.Field == "" ||
		config.Delimiters.Repeat == "" ||
		config.Delimiters.Component == "" ||
		config.Delimiters.Escape == "" {
		config.Delimiters = astmmodels.DefaultDelimiters
	}
	config.TimeLocation, err = time.LoadLocation(config.TimeZone)
	if err != nil {
		return nil, err
	}
	return config, nil
}
