http = require("net/http")


def hello(w, req):
    w.WriteHeader(201)
    w.Write("hello world\n")


http.HandleFunc("/hello", hello)

http.ListenAndServe(":8080", http.Handler)
