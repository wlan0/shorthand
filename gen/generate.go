package gen

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/golang/glog"
	"github.com/koki/shorthand/gen/samples"
	"github.com/kr/pretty"
)

func printFileAST(file string) {
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", file, 0)
	if err != nil {
		glog.Fatal(err)
	}

	// Print the AST.
	_ = ast.Print(fset, f)
	_, _ = pretty.Println(f)
}

func serializeFileAST(file *ast.File) *bytes.Buffer {
	// Create a FileSet for node. Since the node does not come
	// from a real source file, fset will be empty.
	fset := token.NewFileSet()

	var buf bytes.Buffer
	err := format.Node(&buf, fset, file)
	if err != nil {
		glog.Fatal(err)
	}

	return &buf
}

func printAllTokens() {
	for i := 0; i <= int(token.VAR); i++ {
		_, _ = pretty.Printf("%d: %s\n", i, token.Token(i).String())
	}
}

func Generate() {
	printFileAST(samples.PodSrc)
	_, _ = pretty.Println(serializeFileAST(samples.PodAST).String())
}
