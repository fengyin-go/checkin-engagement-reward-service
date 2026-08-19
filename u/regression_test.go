package u_test

import (
	"strings"
	"testing"
	"time"

	"checkin/internal/config"
	"checkin/internal/model"
	"checkin/internal/service"
	"checkin/internal/store"
	"checkin/pkg/logger"
)

func newRegressionService() *service.Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return service.New(store.NewMemoryStore(), log, cfg)
}

func newRegressionFixture() (*service.Service, *store.MemoryStore) {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	st := store.NewMemoryStore()
	return service.New(st, log, cfg), st
}

var _ = strings.Contains
var _ = time.Now

func TestTrimmedUserNameCannotBeCreatedTwice(t *testing.T) {
	s := newRegressionService()
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
