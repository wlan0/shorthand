package mapping

import (
	"go/ast"

	"github.com/golang/glog"
	"github.com/koki/shorthand/inspect"
	"github.com/kr/pretty"
)

// TypeWrapper wrapper for TypeIdents.
type TypeWrapper interface {
	getWrapper() TypeWrapper
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

// Pointer is *Value
type Pointer struct {
	Value TypeWrapper
}

type TypeDefinition interface {
	getTypeDefinition() TypeDefinition
}

type WrapperDefinition struct {
	Value TypeWrapper
}

type Field struct {
	Name string
	Type TypeWrapper
}

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
		return &Pointer{Value: ParseTypeExpr(context, typeExpr.X)}
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

	glog.Fatal(pretty.Printf(
		"unsupported type expr (%# v)",
		typeExpr))
	return nil
}

// DeepParseType is like ParseTypeExpr except it fills in TypeIdent.Definition
func DeepParseType(context *inspect.Context) TypeWrapper {
	return &TypeIdent{PkgPath: context.PackagePath(),
		Name:       context.TypeSpec.Name.Name,
		Definition: ParseTypeDefinition(context)}
}

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
		return &Pointer{Value: DeepParseTypeExpr(context, typeExpr.X)}
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

	glog.Fatal(pretty.Printf(
		"unsupported type expr (%# v)",
		typeExpr))
	return nil
}

func ParseTypeDefinition(context *inspect.Context) TypeDefinition {
	switch typeExpr := context.TypeSpec.Type.(type) {
	case *ast.StructType:
		fields := []Field{}
		for _, field := range typeExpr.Fields.List {
			fieldType := DeepParseTypeExpr(context, field.Type)
			for _, fieldIdent := range field.Names {
				var fieldName string
				if fieldIdent != nil {
					fieldName = fieldIdent.Name
				}

				fields = append(fields, Field{
					Name: fieldName, Type: fieldType})
			}
		}

		return &StructDefinition{Fields: fields}
	default:
		return &WrapperDefinition{Value: ParseTypeExpr(context, typeExpr)}
	}
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

func (t *Pointer) getWrapper() TypeWrapper {
	return t
}
