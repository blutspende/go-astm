package astm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/blutspende/go-astm/v2/functions"
	"github.com/blutspende/go-astm/v2/models"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

const MAX_MESSAGE_COUNT = 44
const MAX_DEPTH = 44

func Unmarshal_Old(messageData []byte, targetStruct interface{}, enc Encoding, tz Timezone) error {

	if len(messageData) == 0 {
		return fmt.Errorf("message has nil value - aborting")
	}
	var (
		messageBytes []byte
		err          error
	)

	switch enc {
	case EncodingUTF8:
		messageBytes = messageData
	case EncodingASCII:
		messageBytes = messageData
	case EncodingDOS866:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.CodePage866, messageData); err != nil {
			return err
		}
	case EncodingDOS855:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.CodePage855, messageData); err != nil {
			return err
		}
	case EncodingDOS852:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.CodePage852, messageData); err != nil {
			return err
		}
	case EncodingWindows1250:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.Windows1250, messageData); err != nil {
			return err
		}
	case EncodingWindows1251:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.Windows1251, messageData); err != nil {
			return err
		}
	case EncodingWindows1252:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.Windows1252, messageData); err != nil {
			return err
		}
	case EncodingISO8859_1:
		if messageBytes, err = EncodeCharsetToUTF8From(charmap.ISO8859_1, messageData); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid Codepage Id='%d' - %w", enc, err)
	}

	// first try to break by 0x0a (non-standard, but used sometimes)
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
			bufferedInputLines = append(bufferedInputLines, bufferedInputLinesWithEmptyLines[i])
		}
	}

	var (
		repeatDelimiter    = "\\"
		componentDelimiter = "^"
		escapeDelimiter    = "&"
	)

	currentInputLine := 0
	currentInputLine, _, err = reflectInputToStruct(
		bufferedInputLines,
		1, /*recursion-depth*/
		currentInputLine,
		targetStruct,
		enc,
		tz,
		&repeatDelimiter,
		&componentDelimiter,
		&escapeDelimiter)

	if err != nil {
		return err
	}

	// if we have reached the end of the first message but not the end of our buffered input
	if currentInputLine < len(bufferedInputLines) {
		// return an error to avoid data loss
		return fmt.Errorf("%d lines of input were skipped. Last line was %d: '%s' ", len(bufferedInputLines)-currentInputLine+1, currentInputLine, bufferedInputLines[currentInputLine])
	}

	return nil
}

type RETV int

const (
	OK         RETV = 1
	UNEXPECTED RETV = 2 // an exit that won't abort processing. used for skipping optional records
	ERROR      RETV = 3 // a definite error that stops the process
)

func EncodeCharsetToUTF8From(charmap *charmap.Charmap, data []byte) ([]byte, error) {
	sr := bytes.NewReader(data)
	e := charmap.NewDecoder().Reader(sr)
	bytes, err := ioutil.ReadAll(e)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

/* This function takes a string and a struct and matches the annotated fields to the string-input */
func reflectInputToStruct(bufferedInputLines []string, depth int, currentInputLine int, targetStruct interface{}, enc Encoding, tz Timezone,
	repeatDelimiter, componentDelimiter, escapeDelimiter *string) (int, RETV, error) {

	if depth > MAX_DEPTH {
		return currentInputLine, ERROR, errors.New(fmt.Sprintf("Maximum recursion depth reached (%d). Too many nested structures ? - aborting", depth))
	}

	if bufferedInputLines[currentInputLine] == "" {
		// Caution : +1 might skip one; .. without could stick in loop
		return currentInputLine + 1, UNEXPECTED, errors.New(fmt.Sprintf("Empty Input"))
	}

	var targetStructType reflect.Type
	var targetStructValue reflect.Value
	if reflect.TypeOf(targetStruct).Kind() == reflect.Struct {
		targetStructType = reflect.TypeOf(targetStruct)
		targetStructValue = reflect.ValueOf(targetStruct)
	} else {
		targetStructType = reflect.TypeOf(targetStruct).Elem()
		targetStructValue = reflect.ValueOf(targetStruct).Elem()
	}
	timeLocation, err := time.LoadLocation(string(tz))
	if err != nil {
		return currentInputLine, ERROR, err
	}

	for i := 0; i < targetStructType.NumField(); i++ {
		currentRecord := targetStructValue.Field(i)
		ftype := targetStructType.Field(i)
		astmTag := ftype.Tag.Get("astm")
		astmTagsList := strings.Split(astmTag, ",")

		if len(astmTagsList) < 1 {
			continue // not annotated = no processing
		}

		// no annotation after astm:.. provided means a nested array with more records or ignore
		if len(astmTagsList[0]) < 1 {

			// Not annotated array. If it's a struct have to recurse, otherwise skip
			if targetStructType.Field(i).Type.Kind() == reflect.Slice {

				// Array of Structs
				if reflect.TypeOf(targetStructValue.Interface()).Kind() == reflect.Struct {

					innerStructureType := targetStructType.Field(i).Type.Elem()

					sliceForNestedStructure := reflect.MakeSlice(targetStructType.Field(i).Type, 0, 0)

					for currentInputLine < len(bufferedInputLines) { // iterate for as long as there is input or an unexpecte dinput type
						allocatedElement := reflect.New(innerStructureType)
						var err error
						var retv RETV
						currentInputLine, retv, err = reflectInputToStruct(bufferedInputLines, depth+1,
							currentInputLine, allocatedElement.Interface(), enc, tz, repeatDelimiter, componentDelimiter, escapeDelimiter)

						if err != nil {
							if retv == UNEXPECTED {
								break
							}
							if retv == ERROR { // a serious error ends the processing
								return currentInputLine, ERROR, err
							}
						}

						sliceForNestedStructure = reflect.Append(sliceForNestedStructure, allocatedElement.Elem())
						reflect.ValueOf(targetStruct).Elem().Field(i).Set(sliceForNestedStructure)
					}
					continue
				}
			} else if targetStructType.Field(i).Type.Kind() == reflect.Struct { // struct without annotation - descending

				var err error
				var retv RETV

				dood := currentRecord.Addr().Interface()

				currentInputLine, retv, err = reflectInputToStruct(bufferedInputLines, depth+1, currentInputLine, dood, enc, tz,
					repeatDelimiter, componentDelimiter, escapeDelimiter)
				if err != nil {
					if retv == UNEXPECTED {
						if depth > 0 {
							// if nested structures abort due to unexpected records that does not create an error
							// as the parse will be continued one level higher
							return currentInputLine, UNEXPECTED, err
						} else {
							return currentInputLine, ERROR, err
						}
					}
					if retv == ERROR { // a serious error ends the processing
						return currentInputLine, ERROR, err
					}
				}

				continue

			} else {
				return currentInputLine, ERROR, errors.New(fmt.Sprintf("Invalid Datatype '%s' - abort unmarshal.", targetStructType.Field(i).Type.Kind()))
			}
		}

		expectInputRecordType := astmTagsList[0][0] // Expected Record type
		expectedInputRecordTypeOptional := false
		if sliceContainsString(astmTagsList, ANNOTATION_OPTIONAL) {
			expectedInputRecordTypeOptional = true
		}

		if currentInputLine >= len(bufferedInputLines) { // premature end ...
			return currentInputLine, ERROR, fmt.Errorf("premature end of input in line %d (Missing Data)", currentInputLine)
		}

		if len(bufferedInputLines[currentInputLine]) == 0 {
			continue // empty lines can only be skipped
		}

		// headers require delimiters to be disregarded
		isHeader := false
		if bufferedInputLines[currentInputLine][0] == 'H' {
			isHeader = true
		}

		if expectInputRecordType == bufferedInputLines[currentInputLine][0] {

			//Special case: it's not an annotated record, it is an array of annotated records here :
			if currentRecord.Kind() == reflect.Slice {
				innerStructureType := targetStructType.Field(i).Type.Elem()
				sliceForNestedStructure := reflect.MakeSlice(targetStructType.Field(i).Type, 0, 0)
				for { // iterate for as long as the same type repeats
					allocatedElement := reflect.New(innerStructureType)

					if err = reflectAnnotatedFields(bufferedInputLines[currentInputLine], allocatedElement.Elem(), timeLocation, isHeader, repeatDelimiter, componentDelimiter, escapeDelimiter); err != nil {
						return currentInputLine, ERROR, errors.New(fmt.Sprintf("Failed to process input line '%s' err:%s", bufferedInputLines[currentInputLine], err))
					}

					sliceForNestedStructure = reflect.Append(sliceForNestedStructure, allocatedElement.Elem())
					reflect.ValueOf(targetStruct).Elem().Field(i).Set(sliceForNestedStructure)

					// keep reading while same elements are up
					currentInputLine = currentInputLine + 1
					if currentInputLine >= len(bufferedInputLines) {
						break
					}
					if expectInputRecordType != bufferedInputLines[currentInputLine][0] {
						break
					}
				}

			} else { // The "normal" case: scanning a string into a structure :
				if err = reflectAnnotatedFields(bufferedInputLines[currentInputLine], currentRecord, timeLocation, isHeader, repeatDelimiter, componentDelimiter, escapeDelimiter); err != nil {
					return currentInputLine, ERROR, errors.New(fmt.Sprintf("Failed to process input line '%s' err:%s", bufferedInputLines[currentInputLine], err))
				}
				currentInputLine = currentInputLine + 1
			}

		} else { // The expected input-record did not occur
			if expectedInputRecordTypeOptional {
				continue // skipping optional record instead of an error
			} else {
				return currentInputLine, UNEXPECTED, errors.New(fmt.Sprintf("Expected Record-Type '%c' input was '%c' in depth (%d) (Abort)", expectInputRecordType, bufferedInputLines[currentInputLine][0], depth))
			}
		}

		if currentInputLine >= len(bufferedInputLines) {
			break
		}
	}

	return currentInputLine, OK, nil
}

func reflectAnnotatedFields(inputStr string, record reflect.Value, timezone *time.Location, isHeader bool,
	repeatDelimiter, componentDelimiter, escapeDelimiter *string) error {

	if reflect.ValueOf(record).Type().Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("invalid type of target: '%s', expecting 'struct'", reflect.ValueOf(record).Type().Kind()))
	}

	inputFields := strings.Split(inputStr, "|")
	if len(inputFields) < 1 {
		return errors.New("Input contains no data")
	}

	for j := 0; j < record.NumField(); j++ {
		recordfield := record.Field(j)
		if !recordfield.CanInterface() {
			return errors.New(fmt.Sprintf("Field %s is not exported - aborting import", recordfield.Type().Name()))
		}
		recordFieldInterface := recordfield.Addr().Interface()

		hasOverrideDelimiterAnnotation := false
		inputIsRequired := false
		astmTag := record.Type().Field(j).Tag.Get("astm")
		if astmTag == "" {
			continue // nothing to process when someone requires astm:
		}
		astmTagsList := strings.Split(astmTag, ",")
		for i := 0; i < len(astmTagsList); i++ {
			astmTagsList[i] = strings.Trim(astmTagsList[i], " ")
		}
		if sliceContainsString(astmTagsList, ANNOTATION_DELIMITER) {
			// the delimiter is instantly replaced with the delimiters from the file for further parsing. By default that is "\^&"
			hasOverrideDelimiterAnnotation = true
		}
		if sliceContainsString(astmTagsList, ANNOTATION_REQUIRED) {
			inputIsRequired = true
		}
		currentInputFieldNo, repeat, component, hasComponents, err := readFieldAddressAnnotation(astmTagsList[0])
		if err != nil {
			return errors.New(fmt.Sprintf("Invalid annotation for field %s. (%s)", record.Type().Field(j).Name, err))
		}
		if currentInputFieldNo >= len(inputFields) || currentInputFieldNo < 0 {
			//TODO: user should be able to toggle whether he wants an exact match = error or bestfit = skip silent
			continue // mapped field is beyond the data
		}

		switch reflect.TypeOf(recordfield.Interface()).Kind() {
		case reflect.String:
			if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo],
				repeat, component, *repeatDelimiter, *componentDelimiter, sliceContainsString(astmTagsList, ANNOTATION_REQUIRED), hasComponents); err == nil {

				// in headers there can be special characters, that is why the value needs to disregard the delimiters:
				if isHeader {
					value = inputFields[currentInputFieldNo]
				}

				reflect.ValueOf(recordFieldInterface).Elem().SetString(reflect.ValueOf(value).String())

				if hasOverrideDelimiterAnnotation { // the first three characters become the new delimiters
					if len(value) >= 1 {
						*repeatDelimiter = value[0:1]
					}
					if len(value) >= 2 {
						*componentDelimiter = value[1:2]
					}
					if len(value) >= 3 {
						*escapeDelimiter = value[2:3]
					}
				}
			} else {
				if inputIsRequired { // by default we ignore missing input
					return errors.New(fmt.Sprintf("Failed to extract index (%d.%d.%d) from input line '%s' : (%s)",
						currentInputFieldNo+1, repeat+1, component+1, inputStr, err))
				}
			}
		case reflect.Int:
			if hasOverrideDelimiterAnnotation {
				return errors.New("delimiter-annotation is only allowed for string-type, not int.")
			}

			if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo], repeat, component,
				*repeatDelimiter, *componentDelimiter, sliceContainsString(astmTagsList, ANNOTATION_REQUIRED), hasComponents); err == nil {

				if num, err := strconv.Atoi(value); err == nil {
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(num))
				} else {
					if inputIsRequired { // by default we ignore missing input
						return errors.New(fmt.Sprintf("Failed to extract index (%d,%d) from field %s(%s)", repeat, component, inputFields[currentInputFieldNo], err))
					}
				}

			} else {
				return err
			}
		case reflect.Float32:
			if hasOverrideDelimiterAnnotation {
				return errors.New("delimiter-annotation is only allowed for string-type, not int.")
			}

			if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo],
				repeat, component, *repeatDelimiter, *componentDelimiter,
				sliceContainsString(astmTagsList, ANNOTATION_REQUIRED), hasComponents); err == nil {

				if num, err := strconv.ParseFloat(value, 32); err == nil {
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(float32(num)))
				} else {
					if inputIsRequired { // by default we ignore missing input
						return errors.New(fmt.Sprintf("Failed to extract index (%d,%d) from field %s(%s)", repeat, component, inputFields[currentInputFieldNo], err))
					}
				}

			} else {
				return err
			}
		case reflect.Float64:
			if hasOverrideDelimiterAnnotation {
				return errors.New("delimiter-annotation is only allowed for string-type, not int.")
			}

			if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo],
				repeat, component, *repeatDelimiter, *componentDelimiter,
				sliceContainsString(astmTagsList, ANNOTATION_REQUIRED), hasComponents); err == nil {

				if num, err := strconv.ParseFloat(value, 64); err == nil {
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(float64(num)))
				} else {
					if inputIsRequired { // by default we ignore missing input
						return errors.New(fmt.Sprintf("Failed to extract index (%d,%d) from field %s(%s)", repeat, component, inputFields[currentInputFieldNo], err))
					}
				}

			} else {
				return err
			}
		case reflect.Slice: // Decoding a slice in a field
			switch recordfield.Type().Elem().Kind() {
			case reflect.String: // []string
				decodedArray := strings.Split(inputFields[currentInputFieldNo], *repeatDelimiter)
				reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(decodedArray))
			case reflect.Int: // []int
				decodedArray := strings.Split(inputFields[currentInputFieldNo], *repeatDelimiter)
				convertedArray := make([]int, len(decodedArray))
				for i := 0; i < len(decodedArray); i++ {
					convertedArray[i], err = strconv.Atoi(decodedArray[i])
					if err != nil {
						return fmt.Errorf("invalid input type for decoding to []int %s in '%s'", decodedArray[i], inputStr)
					}
				}
				reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(convertedArray))
			case reflect.Float32: // []float32
				decodedArray := strings.Split(inputFields[currentInputFieldNo], *repeatDelimiter)
				convertedArray := make([]float32, len(decodedArray))
				for i := 0; i < len(decodedArray); i++ {
					f64, err := strconv.ParseFloat(decodedArray[i], 32)
					if err != nil {
						return fmt.Errorf("invalid input type for decoding to []float32 %s (%v) in '%s'", decodedArray[i], err, inputStr)
					}
					convertedArray[i] = float32(f64)
				}
				reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(convertedArray))
			case reflect.Float64: // []float64
				decodedArray := strings.Split(inputFields[currentInputFieldNo], *repeatDelimiter)
				convertedArray := make([]float64, len(decodedArray))
				for i := 0; i < len(decodedArray); i++ {
					convertedArray[i], err = strconv.ParseFloat(decodedArray[i], 64)
					if err != nil {
						return fmt.Errorf("invalid input type for decoding to []float64 %s (%v) in '%s'", decodedArray[i], err, inputStr)
					}
				}
				reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(convertedArray))
			case reflect.Struct:
				decodedArray := strings.Split(inputFields[currentInputFieldNo], *repeatDelimiter)
				sliceOfStruct := reflect.MakeSlice(recordfield.Type(), len(decodedArray), len(decodedArray))
				for i := 0; i < len(decodedArray); i++ {
					err := reflectAnnotatedFields(decodedArray[i],
						sliceOfStruct.Index(i), timezone, isHeader,
						repeatDelimiter, componentDelimiter, escapeDelimiter)
					if err != nil {
						return fmt.Errorf("unable to decode array of []structs %s (%v) in '%s'", decodedArray[i], err, inputStr)
					}
				}
				reflect.ValueOf(recordFieldInterface).Elem().Set(sliceOfStruct)
			default:
				return fmt.Errorf("when resolving a slice within a record : Invalid type of Field '%s' while trying to unmarshal this string '%s'. This datatype is not implemented.",
					reflect.TypeOf(recordfield.Interface()).Elem().Kind(), inputStr)
			}
		case reflect.Struct:
			switch reflect.TypeOf(recordfield.Interface()).Name() {
			case "Time":
				if hasOverrideDelimiterAnnotation {
					return errors.New("delimiter-annotation is only allowed for string-type, not Time")
				}

				var inputFieldValue string
				if value, err := extractAstmFieldByRepeatAndComponent(inputFields[currentInputFieldNo],
					repeat, component, *repeatDelimiter, *componentDelimiter,
					sliceContainsString(astmTagsList, ANNOTATION_REQUIRED), hasComponents); err == nil {
					inputFieldValue = value
				} else {
					return errors.New(fmt.Sprintf("Error extracting field '%s' tagged: '%s' : %s ", recordfield.Type().Name(), astmTag, err))
				}

				if inputFieldValue == "" {
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(time.Time{}))
				} else if len(inputFieldValue) == 8 { // YYYYMMDD See Section 5.6.2 https://samson-rus.com/wp-content/files/LIS2-A2.pdf
					timeInLocation, err := time.ParseInLocation("20060102", inputFieldValue, timezone)
					if err != nil {
						return errors.New(fmt.Sprintf("Invalid time format <%s>", inputFieldValue))
					}
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(timeInLocation))

				} else if len(inputFieldValue) == 14 { // YYYYMMDDHHMMSS
					timeInLocation, err := time.ParseInLocation("20060102150405", inputFieldValue, timezone)
					if err != nil {
						return errors.New(fmt.Sprintf("Invalid time format <%s>", inputFieldValue))
					}
					reflect.ValueOf(recordFieldInterface).Elem().Set(reflect.ValueOf(timeInLocation.UTC()))
				} else {
					return errors.New(fmt.Sprintf("Unrecognized time format <%s>", inputFieldValue))
				}
			default:
				return errors.New(fmt.Sprintf("Invalid type of Field '%s' while trying to unmarshal this string '%s'. This datatype is a structure type which is not implemented.",
					reflect.TypeOf(recordfield.Interface()).Name(), inputStr))
			}
		default:
			return errors.New(fmt.Sprintf("Invalid type of Field '%s' while trying to unmarshal this string '%s'. This datatype is not implemented.",
				reflect.TypeOf(recordfield.Interface()).Kind(), inputStr))
		}
	}

	return nil
}

// Translating the annotation of a field to field, index/repeat, component
// Input of one value : e.g."4" -> field -> 4
// Input of two values :"4.2" -> field, component -> 4,1,2
// Input of three values "4.1.1" -> field, repeat, component -> 4,1,1
// "whereas field indexes should be 1-99 (check plz)
func readFieldAddressAnnotation(annotation string) (field int, repeat int, component int, hasComponents bool, err error) {

	if annotation == "" { // no annotation will always return the first of everything
		return 0, 0, 0, false, nil
	}
	field = 1
	repeat = 1
	component = 1
	hasComponents = false
	fieldSplitted := strings.Split(annotation, ".")

	if len(fieldSplitted) >= 1 {
		if field, err = strconv.Atoi(fieldSplitted[0]); err != nil {
			return 0, 0, 0, false, err
		}
	}
	if len(fieldSplitted) >= 2 {
		if component, err = strconv.Atoi(fieldSplitted[1]); err != nil {
			return 0, 0, 0, false, err
		}
		// Workaround: for multi component inputs with single field struct definitions to default to the whole field instead of only the first component
		hasComponents = true
	}
	if len(fieldSplitted) >= 3 {
		if repeat, err = strconv.Atoi(fieldSplitted[1]); err != nil {
			return 0, 0, 0, false, err
		}
		if component, err = strconv.Atoi(fieldSplitted[2]); err != nil {
			return 0, 0, 0, false, err
		}
	}

	return field - 1, repeat - 1, component - 1, hasComponents, nil
}

// input is an unpacked field from an astm-file free of the field delimiter ("|")
// this function extracts the field by repeat and component-delimiter
func extractAstmFieldByRepeatAndComponent(text string, repeat int, component int, repeatDelimiter, componentDelimiter string, isRequired bool, hasComponents bool) (string, error) {

	subfield := strings.Split(text, repeatDelimiter)
	if repeat >= len(subfield) {
		if isRequired {
			return "", errors.New(fmt.Sprintf("Index (%d, %d) out of bounds '%s', delimiter '%s'", repeat, component, text, repeatDelimiter))
		}
		return "", nil
	}

	// Workaround: single component value takes the whole subfield instead of just the first component
	if !hasComponents {
		return subfield[repeat], nil
	}

	subsubfield := strings.Split(subfield[repeat], componentDelimiter)
	if component >= len(subsubfield) || component < 0 {
		if isRequired {
			return "", errors.New(fmt.Sprintf("Index (%d, %d) out of bounds '%s' delimiter '%s'", repeat, component, text, componentDelimiter))
		}
		return "", nil
	}

	return subsubfield[component], nil
}

func sliceContainsString(list []string, search string) bool {
	for _, x := range list {
		if x == search {
			return true
		}
	}
	return false
}
