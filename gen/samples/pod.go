package samples

import (
	"go/ast"
	"go/token"
)

const PodSrc = `
package pod
import (
v1 "k8s.io/api/core/v1"
)
func Pod(pod *v1.Pod) *Pod {
var pod1 Pod
pod1.Doot = pod.Spec.Doot
return &pod1
}
`

var PodAST = &ast.File{
	Doc:     (*ast.CommentGroup)(nil),
	Package: 2,
	Name: &ast.Ident{
		Name: "pod",
		Obj:  (*ast.Object)(nil),
	},
	Decls: []ast.Decl{
		&ast.GenDecl{
			Doc: (*ast.CommentGroup)(nil),
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Doc: (*ast.CommentGroup)(nil),
					Name: &ast.Ident{
						Name: "v1",
						Obj:  (*ast.Object)(nil),
					},
					Path:    &ast.BasicLit{Kind: 9, Value: "\"k8s.io/api/core/v1\""},
					Comment: (*ast.CommentGroup)(nil),
				},
			},
		},
		&ast.FuncDecl{
			Doc:  (*ast.CommentGroup)(nil),
			Recv: (*ast.FieldList)(nil),
			Name: &ast.Ident{
				Name: "Pod",
				Obj:  nil,
			},
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Doc: (*ast.CommentGroup)(nil),
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "pod",
									Obj:  nil,
								},
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "v1",
										Obj:  (*ast.Object)(nil),
									},
									Sel: &ast.Ident{
										Name: "Pod",
										Obj:  (*ast.Object)(nil),
									},
								},
							},
							Tag:     (*ast.BasicLit)(nil),
							Comment: (*ast.CommentGroup)(nil),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Doc:   (*ast.CommentGroup)(nil),
							Names: nil,
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Pod",
									Obj:  nil,
								},
							},
							Tag:     (*ast.BasicLit)(nil),
							Comment: (*ast.CommentGroup)(nil),
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.DeclStmt{
						Decl: &ast.GenDecl{
							Doc: (*ast.CommentGroup)(nil),
							Tok: 85,
							Specs: []ast.Spec{
								&ast.ValueSpec{
									Doc: (*ast.CommentGroup)(nil),
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "pod1",
											Obj:  nil,
										},
									},
									Type: &ast.Ident{
										Name: "Pod",
										Obj:  nil,
									},
									Values:  nil,
									Comment: (*ast.CommentGroup)(nil),
								},
							},
						},
					},
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.Ident{
									Name: "pod1",
									Obj:  nil,
								},
								Sel: &ast.Ident{
									Name: "Doot",
									Obj:  (*ast.Object)(nil),
								},
							},
						},
						Tok: 42,
						Rhs: []ast.Expr{
							&ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "pod",
										Obj:  nil,
									},
									Sel: &ast.Ident{
										Name: "Spec",
										Obj:  (*ast.Object)(nil),
									},
								},
								Sel: &ast.Ident{
									Name: "Doot",
									Obj:  (*ast.Object)(nil),
								},
							},
						},
					},
					&ast.ReturnStmt{
						Return: 117,
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: 17,
								X: &ast.Ident{
									Name: "pod1",
									Obj:  nil,
								},
							},
						},
					},
				},
			},
		},
	},
	Scope: &ast.Scope{
		Outer:   (*ast.Scope)(nil),
		Objects: map[string]*ast.Object{},
	},
	Imports:    []*ast.ImportSpec{},
	Unresolved: []*ast.Ident{},
	Comments:   nil,
}
