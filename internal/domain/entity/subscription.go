package entity

import (
	valueobject "github.com/tian841224/stock-bot/internal/domain/valueobject"
)

type Subscription struct {
	ID           uint
	UserID       uint
	Item         valueobject.SubscriptionType
	Active       bool
	ScheduleCron string
	FeatureID    uint
	Feature      *Feature
	User         *User
}

func (s *Subscription) IsActive() bool { return s.Active }

func (s *Subscription) Enable() { s.Active = true }

func (s *Subscription) Disable() { s.Active = false }

func (s *Subscription) HasSchedule() bool { return s.ScheduleCron != "" }
