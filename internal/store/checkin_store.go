package store

import "checkin/internal/model"

func (s *MemoryStore) CreateCheckin(c *model.Checkin) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.checkins {
		if exist.UserID == c.UserID && exist.Date == c.Date && exist.Source == c.Source {
			return ErrConflict
		}
	}
	s.checkins[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCheckin(id string) (*model.Checkin, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.checkins[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetCheckinByUserDate(userID, date string) (*model.Checkin, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.checkins {
		if c.UserID == userID && c.Date == date {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCheckins() []*model.Checkin {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Checkin, 0, len(s.checkins))
	for _, c := range s.checkins {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) ListCheckinsByUser(userID string) []*model.Checkin {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Checkin, 0)
	for _, c := range s.checkins {
		if c.UserID == userID {
			list = append(list, c)
		}
	}
	return list
}

func (s *MemoryStore) DeleteCheckin(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.checkins[id]; !ok {
		return ErrNotFound
	}
	delete(s.checkins, id)
	return nil
}
