package time

import (
	"time"

	"github.com/embly/star"
	"go.starlark.net/starlark"
)

var Duration = starlark.Int{}
var Nanosecond = starlark.MakeInt64(1)
var Microsecond = starlark.MakeInt64(1000).Mul(Nanosecond)
var Millisecond = starlark.MakeInt64(1000).Mul(Microsecond)
var Second = starlark.MakeInt64(1000).Mul(Millisecond)
var Minute = starlark.MakeInt64(60).Mul(Second)
var Hour = starlark.MakeInt64(60).Mul(Minute)

var Sleep = star.Function{
	FunctionName: "Sleep",
	Args:         []starlark.Value{starlark.Int{}},
	Call: func(args []interface{}) (resp []interface{}) {
		time.Sleep(time.Duration(args[0].(int)))
		return
	},
}

var Now = star.Function{
	FunctionName: "Now",
	Returns:      []starlark.Value{Time},
	Call: func(args []interface{}) (resp []interface{}) {
		return []interface{}{time.Now()}
	},
}

var Time = star.Struct{
	TypeName: "Time",
}
