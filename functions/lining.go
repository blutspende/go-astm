package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/blutspende/go-astm/v2/models"
	"strings"
)

func SliceLines(input string, config *models.Configuration) (output []string, err error) {
	// Check for empty input
	if input == "" {
		return nil, errmsg.Lining_ErrNotEnoughLines
	}

	var lines []string
	if config.LineSeparator != "" {
		// Line separator provided in config, split by it
		lines = strings.Split(input, config.LineSeparator)
	} else {
		// No line separator provided: default behavior
		lfCnt := 0
		crCnt := 0
		for _, c := range input {
			if c == rune(constants.LF[0]) {
				lfCnt++
			} else if c == rune(constants.CR[0]) {
				crCnt++
			}
		}
		if lfCnt == 0 && crCnt == 0 {
			return nil, errmsg.Lining_ErrInvalidLinebreak
		}
		if lfCnt > 0 && crCnt > 0 && lfCnt != crCnt {
			return nil, errmsg.Lining_ErrInvalidLinebreak
		}

		if lfCnt == 0 {
			input = strings.ReplaceAll(input, constants.CR, constants.LF)

		} else {
			input = strings.ReplaceAll(input, constants.CR, "")
		}

		lines = strings.Split(input, constants.LF)
	}

	for i := range lines {
		lines[i] = strings.Trim(lines[i], " ")
		if lines[i] != "" {
			output = append(output, lines[i])
		}
	}

	return output, nil
}

func BuildLines(input []string, config *models.Configuration) (output string) {
	linebreak := ""
	if config.LineSeparator != "" {
		linebreak = config.LineSeparator
	} else {
		linebreak = constants.LF
	}

	for i, line := range input {
		output += line
		if i < len(input)-1 {
			output += linebreak
		}
	}

	return output
}
