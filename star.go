package star

import (
	"fmt"
	"os"

	"go.starlark.net/repl"
	"go.starlark.net/starlark"
)

func RunScript(file string) (err error) {
	thread := &starlark.Thread{Name: ""}
	globals, err := starlark.ExecFile(thread, os.Args[1], nil, starlark.StringDict{
		"require": starlark.NewBuiltin("require", Require),
	})
	if err != nil {
		return
	}
	main := globals["main"]
	_, err = starlark.Call(thread, main, starlark.Tuple{}, nil)
	if err != nil {
		if er, ok := err.(*starlark.EvalError); ok {
			fmt.Println(er.Backtrace())
		}
		return
	}
	return
}

func REPL() {
	thread := &starlark.Thread{Name: ""}
	repl.REPL(thread, starlark.StringDict{
		"require": starlark.NewBuiltin("require", Require),
	})
}
