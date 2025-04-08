package constants

//type Annotation string

const ANNOTATION_DELIMITER string = "delimiter" // annotation that triggers the delimiters in the scanner to be reset
const ANNOTATION_REQUIRED string = "require"    // field-annotation: by default all fields are optinal
const ANNOTATION_OPTIONAL string = "optional"   // record-annotation: by default all records are mandatory
const ANNOTATION_SEQUENCE string = "sequence"   // indicating that a sequence number should be generated (output only)
const ANNOTATION_LONGDATE string = "longdate"
const ANNOTATION_LENGTH string = "length" // used for specifying the decimal length of float fields - astm:"1,length:2" (output only)
