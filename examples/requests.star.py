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
        "https://www.embly.run/hello/",
        "https://www.embly.run/hello/",
        "https://www.embly.run/hello/",
    ]
    for url in urls:
        star.go(get_url, url, wg)
    wg.Wait()
