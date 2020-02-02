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
var Request = star.Struct{
	TypeName: "Request",
}

var ResponseWriter = star.Interface{
	Name: "ResponseWriter",
	Methods: map[string]star.Method{
		"Write": star.Method{
			// TODO: should accept lots of things?
			Args:    []starlark.Value{starlark.String("")},
			Returns: []starlark.Value{starlark.Int{}, star.Error{}},
			Call: func(recv interface{}, args []interface{}) (resp []interface{}) {
				resp = make([]interface{}, 2)
				resp[0], resp[1] = recv.(http.ResponseWriter).Write([]byte(args[0].(string)))
				return
			},
		},
		"WriteHeader": star.Method{
			Args: []starlark.Value{starlark.Int{}},
			Call: func(recv interface{}, args []interface{}) (resp []interface{}) {
				recv.(http.ResponseWriter).WriteHeader(args[0].(int))
				return
			},
		},
	},
}

var Handler = star.Struct{
	TypeName: "Handler",
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

var ListenAndServe = star.Function{
	FunctionName: "ListenAndServe",
	Args:         []starlark.Value{starlark.String(""), Handler},
	Returns:      []starlark.Value{star.Error{}},
	Call: func(args []interface{}) (resp []interface{}) {
		var handler http.Handler
		if args[1] != nil {
			handler = args[1].(http.Handler)
		}
		resp = make([]interface{}, 1)
		resp[0] = http.ListenAndServe(args[0].(string), handler)
		return
	},
}

var HandleFunc = star.Function{
	FunctionName: "HandleFunc",
	Args:         []starlark.Value{starlark.String(""), &starlark.Function{}},
	Call: func(args []interface{}) (resp []interface{}) {
		http.HandleFunc(args[0].(string), func(w http.ResponseWriter, req *http.Request) {
			fn := args[1].(*starlark.Function)
			starRequest := Request
			starRequest.Value = req
			starWriter := ResponseWriter
			starWriter.Value = w

			// we want a thread for every request goroutine
			thread := &starlark.Thread{Name: ""}
			_, err := starlark.Call(thread, fn, starlark.Tuple([]starlark.Value{starWriter, starRequest}), []starlark.Tuple{})
			if err != nil {
				panic(err)
			}
		})
		return
	},
}
