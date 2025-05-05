package errmsg

import "errors"

// Encoding
var (
	ErrEncodingInvalidEncoding = errors.New("invalid encoding")
)

// Lining
var (
	ErrLineProcessingEmptyInput       = errors.New("empty input")
	ErrLineProcessingInvalidLinebreak = errors.New("invalid line breaking")
	ErrLineProcessingNoLineSeparator  = errors.New("separator has to be provided if auto-detect is disabled")
)

// AnnotationParsing
var (
	ErrAnnotationParsingMissingAstmAnnotation        = errors.New("astm annotation missing")
	ErrAnnotationParsingInvalidAstmAnnotation        = errors.New("invalid astm annotation")
	ErrAnnotationParsingInvalidAstmAttribute         = errors.New("invalid astm attribute")
	ErrAnnotationParsingInvalidAstmAttributeFormat   = errors.New("invalid astm attribute format")
	ErrAnnotationParsingInvalidInputStruct           = errors.New("invalid input struct")
	ErrAnnotationParsingIllegalComponentArray        = errors.New("component array is not allowed")
	ErrAnnotationParsingIllegalComponentSubstructure = errors.New("component substructure is not allowed")
)

// LineParsing
var (
	ErrLineParsingEmptyInput                  = errors.New("empty input")
	ErrLineParsingHeaderTooShort              = errors.New("header too short")
	ErrLineParsingMandatoryInputFieldsMissing = errors.New("mandatory input fields missing")
	ErrLineParsingLineTypeNameMismatch        = errors.New("line type name mismatch")
	ErrLineParsingSequenceNumberMismatch      = errors.New("sequence number mismatch")
	ErrLineParsingRequiredInputFieldMissing   = errors.New("required input field missing")
	ErrLineParsingInputComponentsMissing      = errors.New("input components missing")
	ErrLineParsingNonSettableField            = errors.New("field is not settable")
	ErrLineParsingDataParsingError            = errors.New("data parsing error")
	ErrLineParsingInvalidDateFormat           = errors.New("invalid date format")
	ErrLineParsingUnsupportedDataType         = errors.New("unsupported data type")
	ErrLineParsingReservedFieldPosReference   = errors.New("field position 1 and 2 are reserved")
)

// StructureParsing
var (
	ErrStructureParsingMaxDepthReached      = errors.New("max depth reached")
	ErrStructureParsingInputLinesDepleted   = errors.New("input lines depleted")
	ErrStructureParsingLineTypeNameMismatch = errors.New("line type name mismatch")
)

// LineBuilding
var (
	ErrLineBuildingInvalidDateFormat           = errors.New("invalid date format")
	ErrLineBuildingUsupportedDataType          = errors.New("unsupported data type")
	ErrLineBuildingReservedFieldPosReference   = errors.New("field position 1 and 2 are reserved")
	ErrLineBuildingInvalidLengthAttributeValue = errors.New("invalid length attribute value")
)
