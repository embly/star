package star

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

func TestRunScript(t *testing.T) {

	{
		err := RunScript("")
		if !strings.Contains(err.Error(), "no such file") {
			t.Fatal("file shouldn't exist")
		}
	}

	if err := RunScript("./print.star.py"); err != nil {
		t.Fatal(err)
	}

	{
		err := RunScript("./error.star.py")
		er := err.(*starlark.EvalError)
		if !strings.Contains(er.Msg, "invalid literal with base 10") {
			t.Fatal("shold error")
		}
	}

}
