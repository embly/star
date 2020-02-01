package star

import (
	"os"

	"go.starlark.net/starlark"
)

func RunScript(file string) {
	thread := &starlark.Thread{Name: ""}
	globals, err := starlark.ExecFile(thread, os.Args[1], nil, starlark.StringDict{
		"require": starlark.NewBuiltin("require", Require),
	})
	if err != nil {
		panic(err)
	}
	main := globals["main"]
	_, err = starlark.Call(thread, main, starlark.Tuple{}, nil)
	if err != nil {
		panic(err)
	}
}
