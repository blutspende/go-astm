package functions

import (
	"github.com/blutspende/go-astm/v2/constants"
	"github.com/blutspende/go-astm/v2/errmsg"
	"github.com/blutspende/go-astm/v2/models"
	"reflect"
)

func ParseStruct(inputLines []string, targetStruct interface{}, lineIndex *int, depth int, config *models.Configuration) (err error) {
	// Check for maximum depth
	if depth >= constants.MAX_DEPTH {
		return errmsg.StructureParsing_ErrMaxDepthReached
	}
	// Check for enough input lines
	if *lineIndex >= len(inputLines) {
		return errmsg.StructureParsing_ErrInputLinesDepleted
	}

	// Process the target structure
	targetTypes, targetValues, _, _ := ProcessStructReflection(targetStruct)

	// Iterate over the inputFields of the targetStruct struct
	for i, targetType := range targetTypes {

		// Parse the targetStruct field targetFieldAnnotation
		targetStructAnnotation, _ := ParseAstmStructAnnotation(targetType)

		// Save the target value pointer
		targetValue := targetValues[i].Addr().Interface()

		//TODO: handle optional elements

		// Target is an array it is iterated with conditional break (unknown length)
		if targetStructAnnotation.IsArray {
			// Create the array structure
			sliceType := reflect.SliceOf(targetValues[i].Type().Elem())
			targetValues[i].Set(reflect.MakeSlice(sliceType, 0, 0))

			// Iterate as long as we have matching input structure and still have input lines
			for j := 1; *lineIndex < len(inputLines); j++ {
				// Create a new element for the slice to parse into
				elem := reflect.New(targetValues[i].Type().Elem()).Elem()

				if targetStructAnnotation.IsComposite {
					// Composite target: recursively parse the composite structure
					err = ParseStruct(inputLines, elem.Addr().Interface(), lineIndex, depth+1, config)
				} else {
					// Non-composite target: parse the line into the new element
					err = ParseLine(inputLines[*lineIndex], elem.Addr().Interface(), targetStructAnnotation.StructName, j, config)
					// Increment the line index
					*lineIndex++
				}

				if err != nil {
					// TODO: handle structure change better than a specific error
					if err == errmsg.LineParsing_ErrLineTypeNameMismatch {
						// If the error is a line type name mismatch, it means the end of the array
						break
					} else {
						// Other error
						return err
					}
				}
				// If no error, add the new element to the slice
				targetValues[i].Set(reflect.Append(targetValues[i], elem))
			}
		} else {
			// Single element structure
			if targetStructAnnotation.IsComposite {
				// Composite target: go further down the rabbit hole
				err = ParseStruct(inputLines, targetValue, lineIndex, depth+1, config)
				if err != nil {
					return err
				}
			} else {
				// Non-composite target: there is a single line to parse
				if *lineIndex >= len(inputLines) {
					return errmsg.StructureParsing_ErrInputLinesDepleted
				}
				err = ParseLine(inputLines[*lineIndex], targetValue, targetStructAnnotation.StructName, 1, config)
				*lineIndex++
				if err != nil {
					return err
				}
			}
		}
	}
	// Return nil if everything went well
	return nil
}
