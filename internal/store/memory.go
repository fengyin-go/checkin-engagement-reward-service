package store

import (
	"sync"

	"checkin/internal/model"
)

// MemoryStore 基于内存的 Store 实现，线程安全。
type MemoryStore struct {
	mu         sync.RWMutex
	users      map[string]*model.User
	checkins   map[string]*model.Checkin
	rewards    map[string]*model.Reward
	makeupCards map[string]*model.MakeupCard // key = userID
	badges     map[string]*model.Badge
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:       make(map[string]*model.User),
		checkins:    make(map[string]*model.Checkin),
		rewards:     make(map[string]*model.Reward),
		makeupCards: make(map[string]*model.MakeupCard),
		badges:      make(map[string]*model.Badge),
	}
}

var _ Store = (*MemoryStore)(nil)
