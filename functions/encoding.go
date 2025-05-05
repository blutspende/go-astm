package functions

import (
	"bytes"
	"fmt"
	encodingconst "github.com/blutspende/go-astm/enums/encoding"
	"github.com/blutspende/go-astm/errmsg"
	"github.com/blutspende/go-astm/models/astmmodels"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
)

func ConvertFromEncodingToUtf8(input []byte, config *astmmodels.Configuration) (output string, err error) {
	cmap, err := findCharmapForEncoding(config.Encoding)
	if err != nil {
		return "", err
	}
	if cmap == nil {
		return string(input), nil
	}
	encoded, err := io.ReadAll(cmap.NewDecoder().Reader(bytes.NewReader(input)))
	return string(encoded), err
}

func ConvertFromUtf8ToEncoding(input string, config *astmmodels.Configuration) (output []byte, err error) {
	cmap, err := findCharmapForEncoding(config.Encoding)
	if err != nil {
		return []byte{}, err
	}
	if cmap == nil {
		return []byte(input), nil
	}
	output, _, err = transform.Bytes(cmap.NewEncoder(), []byte(input))
	return output, err
}

func ConvertArrayFromUtf8ToEncoding(input []string, config *astmmodels.Configuration) (output [][]byte, err error) {
	output = make([][]byte, len(input))
	for i, line := range input {
		output[i], err = ConvertFromUtf8ToEncoding(line, config)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func findCharmapForEncoding(encoding string) (*charmap.Charmap, error) {
	switch encoding {
	case encodingconst.UTF8:
		return nil, nil
	case encodingconst.ASCII:
		return nil, nil
	case encodingconst.Windows1250:
		return charmap.Windows1250, nil
	case encodingconst.Windows1251:
		return charmap.Windows1251, nil
	case encodingconst.Windows1252:
		return charmap.Windows1252, nil
	case encodingconst.DOS852:
		return charmap.CodePage852, nil
	case encodingconst.DOS855:
		return charmap.CodePage855, nil
	case encodingconst.DOS866:
		return charmap.CodePage866, nil
	case encodingconst.ISO8859_1:
		return charmap.ISO8859_1, nil
	case "IBM037":
		return charmap.CodePage037, nil
	case "IBM437":
		return charmap.CodePage437, nil
	case "IBM850":
		return charmap.CodePage850, nil
	case "IBM852":
		return charmap.CodePage852, nil
	case "IBM855":
		return charmap.CodePage855, nil
	case "IBM858":
		return charmap.CodePage858, nil
	case "IBM860":
		return charmap.CodePage860, nil
	case "IBM862":
		return charmap.CodePage862, nil
	case "IBM863":
		return charmap.CodePage863, nil
	case "IBM865":
		return charmap.CodePage865, nil
	case "IBM866":
		return charmap.CodePage866, nil
	case "IBM1047":
		return charmap.CodePage1047, nil
	case "IBM1140":
		return charmap.CodePage1140, nil
	case "ISO8859-2":
		return charmap.ISO8859_2, nil
	case "ISO8859-3":
		return charmap.ISO8859_3, nil
	case "ISO8859-4":
		return charmap.ISO8859_4, nil
	case "ISO8859-5":
		return charmap.ISO8859_5, nil
	case "ISO8859-6":
		return charmap.ISO8859_6, nil
	case "ISO8859-6E":
		return nil, fmt.Errorf("ISO8859-6E is not supported as *charmap.Charmap")
	case "ISO8859-6I":
		return nil, fmt.Errorf("ISO8859-6I is not supported as *charmap.Charmap")
	case "ISO8859-7":
		return charmap.ISO8859_7, nil
	case "ISO8859-8":
		return charmap.ISO8859_8, nil
	case "ISO8859-8E":
		return nil, fmt.Errorf("ISO8859-8E is not supported as *charmap.Charmap")
	case "ISO8859-8I":
		return nil, fmt.Errorf("ISO8859-8I is not supported as *charmap.Charmap")
	case "ISO8859-9":
		return charmap.ISO8859_9, nil
	case "ISO8859-10":
		return charmap.ISO8859_10, nil
	case "ISO8859-13":
		return charmap.ISO8859_13, nil
	case "ISO8859-14":
		return charmap.ISO8859_14, nil
	case "ISO8859-15":
		return charmap.ISO8859_15, nil
	case "ISO8859-16":
		return charmap.ISO8859_16, nil
	case "KOI8-R":
		return charmap.KOI8R, nil
	case "KOI8-U":
		return charmap.KOI8U, nil
	case "Macintosh":
		return charmap.Macintosh, nil
	case "MacintoshCyrillic":
		return charmap.MacintoshCyrillic, nil
	case "Windows-874":
		return charmap.Windows874, nil
	case "Windows-1250":
		return charmap.Windows1250, nil
	case "Windows-1251":
		return charmap.Windows1251, nil
	case "Windows-1252":
		return charmap.Windows1252, nil
	case "Windows-1253":
		return charmap.Windows1253, nil
	case "Windows-1254":
		return charmap.Windows1254, nil
	case "Windows-1255":
		return charmap.Windows1255, nil
	case "Windows-1256":
		return charmap.Windows1256, nil
	case "Windows-1257":
		return charmap.Windows1257, nil
	case "Windows-1258":
		return charmap.Windows1258, nil
	default:
		return nil, fmt.Errorf("%s: %w", encoding, errmsg.ErrEncodingInvalidEncoding)
	}
}

// List of encodings from skeleton
//ASCII
//DOS852
//DOS855
//DOS866
//ISOLatin1
//ISOLatin2
//ISOLatin3
//ISOLatin4
//ISOLatinCyrillic
//ISOLatinArabic
//ISOLatinGreek
//ISOLatinHebrew
//ISOLatin5
//ISOLatin6
//ISOTextComm
//HalfWidthKatakana
//JISEncoding
//ShiftJIS
//EUCPkdFmtJapanese
//EUCFixWidJapanese
//ISO4UnitedKingdom
//ISO11SwedishForNames
//ISO15Italian
//ISO17Spanish
//ISO21German
//ISO60Norwegian1
//ISO69French
//ISO10646UTF1
//ISO646basic1983
//INVARIANT
//ISO2IntlRefVersion
//NATSSEFI
//NATSSEFIADD
//NATSDANO
//NATSDANOADD
//ISO10Swedish
//KSC56011987
//ISO2022KR
//EUCKR
//ISO2022JP
//ISO2022JP2
//ISO13JISC6220jp
//ISO14JISC6220ro
//ISO16Portuguese
//ISO18Greek7Old
//ISO19LatinGreek
//ISO25French
//ISO27LatinGreek1
//ISO5427Cyrillic
//ISO42JISC62261978
//ISO47BSViewdata
//ISO49INIS
//ISO50INIS8
//ISO51INISCyrillic
//ISO54271981
//ISO5428Greek
//ISO57GB1988
//ISO58GB231280
//ISO61Norwegian2
//ISO70VideotexSupp1
//ISO84Portuguese2
//ISO85Spanish2
//ISO86Hungarian
//ISO87JISX0208
//ISO88Greek7
//ISO89ASMO449
//ISO90
//ISO91JISC62291984a
//ISO92JISC62991984b
//ISO93JIS62291984badd
//ISO94JIS62291984hand
//ISO95JIS62291984handadd
//ISO96JISC62291984kana
//ISO2033
//ISO99NAPLPS
//ISO102T617bit
//ISO103T618bit
//ISO111ECMACyrillic
//ISO121Canadian1
//ISO122Canadian2
//ISO123CSAZ24341985gr
//ISO88596E
//ISO88596I
//ISO128T101G2
//ISO88598E
//ISO88598I
//ISO139CSN369103
//ISO141JUSIB1002
//ISO143IECP271
//ISO146Serbian
//ISO147Macedonian
//ISO150GreekCCITT
//ISO151Cuba
//ISO6937Add
//ISO153GOST1976874
//ISO8859Supp
//ISO10367Box
//ISO158Lap
//ISO159JISX02121990
//ISO646Danish
//USDK
//DKUS
//KSC5636
//Unicode11UTF7
//ISO2022CN
//ISO2022CNEXT
//UTF8
//ISO885913
//ISO885914
//ISO885915
//ISO885916
//GBK
//GB18030
//OSDEBCDICDF0415
//OSDEBCDICDF03IRV
//OSDEBCDICDF041
//ISO115481
//KZ1048
//Unicode
//UCS4
//UnicodeASCII
//UnicodeLatin1
//UnicodeJapanese
//UnicodeIBM1261
//UnicodeIBM1268
//UnicodeIBM1276
//UnicodeIBM1264
//UnicodeIBM1265
//Unicode11
//SCSU
//UTF7
//UTF16BE
//UTF16LE
//UTF16
//CESU8
//UTF32
//UTF32BE
//UTF32LE
//BOCU1
//UTF7IMAP
//Windows30Latin1
//Windows31Latin1
//Windows31Latin2
//Windows31Latin5
//HPRoman8
//AdobeStandardEncoding
//VenturaUS
//VenturaInternational
//DECMCS
//PC850Multilingual
//PC8DanishNorwegian
//PC862LatinHebrew
//PC8Turkish
//IBMSymbols
//IBMThai
//HPLegal
//HPPiFont
//HPMath8
//HPPSMath
//HPDesktop
//VenturaMath
//MicrosoftPublishing
//Windows31J
//GB2312
//Big5
//Macintosh
//IBM037
//IBM038
//IBM273
//IBM274
//IBM275
//IBM277
//IBM278
//IBM280
//IBM281
//IBM284
//IBM285
//IBM290
//IBM297
//IBM420
//IBM423
//IBM424
//PC8CodePage437
//IBM500
//IBM851
//PCp852
//IBM855
//IBM857
//IBM860
//IBM861
//IBM863
//IBM864
//IBM865
//IBM868
//IBM869
//IBM870
//IBM871
//IBM880
//IBM891
//IBM903
//IBBM904
//IBM905
//IBM918
//IBM1026
//IBMEBCDICATDE
//EBCDICATDEA
//EBCDICCAFR
//EBCDICDKNO
//EBCDICDKNOA
//EBCDICFISE
//EBCDICFISEA
//EBCDICFR
//EBCDICIT
//EBCDICPT
//EBCDICES
//EBCDICESA
//EBCDICESS
//EBCDICUK
//EBCDICUS
//Unknown8BiT
//Mnemonic
//Mnem
//VISCII
//VIQR
//KOI8R
//HZGB2312
//IBM866
//PC775Baltic
//KOI8U
//IBM00858
//IBM00924
//IBM01140
//IBM01141
//IBM01142
//IBM01143
//IBM01144
//IBM01145
//IBM01146
//IBM01147
//IBM01148
//IBM01149
//Big5HKSCS
//IBM1047
//PTCP154
//Amiga1251
//KOI7switched
//BRF
//TSCII
//CP51932
//Windows874
//Windows1250
//Windows1251
//Windows1252
//Windows1253
//Windows1254
//Windows1255
//Windows1256
//Windows1257
//Windows1258
//TIS620
//CP50220
//ISO8859-1
