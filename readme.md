# Star

Star uses [starlark](https://github.com/google/starlark-go) to provide a python-like environment to run Go code. This is a quick and dirty experiment and should not be taken seriously.

The [intro blog post](https://embly.run/star) is likely the best place to get started. Or you can play around in the repl: https://embly.github.com/star

## How to use


Use Go's concurrency model:
```python
http = require("net/http")
ioutil = require("io/ioutil")
sync = require("sync")
star = require("star")
time = require("time")

def get_url(url, wg):
    resp, err = http.Get(url)
    if err:
        return print(err)
    b, err = ioutil.ReadAll(resp.Body)
    if err:
        return print(err)
    body, err = star.bytes_to_string(b)
    if err:
        return print(err)
    time.Sleep(time.Second * 2)
    wg.Done()


def main():
    wg = sync.WaitGroup()
    wg.Add(3)
    urls = [
        "https://api.exchangeratesapi.io/latest",
        "https://api.exchangeratesapi.io/latest",
        "https://api.exchangeratesapi.io/latest",
    ]
    for url in urls:
        star.go(get_url, url, wg)
    wg.Wait()
```


Run a web server:

```python
http = require("net/http")


def hello(w, req):
    w.WriteHeader(201)
    w.Write("hello world\n")


http.HandleFunc("/hello", hello)

http.ListenAndServe(":8080", http.Hand
```
