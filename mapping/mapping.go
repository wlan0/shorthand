package mapping

import (
	"go/ast"

	"github.com/koki/shorthand/inspect"
)

// SourceField is the name of a source field and its type information.
type SourceField struct {
	// Name of the field in the source object.
	Name string
	// The type information of this field.
	// The important information here is whether the type needs a
	//   nil check before its subfields are accessed.
	TypeExpr ast.Expr
	// A Context focused on the struct this field is part of.
	// This is used to interpret TypeExpr properly.
	Context *inspect.Context
}

// MappedField contains either a new struct type or a value extracted from the
//   source object.
type MappedField struct {
	// Name is the name of the new field.
	Name string
	// NewStruct is nil if this MappedField is just mapped from a single value
	//   from the source object.
	// Otherwise, it's a whole new struct type.
	NewStruct *MappedStruct
	// OriginPath tells us all the field names and types from the root
	//   source object to the value we want to insert at Name in this
	//   destination object (which may not be the root).
	// Only used if NewStruct is nil.
	OriginPath []*SourceField
}

// MappedStruct is a new struct type created by moving around the fields of
//   a source struct.
type MappedStruct struct {
	// Name of the new struct type.
	Name string
	// The fields of this struct.
	Fields []*MappedField
}
