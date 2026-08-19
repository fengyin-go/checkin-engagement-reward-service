package service

import (
	"sort"
	"time"

	"checkin/internal/model"
	"checkin/pkg/idgen"
)

// CreateUser 创建用户。
func (s *Service) CreateUser(name string) (*model.User, error) {
	u := &model.User{Name: name}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetUserByName(u.Name); err == nil {
		return nil, model.NewValidationError("name", "用户名称已存在")
	}
	u.ID = idgen.Hex()
	u.CreatedAt = time.Now()
	if err := s.store.CreateUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) GetUser(id string) (*model.User, error) {
	return s.store.GetUser(id)
}

func (s *Service) ListUsers(filter model.UserFilter, page, size int) ([]*model.User, int, error) {
	all := s.store.ListUsers()
	matched := make([]*model.User, 0, len(all))
	for _, u := range all {
		if filter.Match(u) {
			matched = append(matched, u)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.User{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateUser 更新用户名称。
func (s *Service) UpdateUser(id, name string) (*model.User, error) {
	u, err := s.store.GetUser(id)
	if err != nil {
		return nil, err
	}
	// 在独立对象上校验新名称，避免覆写存储中的用户而干扰同名检查。
	candidate := &model.User{Name: name}
	if err := candidate.Validate(); err != nil {
		return nil, err
	}
	if exist, err := s.store.GetUserByName(candidate.Name); err == nil && exist.ID != u.ID {
		return nil, model.NewValidationError("name", "用户名称已存在")
	}
	u.Name = candidate.Name
	if err := s.store.UpdateUser(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) DeleteUser(id string) error {
	return s.store.DeleteUser(id)
}
