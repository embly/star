package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"go.starlark.net/starlark"
)

func BenchmarkThing(b *testing.B) {
	// Execute Starlark program in a file.
	thread := &starlark.Thread{Name: ""}
	globals, err := starlark.ExecFile(thread, "script.star.py", nil, nil)
	if err != nil {
		panic(err)
	}
	f, err := ioutil.TempFile("", "starlark-prof")
	if err != nil {
		panic(err)
	}
	// Retrieve a module global.
	fibonacci := globals["hello"]

	for i := 0; i < b.N; i++ {
		// Call Starlark function from Go.
		v, err := starlark.Call(thread, fibonacci, starlark.Tuple{}, nil)
		if err != nil {
			panic(err)
		}
		fmt.Println(f.Name(), len(fmt.Sprint(v)))
	}

}
