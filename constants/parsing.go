package constants

const ATTRIBUTE_REQUIRED string = "required" // field-annotation: by default all fields are optinal
const ATTRIBUTE_OPTIONAL string = "optional" // record-annotation: by default all records are mandatory
const ATTRIBUTE_SEQUENCE string = "sequence" // indicating that a sequence number should be generated (output only)
const ATTRIBUTE_LONGDATE string = "longdate" // Indicating that the date should be formatted as longdate (output only)
const ATTRIBUTE_LENGTH string = "length"     // used for specifying the decimal length of float fields - astm:"1,length:2" (output only)

const MAX_DEPTH int = 44
