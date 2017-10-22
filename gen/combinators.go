package gen

import (
	"fmt"
	"go/ast"
	"go/token"
)

// BlockOf is "{statements...}"
func BlockOf(statements ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: statements}
}

// PointerOf is "*typeExpr"
func PointerOf(typeExpr ast.Expr) *ast.StarExpr {
	return &ast.StarExpr{X: typeExpr}
}

// ImportsOf is "import (imports...)".
func ImportsOf(imports ...*ast.ImportSpec) *ast.GenDecl {
	specs := make([]ast.Spec, len(imports))
	for ix, imprt := range imports {
		specs[ix] = imprt
	}

	return &ast.GenDecl{Tok: token.IMPORT, Specs: specs}
}

// ImportOf is "localName \"pkg\"".
func ImportOf(localName string, pkg string) *ast.ImportSpec {
	var localIdent *ast.Ident
	if len(localName) > 0 {
		localIdent = IdentFor(localName)
	}

	return &ast.ImportSpec{
		Name: localIdent,
		Path: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", pkg)},
	}
}

// FieldOf is "name typeExpr".
func FieldOf(name string, typeExpr ast.Expr) *ast.Field {
	return &ast.Field{
		Names: []*ast.Ident{IdentFor(name)},
		Type:  typeExpr,
	}
}

// SimpleFuncType is "func(paramName paramType) returnType".
func SimpleFuncType(paramName string, paramType ast.Expr, returnType ast.Expr) *ast.FuncType {
	return &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				FieldOf(paramName, paramType),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Type: returnType,
				},
			},
		},
	}
}

// VarOfType is "var name typeSelectorOrIdent".
func VarOfType(name string, typeSelectorOrIdent ast.Expr) *ast.DeclStmt {
	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{IdentFor(name)},
					Type:  typeSelectorOrIdent,
				},
			},
		},
	}
}

// ReturnAddressOf is "return &thingName".
func ReturnAddressOf(thingName string) ast.Stmt {
	return &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X:  IdentFor(thingName),
			},
		},
	}
}

// AssignmentOf is "lhs = rhs".
func AssignmentOf(lhs, rhs ast.Expr) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{rhs},
	}
}

// IdentFor is "thingName".
func IdentFor(thingName string) *ast.Ident {
	return &ast.Ident{Name: thingName}
}

// SelectorOrIdentForV is the variadic version of SelectorOrIdentFor.
func SelectorOrIdentForV(segments ...string) ast.Expr {
	return SelectorOrIdentFor(segments)
}

// SelectorOrIdentFor is "segment0.segment1.segment..."
func SelectorOrIdentFor(segments []string) ast.Expr {
	var result ast.Expr
	for _, segment := range segments {
		if result == nil {
			result = IdentFor(segment)
		} else {
			result = &ast.SelectorExpr{
				X:   result,
				Sel: IdentFor(segment),
			}
		}
	}

	return result
}

// SimpleCallOf is "funcName(argName)".
func SimpleCallOf(funcName string, argName string) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  IdentFor(funcName),
		Args: []ast.Expr{IdentFor(argName)},
	}
}
