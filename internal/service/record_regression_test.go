package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestRecord007DeletedUserDataIsNotReportedInMonthlyStats(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("deleted-user")
	s.store.CreateCheckin(&model.Checkin{ID: "leftover", UserID: u.ID, Date: "2026-08-09", Source: model.SourceNormal})
	if err := s.DeleteUser(u.ID); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	stats, err := s.MonthlyStats("2026-08")
	if err != nil {
		t.Fatalf("MonthlyStats: %v", err)
	}
	if stats.TotalChecks != 0 || stats.ActiveUsers != 0 {
		t.Fatalf("stats = %+v, want deleted user ignored", stats)
	}
}
