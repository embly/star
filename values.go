package star

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

type ValueSetter interface {
	SetValue(v interface{})
	GetValue() (v interface{})
}

func makeOut(returns []starlark.Value, resp []interface{}) []starlark.Value {
	out := []starlark.Value{}
	for i, kind := range returns {
		switch t := kind.(type) {
		case Struct:
			t.Value = resp[i]
			out = append(out, t)
		case Interface:
			t.Value = resp[i]
			out = append(out, t)
		case Error:
			var err error
			if er, ok := resp[i].(error); ok {
				err = er
			}
			out = append(out, Error{err: err})
		case ByteArray:
			b, _ := resp[i].([]byte)
			out = append(out, ByteArray{b: b})
		case starlark.Int:
			out = append(out, starlark.MakeInt(resp[i].(int)))
		default:
			fmt.Println(kind.Type(), resp[i])
			panic("no match")
		}
	}
	if len(out) == 0 || out[len(out)-1].Type() != "error" {
		out = append(out, Error{})
	}
	return out
}

type Struct struct {
	Name       string
	Value      interface{}
	Methods    map[string]Method
	Attributes map[string]starlark.Value
}

func (s Struct) String() string        { return s.Name }
func (s Struct) Type() string          { return s.Name }
func (s Struct) Freeze()               {}
func (s Struct) Truth() starlark.Bool  { return starlark.True }
func (s Struct) Hash() (uint32, error) { return 0, errors.New("not hashable") }

func (s Struct) AttrNames() []string {
	out := []string{}
	for name := range s.Methods {
		out = append(out, name)
	}
	for name := range s.Attributes {
		out = append(out, name)
	}
	return out
}

func (s Struct) Attr(name string) (v starlark.Value, err error) {
	method, ok := s.Methods[name]
	if ok {
		v = starlark.NewBuiltin(name, func(thread *starlark.Thread,
			fn *starlark.Builtin, args starlark.Tuple,
			kwargs []starlark.Tuple) (v starlark.Value, err error) {
			values, err := ValidateArgs(method.Args, args)
			if err != nil {
				out := makeOut(method.Returns, make([]interface{}, len(method.Returns)))
				out[len(out)-1] = Error{err: err}
				return
			}
			resp := method.Call(s.Value, values)
			return starlark.Tuple(makeOut(method.Returns, resp)), nil
		})
	}

	attrType, ok := s.Attributes[name]
	if ok {
		s := reflect.ValueOf(s.Value).Elem()
		switch t := attrType.(type) {
		case Struct:
			t.Value = s.FieldByName(name).Interface()
			v = t
			return
		case Interface:
			t.Value = s.FieldByName(name).Interface()
			v = t
			return
		default:
			panic("NO ATTRIB")
		}
	}
	return
}

type Method struct {
	Args    []starlark.Value
	Returns []starlark.Value
	Call    MethodCall
}

type MethodCall func(interface{}, []interface{}) []interface{}

type Interface struct {
	Struct
}

type Function struct {
	FunctionName string
	Args         []starlark.Value
	Returns      []starlark.Value
	Call         func([]interface{}) []interface{}
}

func (f Function) String() string        { return f.FunctionName }
func (f Function) Name() string          { return f.FunctionName }
func (f Function) Type() string          { return f.FunctionName } // TODO: return the actual type signature?
func (f Function) Freeze()               {}
func (f Function) Truth() starlark.Bool  { return starlark.True }
func (f Function) Hash() (uint32, error) { return 0, errors.New("not hashable") }

func (f Function) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	values, err := ValidateArgs(f.Args, args)
	if err != nil {
		out := makeOut(f.Returns, make([]interface{}, len(f.Returns)))
		out[len(out)-1] = Error{err: err}
		return
	}
	resp := f.Call(values)
	return starlark.Tuple(makeOut(f.Returns, resp)), nil
}

func (f Function) Builtin() *starlark.Builtin {
	return starlark.NewBuiltin(f.Name(), func(thread *starlark.Thread,
		fn *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (v starlark.Value, err error) {
		return f.CallInternal(thread, args, kwargs)
	})
}

// ValidateArgs
func ValidateArgs(types []starlark.Value, args starlark.Tuple) (values []interface{}, err error) {
	// TODO: could take kwargs and check against type name
	var i int

	if len(types) != args.Len() {
		// TODO: better error message
		err = errors.Errorf(`not enough arguments`)
		return
	}
	if args.Len() == 0 {
		return
	}

	for {
		var val starlark.Value
		done := args.Iterate().Next(&val)

		in, ok := val.(Interface)
		in2, ok2 := types[i].(Interface)
		// check if interface satisfies another
		if ok && ok2 {
			for name := range in2.Methods {
				_, ok := in.Methods[name]
				if !ok {
					err = errors.Errorf(
						`argument %d was passed interface %s which doesn't `+
							`implement "%s". missing method "%s"`,
						i+1,
						val.Type(),
						types[i],
						name,
					)
					return
				}
			}
		} else if types[i].Type() != val.Type() {
			err = errors.Errorf(
				`argument %d is the wrong type, must `+
					`be "%s" but got "%s"`,
				i+1, types[i], val.Type())
			return
		}

		values = append(values, underlyingValue(val))
		if done {
			return
		}
	}
}

func underlyingValue(val starlark.Value) interface{} {
	switch v := val.(type) {
	case Struct:
		return v.Value
	case Interface:
		return v.Value
	case starlark.String:
		return string(v)
	default:
		panic(v)
	}
}
