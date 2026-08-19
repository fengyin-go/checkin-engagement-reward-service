package store

import "checkin/internal/model"

// CreateBadge 创建徽章规则。
func (s *MemoryStore) CreateBadge(b *model.Badge) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.badges[b.ID] = b
	return nil
}

// GetBadge 按 ID 查询徽章规则。
func (s *MemoryStore) GetBadge(id string) (*model.Badge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.badges[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

// ListBadges 列出全部徽章规则。
func (s *MemoryStore) ListBadges() []*model.Badge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Badge, 0)
	for _, b := range s.badges {
		list = append(list, b)
	}
	return list
}

// UpdateBadge 更新徽章规则。
func (s *MemoryStore) UpdateBadge(b *model.Badge) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.badges[b.ID]; !ok {
		return ErrNotFound
	}
	s.badges[b.ID] = b
	return nil
}

// DeleteBadge 删除徽章规则。
func (s *MemoryStore) DeleteBadge(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.badges[id]; !ok {
		return ErrNotFound
	}
	delete(s.badges, id)
	return nil
}
