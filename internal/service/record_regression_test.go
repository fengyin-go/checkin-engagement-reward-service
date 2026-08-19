package service

import (
	"testing"
	"time"

	"checkin/internal/model"
)

var _ = time.Now
var _ = model.SourceNormal

func TestTrimmedUserNameCannotBeCreatedTwice(t *testing.T) {
	s := newTestService()
	if _, err := s.CreateUser("alice"); err != nil {
		t.Fatalf("first CreateUser: %v", err)
	}
	if _, err := s.CreateUser("  alice  "); !model.IsValidationError(err) {
		t.Fatalf("expected duplicate name validation, got %v", err)
	}
	users, total, err := s.ListUsers(model.UserFilter{Keyword: "alice"}, 1, 10)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	if total != 1 || len(users) != 1 {
		t.Fatalf("users = %d/%d, want one alice", len(users), total)
	}
}
