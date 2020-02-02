package http

import (
	"net/http"

	"github.com/embly/star"
	"github.com/embly/star/src/io"
	"go.starlark.net/starlark"
)

var Response = star.Struct{
	TypeName: "Response",
	Attributes: map[string]starlark.Value{
		"Body": io.ReadCloser,
	},
}

var Get = star.Function{
	FunctionName: "Get",
	Args:         []starlark.Value{starlark.String("")},
	Returns:      []starlark.Value{Response, star.Error{}},
	Call: func(args []interface{}) (resp []interface{}) {
		resp = make([]interface{}, 2)
		resp[0], resp[1] = http.Get(args[0].(string))
		return
	},
}
