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
	sortRewardsByDaysThenPoints(all)
	return all, nil
}

// rewardLess 定义奖励列表的展示顺序：先按所需连续天数升序，
// 天数相同时再按积分降序（积分高的排在前面）。
func rewardLess(a, b *model.Reward) bool {
	if a.RequiredDays != b.RequiredDays {
		return a.RequiredDays < b.RequiredDays
	}
	return a.Points > b.Points
}

// sortRewardsByDaysThenPoints 按展示顺序对奖励切片进行稳定排序。
func sortRewardsByDaysThenPoints(rs []*model.Reward) {
	sort.SliceStable(rs, func(i, j int) bool {
		return rewardLess(rs[i], rs[j])
	})
}

func (s *Service) UpdateReward(id string, r model.Reward) (*model.Reward, error) {
	exist, err := s.store.GetReward(id)
	if err != nil {
		return nil, err
	}
	exist.Name = r.Name
	exist.Points = r.Points
	exist.RequiredDays = r.RequiredDays
	exist.Description = r.Description
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateReward(exist); err != nil {
		return nil, err
	}
	return exist, nil
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
	sortRewardsByDaysThenPoints(claimable)
	return claimable, nil
}
