package store

import (
	"errors"
	"testing"

	"checkin/internal/model"
)

func TestUserStoreCRUD(t *testing.T) {
	s := NewMemoryStore()

	u := &model.User{ID: "u1", Name: "alice"}
	if err := s.CreateUser(u); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	got, err := s.GetUser("u1")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if got.Name != "alice" {
		t.Fatalf("name = %q, want %q", got.Name, "alice")
	}

	if _, err := s.GetUser("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if list := s.ListUsers(); len(list) != 1 {
		t.Fatalf("ListUsers len = %d, want 1", len(list))
	}

	u.Name = "bob"
	if err := s.UpdateUser(u); err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	got, _ = s.GetUser("u1")
	if got.Name != "bob" {
		t.Fatalf("name after update = %q, want bob", got.Name)
	}

	if err := s.UpdateUser(&model.User{ID: "nope"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound on update missing, got %v", err)
	}

	if err := s.DeleteUser("u1"); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if err := s.DeleteUser("u1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound on delete missing, got %v", err)
	}
}

func TestCheckinStoreCRUDAndUniqueness(t *testing.T) {
	s := NewMemoryStore()

	c1 := &model.Checkin{ID: "c1", UserID: "u1", Date: "2026-08-16"}
	if err := s.CreateCheckin(c1); err != nil {
		t.Fatalf("CreateCheckin: %v", err)
	}

	// 同用户同日期冲突
	c2 := &model.Checkin{ID: "c2", UserID: "u1", Date: "2026-08-16"}
	if err := s.CreateCheckin(c2); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict for duplicate checkin, got %v", err)
	}

	// 不同用户同日期不冲突
	c3 := &model.Checkin{ID: "c3", UserID: "u2", Date: "2026-08-16"}
	if err := s.CreateCheckin(c3); err != nil {
		t.Fatalf("CreateCheckin c3: %v", err)
	}

	if _, err := s.GetCheckinByUserDate("u1", "2026-08-16"); err != nil {
		t.Fatalf("GetCheckinByUserDate: %v", err)
	}
	if _, err := s.GetCheckinByUserDate("u1", "2026-08-15"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if list := s.ListCheckinsByUser("u1"); len(list) != 1 {
		t.Fatalf("ListCheckinsByUser(u1) len = %d, want 1", len(list))
	}
	if list := s.ListCheckins(); len(list) != 2 {
		t.Fatalf("ListCheckins len = %d, want 2", len(list))
	}

	if _, err := s.GetCheckin("c1"); err != nil {
		t.Fatalf("GetCheckin: %v", err)
	}
	if err := s.DeleteCheckin("c1"); err != nil {
		t.Fatalf("DeleteCheckin: %v", err)
	}
	if err := s.DeleteCheckin("c1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRewardStoreCRUD(t *testing.T) {
	s := NewMemoryStore()

	r := &model.Reward{ID: "r1", Name: "连续签到奖", Points: 100, RequiredDays: 7}
	if err := s.CreateReward(r); err != nil {
		t.Fatalf("CreateReward: %v", err)
	}

	got, err := s.GetReward("r1")
	if err != nil {
		t.Fatalf("GetReward: %v", err)
	}
	if got.Points != 100 {
		t.Fatalf("points = %d, want 100", got.Points)
	}

	if _, err := s.GetReward("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if list := s.ListRewards(); len(list) != 1 {
		t.Fatalf("ListRewards len = %d, want 1", len(list))
	}

	r.Points = 200
	if err := s.UpdateReward(r); err != nil {
		t.Fatalf("UpdateReward: %v", err)
	}
	if err := s.UpdateReward(&model.Reward{ID: "nope"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if err := s.DeleteReward("r1"); err != nil {
		t.Fatalf("DeleteReward: %v", err)
	}
	if err := s.DeleteReward("r1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
