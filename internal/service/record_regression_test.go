package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestRecord004BadgeProgressPutsAchievedBadgesFirst(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("badge-count")
	s.store.CreateCheckin(&model.Checkin{ID: "normal", UserID: u.ID, Date: model.DateOf(time.Now()), Source: model.SourceNormal})
	s.CreateBadge(model.Badge{Name: "done", Type: model.BadgeTotalDays, Threshold: 1})
	s.CreateBadge(model.Badge{Name: "later", Type: model.BadgeTotalDays, Threshold: 5})
	progress, err := s.UserBadges(u.ID)
	if err != nil {
		t.Fatalf("UserBadges: %v", err)
	}
	if len(progress.Badges) < 2 || !progress.Badges[0].Achieved || progress.Badges[0].Badge.Name != "done" {
		t.Fatalf("first progress = %+v, want achieved badge first", progress.Badges)
	}
}
