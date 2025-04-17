package functions

import (
	"github.com/blutspende/go-astm/v2/constants/astmconst"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/blutspende/go-astm/v2/models"
	"strings"
)

func SliceLines(input string, config *models.Configuration) (output []string, err error) {
	// Check for empty input
	if input == "" {
		return nil, errmsg.Lining_ErrNotEnoughLines
	}

	// A line separator has to be provided if auto-detect is disabled
	if !config.AutoDetectLineSeparator && config.LineSeparator == "" {
		return nil, errmsg.Lining_ErrNoLineSeparator
	}

	var lines []string
	if !config.AutoDetectLineSeparator {
		// Line separator provided in config, no auto-detect
		lines = strings.Split(input, config.LineSeparator)
	} else {
		// Auto-detect line separator
		lfCnt := 0
		crCnt := 0
		for _, c := range input {
			if c == rune(astmconst.LF[0]) {
				lfCnt++
			} else if c == rune(astmconst.CR[0]) {
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
			input = strings.ReplaceAll(input, astmconst.CR, astmconst.LF)

		} else {
			input = strings.ReplaceAll(input, astmconst.CR, "")
		}

		lines = strings.Split(input, astmconst.LF)
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
	linebreak := astmconst.LF
	if config.LineSeparator != "" && !config.AutoDetectLineSeparator {
		linebreak = config.LineSeparator
	}

	for i, line := range input {
		output += line
		if i < len(input)-1 {
			output += linebreak
		}
	}

	return output
}
