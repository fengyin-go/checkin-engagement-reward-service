# 签到打卡服务（checkin）

纯 Go 标准库实现的签到打卡后端，零第三方依赖。

## 运行

```bash
go run ./cmd/server
# 默认监听 :8080，可通过 PORT / ADDR / MAX_PAGE_SIZE / LOG_LEVEL 环境变量配置
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/users | 创建用户 |
| GET | /api/users?keyword=&page=&size= | 分页查询用户 |
| GET | /api/users/{id} | 查询用户 |
| PUT | /api/users/{id} | 更新用户名称 |
| DELETE | /api/users/{id} | 删除用户 |
| POST | /api/checkins | 签到（body: `{"user_id": "..."}`，重复签到返回 409） |
| DELETE | /api/checkins/{id} | 撤销一条签到记录 |
| GET | /api/users/{id}/summary | 签到概览（连续天数/累计天数/今日是否已签到） |
| GET | /api/users/{id}/checkins?month=YYYY-MM&page=&size= | 分页查询某用户签到记录 |
| GET | /api/stats/checkins?month=YYYY-MM | 全站月度统计（签到次数/活跃用户数） |
| POST | /api/rewards | 创建奖励规则 |
| GET | /api/rewards | 列出奖励规则 |
| GET | /api/rewards/{id} | 查询奖励规则 |
| PUT | /api/rewards/{id} | 更新奖励规则 |
| DELETE | /api/rewards/{id} | 删除奖励规则 |
| GET | /api/users/{id}/rewards | 查询当前连续天数可领取的奖励 |
| POST | /api/makeup-cards | 发放补签卡（body: `{"user_id":"...","amount":3}`） |
| GET | /api/users/{id}/makeup-card | 查询补签卡余额 |
| POST | /api/makeups | 使用补签卡补签（body: `{"user_id":"...","date":"YYYY-MM-DD"}`，须早于今天） |
| GET | /api/users/{id}/makeups | 查询某用户补签记录 |
| POST | /api/badges | 创建成就徽章规则 |
| GET | /api/badges | 列出徽章规则 |
| GET | /api/badges/{id} | 查询徽章规则 |
| PUT | /api/badges/{id} | 更新徽章规则 |
| DELETE | /api/badges/{id} | 删除徽章规则 |
| GET | /api/users/{id}/badges | 查询用户签到数据在全部徽章上的达成进度 |
| GET | /api/leaderboard/streak?limit= | 连续签到天数排行榜 |
| GET | /api/leaderboard/total?limit= | 累计签到天数排行榜 |
| GET | /api/leaderboard/monthly?month=&limit= | 月度签到次数排行榜 |

## 业务说明

- **签到**：同一用户一天只能签到一次，重复签到返回 409。
- **连续签到**：从今天（或昨天，若今天尚未签到）起往前连续计算，断签则重置。
- **奖励规则**：按「所需连续签到天数」判定是否可领取。
- **补签卡**：可为用户发放补签卡，消耗一张即可为过去某一天补签；补签日期须早于今天且该日未签到。
- **成就徽章**：按「累计天数 / 连续天数 / 累计次数」三种维度设定阈值，用户数据达标即获得。
- **排行榜**：提供连续、累计、月度三种全站排行。

## 分层

- `cmd/server` 入口；`internal/app` 依赖装配；`internal/config` 环境变量配置。
- `internal/model` 领域模型与校验；`internal/store` 内存存储；`internal/service` 业务逻辑；`internal/handler` HTTP 路由。
- `pkg/httpx`、`pkg/idgen`、`pkg/logger` 通用工具。
