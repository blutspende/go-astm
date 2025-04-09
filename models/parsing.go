package models

import "time"

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

type Configuration struct {
	Delimiters   *Delimiters
	TimeLocation *time.Location
}
