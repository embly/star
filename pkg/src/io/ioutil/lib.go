package ioutil

import (
	go_io "io"
	"io/ioutil"

	"github.com/embly/star/pkg/star"
	"github.com/embly/star/pkg/src/io"
	"go.starlark.net/starlark"
)

var NopCloser = star.Function{
	FunctionName: "NopCloser",
	Args:         []starlark.Value{io.Reader},
	Returns:      []starlark.Value{io.ReadCloser},
	Call: func(args []interface{}) (resp []interface{}) {
		resp = make([]interface{}, 1)
		resp[0] = ioutil.NopCloser(args[0].(go_io.Reader))
		return
	},
}

var ReadAll = star.Function{
	FunctionName: "ReadAll",
	Args:         []starlark.Value{io.Reader},
	Returns:      []starlark.Value{star.ByteArray{}, star.Error{}},
	Call: func(args []interface{}) (resp []interface{}) {
		resp = make([]interface{}, 2)
		resp[0], resp[1] = ioutil.ReadAll(args[0].(go_io.Reader))
		return
	},
}
