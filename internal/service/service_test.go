package service

import (
	"strings"
	"testing"
	"time"

	"checkin/internal/config"
	"checkin/internal/model"
	"checkin/internal/store"
	"checkin/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestCreateAndGetUser(t *testing.T) {
	s := newTestService()

	u, err := s.CreateUser("alice")
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if u.ID == "" {
		t.Fatal("expected generated id")
	}

	got, err := s.GetUser(u.ID)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.Name != "alice" {
		t.Fatalf("name = %q, want alice", got.Name)
	}

	// 空名称校验
	if _, err := s.CreateUser("  "); !model.IsValidationError(err) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestUpdateAndDeleteUser(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("alice")

	updated, err := s.UpdateUser(u.ID, "bob")
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if updated.Name != "bob" {
		t.Fatalf("name = %q, want bob", updated.Name)
	}

	if _, err := s.UpdateUser("missing", "x"); err == nil {
		t.Fatal("expected error updating missing user")
	}

	if err := s.DeleteUser(u.ID); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if _, err := s.GetUser(u.ID); err == nil {
		t.Fatal("expected error getting deleted user")
	}
}

func TestListUsersFilterAndPagination(t *testing.T) {
	s := newTestService()
	s.CreateUser("alice")
	s.CreateUser("bob")
	s.CreateUser("charlie")

	items, total, err := s.ListUsers(model.UserFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if total != 3 {
		t.Fatalf("total = %d, want 3", total)
	}
	if len(items) != 2 {
		t.Fatalf("page size = %d, want 2", len(items))
	}

	items, total, _ = s.ListUsers(model.UserFilter{Keyword: "ali"}, 1, 10)
	if total != 1 || items[0].Name != "alice" {
		t.Fatalf("filtered total = %d, want 1", total)
	}

	// 越界页
	items, total, _ = s.ListUsers(model.UserFilter{}, 99, 10)
	if len(items) != 0 {
		t.Fatalf("expected empty page, got %d items", len(items))
	}
}

func TestCheckInConflictAndSummary(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("alice")

	if _, err := s.CheckIn(u.ID); err != nil {
		t.Fatalf("first CheckIn: %v", err)
	}
	if _, err := s.CheckIn(u.ID); !model.IsValidationError(err) {
		t.Fatalf("expected ValidationError on duplicate checkin, got %v", err)
	}

	summary, err := s.GetSummary(u.ID)
	if err != nil {
		t.Fatalf("GetSummary: %v", err)
	}
	if !summary.TodaySigned {
		t.Fatal("expected today signed")
	}
	if summary.TotalDays != 1 {
		t.Fatalf("total days = %d, want 1", summary.TotalDays)
	}

	// 不存在的用户签到
	if _, err := s.CheckIn("missing"); err == nil {
		t.Fatal("expected error checking in missing user")
	}
}

func TestComputeStreak(t *testing.T) {
	now := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name  string
		dates []string
		want  int
	}{
		{"无签到", nil, 0},
		{"仅今天", []string{"2026-08-16"}, 1},
		{"连续三天到今", []string{"2026-08-14", "2026-08-15", "2026-08-16"}, 3},
		{"断签重置", []string{"2026-08-14", "2026-08-16"}, 1},
		{"今天未签但昨天连续", []string{"2026-08-14", "2026-08-15"}, 2},
		{"昨天也没签则从最近回溯", []string{"2026-08-10", "2026-08-11", "2026-08-12"}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := computeStreak(tc.dates, now); got != tc.want {
				t.Fatalf("streak = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestListCheckinsByMonth(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("alice")

	for _, d := range []string{"2026-08-01", "2026-08-05", "2026-07-31"} {
		s.store.CreateCheckin(&model.Checkin{ID: d, UserID: u.ID, Date: d, CreatedAt: time.Now()})
	}

	items, total, err := s.ListCheckins(u.ID, "2026-08", 1, 10)
	if err != nil {
		t.Fatalf("ListCheckins: %v", err)
	}
	if total != 2 {
		t.Fatalf("aug total = %d, want 2", total)
	}
	if len(items) != 2 {
		t.Fatalf("aug items = %d, want 2", len(items))
	}

	// 全部月份
	_, total, _ = s.ListCheckins(u.ID, "", 1, 10)
	if total != 3 {
		t.Fatalf("all total = %d, want 3", total)
	}
}

func TestMonthlyStats(t *testing.T) {
	s := newTestService()
	u1, _ := s.CreateUser("alice")
	u2, _ := s.CreateUser("bob")

	s.store.CreateCheckin(&model.Checkin{ID: "a", UserID: u1.ID, Date: "2026-08-01"})
	s.store.CreateCheckin(&model.Checkin{ID: "b", UserID: u1.ID, Date: "2026-08-02"})
	s.store.CreateCheckin(&model.Checkin{ID: "c", UserID: u2.ID, Date: "2026-08-01"})
	s.store.CreateCheckin(&model.Checkin{ID: "d", UserID: u2.ID, Date: "2026-07-01"})

	stats, err := s.MonthlyStats("2026-08")
	if err != nil {
		t.Fatalf("MonthlyStats: %v", err)
	}
	if stats.TotalChecks != 3 {
		t.Fatalf("aug checks = %d, want 3", stats.TotalChecks)
	}
	if stats.ActiveUsers != 2 {
		t.Fatalf("aug active users = %d, want 2", stats.ActiveUsers)
	}
}

func TestRewardCRUDAndClaimable(t *testing.T) {
	s := newTestService()

	r, err := s.CreateReward(model.Reward{Name: "七日礼", Points: 100, RequiredDays: 7})
	if err != nil {
		t.Fatalf("CreateReward: %v", err)
	}
	if r.ID == "" {
		t.Fatal("expected generated reward id")
	}

	// 校验
	if _, err := s.CreateReward(model.Reward{Name: "", Points: 0, RequiredDays: 0}); !model.IsValidationError(err) {
		t.Fatalf("expected ValidationError, got %v", err)
	}

	updated, err := s.UpdateReward(r.ID, model.Reward{Name: "七日礼", Points: 200, RequiredDays: 7})
	if err != nil {
		t.Fatalf("UpdateReward: %v", err)
	}
	if updated.Points != 200 {
		t.Fatalf("points = %d, want 200", updated.Points)
	}

	rewards, err := s.ListRewards()
	if err != nil || len(rewards) != 1 {
		t.Fatalf("ListRewards: %v", err)
	}

	// 用户连续签到 3 天（到今天的前一天为止），可领取 required_days<=3 的奖励。
	// 使用相对日期构造，避免依赖运行当天日期。
	u, _ := s.CreateUser("alice")
	now := time.Now()
	s.store.CreateCheckin(&model.Checkin{ID: "1", UserID: u.ID, Date: model.DateOf(now.AddDate(0, 0, -3))})
	s.store.CreateCheckin(&model.Checkin{ID: "2", UserID: u.ID, Date: model.DateOf(now.AddDate(0, 0, -2))})
	s.store.CreateCheckin(&model.Checkin{ID: "3", UserID: u.ID, Date: model.DateOf(now.AddDate(0, 0, -1))})

	claimable, err := s.ClaimableRewards(u.ID)
	if err != nil {
		t.Fatalf("ClaimableRewards: %v", err)
	}
	// required_days=7 > 3，不可领取
	if len(claimable) != 0 {
		t.Fatalf("claimable len = %d, want 0", len(claimable))
	}

	// 添加一个 2 天即可领取的奖励
	s.CreateReward(model.Reward{Name: "两日礼", Points: 10, RequiredDays: 2})
	claimable, _ = s.ClaimableRewards(u.ID)
	if len(claimable) != 1 || !strings.Contains(claimable[0].Name, "两日礼") {
		t.Fatalf("claimable = %+v, want 两日礼 only", claimable)
	}

	if err := s.DeleteReward(r.ID); err != nil {
		t.Fatalf("DeleteReward: %v", err)
	}
}
