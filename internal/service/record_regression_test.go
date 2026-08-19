package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestRecord001MakeupRejectsAlreadySignedDate(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("makeup-repeat")
	day := model.DateOf(time.Now().AddDate(0, 0, -1))
	if err := s.store.CreateCheckin(&model.Checkin{ID: "normal-day", UserID: u.ID, Date: day, Source: model.SourceNormal}); err != nil {
		t.Fatalf("seed checkin: %v", err)
	}
	if _, err := s.GrantMakeupCard(u.ID, 1); err != nil {
		t.Fatalf("grant card: %v", err)
	}
	if _, err := s.MakeupCheckIn(u.ID, day); !model.IsValidationError(err) {
		t.Fatalf("expected duplicate day validation, got %v", err)
	}
	items, _, err := s.ListCheckins(u.ID, "", 1, 10)
	if err != nil {
		t.Fatalf("list checkins: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("checkins len = %d, want 1", len(items))
	}
}
