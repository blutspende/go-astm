package astm

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type MessageType int

const MessageTypeUnkown MessageType = -1
const MessageTypeQuery MessageType = 1
const MessageTypeOrdersOnly MessageType = 2
const MessageTypeOrdersAndResults MessageType = 3

func IdentifyMessage(messageEncoded []byte, enc Encoding) (MessageType, error) {

	messageBytes, err := utilityConvertByteArrayToUTF(messageEncoded, enc)
	if err != nil {
		return MessageTypeUnkown, err
	}

	bufferedInputLinesWithEmptyLines := strings.Split(string(messageBytes), string([]byte{0x0A})) // copy
	if len(bufferedInputLinesWithEmptyLines) <= 1 {                                               // if it was not possible to break with non-standard 0x0a line-break try 0d (standard)
		bufferedInputLinesWithEmptyLines = strings.Split(string(messageBytes), string([]byte{0x0D}))
	}
	// strip the remaining 0A and 0D Linefeed at the end
	for i := 0; i < len(bufferedInputLinesWithEmptyLines); i++ {
		// 0d,0a then again as there have been files observed which had 0a0d (0d0a would be normal)
		bufferedInputLinesWithEmptyLines[i] = strings.Trim(bufferedInputLinesWithEmptyLines[i], string([]byte{0x0A}))
		bufferedInputLinesWithEmptyLines[i] = strings.Trim(bufferedInputLinesWithEmptyLines[i], string([]byte{0x0D}))
		bufferedInputLinesWithEmptyLines[i] = strings.Trim(bufferedInputLinesWithEmptyLines[i], string([]byte{0x0A}))
		bufferedInputLinesWithEmptyLines[i] = strings.Trim(bufferedInputLinesWithEmptyLines[i], string([]byte{0x0D}))
	}

	// remove empty lines
	bufferedInputLines := []string{}
	for i := range bufferedInputLinesWithEmptyLines {
		if strings.Trim(bufferedInputLinesWithEmptyLines[i], " ") != "" {
			bufferedInputLines = append(bufferedInputLines, strings.TrimSpace(bufferedInputLinesWithEmptyLines[i]))
		}
	}

	genome := ""
	for _, line := range bufferedInputLines {
		if len(line) > 1 {
			genome = genome + string(line[0])
		}
	}

	expressionQuery := "^(HQ+)+L?$"
	expressionOrder := "^(H(PM?C?M?OM?C?M?)+)+L?$"
	expressionOrderAndResult := "^(H(PM*C?M*OM*C?M*(RM*C?M*)+)+)+L?$"
	expressionManyOrderAndResult := "^(H(PM*C?M*(OM*C?M*(RM*C?M*)+)*)+)L?$"

	if match, _ := regexp.MatchString(expressionQuery, genome); match {
		return MessageTypeQuery, nil
	}

	if match, _ := regexp.MatchString(expressionOrder, genome); match {
		return MessageTypeOrdersOnly, nil
	}

	if match, _ := regexp.MatchString(expressionOrderAndResult, genome); match {
		return MessageTypeOrdersAndResults, nil
	}

	if match, _ := regexp.MatchString(expressionManyOrderAndResult, genome); match {
		return MessageTypeOrdersAndResults, nil
	}

	return MessageTypeUnkown, nil
}

func utilityConvertByteArrayToUTF(messageData []byte, fromEncoding Encoding) (string, error) {

	var (
		messageBytes []byte
		err          error
	)

	switch fromEncoding {
	case EncodingUTF8:
		messageBytes = messageData
	case EncodingASCII:
		messageBytes = messageData
	case EncodingDOS866:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.CodePage866, messageData); err != nil {
			return "", err
		}
	case EncodingDOS855:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.CodePage855, messageData); err != nil {
			return "", err
		}
	case EncodingDOS852:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.CodePage852, messageData); err != nil {
			return "", err
		}
	case EncodingWindows1250:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.Windows1250, messageData); err != nil {
			return "", err
		}
	case EncodingWindows1251:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.Windows1251, messageData); err != nil {
			return "", err
		}
	case EncodingWindows1252:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.Windows1252, messageData); err != nil {
			return "", err
		}
	case EncodingISO8859_1:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.ISO8859_1, messageData); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("invalid Codepage Id='%d' - %w", fromEncoding, err)
	}

	return string(messageBytes), nil
}
