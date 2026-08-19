package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestRecord002ListCheckinsUsesExactMonthAndNewestFirst(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("month-list")
	for _, d := range []string{"2026-08-01", "2026-08-30", "2026-09-08", "2026-07-08"} {
		s.store.CreateCheckin(&model.Checkin{ID: d, UserID: u.ID, Date: d, Source: model.SourceNormal})
	}
	items, total, err := s.ListCheckins(u.ID, "2026-08", 1, 10)
	if err != nil {
		t.Fatalf("ListCheckins: %v", err)
	}
	if total != 2 || len(items) != 2 {
		t.Fatalf("august total/items = %d/%d, want 2/2", total, len(items))
	}
	if items[0].Date != "2026-08-30" || items[1].Date != "2026-08-01" {
		t.Fatalf("dates = %v, want newest August dates only", []string{items[0].Date, items[1].Date})
	}
}
