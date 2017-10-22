package cmd

// This file contains methods for recursively inspecting a type definition
//   and a Print() method to test their implementation.

import (
	"github.com/golang/glog"
	"github.com/koki/shorthand/inspect"
	"golang.org/x/tools/go/loader"

	// We're importing "v1" so it and its dependencies are added to vendor/.
	// This package doesn't need "v1" to build, but it does need access
	//   to the source files at runtime.
	// TODO(ublubu): Make sure "go/loader" is configured to find the vendor
	//   folder even if we run "shorthand" outside its source directory.
	_ "k8s.io/api/core/v1"
)

// load a package and traverse all its types.
func loadAndPrint(typeName string) {
	var conf loader.Config

	// We're just loading "v1" (and then all its dependencies).
	conf.Import("k8s.io/api/core/v1")
	program, err := conf.Load()

	// If we're missing dependencies or our "v1" is otherwise broken, quit.
	if err != nil {
		glog.Error(err)
		return
	}

	v1pkg := program.InitialPackages()[0]
	var contexts []*inspect.Context

	if len(typeName) > 0 {
		// Create a traversal context for the named type.
		context := inspect.ContextForPackageAndType(
			program, v1pkg, typeName)
		if context == nil {
			glog.Errorf("couldn't find v1 type %s", typeName)
			return
		}

		contexts = []*inspect.Context{context}
	} else {
		// Create a traversal context for each type definition in "v1".
		contexts = inspect.ContextsForPackage(program, v1pkg)
	}

	// Test the traversal context by printing all fields "recursively".
	for _, context := range contexts {
		context.Print(0)
	}
}
