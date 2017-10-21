package inspect

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/golang/glog"
	"github.com/kr/pretty"
)

// These (and init()) are just for formatting the Print methods.
var maxDepth = 40
var indents = make([]string, maxDepth)

func init() {
	for index := range indents {
		indents[index] = strings.Repeat("  ", index)
	}
}

// Print recursively traverses all the fields of a type and prints them.
// "depth" is the initial indentation depth.
func (context *Context) Print(depth int) {
	if context == nil {
		glog.Fatal()
	}

	fmt.Println(indents[depth], context.TypeSpec.Name)
	context.PrintType(depth+1, context.TypeSpec.Type)
}

// PrintType prints the contents (RHS) of a type declaration.
// "depth" is the initial indentation depth.
func (context *Context) PrintType(depth int, root ast.Expr) {
	switch root := root.(type) {
	case *ast.ParenExpr:
		// Strip parens
		context.PrintType(depth, root.X)
	case *ast.Ident:
		// Print name and then associated object.
		context.PrintTypeIdent(depth, root)
	case *ast.SelectorExpr:
		context.PrintSelector(depth, root)
	case *ast.StarExpr:
		fmt.Println(indents[depth], "*")
		context.PrintType(depth, root.X)
	case *ast.FuncType:
		fmt.Println(indents[depth], "func")
	case *ast.ChanType:
		fmt.Println(indents[depth], "chan")
	case *ast.ArrayType:
		fmt.Println(indents[depth], "[]")
		context.PrintType(depth, root.Elt)
	case *ast.StructType:
		for _, field := range root.Fields.List {
			context.PrintField(depth, field)
		}
	case *ast.MapType:
		fmt.Println(indents[depth], "map")
		context.PrintType(depth, root.Key)
		context.PrintType(depth, root.Value)
	case *ast.InterfaceType:
		fmt.Println(indents[depth], "interface")
	default:
		glog.Fatal(pretty.Printf("non-type expr (%# v)", root))
	}
}

// PrintSelector prints a Selector and the details of the type it represents.
// "depth" is the initial indentation depth.
func (context *Context) PrintSelector(depth int, root *ast.SelectorExpr) {
	var pkgName string
	switch expr := root.X.(type) {
	case *ast.Ident:
		pkgName = expr.Name
	default:
		glog.Fatal(pretty.Sprint(root))
	}

	typeName := root.Sel.Name

	fmt.Printf("%s%s.%s\n", indents[depth], pkgName, typeName)
	context.RefocusedWithSelector(root).Print(depth)
}

// PrintTypeIdent prints a type identifier along with its contents.
// An Ident either refers to a built-in type or a type within the same Package.
// "depth" is the initial indentation depth.
func (context *Context) PrintTypeIdent(depth int, root *ast.Ident) {
	obj := root.Obj
	if obj != nil {
		switch obj.Kind {
		case ast.Typ:
			switch decl := obj.Decl.(type) {
			case *ast.TypeSpec:
				context.RefocusedWithinPackage(decl).Print(depth)
			default:
				fmt.Println(indents[depth], root.Name)
				fmt.Println(indents[depth+1], "Typ but no TypeSpec")
			}

			return

		default:
			fmt.Println(indents[depth], root.Name)
			fmt.Println(indents[depth+1], obj)
		}

		return
	}

	fmt.Println(indents[depth], root.Name)
}

// PrintField prints a struct Field and its type information.
// "depth" is the initial indentation depth.
func (context *Context) PrintField(depth int, root *ast.Field) {
	l := len(root.Names)
	if l > 0 {
		names := make([]string, l)
		for ix, ident := range root.Names {
			names[ix] = ident.Name
		}

		fmt.Println(indents[depth], "-", strings.Join(names, ", "))
	} else {
		fmt.Println(indents[depth], "<anonymous>")
	}

	context.PrintType(depth+1, root.Type)
}
