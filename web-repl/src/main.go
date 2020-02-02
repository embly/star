package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/embly/star"
	"github.com/embly/star/src"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func main() {
	registerCallbacks()
	c := make(chan struct{}, 0)
	star.AddPackages(src.Packages)
	thread := &starlark.Thread{Name: ""}
	globals := starlark.StringDict{
		"require": starlark.NewBuiltin("require", star.Require),
	}
	js.Global().Get("window").Call("go_ready")
	REPL(thread, globals)
	<-c

}

var lineChan = make(chan string)

func registerCallbacks() {
	js.Global().Set("newLine", js.FuncOf(newLine))
}

func newLine(this js.Value, args []js.Value) interface{} {
	lineChan <- args[0].String()
	return nil
}

func REPL(thread *starlark.Thread, globals starlark.StringDict) {
	for {
		if err := rep(thread, globals); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println()
}

// rep reads, evaluates, and prints one item.
//
// It returns an error (possibly readline.ErrInterrupt)
// only if readline failed. Starlark errors are printed.
func rep(thread *starlark.Thread, globals starlark.StringDict) error {
	// Each item gets its own context,
	// which is cancelled by a SIGINT.
	//
	// Note: during Readline calls, Control-C causes Readline to return
	// ErrInterrupt but does not generate a SIGINT.

	eof := false
	js.Global().Get("window").Call("set_prompt", js.ValueOf(">>> "))
	// readline returns EOF, ErrInterrupted, or a line including "\n".
	// rl.SetPrompt(">>> ")
	readline := func() ([]byte, error) {
		line := <-lineChan
		// line, err := rl.Readline()
		js.Global().Get("window").Call("set_prompt", js.ValueOf("... "))
		// rl.SetPrompt("... ")
		// if err != nil {
		// 	if err == io.EOF {
		// 		eof = true
		// 	}
		// 	return nil, err
		// }
		return []byte(line + "\n"), nil
	}

	// parse
	f, err := syntax.ParseCompoundStmt("<stdin>", readline)
	if err != nil {
		if eof {
			return io.EOF
		}
		PrintError(err)
		return nil
	}

	// Treat load bindings as global (like they used to be) in the REPL.
	// This is a workaround for github.com/google/starlark-go/issues/224.
	// TODO(adonovan): not safe wrt concurrent interpreters.
	// Come up with a more principled solution (or plumb options everywhere).
	defer func(prev bool) { resolve.LoadBindsGlobally = prev }(resolve.LoadBindsGlobally)
	resolve.LoadBindsGlobally = true

	if expr := soleExpr(f); expr != nil {
		// eval
		v, err := starlark.EvalExpr(thread, expr, globals)
		if err != nil {
			PrintError(err)
			return nil
		}

		// print
		if v != starlark.None {
			Println(v)
		}
	} else if err := starlark.ExecREPLChunk(f, thread, globals); err != nil {
		PrintError(err)
		return nil
	}

	return nil
}

func soleExpr(f *syntax.File) syntax.Expr {
	if len(f.Stmts) == 1 {
		if stmt, ok := f.Stmts[0].(*syntax.ExprStmt); ok {
			return stmt.X
		}
	}
	return nil
}

func Println(v interface{}) {
	fmt.Println(strconv.Quote(fmt.Sprint(v)))

	js.Global().Get("window").Call(
		"term_write",
		js.ValueOf(
			strings.Replace(
				fmt.Sprintln(v), "\n", "\n\r", -1),
		))
}

// PrintError prints the error to stderr,
// or its backtrace if it is a Starlark evaluation error.
func PrintError(err error) {
	if evalErr, ok := err.(*starlark.EvalError); ok {
		Println(evalErr.Backtrace())
	} else {
		Println(err)
	}
}
