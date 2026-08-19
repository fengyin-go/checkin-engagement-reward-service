package store

import "checkin/internal/model"

func (s *MemoryStore) CreateReward(r *model.Reward) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rewards[r.ID] = r
	return nil
}

func (s *MemoryStore) GetReward(id string) (*model.Reward, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rewards[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRewards() []*model.Reward {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Reward, 0, len(s.rewards))
	for _, r := range s.rewards {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateReward(r *model.Reward) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rewards[r.ID]; !ok {
		return ErrNotFound
	}
	s.rewards[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteReward(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rewards[id]; !ok {
		return ErrNotFound
	}
	delete(s.rewards, id)
	return nil
}
