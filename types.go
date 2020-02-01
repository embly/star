package star

import (
	"fmt"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

type ByteArray struct {
	b []byte
}

func (ba ByteArray) String() string             { return fmt.Sprint(ba.b) }
func (ba ByteArray) Type() string               { return "ByteArray" }
func (ba ByteArray) Freeze()                    {}
func (ba ByteArray) Len() int                   { return len(ba.b) }
func (ba ByteArray) Truth() starlark.Bool       { return starlark.True }
func (ba ByteArray) Iterate() starlark.Iterator { return &byteIterator{ByteArray: &ba} }
func (ba ByteArray) Hash() (uint32, error)      { return 0, errors.New("not hashable") }

type byteIterator struct {
	*ByteArray
	i int
}

func (r *byteIterator) Next(p *starlark.Value) bool {
	i := starlark.MakeInt(int(r.ByteArray.b[r.i]))
	*p = &i
	r.i++
	return !(r.i > len(r.ByteArray.b)-1)
}

func (r *byteIterator) Done() {}

type Error struct {
	err error
}

func (err Error) Type() string          { return "error" }
func (err Error) Freeze()               {}
func (err Error) Truth() starlark.Bool  { return starlark.Bool(err.err != nil) }
func (err Error) Hash() (uint32, error) { return 0, errors.New("not hashable") }
func (err Error) AttrNames() []string    { return []string{"stacktrace"} }

func (err Error) String() string {
	if err.err != nil {
		return err.err.Error()
	} else {
		return ""
	}
}

func (r *Error) Attr(name string) (v starlark.Value, err error) {
	if name != "stacktrace" {
		return
	}
	v = starlark.NewBuiltin("stacktrace", func(thread *starlark.Thread,
		fn *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (v starlark.Value, err error) {
		v = starlark.String(fmt.Sprint(thread.CallStack()))
		return
	})
	return
}
