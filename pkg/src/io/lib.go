package io

import (
	"io"

	"github.com/embly/star/pkg/star"
	"go.starlark.net/starlark"
)

var Reader = star.Interface{
	Name: "io.Reader",
	Methods: map[string]star.Method{
		"Read": star.Method{
			Args:    []starlark.Value{star.ByteArray{}},
			Returns: []starlark.Value{starlark.Int{}, star.Error{}},
			Call: func(recv interface{}, args []interface{}) (resp []interface{}) {
				resp = make([]interface{}, 2)
				resp[0], resp[1] = recv.(io.Reader).Read(args[0].([]byte))
				return
			},
		},
	},
}

var ReadCloser = star.Interface{
	Name: "io.ReadCloser",
	Methods: map[string]star.Method{
		"Read": star.Method{
			Args:    []starlark.Value{star.ByteArray{}},
			Returns: []starlark.Value{starlark.Int{}, star.Error{}},
			Call: func(recv interface{}, args []interface{}) (resp []interface{}) {
				resp = make([]interface{}, 2)
				resp[0], resp[1] = recv.(io.ReadCloser).Read(args[0].([]byte))
				return
			},
		},
		"Close": star.Method{
			Returns: []starlark.Value{star.Error{}},
			Call: func(recv interface{}, args []interface{}) (resp []interface{}) {
				resp = make([]interface{}, 1)
				resp[0] = recv.(io.ReadCloser).Close()
				return
			},
		},
	},
}
