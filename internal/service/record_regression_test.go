package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestRecord005GrantMakeupCardKeepsExistingBalance(t *testing.T) {
	s := newTestService()
	u, _ := s.CreateUser("makeup-balance")
	if _, err := s.GrantMakeupCard(u.ID, 2); err != nil {
		t.Fatalf("first grant: %v", err)
	}
	card, err := s.GrantMakeupCard(u.ID, 3)
	if err != nil {
		t.Fatalf("second grant: %v", err)
	}
	if card.Balance != 5 {
		t.Fatalf("balance = %d, want 5", card.Balance)
	}
}
