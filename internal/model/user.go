package model

import (
	"strings"
	"time"
)

// User 参与签到的用户。
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	u.Name = strings.TrimSpace(u.Name)
	if u.Name == "" {
		return NewValidationError("name", "用户名称不能为空")
	}
	return nil
}

type UserFilter struct {
	Keyword string
}

func (f UserFilter) Match(u *User) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(u.Name), k) {
			return false
		}
	}
	return true
}
