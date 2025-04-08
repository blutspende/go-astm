package models

type AstmFieldAnnotation struct {
	Raw               string
	IsArray           bool
	FieldPos          int
	HasComponent      bool
	ComponentPos      int
	HasAttribute      bool
	Attribute         string
	HasAttributeValue bool
	AttributeValue    int
}

type AstmStructAnnotation struct {
	Raw          string
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
