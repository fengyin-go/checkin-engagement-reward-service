package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestEqualDayRewardsUsePointTieBreaker(t *testing.T) {
	s := newTestService()
	s.CreateReward(model.Reward{Name: "small", Points: 10, RequiredDays: 3})
	s.CreateReward(model.Reward{Name: "large", Points: 30, RequiredDays: 3})
	rewards, err := s.ListRewards()
	if err != nil {
		t.Fatalf("ListRewards: %v", err)
	}
	if len(rewards) != 2 || rewards[0].Name != "large" || rewards[1].Name != "small" {
		t.Fatalf("rewards = %+v, want same-day rewards by points desc", rewards)
	}
}
