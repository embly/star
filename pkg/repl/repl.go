package repl

// This has to be in its own package so that we don't have to compile
// github.com/chzyer/readline when compiling with webassembly

import (
	"github.com/embly/star/pkg/star"
	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

func REPL() {
	thread := &starlark.Thread{Name: ""}
	repl.REPL(thread, starlark.StringDict{
		"require": starlark.NewBuiltin("require", star.Require),
	})
}
