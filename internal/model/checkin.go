package model

import (
	"regexp"
	"time"
)

// dateLayout 签到日期的统一格式 YYYY-MM-DD。
const dateLayout = "2006-01-02"

var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// 签到来源。
const (
	SourceNormal = "normal" // 正常签到
	SourceMakeup = "makeup" // 补签
	SourceRetry  = "retry"  // 重试补写
)

// Checkin 一条签到记录，一个用户一天最多一条。
type Checkin struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Date      string    `json:"date"`   // YYYY-MM-DD
	Source    string    `json:"source"` // normal / makeup
	CreatedAt time.Time `json:"created_at"`
}

func (c *Checkin) Validate() error {
	if c.UserID == "" {
		return NewValidationError("user_id", "用户 ID 不能为空")
	}
	if !dateRe.MatchString(c.Date) {
		return NewValidationError("date", "日期格式必须为 YYYY-MM-DD")
	}
	return nil
}

// DateOf 返回 t 对应日期的字符串表示。
func DateOf(t time.Time) string {
	return t.Format(dateLayout)
}

// IsValidDate 判断字符串是否为合法的 YYYY-MM-DD 日期。
func IsValidDate(s string) bool {
	return dateRe.MatchString(s)
}
