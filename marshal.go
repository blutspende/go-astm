package astm

import (
	"bytes"
	"fmt"
	"github.com/blutspende/go-astm/v2/functions"
	"github.com/blutspende/go-astm/v2/models"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func Marshal(sourceStruct interface{}, configuration ...*models.Configuration) (result [][]byte, err error) {
	// Set up the configuration
	var config *models.Configuration
	if len(configuration) > 0 {
		config = configuration[0]
	} else {
		config = &models.DefaultConfiguration
		config.Internal.Delimiters = models.DefaultDelimiters
	}
	config.Internal.TimeLocation, err = time.LoadLocation(config.TimeZone)
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

/** Marshal - wrap datastructure to code
**/
func Marshal_old(message interface{}, enc Encoding, tz Timezone, notation Notation) ([][]byte, error) {

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

	// default delimiters. These will be overwritten by the first occurrence of "delimiter"-annotation
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

		if len(recordAstmTag) == 0 { // no annotation = Descend if it's an array or a struct of such

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
	var err error
	for i := 0; i < currentRecord.NumField(); i++ {
		field := currentRecord.Field(i)
		fieldAstmTag := currentRecord.Type().Field(i).Tag.Get("astm")
		if fieldAstmTag == "" {
			continue
		}
		fieldList, err = getOutputRecords(
			field,
			fieldAstmTag,
			fieldList,
			currentRecord,
			generatedSequenceNumber,
			location,
			repeatDelimiter,
			componentDelimiter,
			escapeDelimiter,
			notation)
		if err != nil {
			return "", err
		}

	}

	return generateOutputString(recordType, fieldList, *repeatDelimiter, *componentDelimiter, *escapeDelimiter), nil
}

func getOutputRecords(
	field reflect.Value,
	fieldAstmTag string,
	fieldList []OutputRecord,
	currentRecord reflect.Value,
	generatedSequenceNumber int,
	location *time.Location,
	repeatDelimiter, componentDelimiter, escapeDelimiter *string,
	notation Notation,
) ([]OutputRecord, error) {
	fieldAstmTagsList := strings.Split(fieldAstmTag, ",")
	//TODO: this should not depend on the unmarshal code, or make it explicit
	fieldIdx, repeatIdx, componentIdx, _, err := readFieldAddressAnnotation(fieldAstmTagsList[0])
	if err != nil {
		return nil, fmt.Errorf("invalid annotation for field %s : (%w)", reflect.TypeOf(field).Name(), err)
	}

	switch field.Type().Kind() {
	case reflect.String:
		value := ""

		if sliceContainsString(fieldAstmTagsList, ANNOTATION_SEQUENCE) {
			return nil, fmt.Errorf("invalid annotation %s for string-field", ANNOTATION_SEQUENCE)
		}

		// if no delimiters are given, default is \^&
		if sliceContainsString(fieldAstmTagsList, ANNOTATION_DELIMITER) && field.String() == "" {
			value = *repeatDelimiter + *componentDelimiter + *escapeDelimiter
		} else {
			value = field.String()
		}

		fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx, componentIdx, value)
	case reflect.Int:
		value := fmt.Sprintf("%d", field.Int())
		if sliceContainsString(fieldAstmTagsList, ANNOTATION_SEQUENCE) {
			value = fmt.Sprintf("%d", generatedSequenceNumber)
			generatedSequenceNumber = generatedSequenceNumber + 1
		}

		fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx, componentIdx, value)
	case reflect.Float32, reflect.Float64:
		decimalLength := getDecimalLengthByASTMTagList(fieldAstmTagsList)
		format := fmt.Sprintf("%%.%df", decimalLength)
		value := fmt.Sprintf(format, field.Float())
		fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx, componentIdx, value)
	case reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			switch field.Type().Elem().Kind() {
			case reflect.Struct:
				for j := 0; j < field.Index(i).NumField(); j++ {
					subfieldList := make([]OutputRecord, 0)
					subfield := field.Index(i).Field(j)
					subfieldAstmTag := field.Index(i).Type().Field(j).Tag.Get("astm")
					if subfieldAstmTag == "" {
						continue
					}
					subfieldList, err = getOutputRecords(subfield, subfieldAstmTag, subfieldList, field.Index(i), generatedSequenceNumber, location, repeatDelimiter, componentDelimiter, escapeDelimiter, notation)
					for k := range subfieldList {
						subfieldList[k].Field += fieldIdx
						subfieldList[k].Repeat += i
						subfieldList[k].Component += componentIdx
					}
					fieldList = append(fieldList, subfieldList...)
				}
			case reflect.String:
				fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx+i, componentIdx, field.Index(i).String())
			case reflect.Int:
				fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx+i, componentIdx, fmt.Sprintf("%d", field.Index(i).Int()))
			case reflect.Float32, reflect.Float64:
				decimalLength := getDecimalLengthByASTMTagList(fieldAstmTagsList)
				format := fmt.Sprintf("%%.%df", decimalLength)
				value := fmt.Sprintf(format, field.Index(i).Float())
				fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx+i, componentIdx, value)
			default:
				return nil, fmt.Errorf("while resolving slice: invalid type '%s'", field.Type().Elem().Kind().String())
			}

		}
	case reflect.Struct:
		switch field.Type().Name() {
		case "Time":
			time := field.Interface().(time.Time)

			if !time.IsZero() {

				if sliceContainsString(fieldAstmTagsList, ANNOTATION_LONGDATE) {
					value := time.In(location).Format("20060102150405")
					fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx, componentIdx, value)
				} else { // short date
					value := time.In(location).Format("20060102")
					fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx, componentIdx, value)
				}
			} else {
				fieldList = addASTMFieldToList(fieldList, fieldIdx, repeatIdx, componentIdx, "")
			}
		default:
			return nil, fmt.Errorf("invalid field type '%s' in struct '%s', input not processed", field.Type().Name(), currentRecord.Type().Name())
		}
	default:
		return nil, fmt.Errorf("invalid field type '%s' in struct '%s', input not processed", field.Type().Name(), currentRecord.Type().Name())
	}

	if notation == ShortNotation {
		for i := len(fieldList) - 1; i >= 0; i-- {
			if fieldList[i].Value != "" {
				break
			}
			fieldList = fieldList[:i]
		}
	}

	return fieldList, nil
}

func addASTMFieldToList(data []OutputRecord, field, repeat, component int, value string) []OutputRecord {

	or := OutputRecord{
		Field:     field,
		Repeat:    repeat,
		Component: component,
		Value:     value,
	}

	data = append(data, or)
	return data
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

func getDecimalLengthByASTMTagList(astmTagList []string) int {
	decimalLength := 3
	for _, tag := range astmTagList {
		if strings.Contains(tag, ANNOTATION_LENGTH) {
			decimalLengthTags := strings.Split(tag, ":")
			if len(decimalLengthTags) < 2 {
				break
			}
			customLength, err := strconv.Atoi(decimalLengthTags[1])
			if err != nil || customLength == 0 {
				break
			}
			decimalLength = customLength
			break
		}
	}

	return decimalLength
}

func generateOutputString(recordtype string, fieldList OutputRecords, REPEAT_DELIMITER, COMPONENT_DELIMITER, ESCAPE_DELMITER string) string {

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
