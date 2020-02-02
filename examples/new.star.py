http = require("net/http")
ioutil = require("io/ioutil")
sync = require("sync")
star = require("star")
time = require("time")


def get_url(url, wg):
    resp, err = http.Get(url)
    if err:
        return print(err)
    print("hi")
    b, err = ioutil.ReadAll(resp.Body)

    if err:
        return print(err)

    body, err = star.bytes_to_string(b)
    if err:
        return print(err)

    print(body)
    time.Sleep(time.Second * 2)
    wg.Done()


def sleep():
    print("hello")


def main():
    wg = sync.WaitGroup()
    wg.Add(3)

    start = time.Now()
    print(start)
    urls = [
        "https://api.exchangeratesapi.io/latest",
        "https://api.exchangeratesapi.io/latest",
        "https://api.exchangeratesapi.io/latest",
    ]
    for url in urls:
        star.go(get_url, url, wg)

    wg.Wait()
