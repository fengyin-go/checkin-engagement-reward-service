# BUG_REPRO

## Bug 是什么
刚好达到奖励门槛时不可领取

## 如何触发
运行定向测试，用户连续签到两天并配置 required_days=2 的奖励。

## 错误信息
claimable = [], want the two-day reward
