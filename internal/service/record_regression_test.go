package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestMakeupYesterdayExtendsCurrentStreak(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("bridge-streak")
	now := time.Now()
	s.store.CreateCheckin(&model.Checkin{ID: "today", UserID: u.ID, Date: model.DateOf(now), Source: model.SourceNormal})
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
