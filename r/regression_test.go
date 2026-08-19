package r_test

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

func TestEqualDayRewardsUsePointTieBreaker(t *testing.T) {
	s := newRegressionService()
	s.CreateReward(model.Reward{Name: "small", Points: 10, RequiredDays: 3})
	s.CreateReward(model.Reward{Name: "large", Points: 30, RequiredDays: 3})
	rewards, err := s.ListRewards()
	if err != nil {
		t.Fatalf("ListRewards: %v", err)
	}
	if len(rewards) != 2 || rewards[0].Name != "large" || rewards[1].Name != "small" {
		t.Fatalf("rewards = %+v, want same-day rewards by points desc", rewards)
	}
	if !strings.Contains(rewards[0].Name, "large") {
		t.Fatalf("top reward name = %q", rewards[0].Name)
	}
}
