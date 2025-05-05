package models

// Annotation types for ASTM fields and structures
type AstmFieldAnnotation struct {
	Raw            string
	FieldPos       int
	IsArray        bool
	IsComponent    bool
	ComponentPos   int
	IsSubstructure bool
	Attributes     map[string]string
}
type AstmStructAnnotation struct {
	Raw         string
	StructName  string
	IsComposite bool
	IsArray     bool
	Attributes  map[string]string
}
