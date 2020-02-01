package src

import (
	"github.com/embly/star/src/io"
	"github.com/embly/star/src/io/ioutil"
	"github.com/embly/star/src/net/http"
	"go.starlark.net/starlark"
)

var Packages = map[string]map[string]starlark.Value{
	"io": map[string]starlark.Value{
		"Reader":     io.Reader,
		"ReadCloser": io.ReadCloser,
	},
	"io/ioutil": map[string]starlark.Value{
		"NopCloser": ioutil.NopCloser,
		"ReadAll":   ioutil.ReadAll,
	},
	"net/http": map[string]starlark.Value{
		"Get":      http.Get,
		"Response": http.Response,
	},
}
