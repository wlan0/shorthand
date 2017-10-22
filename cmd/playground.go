package cmd

// This file contains methods for recursively inspecting a type definition
//   and a Print() method to test their implementation.

import (
	"github.com/koki/shorthand/gen"
	"github.com/koki/shorthand/gen/samples"
	"github.com/kr/pretty"
)

func genThings() {
	//gen.PrintFileAST(samples.PodSrc)
	_, _ = pretty.Println(gen.SerializeFileAST(samples.PodAST).String())
}
