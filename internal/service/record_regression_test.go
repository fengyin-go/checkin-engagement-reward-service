package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestRecord006UpdateRewardInvalidInputDoesNotCorruptStoredRule(t *testing.T) {
	s := newTestService()
	r, _ := s.CreateReward(model.Reward{Name: "weekly", Points: 70, RequiredDays: 7})
	if _, err := s.UpdateReward(r.ID, model.Reward{Name: "", Points: 1, RequiredDays: 1}); !model.IsValidationError(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
	got, err := s.GetReward(r.ID)
	if err != nil {
		t.Fatalf("GetReward: %v", err)
	}
	if got.Name != "weekly" || got.Points != 70 || got.RequiredDays != 7 {
		t.Fatalf("stored reward changed to %+v", got)
	}
}
