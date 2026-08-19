package service

import (
	"sort"
	"strings"
	"time"
)

// LeaderboardEntry 排行榜条目。
type LeaderboardEntry struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Value  int    `json:"value"`
	Rank   int    `json:"rank"`
}

// StreakLeaderboard 全站连续签到天数排行榜。
func (s *Service) StreakLeaderboard(limit int) ([]LeaderboardEntry, error) {
	users := s.store.ListUsers()
	now := time.Now()
	entries := make([]LeaderboardEntry, 0, len(users))
	for _, u := range users {
		records := s.store.ListCheckinsByUser(u.ID)
		entries = append(entries, LeaderboardEntry{
			UserID: u.ID,
			Name:   u.Name,
			Value:  computeStreak(recordCheckinDates(records), now),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})
	return trimAndRank(entries, limit), nil
}

// TotalDaysLeaderboard 全站累计签到天数排行榜。
func (s *Service) TotalDaysLeaderboard(limit int) ([]LeaderboardEntry, error) {
	users := s.store.ListUsers()
	entries := make([]LeaderboardEntry, 0, len(users))
	for _, u := range users {
		records := s.store.ListCheckinsByUser(u.ID)
		unique := map[string]bool{}
		for _, c := range records {
			unique[c.Date] = true
		}
		entries = append(entries, LeaderboardEntry{
			UserID: u.ID,
			Name:   u.Name,
			Value:  len(unique),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})
	return trimAndRank(entries, limit), nil
}

// MonthlyLeaderboard 全站某月签到次数排行榜。
func (s *Service) MonthlyLeaderboard(month string, limit int) ([]LeaderboardEntry, error) {
	users := s.store.ListUsers()
	entries := make([]LeaderboardEntry, 0, len(users))
	for _, u := range users {
		records := s.store.ListCheckinsByUser(u.ID)
		cnt := 0
		for _, c := range records {
			if month == "" || strings.HasPrefix(c.Date, month) {
				cnt++
			}
		}
		entries = append(entries, LeaderboardEntry{
			UserID: u.ID,
			Name:   u.Name,
			Value:  cnt,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})
	return trimAndRank(entries, limit), nil
}

// trimAndRank 截断到 limit 条并填充名次。
func trimAndRank(entries []LeaderboardEntry, limit int) []LeaderboardEntry {
	if limit <= 0 {
		limit = 10
	}
	if len(entries) > limit {
		entries = entries[:limit]
	}
	for i := range entries {
		entries[i].Rank = i + 1
	}
	return entries
}
