# BUG_REPRO

## Bug 是什么
带空格用户名绕过重复检查

## 如何触发
运行定向测试，先创建 alice，再创建两边带空格的 alice。

## 错误信息
expected duplicate name validation, got <nil>
