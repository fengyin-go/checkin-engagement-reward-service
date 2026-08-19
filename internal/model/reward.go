package model

import "strings"

// Reward 签到奖励规则，连续签到达到指定天数即可领取。
type Reward struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Points       int    `json:"points"`
	RequiredDays int    `json:"required_days"` // 需要连续签到天数
	Description  string `json:"description"`
}

func (r *Reward) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)
	if r.Name == "" {
		return NewValidationError("name", "奖励名称不能为空")
	}
	if r.Points < 0 {
		return NewValidationError("points", "奖励积分不能为负数")
	}
	if r.RequiredDays < 1 {
		return NewValidationError("required_days", "所需连续天数必须大于 0")
	}
	return nil
}
