package store

import "checkin/internal/model"

func cloneRecordReward(r *model.Reward) *model.Reward {
	if r == nil {
		return nil
	}
	copy := *r
	return &copy
}

func cloneRecordCheckin(c *model.Checkin) *model.Checkin {
	if c == nil {
		return nil
	}
	copy := *c
	return &copy
}
