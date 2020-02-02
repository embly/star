package src

import (
	"github.com/embly/star/src/io"
	"github.com/embly/star/src/io/ioutil"
	"github.com/embly/star/src/net/http"
	"github.com/embly/star/src/sync"
	"github.com/embly/star/src/time"
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
	"sync": map[string]starlark.Value{
		"WaitGroup": sync.WaitGroup,
	},
	"time": map[string]starlark.Value{
		"Sleep":       time.Sleep,
		"Duration":    time.Duration,
		"Nanosecond":  time.Nanosecond,
		"Microsecond": time.Microsecond,
		"Millisecond": time.Millisecond,
		"Second":      time.Second,
		"Minute":      time.Minute,
		"Hour":        time.Hour,
		"Now":         time.Now,
		"Time":        time.Time,
	},
}
