package constants

const ATTRIBUTE_DELIMITER string = "delimiter" // annotation that triggers the delimiters in the scanner to be reset
const ATTRIBUTE_REQUIRE string = "require"     // field-annotation: by default all fields are optinal
const ATTRIBUTE_OPTIONAL string = "optional"   // record-annotation: by default all records are mandatory
const ATTRIBUTE_SEQUENCE string = "sequence"   // indicating that a sequence number should be generated (output only)
const ATTRIBUTE_LONGDATE string = "longdate"
const ATTRIBUTE_LENGTH string = "length" // used for specifying the decimal length of float fields - astm:"1,length:2" (output only)

const MAX_DEPTH int = 44
