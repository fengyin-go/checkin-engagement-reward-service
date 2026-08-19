package model

import "strings"

// BadgeType 成就徽章的判定维度。
type BadgeType string

const (
	BadgeTotalDays   BadgeType = "total_days"   // 累计签到天数达标
	BadgeStreak      BadgeType = "streak"       // 连续签到天数达标
	BadgeTotalChecks BadgeType = "total_checks" // 累计签到次数达标
)

// validBadgeTypes 允许的徽章类型集合。
var validBadgeTypes = map[BadgeType]bool{
	BadgeTotalDays:   true,
	BadgeStreak:      true,
	BadgeTotalChecks: true,
}

// Badge 成就徽章规则，用户签到数据达到阈值即可获得。
type Badge struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        BadgeType `json:"type"`      // total_days / streak / total_checks
	Threshold   int       `json:"threshold"` // 达到该数值即获得徽章
	Description string    `json:"description"`
}

func (b *Badge) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	b.Description = strings.TrimSpace(b.Description)
	if b.Name == "" {
		return NewValidationError("name", "徽章名称不能为空")
	}
	if !validBadgeTypes[b.Type] {
		return NewValidationError("type", "徽章类型必须为 total_days / streak / total_checks 之一")
	}
	if b.Threshold < 1 {
		return NewValidationError("threshold", "徽章阈值必须大于 0")
	}
	return nil
}

// TypeLabel 返回徽章类型的中文描述。
func (b *Badge) TypeLabel() string {
	switch b.Type {
	case BadgeTotalDays:
		return "累计签到天数"
	case BadgeStreak:
		return "连续签到天数"
	case BadgeTotalChecks:
		return "累计签到次数"
	default:
		return string(b.Type)
	}
}
