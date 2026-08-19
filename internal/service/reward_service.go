package service

import (
	"sort"

	"checkin/internal/model"
	"checkin/pkg/idgen"
)

// CreateReward 创建奖励规则。
func (s *Service) CreateReward(r model.Reward) (*model.Reward, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	r.ID = idgen.Hex()
	if err := s.store.CreateReward(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Service) GetReward(id string) (*model.Reward, error) {
	return s.store.GetReward(id)
}

func (s *Service) ListRewards() ([]*model.Reward, error) {
	all := s.store.ListRewards()
	sort.Slice(all, func(i, j int) bool {
		return all[i].RequiredDays < all[j].RequiredDays
	})
	return all, nil
}

func (s *Service) UpdateReward(id string, r model.Reward) (*model.Reward, error) {
	exist, err := s.store.GetReward(id)
	if err != nil {
		return nil, err
	}
	updated := *exist
	updated.Name = r.Name
	updated.Points = r.Points
	updated.RequiredDays = r.RequiredDays
	updated.Description = r.Description
	if err := updated.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateReward(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *Service) DeleteReward(id string) error {
	return s.store.DeleteReward(id)
}

// ClaimableRewards 返回当前连续天数已满足领取条件的奖励。
func (s *Service) ClaimableRewards(userID string) ([]*model.Reward, error) {
	summary, err := s.GetSummary(userID)
	if err != nil {
		return nil, err
	}
	rewards := s.store.ListRewards()
	claimable := make([]*model.Reward, 0)
	for _, r := range rewards {
		if summary.Streak >= r.RequiredDays {
			claimable = append(claimable, r)
		}
	}
	sort.Slice(claimable, func(i, j int) bool {
		return claimable[i].RequiredDays < claimable[j].RequiredDays
	})
	return claimable, nil
}
