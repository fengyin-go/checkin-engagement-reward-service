package s_test

import (
	"strings"
	"testing"
	"time"

	"checkin/internal/config"
	"checkin/internal/model"
	"checkin/internal/service"
	"checkin/internal/store"
	"checkin/pkg/logger"
)

func newRegressionService() *service.Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return service.New(store.NewMemoryStore(), log, cfg)
}

func newRegressionFixture() (*service.Service, *store.MemoryStore) {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	st := store.NewMemoryStore()
	return service.New(st, log, cfg), st
}

var _ = strings.Contains
var _ = time.Now

func TestMakeupYesterdayExtendsCurrentStreak(t *testing.T) {
	s, st := newRegressionFixture()
	u, _ := s.CreateUser("bridge-streak")
	now := time.Now()
	st.CreateCheckin(&model.Checkin{ID: "today", UserID: u.ID, Date: model.DateOf(now), Source: model.SourceNormal})
	if _, err := s.GrantMakeupCard(u.ID, 1); err != nil {
		t.Fatalf("grant card: %v", err)
	}
	if _, err := s.MakeupCheckIn(u.ID, model.DateOf(now.AddDate(0, 0, -1))); err != nil {
		t.Fatalf("makeup yesterday: %v", err)
	}
	summary, err := s.GetSummary(u.ID)
	if err != nil {
		t.Fatalf("GetSummary: %v", err)
	}
	if summary.Streak != 2 {
		t.Fatalf("streak = %d, want 2 after makeup bridge", summary.Streak)
	}
}
