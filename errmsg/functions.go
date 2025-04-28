package errmsg

import "errors"

// Encoding
const (
	Encoding_MsgInvalidEncoding = "invalid encoding: %s"
)

// Lining
var (
	Lining_ErrEmptyInput       = errors.New("empty input")
	Lining_ErrInvalidLinebreak = errors.New("invalid line breaking")
	Lining_ErrNoLineSeparator  = errors.New("separator has to be provided if auto-detect is disabled")
)

// AnnotationParsing
var (
	AnnotationParsing_ErrMissingAstmAnnotation        = errors.New("astm annotation missing")
	AnnotationParsing_ErrInvalidAstmAnnotation        = errors.New("invalid astm annotation")
	AnnotationParsing_ErrInvalidAstmAttribute         = errors.New("invalid astm attribute")
	AnnotationParsing_ErrInvalidAstmAttributeFormat   = errors.New("invalid astm attribute format")
	AnnotationParsing_ErrInvalidInputStruct           = errors.New("invalid input struct")
	AnnotationParsing_ErrIllegalComponentArray        = errors.New("component array is not allowed")
	AnnotationParsing_ErrIllegalComponentSubstructure = errors.New("component substructure is not allowed")
)

// LineParsing
var (
	LineParsing_ErrEmptyInput                  = errors.New("empty input")
	LineParsing_ErrHeaderTooShort              = errors.New("header too short")
	LineParsing_ErrMandatoryInputFieldsMissing = errors.New("mandatory input fields missing")
	LineParsing_ErrLineTypeNameMismatch        = errors.New("line type name mismatch")
	LineParsing_ErrSequenceNumberMismatch      = errors.New("sequence number mismatch")
	LineParsing_ErrRequiredInputFieldMissing   = errors.New("required input field missing")
	LineParsing_ErrInputComponentsMissing      = errors.New("input components missing")

	LineParsing_ErrNonSettableField          = errors.New("field is not settable")
	LineParsing_ErrDataParsingError          = errors.New("data parsing error")
	LineParsing_ErrInvalidDateFormat         = errors.New("invalid date format")
	LineParsing_ErrUsupportedDataType        = errors.New("unsupported data type")
	LineParsing_ErrReservedFieldPosReference = errors.New("field position 1 and 2 are reserved")
)

// StructureParsing
var (
	StructureParsing_ErrMaxDepthReached      = errors.New("max depth reached")
	StructureParsing_ErrInputLinesDepleted   = errors.New("input lines depleted")
	StructureParsing_ErrLineTypeNameMismatch = errors.New("line type name mismatch")
)

// LineBuilding
var (
	LineBuilding_ErrInvalidDateFormat           = errors.New("invalid date format")
	LineBuilding_ErrUsupportedDataType          = errors.New("unsupported data type")
	LineBuilding_ErrReservedFieldPosReference   = errors.New("field position 1 and 2 are reserved")
	LineBuilding_ErrInvalidLengthAttributeValue = errors.New("invalid length attribute value")
)
