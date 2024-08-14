package astm

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

/** Marshal - wrap datastructure to code
**/
func Marshal(message interface{}, enc Encoding, tz Timezone, notation Notation) ([][]byte, error) {

	// dereference for as long as we deal with pointers
	if reflect.TypeOf(message).Kind() == reflect.Ptr {
		// return Marshal(reflect.ValueOf(message).Elem(), enc, tz, notation)
		return [][]byte{}, fmt.Errorf("marshal can not be used with pointers")
	}

	if reflect.ValueOf(message).Kind() != reflect.Struct {
		return [][]byte{}, fmt.Errorf("can only marshal annotated structs (see readme)")
	}

	location, err := time.LoadLocation(string(tz))
	if err != nil {
		return [][]byte{}, err
	}

	// default delmiters. These will be overwritten by the first occurence of "delimter"-annotation
	repeatDelimiter := "\\"
	componentDelimiter := "^"
	escapeDelimiter := "&"

	buffer, err := iterateStructFieldsAndBuildOutput(message, 1, 1, enc, location, notation, &repeatDelimiter, &componentDelimiter, &escapeDelimiter)

	return buffer, err
}

type OutputRecord struct {
	Field, Repeat, Component int
	Value                    string
}

type OutputRecords []OutputRecord

func iterateStructFieldsAndBuildOutput(message interface{}, depth, sequence int, enc Encoding, location *time.Location, notation Notation,
	repeatDelimiter, componentDelimiter, escapeDelimiter *string) ([][]byte, error) {

	buffer := make([][]byte, 0)

	messageValue := reflect.ValueOf(message)
	messageType := reflect.TypeOf(message)

	for i := 0; i < messageValue.NumField(); i++ {

		currentRecord := messageValue.Field(i)
		recordAstmTag := messageType.Field(i).Tag.Get("astm")
		recordAstmTagsList := strings.Split(recordAstmTag, ",")

		if len(recordAstmTag) == 0 { // no annotation = Descend if its an array or a struct of such

			if currentRecord.Kind() == reflect.Slice { // array of something = iterate and recurse
				for x := 0; x < currentRecord.Len(); x++ {
					dood := currentRecord.Index(x).Interface()

					if bytes, err := iterateStructFieldsAndBuildOutput(dood, depth+1, (x + 1), enc, location, notation, repeatDelimiter, componentDelimiter, escapeDelimiter); err != nil {
						return nil, err
					} else {
						for line := 0; line < len(bytes); line++ {
							buffer = append(buffer, bytes[line])
						}
					}
				}
			} else if currentRecord.Kind() == reflect.Struct { // got the struct straignt = recurse directly

				if bytes, err := iterateStructFieldsAndBuildOutput(currentRecord.Interface(), depth+1, 1, enc, location, notation, repeatDelimiter, componentDelimiter, escapeDelimiter); err != nil {
					return nil, err
				} else {
					for line := 0; line < len(bytes); line++ {
						buffer = append(buffer, bytes[line])
					}
				}

			} else {
				return nil, fmt.Errorf("ivalid Datatype without any annotation '%s' - you can use struct or slices of structs", currentRecord.Kind())
			}

		} else {

			recordType := recordAstmTagsList[0]

			if currentRecord.Kind() == reflect.Slice { // it is an annotated slice
				if !currentRecord.IsNil() {
					for x := 0; x < currentRecord.Len(); x++ {
						outs, err := processOneRecord(recordType, currentRecord.Index(x), x+1, location, repeatDelimiter, componentDelimiter, escapeDelimiter, notation) // fmt.Println(outp)
						if err != nil {
							return nil, err
						}
						buffer = append(buffer, []byte(outs))
					}
				}
			} else {
				outs, err := processOneRecord(recordType, currentRecord, sequence, location, repeatDelimiter, componentDelimiter, escapeDelimiter, notation) // fmt.Println(outp)
				if err != nil {
					return nil, err
				}
				buffer = append(buffer, []byte(outs))
			}
		}

	}

	switch enc {
	case EncodingUTF8:
		// nothing
	case EncodingASCII:
		// nothing
	case EncodingDOS866:
		for i, x := range buffer {
			buffer[i] = EncodeUTF8ToCharset(charmap.CodePage866, x)
		}
	case EncodingDOS855:
		for i, x := range buffer {
			buffer[i] = EncodeUTF8ToCharset(charmap.CodePage855, x)
		}
	case EncodingDOS852:
		for i, x := range buffer {
			buffer[i] = EncodeUTF8ToCharset(charmap.CodePage852, x)
		}
	case EncodingWindows1250:
		for i, x := range buffer {
			buffer[i] = EncodeUTF8ToCharset(charmap.Windows1250, x)
		}
	case EncodingWindows1251:
		for i, x := range buffer {
			buffer[i] = EncodeUTF8ToCharset(charmap.Windows1251, x)
		}
	case EncodingWindows1252:
		for i, x := range buffer {
			buffer[i] = EncodeUTF8ToCharset(charmap.Windows1252, x)
		}
	case EncodingISO8859_1:
		for i, x := range buffer {
			buffer[i] = EncodeUTF8ToCharset(charmap.ISO8859_1, x)
		}
	default:
		return nil, fmt.Errorf("invalid Codepage Id='%d' in marshalling message", enc)
	}

	return buffer, nil
}

func EncodeUTF8ToCharset(charmap *charmap.Charmap, data []byte) []byte {
	e := charmap.NewEncoder()
	var b bytes.Buffer
	writer := transform.NewWriter(&b, e)
	writer.Write([]byte(data))
	resultdata := b.Bytes()
	writer.Close()
	return resultdata
}

func processOneRecord(recordType string, currentRecord reflect.Value, generatedSequenceNumber int, location *time.Location,
	repeatDelimiter, componentDelimiter, escapeDelimiter *string, notation Notation) (string, error) {

	if currentRecord.Kind() != reflect.Struct {
		return "", nil // beeing not a struct is not an error
	}

	fieldList := make(OutputRecords, 0)
	emptyRecords := make(OutputRecords, 0)
	isShortNotation := notation == ShortNotation

	for i := 0; i < currentRecord.NumField(); i++ {

		field := currentRecord.Field(i)
		fieldAstmTag := currentRecord.Type().Field(i).Tag.Get("astm")

		if fieldAstmTag == "" {
			continue
		}

		fieldAstmTagsList := strings.Split(fieldAstmTag, ",")

		fieldIdx, repeatIdx, componentIdx, err := readFieldAddressAnnotation(fieldAstmTagsList[0])
		if err != nil {
			return "", fmt.Errorf("invalid annotation for field %s : (%w)", currentRecord.Type().Field(i).Name, err)
		}

		switch field.Type().Kind() {
		case reflect.String:
			value := ""

			if sliceContainsString(fieldAstmTagsList, ANNOTATION_SEQUENCE) {
				return "", fmt.Errorf("invalid annotation %s for string-field", ANNOTATION_SEQUENCE)
			}

			// if no delimiters are given, default is \^&
			if sliceContainsString(fieldAstmTagsList, ANNOTATION_DELIMITER) && field.String() == "" {
				value = *repeatDelimiter + *componentDelimiter + *escapeDelimiter
			} else {
				value = field.String()
			}

			fieldList, emptyRecords = addASTMFieldToList(fieldList, emptyRecords, fieldIdx, repeatIdx, componentIdx, value, isShortNotation)
		case reflect.Int:
			value := fmt.Sprintf("%d", field.Int())
			if sliceContainsString(fieldAstmTagsList, ANNOTATION_SEQUENCE) {
				value = fmt.Sprintf("%d", generatedSequenceNumber)
				generatedSequenceNumber = generatedSequenceNumber + 1
			}

			fieldList, emptyRecords = addASTMFieldToList(fieldList, nil, fieldIdx, repeatIdx, componentIdx, value, isShortNotation)
		case reflect.Float32:
		case reflect.Float64:
			//TODO: add annotation for decimal length
			value := fmt.Sprintf("%.3f", field.Float())
			fieldList, emptyRecords = addASTMFieldToList(fieldList, nil, fieldIdx, repeatIdx, componentIdx, value, isShortNotation)
		case reflect.Struct:
			switch field.Type().Name() {
			case "Time":
				time := field.Interface().(time.Time)

				if !time.IsZero() {

					if sliceContainsString(fieldAstmTagsList, ANNOTATION_LONGDATE) {
						value := time.In(location).Format("20060102150405")
						fieldList, emptyRecords = addASTMFieldToList(fieldList, nil, fieldIdx, repeatIdx, componentIdx, value, isShortNotation)
					} else { // short date
						value := time.In(location).Format("20060102")
						fieldList, emptyRecords = addASTMFieldToList(fieldList, nil, fieldIdx, repeatIdx, componentIdx, value, isShortNotation)
					}
				} else {
					fieldList, emptyRecords = addASTMFieldToList(fieldList, nil, fieldIdx, repeatIdx, componentIdx, "", isShortNotation)
				}
			default:
				return "", fmt.Errorf("invalid field type '%s' in struct '%s', input not processed", field.Type().Name(), currentRecord.Type().Name())
			}
		default:
			return "", fmt.Errorf("invalid field type '%s' in struct '%s', input not processed", field.Type().Name(), currentRecord.Type().Name())
		}

	}

	return generateOutputRecord(recordType, fieldList, *repeatDelimiter, *componentDelimiter, *escapeDelimiter), nil
}

func addASTMFieldToList(data []OutputRecord, emptyRecords []OutputRecord, field, repeat, component int, value string, isShortNotation bool) ([]OutputRecord, []OutputRecord) {

	or := OutputRecord{
		Field:     field,
		Repeat:    repeat,
		Component: component,
		Value:     value,
	}

	if !isShortNotation {
		data = append(data, or)
		return data, emptyRecords
	}

	if value == "" {
		emptyRecords = append(emptyRecords, or)
		return data, emptyRecords
	}

	if len(emptyRecords) != 0 {
		data = append(data, emptyRecords...)
		emptyRecords = make(OutputRecords, 0)
	}
	data = append(data, or)

	return data, emptyRecords
}

// used for sorting
func (or OutputRecords) Len() int { return len(or) }
func (or OutputRecords) Less(i, j int) bool {
	if or[i].Field == or[j].Field {
		if or[i].Repeat == or[j].Repeat {
			return or[i].Component < or[j].Component
		} else {
			return or[i].Repeat < or[j].Repeat
		}
	} else {
		return or[i].Field < or[j].Field
	}
}
func (or OutputRecords) Swap(i, j int) { or[i], or[j] = or[j], or[i] }

/* Converting a list of values (all string already) to the astm format. this funciton works only for one record
   example:
    (0, 0, 2) = first-arr1
    (0, 0, 0) = third-arr1
    (0, 1, 0) = first-arr2
    (0, 1, 1) = second-arr2

	-> .... "|first-arr1^^third-arr1\fist-arr2^second-arr2|"

	returns the full record for output to astm file
*/

func generateOutputRecord(recordtype string, fieldList OutputRecords, REPEAT_DELIMITER, COMPONENT_DELIMITER, ESCAPE_DELMITER string) string {

	var output = ""

	// Record-ID, typical "H", "R", "O", .....
	output += recordtype

	// render fields - concat arrays
	sort.Sort(fieldList)

	var componentbuffer []string
	var lastComponentIdx = -1

	var currFieldGroup = -1
	var prevFieldGroup = -1
	var currFieldRepeat = -1
	var prevFieldRepeat = -1
	for _, field := range fieldList {

		prevFieldGroup = currFieldGroup
		currFieldGroup = field.Field
		var newFieldGroup = prevFieldGroup != currFieldGroup

		prevFieldRepeat = currFieldRepeat
		currFieldRepeat = field.Repeat
		var newRepeatGroup = prevFieldRepeat != currFieldRepeat

		if newFieldGroup || newRepeatGroup {

			// render all in component buffer
			if lastComponentIdx > -1 {
				output += componentbuffer[0]
				for i := 1; i <= lastComponentIdx; i++ {
					output += COMPONENT_DELIMITER + componentbuffer[i]
				}
			}

			if newFieldGroup {
				if prevFieldGroup > 0 { // for the first iteration we don't know on which number the field numbering starts
					for i := 0; i < currFieldGroup-prevFieldGroup; i++ {
						output += "|"
					}
				} else {
					output += "|"
				}
			} else if newRepeatGroup {
				output += REPEAT_DELIMITER
			}

			componentbuffer = make([]string, 100)
			lastComponentIdx = -1
		}

		componentbuffer[field.Component] = field.Value

		if field.Component > lastComponentIdx {
			lastComponentIdx = field.Component
		}
	}

	// render last field in component buffer
	if lastComponentIdx > -1 {
		output += componentbuffer[0]
		for i := 1; i <= lastComponentIdx; i++ {
			output += COMPONENT_DELIMITER + componentbuffer[i]
		}
	}

	return output
}
