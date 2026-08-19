// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"checkin/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	CreateUser(u *model.User) error
	GetUser(id string) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
	ListUsers() []*model.User
	UpdateUser(u *model.User) error
	DeleteUser(id string) error

	CreateCheckin(c *model.Checkin) error
	GetCheckin(id string) (*model.Checkin, error)
	GetCheckinByUserDate(userID, date string) (*model.Checkin, error)
	ListCheckins() []*model.Checkin
	ListCheckinsByUser(userID string) []*model.Checkin
	DeleteCheckin(id string) error

	CreateReward(r *model.Reward) error
	GetReward(id string) (*model.Reward, error)
	ListRewards() []*model.Reward
	UpdateReward(r *model.Reward) error
	DeleteReward(id string) error

	GetMakeupCard(userID string) (*model.MakeupCard, error)
	CreateMakeupCard(m *model.MakeupCard) error
	UpdateMakeupCard(m *model.MakeupCard) error

	CreateBadge(b *model.Badge) error
	GetBadge(id string) (*model.Badge, error)
	ListBadges() []*model.Badge
	UpdateBadge(b *model.Badge) error
	DeleteBadge(id string) error
}
