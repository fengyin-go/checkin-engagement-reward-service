package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestRecord003RewardAvailableWhenStreakEqualsRequiredDays(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("reward-edge")
	now := time.Now()
	s.store.CreateCheckin(&model.Checkin{ID: "yesterday", UserID: u.ID, Date: model.DateOf(now.AddDate(0, 0, -1)), Source: model.SourceNormal})
	s.store.CreateCheckin(&model.Checkin{ID: "today", UserID: u.ID, Date: model.DateOf(now), Source: model.SourceNormal})
	if _, err := s.CreateReward(model.Reward{Name: "two day", Points: 20, RequiredDays: 2}); err != nil {
		t.Fatalf("CreateReward: %v", err)
	}
	claimable, err := s.ClaimableRewards(u.ID)
	if err != nil {
		t.Fatalf("ClaimableRewards: %v", err)
	}
	if len(claimable) != 1 || claimable[0].RequiredDays != 2 {
		t.Fatalf("claimable = %+v, want the two-day reward", claimable)
	}
}
