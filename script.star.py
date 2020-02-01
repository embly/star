"""
hello it is
"""

require("ioutil")

# star = require("star")
# print(star.channel)


def handle_err(args):
    for a in args:
        if type(a) == "error" and a:
            print(a, a.stacktrace())
    return args


def hello(w, req):
    w.write(req.content_type)
    w.write(req.path)

    b, err = handle_err(req.sample_resp)
    print(b, err)

    print(req.Body)

    print(req.some_bytes)
    for b in req.some_bytes:
        print(b)

    # print(str(req.Body))
    # print(str(1 * 4))
    # print(req.error)
    # handle_err(req.error)
    # print(len(req.some_bytes))
    # ords = str(list("hi".elem_ords()))
    # w.write("Hello World\n" + str([x for x in range(10)]) + "\n" + ords)

    return
