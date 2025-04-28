package models

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
	Raw               string
	StructName        string
	IsComposite       bool
	IsArray           bool
	HasAttribute      bool
	Attribute         string
	HasAttributeValue bool
	AttributeValue    string
}
