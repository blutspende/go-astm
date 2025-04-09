package errmsg

import "errors"

// Encoding
const (
	Encoding_MsgInvalidEncoding = "invalid encoding: %s"
)

// Lining
var (
	Lining_ErrNotEnoughLines   = errors.New("not enough lines")
	Lining_ErrInvalidLinebreak = errors.New("invalid line breaking")
)

// AnnotationParsing
var (
	AnnotationParsing_ErrMissingAstmAnnotation = errors.New("astm annotation missing")
	AnnotationParsing_ErrInvalidAstmAnnotation = errors.New("invalid astm annotation")
	AnnotationParsing_ErrTooManyAttributes     = errors.New("only one astm attribute is allowed")
	AnnotationParsing_ErrInvalidAstmAttribute  = errors.New("invalid astm attribute")
	AnnotationParsing_ErrInvalidTargetStruct   = errors.New("invalid target struct")
)

// LineParsing
var (
	LineParsing_ErrLineTypeNameMismatch   = errors.New("line type name mismatch")
	LineParsing_ErrSequenceNumberMismatch = errors.New("sequence number mismatch")
)

// StructureParsing
var (
	StructureParsing_ErrMaxDepthReached    = errors.New("max depth reached")
	StructureParsing_ErrInputLinesDepleted = errors.New("input lines depleted")
)
