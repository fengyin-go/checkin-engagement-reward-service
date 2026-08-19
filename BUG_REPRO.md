# BUG_REPRO

## Bug 是什么
补签昨天不能接上当前连续天数

## 如何触发
运行定向测试，今天签到后补昨天，再查询 summary。

## 错误信息
streak = 1, want 2 after makeup bridge
