# BUG_REPRO

## Bug 是什么
重复补签没有被拦截

## 如何触发
运行定向测试，先普通签到昨天，再用补签卡补同一天。

## 错误信息
expected duplicate day validation, got <nil>
