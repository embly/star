http = require("net/http")
ioutil = require("io/ioutil")


def main():
    resp, err = http.Get("http://www.google.com")
    if err:
        return print(err)
    print("hi")
    b, err = ioutil.ReadAll(resp.Body)
    print("hi")
    if err:
        return print(err)

    print(b)


# def health(w, req):
#     print(w, req)


# http.HandleFunc("/health", health)
