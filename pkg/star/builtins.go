package star

import (
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func AddPackages(packages map[string]map[string]starlark.Value) {
	for name, members := range packages {
		mappings[name] = members
	}
}

var mappings = map[string]map[string]starlark.Value{
	"star": Star,
}

func Require(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	if args.Len() != 1 || args.Index(0).Type() != "string" {
		err = errors.New("require takes one string argument")
		return
	}

	packageName := string(args.Index(0).(starlark.String))
	values, ok := mappings[packageName]
	if !ok {
		err = errors.Errorf(`package doesn't exist with name "%s"`, packageName)
		return
	}
	members := []starlark.Tuple{}
	for name, value := range values {
		members = append(members, starlark.Tuple{starlark.String(name), value})
	}
	v, err = starlarkstruct.Make(nil, nil, nil, members)
	return
}

var Star = map[string]starlark.Value{
	"bytes_to_string": starlark.NewBuiltin("bytes_to_string", func(thread *starlark.Thread,
		fn *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (v starlark.Value, err error) {
		out := []starlark.Value{starlark.None, Error{}}
		if args.Len() != 1 || args.Index(0).Type() != "[]byte" {
			out[1] = Error{err: errors.New("bytes_to_string takes one argument of type []byte")}
		} else {
			out[0] = starlark.String(string(args.Index(0).(ByteArray).b))
		}
		v = starlark.Tuple(out)
		return
	}),
	"chan": starlark.NewBuiltin("chan", func(thread *starlark.Thread,
		fn *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (v starlark.Value, err error) {
		out := []starlark.Value{starlark.None, Error{}}
		if args.Len() != 1 || args.Index(0).Type() != "[]byte" {
			out[1] = Error{err: errors.New("bytes_to_string takes one argument of type []byte")}
		} else {
			out[0] = starlark.String(string(args.Index(0).(ByteArray).b))
		}
		v = starlark.Tuple(out)
		return
	}),
	"go": starlark.NewBuiltin("go", func(thread *starlark.Thread,
		fn *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (v starlark.Value, err error) {

		go func() {
			argsToPass := []starlark.Value(args)[1:]
			thread := &starlark.Thread{Name: ""}
			_, err := starlark.Call(thread, args.Index(0), starlark.Tuple(argsToPass), nil)
			if err != nil {
				panic(err)
			}
		}()
		v = starlark.None
		return
	}),
}
