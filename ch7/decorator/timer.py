import time


def timeit(func):
    def wrapper(*args, **kwargs):
        start = time.time()
        result = func(*args, **kwargs)
        end = time.time()
        print(f"{func.__name__} took {end - start:.4f}s")
        return result

    return wrapper


@timeit
def foo():
    time.sleep(1)
    print("done")


foo()
