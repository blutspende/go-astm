package functions

import (
	"bytes"
	"fmt"
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
)

func ConvertFromEncodingToUtf8(input []byte, encoding constants.Encoding) (output string, err error) {
	cmap, err := findCharmapForEncoding(encoding)
	if err != nil {
		return "", err
	}
	if cmap == nil {
		return string(input), nil
	}
	encoded, err := io.ReadAll(cmap.NewDecoder().Reader(bytes.NewReader(input)))
	return string(encoded), err
}

func ConvertFromUtf8ToEncoding(input string, encoding constants.Encoding) (output []byte, err error) {
	cmap, err := findCharmapForEncoding(encoding)
	if err != nil {
		return []byte{}, err
	}
	if cmap == nil {
		return []byte(input), nil
	}
	output, _, err = transform.Bytes(cmap.NewEncoder(), []byte(input))
	return output, err
}

func findCharmapForEncoding(encoding constants.Encoding) (*charmap.Charmap, error) {
	switch encoding {
	case constants.ENCODING_UTF8:
		return nil, nil
	case constants.ENCODING_ASCII:
		return nil, nil
	case constants.ENCODING_WINDOWS1250:
		return charmap.Windows1250, nil
	case constants.ENCODING_WINDOWS1251:
		return charmap.Windows1251, nil
	case constants.ENCODING_WINDOWS1252:
		return charmap.Windows1252, nil
	case constants.ENCODING_DOS852:
		return charmap.CodePage852, nil
	case constants.ENCODING_DOS855:
		return charmap.CodePage855, nil
	case constants.ENCODING_DOS866:
		return charmap.CodePage866, nil
	case constants.ENCODING_ISO8859_1:
		return charmap.ISO8859_1, nil
	default:
		return nil, fmt.Errorf(errmsg.Encoding_MsgInvalidEncoding, encoding)
	}
}
