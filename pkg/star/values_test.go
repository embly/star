package star

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

func TestValuesTest(t *testing.T) {
	a := assert.New(t)

	// incorrect types
	{
		types := []starlark.Value{
			starlark.Int{},
		}
		args := starlark.Tuple{starlark.String("Hello World")}
		_, err := ValidateArgs(types, args)
		if !strings.Contains(err.Error(), "wrong type") {
			t.Fatal("should be incorrect type")
		}
	}

	// incorrect argument lengths
	{
		types := []starlark.Value{}
		args := starlark.Tuple{starlark.String("Hello World")}
		_, err := ValidateArgs(types, args)
		if !strings.Contains(err.Error(), "not enough arguments") {
			t.Fatal("should not enough arguments")
		}
	}

	// basic boolean support
	{
		types := []starlark.Value{
			starlark.False,
		}
		args := starlark.Tuple{starlark.False}
		values, err := ValidateArgs(types, args)
		panicOnErr(err)
		a.Equal(values, []interface{}{false})
	}

	// basic int support
	{
		i := 213123421
		types := []starlark.Value{
			starlark.Int{},
		}
		args := starlark.Tuple{starlark.MakeInt(i)}
		values, err := ValidateArgs(types, args)
		panicOnErr(err)
		a.Equal(values, []interface{}{i})
	}

	// big int support
	{
		types := []starlark.Value{
			starlark.Int{},
		}
		args := starlark.Tuple{starlark.MakeInt(MaxInt)}
		values, err := ValidateArgs(types, args)
		panicOnErr(err)
		a.Equal(values, []interface{}{MaxInt})
	}

	// basic string support
	{
		s := "Hello World"
		types := []starlark.Value{
			starlark.String(""),
		}
		args := starlark.Tuple{starlark.String(s)}
		values, err := ValidateArgs(types, args)
		panicOnErr(err)
		a.Equal(values, []interface{}{s})
	}

}
