def make_cfg(base_url="http://localhost", timeout=5, **kwargs):
    return {"base_url": base_url, "timeout": timeout, **kwargs}


c1 = make_cfg()  # 使用默认值
c2 = make_cfg("https://jianghushinian.cn/")  # 仅修改 URL
c3 = make_cfg(timeout=10, retry=3)  # 修改超时时间和添加重试次数

print(c1)
print(c2)
print(c3)
