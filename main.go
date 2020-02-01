package star

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

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

func NewHandler(fn starlark.Value) (handle http.Handler) {

	handle = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		thread := &starlark.Thread{Name: ""}
		request, err := starlarkstruct.Make(nil, nil, nil, []starlark.Tuple{
			{starlark.String("path"), starlark.String(req.URL.Path)},
			{starlark.String("content_type"), starlark.String(req.Header.Get("Content-Type"))},
			{starlark.String("method"), starlark.String(req.Method)},
			{starlark.String("Body"), &reader{reader: req.Body}},
			{starlark.String("some_ByteArray"), &ByteArray{b: []byte("hello world")}},
			{starlark.String("error"), Error{err: errors.Wrap(errors.New("this is the error"), "this is it wrapped")}},
			{starlark.String("sample_resp"), starlark.Tuple{
				&ByteArray{b: []byte("hello world")},
				Error{err: errors.Wrap(errors.New("this is the error"), "this is it wrapped")}}},
			{starlark.String("nil_err"), starlark.Tuple{
				&ByteArray{b: []byte("hello world")},
				Error{err: nil}}},
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
