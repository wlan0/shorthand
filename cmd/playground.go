package cmd

// This file is just for trying things out during development.

import (
	"github.com/koki/shorthand/gen"
	"github.com/koki/shorthand/gen/samples"
	"github.com/kr/pretty"
)

func genThings() {
	//gen.PrintFileAST(samples.PodSrc)
	_, _ = pretty.Println(gen.SerializeFileAST(samples.PodAST).String())
}
