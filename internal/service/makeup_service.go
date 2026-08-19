package service

import (
	"sort"
	"time"

	"checkin/internal/model"
	"checkin/pkg/idgen"
)

// GrantMakeupCard 为用户发放补签卡，无账户则先创建。
func (s *Service) GrantMakeupCard(userID string, amount int) (*model.MakeupCard, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, model.NewValidationError("amount", "发放数量必须大于 0")
	}
	card, err := s.store.GetMakeupCard(userID)
	if err != nil {
		card = &model.MakeupCard{ID: idgen.Hex(), UserID: userID, Balance: 0, UpdatedAt: time.Now().UTC()}
		if err := s.store.CreateMakeupCard(card); err != nil {
			return nil, err
		}
	}
	card.Balance = amount
	card.UpdatedAt = time.Now()
	if err := s.store.UpdateMakeupCard(card); err != nil {
		return nil, err
	}
	return card, nil
}

// GetMakeupCard 查询用户补签卡余额，无账户时返回余额为 0 的空账户。
func (s *Service) GetMakeupCard(userID string) (*model.MakeupCard, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, err
	}
	card, err := s.store.GetMakeupCard(userID)
	if err != nil {
		return &model.MakeupCard{UserID: userID, Balance: 0}, nil
	}
	return card, nil
}

// MakeupCheckIn 消耗一张补签卡，为过去某一天补签。
func (s *Service) MakeupCheckIn(userID, date string) (*model.Checkin, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, err
	}
	if !model.IsValidDate(date) {
		return nil, model.NewValidationError("date", "日期格式必须为 YYYY-MM-DD")
	}
	today := model.DateOf(time.Now())
	if date >= today {
		return nil, model.NewValidationError("date", "补签日期必须早于今天")
	}
	if _, err := s.store.GetCheckinByUserDate(userID, date); err == nil {
		return nil, model.NewValidationError("date", "该日期已签到，无需补签")
	}
	card, err := s.store.GetMakeupCard(userID)
	if err != nil {
		return nil, model.NewValidationError("makeup_card", "暂无补签卡")
	}
	if card.Balance < 1 {
		return nil, model.NewValidationError("makeup_card", "补签卡余额不足")
	}
	card.Balance--
	card.UpdatedAt = time.Now()
	if err := s.store.UpdateMakeupCard(card); err != nil {
		return nil, err
	}
	c := &model.Checkin{
		ID:        idgen.Hex(),
		UserID:    userID,
		Date:      date,
		Source:    model.SourceMakeup,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateCheckin(c); err != nil {
		return nil, err
	}
	return c, nil
}

// ListMakeups 列出用户的补签记录，按日期倒序。
func (s *Service) ListMakeups(userID string) ([]*model.Checkin, error) {
	if _, err := s.store.GetUser(userID); err != nil {
		return nil, err
	}
	all := s.store.ListCheckinsByUser(userID)
	makeups := make([]*model.Checkin, 0)
	for _, c := range all {
		if c.Source == model.SourceMakeup {
			makeups = append(makeups, c)
		}
	}
	sort.Slice(makeups, func(i, j int) bool {
		return makeups[i].Date > makeups[j].Date
	})
	return makeups, nil
}
