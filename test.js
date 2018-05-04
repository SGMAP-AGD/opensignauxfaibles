a = {
    "test": {
    "a": 1,
    "b": 2
    }
}

b = {
    "test": {
    "a": 3,
    "c": 4
    }
}
c = {"test": {}}
Object.assign(a.test, b.test)
print(JSON.stringify(a, null, 2))