package samples

import (
	"go/ast"

	"github.com/koki/shorthand/gen"
)

// PodSrc is Go code that corresponds to PodAST below.
const PodSrc = `
package pod
import (
v1 "k8s.io/api/core/v1"
)
func Pod(pod *v1.Pod) *Pod {
var pod1 Pod
pod1.Doot = pod.Spec.Doot
pod1.Beep = Beep(pod)
return &pod1
}
`

// PodAST is a sample artificial AST to serialize into Go code.
var PodAST = &ast.File{
	Name: gen.IdentFor("pod"),
	Decls: []ast.Decl{
		gen.ImportsOf(gen.ImportOf("v1", "k8s.io/api/core/v1")),
		&ast.FuncDecl{
			Name: gen.IdentFor("Pod"),
			Type: gen.SimpleFuncType(
				"pod",
				gen.PointerOf(gen.SelectorOrIdentForV("v1", "Pod")),
				gen.PointerOf(gen.IdentFor("Pod"))),
			Body: gen.BlockOf(
				gen.VarOfType("pod1", gen.IdentFor("Pod")),
				gen.AssignmentOf(
					gen.SelectorOrIdentForV("pod1", "Doot"),
					gen.SelectorOrIdentForV("pod", "Spec", "Doot")),
				gen.AssignmentOf(
					gen.SelectorOrIdentForV("pod1", "Beep"),
					gen.SimpleCallOf("Beep", "pod")),
				gen.ReturnAddressOf("pod1"),
			),
		},
	},
}
