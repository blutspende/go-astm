package models

import "time"

// Annotation types for ASTM fields and structures
type AstmFieldAnnotation struct {
	Raw               string
	IsArray           bool
	FieldPos          int
	IsComponent       bool
	ComponentPos      int
	HasAttribute      bool
	Attribute         string
	HasAttributeValue bool
	AttributeValue    int
}
type AstmStructAnnotation struct {
	Raw          string
	IsComposite  bool
	IsArray      bool
	StructName   string
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
