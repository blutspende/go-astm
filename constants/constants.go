package constants

// Maximum recurtion depth for struct parsing and building
const MaxDepth int = 44

// Attributes for annotations
const AttributeRequired string = "required" // field-annotation: by default all fields are optinal
const AttributeOptional string = "optional" // record-annotation: by default all records are mandatory
const AttributeLongdate string = "longdate" // Indicating that the date should be formatted as date and time (output only)
const AttributeLength string = "length"     // used for specifying the decimal length of float fields - astm:"1,length:2" (output only)
const AttributeSubname string = "subname"   // used for specifying a subname for a record - astm:"M,subname:MATRIX"
