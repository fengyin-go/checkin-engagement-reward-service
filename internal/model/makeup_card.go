package model

import "time"

// MakeupCard 用户的补签卡，剩余额度以 Balance 表示。
// 用户可通过消耗补签卡为过去漏签的某一天补签。
type MakeupCard struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   int       `json:"balance"` // 剩余补签卡张数
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *MakeupCard) Validate() error {
	if m.UserID == "" {
		return NewValidationError("user_id", "用户 ID 不能为空")
	}
	if m.Balance < 0 {
		return NewValidationError("balance", "补签卡余额不能为负数")
	}
	return nil
}
