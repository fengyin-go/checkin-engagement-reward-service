package store

import "checkin/internal/model"

// GetMakeupCard 按用户 ID 查询补签卡，不存在返回 ErrNotFound。
func (s *MemoryStore) GetMakeupCard(userID string) (*model.MakeupCard, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.makeupCards[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

// CreateMakeupCard 为用户创建补签卡账户，已存在则返回 ErrConflict。
func (s *MemoryStore) CreateMakeupCard(m *model.MakeupCard) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.makeupCards[m.UserID]; ok {
		return ErrConflict
	}
	s.makeupCards[m.UserID] = m
	return nil
}

// UpdateMakeupCard 更新补签卡账户，不存在返回 ErrNotFound。
func (s *MemoryStore) UpdateMakeupCard(m *model.MakeupCard) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.makeupCards[m.UserID]; !ok {
		return ErrNotFound
	}
	s.makeupCards[m.UserID] = m
	return nil
}
