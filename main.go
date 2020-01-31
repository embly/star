package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func main() {
	thing()
}

func thing() {
	thread := &starlark.Thread{Name: ""}
	globals, err := starlark.ExecFile(thread, "script.star.py", nil, starlark.StringDict{
		"struct":  starlark.NewBuiltin("struct", starlarkstruct.Make),
		"require": starlark.NewBuiltin("require", require),
	})

	if err != nil {
		panic(err)
	}
	// Retrieve a module global.
	handler := globals["hello"]

	fmt.Println("localhost:8650")
	panic(http.ListenAndServe("localhost:8650", NewHandler(handler)))
}

func require(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	v, err = starlarkstruct.Make(nil, nil, nil, []starlark.Tuple{
		{starlark.String("channel"), starlark.String("channel")},
		{starlark.String("ioutil"), starlark.String("channel")},
	})
	return
}

func writeHeader(w http.ResponseWriter) *starlark.Builtin {
	return starlark.NewBuiltin("write_header", func(thread *starlark.Thread,
		fn *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (v starlark.Value, err error) {
		in, err := strconv.Atoi(args.Index(0).String())
		if err != nil {
			return
		}
		w.WriteHeader(in)
		v = starlark.None
		return
	})
}
func write(w http.ResponseWriter) *starlark.Builtin {
	return starlark.NewBuiltin("write", func(thread *starlark.Thread,
		fn *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (v starlark.Value, err error) {
		_, err = w.Write([]byte(args.Index(0).(starlark.String)))
		v = starlark.None
		return
	})
}

type reader struct {
	reader io.Reader
}

func (r *reader) String() string        { return "io.Reader" }
func (r *reader) Type() string          { return "io.Reader" }
func (r *reader) Freeze()               {}
func (r *reader) Truth() starlark.Bool  { return starlark.True }
func (r *reader) Hash() (uint32, error) { return 0, errors.New("not hashable") }

type bytes struct {
	b []byte
}

func (r *bytes) String() string             { return fmt.Sprint(r.b) }
func (r *bytes) Type() string               { return "bytes" }
func (r *bytes) Freeze()                    {}
func (r *bytes) Len() int                   { return len(r.b) }
func (r *bytes) Truth() starlark.Bool       { return starlark.True }
func (r *bytes) Iterate() starlark.Iterator { return &byteIterator{bytes: r} }
func (r *bytes) Hash() (uint32, error)      { return 0, errors.New("not hashable") }

type byteIterator struct {
	*bytes
	i int
}

func (r *byteIterator) Next(p *starlark.Value) bool {
	if r.i > len(r.bytes.b)-1 {
		return false
	}
	i := starlark.MakeInt(int(r.bytes.b[r.i]))
	*p = &i
	r.i++
	return true
}
func (r *byteIterator) Done() {}

func NewHandler(fn starlark.Value) (handle http.Handler) {
	handle = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		thread := &starlark.Thread{Name: ""}
		request, err := starlarkstruct.Make(nil, nil, nil, []starlark.Tuple{
			{starlark.String("path"), starlark.String(req.URL.Path)},
			{starlark.String("content_type"), starlark.String(req.Header.Get("Content-Type"))},
			{starlark.String("method"), starlark.String(req.Method)},
			{starlark.String("Body"), &reader{reader: req.Body}},
			{starlark.String("some_bytes"), &bytes{b: []byte("hello world")}},
		})
		if err != nil {
			panic(err)
		}
		writer, err := starlarkstruct.Make(nil, nil, nil, []starlark.Tuple{{
			starlark.String("write"), write(w),
		}, {
			starlark.String("write_header"), writeHeader(w),
		}})
		if err != nil {
			panic(err)
		}
		_, err = starlark.Call(thread, fn, starlark.Tuple{writer, request}, nil)
		if err != nil {
			panic(err)
		}
	})
	return
}
