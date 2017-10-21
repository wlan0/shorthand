package inspect

// This file contains methods for recursively inspecting a type definition
//   and a Print() method to test their implementation.

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/golang/glog"
	"github.com/kr/pretty"
	"golang.org/x/tools/go/loader"
)

// Context is a cursor that indicates our current position in the program.
// It contains the information needed to understand every part of TypeSpec.
type Context struct {
	// Program is the entire loaded program.
	Program *loader.Program
	// Package is an entire loaded package.
	Package *loader.PackageInfo
	// File is a file in the Package.
	// The file matters because the local name of an import is
	//   file-specific. This is how we resolve Selector expressions.
	//   e.g. metav1.ObjectMeta
	File *ast.File
	// TypeSpec is a type definition in the File.
	TypeSpec *ast.TypeSpec
}

// importMatchesPackage checks if a package corresponds to an import.
// types.Package uses the actual filesystem path (for non-built-in packages)
//   e.g. github.com/koki/shorthand/vendor/k8s.io/apimachinery
//        go/ast
// ast.ImportSpec uses the quoted string in the source file's import statement.
//   e.g. "k8s.io/apimachinery" <- with quotes
//        "go/ast"
func importMatchesPackage(imprt *ast.ImportSpec, pkg *types.Package) bool {
	quoted := imprt.Path.Value
	stripped := quoted[1 : len(quoted)-1]
	return pathMatchesPackage(stripped, pkg)
}

// pathMatchesPackage checks if a path corresponds to a given package.
// pkgPath can be either of two formats:
//   e.g. github.com/koki/shorthand/vendor/k8s.io/apimachinery (types.Package)
//   e.g. k8s.io/apimachinery (unquoted ast.ImportSpec)
func pathMatchesPackage(pkgPath string, pkg *types.Package) bool {
	// If pkgPath is the long format, then match this way.
	if pkgPath == pkg.Path() {
		return true
	}

	// If pkgPath is the short format, then match this way.
	if strings.HasSuffix(pkg.Path(), "vendor/"+pkgPath) {
		return true
	}

	return false
}

// ContextForImportedType constructs a context for an import
//   and the name of a type.
// It traverses the Program to find the right Package and TypeSpec.
func ContextForImportedType(program *loader.Program, imprt *ast.ImportSpec, typeName string) *Context {
	quoted := imprt.Path.Value
	stripped := quoted[1 : len(quoted)-1]
	return ContextForType(program, stripped, typeName)
}

// ContextForType constructs a context for a pkg path (see pathMatchesPackage)
//   and the name of a type.
// It traverses the Program to find the right Package and TypeSpec.
func ContextForType(program *loader.Program, pkgPath, typeName string) *Context {
	for _, pkg := range program.AllPackages {
		if !pathMatchesPackage(pkgPath, pkg.Pkg) {
			continue
		}

		context := ContextForPackageAndType(program, pkg, typeName)
		if context != nil {
			return context
		}
	}

	return nil
}

// ContextForPackageAndType constructs a context for a package and the name of a type.
// It traverses only the given package to find the right TypeSpec.
func ContextForPackageAndType(program *loader.Program, pkg *loader.PackageInfo, typeName string) *Context {
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				if decl.Tok == token.TYPE {
					for _, spec := range decl.Specs {
						switch spec := spec.(type) {
						case *ast.TypeSpec:
							if spec.Name.Name == typeName {
								return &Context{Program: program, Package: pkg, File: file, TypeSpec: spec}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// ContextsForPackage traverses an entire package and creates a context for each
//   type definition it finds.
func ContextsForPackage(program *loader.Program, pkg *loader.PackageInfo) []*Context {
	contexts := []*Context{}
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				if decl.Tok == token.TYPE {
					for _, spec := range decl.Specs {
						switch spec := spec.(type) {
						case *ast.TypeSpec:
							contexts = append(contexts, &Context{Program: program, Package: pkg, File: file, TypeSpec: spec})
						}
					}
				}
			}
		}
	}

	return contexts
}

// RefocusedWithinPackage navigates from one type (TypeSpec) to another within
//   the same package.
// For example, this would be used to go from inspecting v1.Pod to v1.PodSpec
//   When we inspect v1.Pod, the Spec field gives us a TypeSpec object for
//   v1.PodSpec.
func (context *Context) RefocusedWithinPackage(typeSpec *ast.TypeSpec) *Context {
	for _, file := range context.Package.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				if decl.Tok == token.TYPE {
					for _, spec := range decl.Specs {
						if typeSpec == spec {
							return &Context{Program: context.Program, Package: context.Package, File: file, TypeSpec: typeSpec}
						}
					}
				}
			}
		}
	}

	return nil
}

// RefocusedWithSelector navigates from one context to a context in a different package.
// For example, this would be used to go from v1.Pod to metav1.ObjectMeta.
//   When we inspect v1.Pod, an anonymous field gives us a Selector for metav1.ObjectMeta.
func (context *Context) RefocusedWithSelector(selector *ast.SelectorExpr) *Context {
	var pkgName string
	switch expr := selector.X.(type) {
	case *ast.Ident:
		pkgName = expr.Name
	default:
		glog.Fatal(pretty.Sprint(selector))
	}

	typeName := selector.Sel.Name

	return context.refocusedWithPkgAndTypeNames(pkgName, typeName)
}

// getImportedPackage traverses the Program to find the Package that matches
//   a given Import.
// It's a helper function for refocusedWithPkgAndTypeNames.
func (context *Context) getImportedPackage(imprt *ast.ImportSpec) *loader.PackageInfo {
	for _, pkg := range context.Program.AllPackages {
		if importMatchesPackage(imprt, pkg.Pkg) {
			return pkg
		}
	}

	return nil
}

// refocusedWithPkgAndTypeNames traverses the Program to find a package with
//   a given name (not the path, but the name used to prefix imported definitions)
// It's a helper function for RefocusedWithSelector.
func (context *Context) refocusedWithPkgAndTypeNames(pkgName string, typeName string) *Context {
	for _, imprt := range context.File.Imports {
		if imprt.Name != nil {
			// If the local name matches, look up the type here.
			if imprt.Name.Name == pkgName {
				return ContextForImportedType(context.Program, imprt, typeName)
			}
		} else {
			// If the default name matches, look up the type here.
			pkg := context.getImportedPackage(imprt)
			if pkg.Pkg.Name() == pkgName {
				return ContextForPackageAndType(context.Program, pkg, typeName)
			}
		}
	}

	return nil
}
