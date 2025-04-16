package models

import "time"

// Annotation types for ASTM fields and structures
type AstmFieldAnnotation struct {
	Raw               string
	FieldPos          int
	IsArray           bool
	IsComponent       bool
	ComponentPos      int
	IsSubstructure    bool
	HasAttribute      bool
	Attribute         string
	HasAttributeValue bool
	AttributeValue    int
}
type AstmStructAnnotation struct {
	Raw          string
	StructName   string
	IsComposite  bool
	IsArray      bool
	HasAttribute bool
	Attribute    string
}

// Delimiters used in ASTM parsing
type Delimiters struct {
	Field     string
	Repeat    string
	Component string
	Escape    string
}

var DefaultDelimiters = Delimiters{
	Field:     `|`,
	Repeat:    `\`,
	Component: `^`,
	Escape:    `&`,
}

// Internal configuration for parsing (part of the Configuration struct)
type InternalConfiguration struct {
	TimeLocation *time.Location
	Delimiters   Delimiters
}
