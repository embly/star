package star

import (
	"github.com/pkg/errors"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var builtins = map[string]struct{}{
	"string": struct{}{},
	"error":  struct{}{},
	"int":    struct{}{},
}

func InitPackages(packages map[string]map[string]starlark.Value) {
	mappings = packages
}

var mappings = map[string]map[string]starlark.Value{}

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
