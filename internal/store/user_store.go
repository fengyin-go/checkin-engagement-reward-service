package store

import (
	"strings"

	"checkin/internal/model"
)

// normalizeName 将用户名称归一化（去首尾空白并转小写），作为名称唯一性的比较键。
// 与 model.UserFilter.Match 的匹配口径保持一致：归一化后相同的名称视为同名，
// 避免出现 "alice" 与 "  alice  "（乃至 "Alice"）被当成两个人的情况。
func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (s *MemoryStore) CreateUser(u *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := normalizeName(u.Name)
	for _, exist := range s.users {
		if normalizeName(exist.Name) == key {
			return ErrConflict
		}
	}
	s.users[u.ID] = u
	return nil
}

func (s *MemoryStore) GetUser(id string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

// GetUserByName 按归一化名称查找用户，找不到返回 ErrNotFound。
func (s *MemoryStore) GetUserByName(name string) (*model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := normalizeName(name)
	for _, u := range s.users {
		if normalizeName(u.Name) == key {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListUsers() []*model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.User, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, u)
	}
	return list
}

func (s *MemoryStore) UpdateUser(u *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[u.ID]; !ok {
		return ErrNotFound
	}
	key := normalizeName(u.Name)
	for id, exist := range s.users {
		if id != u.ID && normalizeName(exist.Name) == key {
			return ErrConflict
		}
	}
	s.users[u.ID] = u
	return nil
}

func (s *MemoryStore) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[id]; !ok {
		return ErrNotFound
	}
	delete(s.users, id)
	return nil
}
