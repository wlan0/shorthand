package mapping

import (
	"go/ast"

	"github.com/golang/glog"
	"github.com/koki/shorthand/inspect"
	"github.com/kr/pretty"
)

// TypeWrapper simplified representation of (a subset of) Go types.
type TypeWrapper interface {
	getWrapper() TypeWrapper
	IdentityMapping(prefix []Choice) MappedAtom
}

// TypeIdent is a fully qualified type identifier.
type TypeIdent struct {
	// PkgPath is an empty string for built-in types like "string"
	PkgPath string
	// Name e.g. "string", "Pod"
	Name string
	// Definition of the type (if not built-in).
	Definition TypeDefinition
}

// Map is map[Key]Value
type Map struct {
	Key   TypeWrapper
	Value TypeWrapper
}

// Slice is []Value
type Slice struct {
	Value TypeWrapper
}

// Sum type, a.k.a. Enum or Tagged Union.
// TODO: Pointer is actually a subset of Sum (not implemented).
// We can implement Sum types using interfaces.
type Sum struct {
	Choices []TypeWrapper
}

// PointerOf construct a Sum type for a pointer to a type "t".
func PointerOf(t TypeWrapper) *Sum {
	return &Sum{[]TypeWrapper{nil, t}}
}

// TypeDefinition defines the contents of a type.
type TypeDefinition interface {
	getTypeDefinition() TypeDefinition
}

// WrapperDefinition for wrapper types (e.g. type MyInt int).
type WrapperDefinition struct {
	Value TypeWrapper
}

// Field in a StructDefinition. Name is "" for anonymous fields.
type Field struct {
	Name string
	Type TypeWrapper
}

// StructDefinition the fields in a struct.
type StructDefinition struct {
	Fields []Field
}

// RootType the type at the focus of the context.
func RootType(context *inspect.Context) TypeWrapper {
	return &TypeIdent{PkgPath: context.PackagePath(),
		Name: context.TypeSpec.Name.Name}
}

// ParseTypeExpr converts Go AST to a simplified (limited) representation.
func ParseTypeExpr(context *inspect.Context, typeExpr ast.Expr) TypeWrapper {
	switch typeExpr := typeExpr.(type) {
	// Ident and Selector are base cases because this is about names.
	case *ast.Ident:
		if typeExpr.Obj != nil {
			return &TypeIdent{PkgPath: context.PackagePath(),
				Name: typeExpr.Name}
		}

		// If no linked ast.Object, then it's a built-in type.
		return &TypeIdent{PkgPath: "",
			Name: typeExpr.Name}
	case *ast.SelectorExpr:
		selectedContext := context.RefocusedWithSelector(typeExpr)
		return &TypeIdent{PkgPath: selectedContext.PackagePath(),
			Name: selectedContext.TypeSpec.Name.Name}
	// These are the recursive cases.
	case *ast.ParenExpr:
		return ParseTypeExpr(context, typeExpr.X)
	case *ast.StarExpr:
		return PointerOf(ParseTypeExpr(context, typeExpr.X))
	case *ast.ArrayType:
		if typeExpr.Len == nil {
			return &Slice{Value: ParseTypeExpr(
				context,
				typeExpr.Elt)}
		}
	case *ast.MapType:
		return &Map{
			Key:   ParseTypeExpr(context, typeExpr.Key),
			Value: ParseTypeExpr(context, typeExpr.Value)}
	}

	glog.Fatal(pretty.Sprintf(
		"unsupported type expr (%# v)",
		typeExpr))
	return nil
}

// DeepParseType fully parses the focused TypeSpec in context as a TypeWrapper.
func DeepParseType(context *inspect.Context) TypeWrapper {
	return &TypeIdent{PkgPath: context.PackagePath(),
		Name:       context.TypeSpec.Name.Name,
		Definition: ParseTypeDefinition(context)}
}

// DeepParseTypeExpr is like ParseTypeExpr except it fills in
//   TypeIdent.Definition.
func DeepParseTypeExpr(context *inspect.Context, typeExpr ast.Expr) TypeWrapper {
	switch typeExpr := typeExpr.(type) {
	// Ident and Selector are base cases because this is about names.
	case *ast.Ident:
		if typeExpr.Obj != nil {
			typeSpec := inspect.IdentTypeSpec(typeExpr)
			if typeSpec == nil {
				glog.Fatal(pretty.Sprintf("expected typeSpec in (%# v)", typeExpr))
			}

			newContext := context.RefocusedWithinPackage(typeSpec)
			return &TypeIdent{PkgPath: context.PackagePath(),
				Name:       typeExpr.Name,
				Definition: ParseTypeDefinition(newContext)}
		}

		// If no linked ast.Object, then it's a built-in type.
		return &TypeIdent{PkgPath: "",
			Name: typeExpr.Name}
	case *ast.SelectorExpr:
		selectedContext := context.RefocusedWithSelector(typeExpr)
		return &TypeIdent{PkgPath: selectedContext.PackagePath(),
			Name:       selectedContext.TypeSpec.Name.Name,
			Definition: ParseTypeDefinition(selectedContext)}
	// These are the recursive cases.
	case *ast.ParenExpr:
		return DeepParseTypeExpr(context, typeExpr.X)
	case *ast.StarExpr:
		return PointerOf(DeepParseTypeExpr(context, typeExpr.X))
	case *ast.ArrayType:
		if typeExpr.Len == nil {
			return &Slice{Value: DeepParseTypeExpr(
				context,
				typeExpr.Elt)}
		}
	case *ast.MapType:
		return &Map{
			Key:   DeepParseTypeExpr(context, typeExpr.Key),
			Value: DeepParseTypeExpr(context, typeExpr.Value)}
	}

	glog.Fatal(pretty.Sprintf(
		"unsupported type expr (%# v)",
		typeExpr))
	return nil
}

// ParseTypeDefinition parses the TypeSpec at the context's focus.
func ParseTypeDefinition(context *inspect.Context) TypeDefinition {
	switch typeExpr := context.TypeSpec.Type.(type) {
	case *ast.StructType:
		fields := []Field{}
		for _, field := range typeExpr.Fields.List {
			fieldType := DeepParseTypeExpr(context, field.Type)
			if field.Names == nil {
				// Anonymous field.
				fields = append(fields, Field{
					Name: "", Type: fieldType})
			} else {
				for _, fieldIdent := range field.Names {
					var fieldName string
					if fieldIdent != nil {
						fieldName = fieldIdent.Name
					}

					fields = append(fields, Field{
						Name: fieldName, Type: fieldType})
				}
			}
		}

		return &StructDefinition{Fields: fields}
	default:
		return &WrapperDefinition{Value: ParseTypeExpr(context, typeExpr)}
	}
}

// IdentityMapping identity.
func (t *TypeIdent) IdentityMapping(prefix []Choice) MappedAtom {
	switch d := t.Definition.(type) {
	case *StructDefinition:
		fields := make([]*MappedField, len(d.Fields))
		for ix, field := range d.Fields {
			newPrefix := AppendedChoice(prefix, &StructChoice{field.Name})
			fields[ix] = &MappedField{
				SourcePath: newPrefix,
				Name:       field.Name,
				Atom:       field.Type.IdentityMapping(newPrefix),
			}
		}
		return &MappedStruct{
			SourcePath: prefix,
			Fields:     fields,
		}
	case *WrapperDefinition:
		// Don't bother transforming wrapper types.
		return nil
	}

	// Don't bother transforming built-in types either.
	return nil
}

// IdentityMapping identity.
func (t *Map) IdentityMapping(prefix []Choice) MappedAtom {
	return &MappedMap{t.Value.IdentityMapping(prefix)}
}

// IdentityMapping identity.
func (t *Slice) IdentityMapping(prefix []Choice) MappedAtom {
	return &MappedSlice{t.Value.IdentityMapping(prefix)}
}

// IdentityMapping identity.
func (t *Sum) IdentityMapping(prefix []Choice) MappedAtom {
	if len(t.Choices) != 2 {
		glog.Fatal("Only Pointer is supported as a Sum type.")
	}

	newPrefix := AppendedChoice(prefix, &SumChoice{1})
	return t.Choices[1].IdentityMapping(newPrefix)
}

func (d *WrapperDefinition) getTypeDefinition() TypeDefinition {
	return d
}

func (d *StructDefinition) getTypeDefinition() TypeDefinition {
	return d
}

func (t *TypeIdent) getWrapper() TypeWrapper {
	return t
}

func (t *Map) getWrapper() TypeWrapper {
	return t
}

func (t *Slice) getWrapper() TypeWrapper {
	return t
}

func (t *Sum) getWrapper() TypeWrapper {
	return t
}
