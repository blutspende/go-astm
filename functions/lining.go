package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"strings"
)

func SliceLines(input string) (output []string, err error) {
	if input == "" {
		return nil, errmsg.Parsing_ErrNotEnoughLines
	}

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
		return nil, errmsg.Parsing_ErrInvalidLinebreak
	}
	if lfCnt > 0 && crCnt > 0 && lfCnt != crCnt {
		return nil, errmsg.Parsing_ErrInvalidLinebreak
	}

	if lfCnt == 0 {
		input = strings.ReplaceAll(input, constants.CR, constants.LF)

	} else {
		input = strings.ReplaceAll(input, constants.CR, "")
	}

	lines := strings.Split(input, constants.LF)

	for i := range lines {
		lines[i] = strings.Trim(lines[i], " ")
		if lines[i] != "" {
			output = append(output, lines[i])
		}
	}

	return output, nil
}

func BuildLines(input []string, linebreak string) (output string) {
	for i, line := range input {
		output += line
		if i < len(input)-1 {
			output += linebreak
		}
	}
	return output
}
