package sync

import (
	"sync"

	"github.com/embly/star"
	"go.starlark.net/starlark"
)

var WaitGroup = star.Struct{
	TypeName:   "WaitGroup",
	Initialize: func() interface{} { return &sync.WaitGroup{} },
	Methods: map[string]star.Method{
		"Add": star.Method{
			Args: []starlark.Value{starlark.Int{}},
			Call: func(recv interface{}, args []interface{}) (resp []interface{}) {
				recv.(*sync.WaitGroup).Add(args[0].(int))
				return
			},
		},
		"Done": star.Method{
			Call: func(recv interface{}, args []interface{}) (resp []interface{}) {
				recv.(*sync.WaitGroup).Done()
				return
			},
		},
		"Wait": star.Method{
			Call: func(recv interface{}, args []interface{}) (resp []interface{}) {
				recv.(*sync.WaitGroup).Wait()
				return
			},
		},
	},
}
