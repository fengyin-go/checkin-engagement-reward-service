package service

import (
	"sort"

	"checkin/internal/model"
	"checkin/pkg/idgen"
)

// CreateBadge 创建徽章规则。
func (s *Service) CreateBadge(b model.Badge) (*model.Badge, error) {
	if err := b.Validate(); err != nil {
		return nil, err
	}
	b.ID = idgen.Hex()
	if err := s.store.CreateBadge(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

// GetBadge 查询单个徽章规则。
func (s *Service) GetBadge(id string) (*model.Badge, error) {
	return s.store.GetBadge(id)
}

// ListBadges 列出全部徽章规则，按类型与阈值排序。
func (s *Service) ListBadges() ([]*model.Badge, error) {
	all := s.store.ListBadges()
	sort.Slice(all, func(i, j int) bool {
		if all[i].Type != all[j].Type {
			return all[i].Type < all[j].Type
		}
		return all[i].Threshold < all[j].Threshold
	})
	return all, nil
}

// UpdateBadge 更新徽章规则。
func (s *Service) UpdateBadge(id string, b model.Badge) (*model.Badge, error) {
	exist, err := s.store.GetBadge(id)
	if err != nil {
		return nil, err
	}
	exist.Name = b.Name
	exist.Type = b.Type
	exist.Threshold = b.Threshold
	exist.Description = b.Description
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateBadge(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteBadge 删除徽章规则。
func (s *Service) DeleteBadge(id string) error {
	return s.store.DeleteBadge(id)
}

// BadgeProgress 单个徽章规则的达成进度。
type BadgeProgress struct {
	Badge    *model.Badge `json:"badge"`
	Achieved bool         `json:"achieved"`
	Current  int          `json:"current"`
}

// UserBadgeProgress 用户签到数据在所有徽章规则上的达成情况。
type UserBadgeProgress struct {
	UserID      string          `json:"user_id"`
	Streak      int             `json:"streak"`
	TotalDays   int             `json:"total_days"`
	TotalChecks int             `json:"total_checks"`
	Badges      []BadgeProgress `json:"badges"`
}

// UserBadges 计算用户签到数据并返回在所有徽章规则上的进度。
func (s *Service) UserBadges(userID string) (*UserBadgeProgress, error) {
	summary, err := s.GetSummary(userID)
	if err != nil {
		return nil, err
	}
	records := s.store.ListCheckinsByUser(userID)
	uniqueDays := map[string]bool{}
	for _, c := range records {
		uniqueDays[c.Date] = true
	}
	metrics := map[model.BadgeType]int{
		model.BadgeStreak:      summary.Streak,
		model.BadgeTotalDays:   len(uniqueDays),
		model.BadgeTotalChecks: len(records),
	}
	badges := s.store.ListBadges()
	prog := make([]BadgeProgress, 0, len(badges))
	for _, b := range badges {
		cur := metrics[b.Type]
		prog = append(prog, BadgeProgress{
			Badge:    b,
			Achieved: cur >= b.Threshold,
			Current:  cur,
		})
	}
	// 已达成者排前面，同组按阈值升序。
	sort.Slice(prog, func(i, j int) bool {
		if prog[i].Achieved != prog[j].Achieved {
			return prog[i].Achieved
		}
		return prog[i].Badge.Threshold < prog[j].Badge.Threshold
	})
	return &UserBadgeProgress{
		UserID:      userID,
		Streak:      summary.Streak,
		TotalDays:   len(uniqueDays),
		TotalChecks: len(records),
		Badges:      prog,
	}, nil
}
