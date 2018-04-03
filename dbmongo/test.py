def test(blu):
    if blu == "test":
        yield "ok"

a = ["test", "zinzin"]

b = [test(x) for x in a]

print(b)