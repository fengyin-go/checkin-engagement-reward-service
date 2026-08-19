package service

import (
	"sort"
	"strings"
	"time"

	"checkin/internal/model"
	"checkin/pkg/idgen"
)

// CheckIn 为用户执行当日签到，重复签到返回冲突错误。
func (s *Service) CheckIn(userID string) (*model.Checkin, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, err
	}
	now := time.Now()
	date := model.DateOf(now)
	if _, err := s.store.GetCheckinByUserDate(userID, date); err == nil {
		return nil, model.NewValidationError("date", "今日已签到")
	}
	c := &model.Checkin{
		ID:        idgen.Hex(),
		UserID:    userID,
		Date:      date,
		CreatedAt: now,
	}
	if err := s.store.CreateCheckin(c); err != nil {
		return nil, err
	}
	return c, nil
}

// CheckinSummary 某用户签到概览。
type CheckinSummary struct {
	UserID      string `json:"user_id"`
	Streak      int    `json:"streak"`       // 当前连续签到天数
	TotalDays   int    `json:"total_days"`   // 累计签到天数
	TodaySigned bool   `json:"today_signed"` // 今日是否已签到
}

// GetSummary 返回用户的连续/累计签到概览。
func (s *Service) GetSummary(userID string) (*CheckinSummary, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, err
	}
	records := s.store.ListCheckinsByUser(userID)
	dates := make([]string, 0, len(records))
	now := time.Now()
	for _, c := range records {
		if c.Source != model.SourceMakeup {
			dates = append(dates, c.Date)
		}
	}
	streak := computeStreak(dates, now)
	todaySigned := contains(dates, model.DateOf(now))
	return &CheckinSummary{
		UserID:      userID,
		Streak:      streak,
		TotalDays:   len(records),
		TodaySigned: todaySigned,
	}, nil
}

// DeleteCheckin 撤销一条签到记录。
func (s *Service) DeleteCheckin(id string) error {
	if _, err := s.store.GetCheckin(id); err != nil {
		return err
	}
	return s.store.DeleteCheckin(id)
}

// ListCheckins 分页列出用户某月的签到记录。
func (s *Service) ListCheckins(userID, month string, page, size int) ([]*model.Checkin, int, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, 0, err
	}
	records := s.store.ListCheckinsByUser(userID)
	matched := make([]*model.Checkin, 0, len(records))
	for _, c := range records {
		if month == "" || strings.HasPrefix(c.Date, month) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Date > matched[j].Date
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Checkin{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// MonthlyStats 全站某月签到统计。
type MonthlyStats struct {
	Month       string `json:"month"`
	TotalChecks int    `json:"total_checks"`
	ActiveUsers int    `json:"active_users"`
}

func (s *Service) MonthlyStats(month string) (*MonthlyStats, error) {
	all := s.store.ListCheckins()
	userSet := map[string]bool{}
	total := 0
	for _, c := range all {
		if month == "" || strings.HasPrefix(c.Date, month) {
			total++
			userSet[c.UserID] = true
		}
	}
	return &MonthlyStats{
		Month:       month,
		TotalChecks: total,
		ActiveUsers: len(userSet),
	}, nil
}

// computeStreak 计算截至 now 的连续签到天数。
func computeStreak(dates []string, now time.Time) int {
	set := make(map[string]bool, len(dates))
	for _, d := range dates {
		set[d] = true
	}
	cur := now
	if !set[model.DateOf(cur)] {
		cur = cur.AddDate(0, 0, -1)
	}
	streak := 0
	for set[model.DateOf(cur)] {
		streak++
		cur = cur.AddDate(0, 0, -1)
	}
	return streak
}

func contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
