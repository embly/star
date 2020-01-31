"""
hello it is
"""


star = require("star")
print(star.channel)


def hello(w, req):
    w.write(req.content_type)
    w.write(req.path)

    print(str(req.Body))
    print(str(1 * 4))
    print(req.some_bytes)
    print(len(req.some_bytes))
    for b in req.some_bytes:
        print(b)
    ords = str(list("hi".elem_ords()))
    w.write("Hello World\n" + str([x for x in range(10)]) + "\n" + ords)

    return
