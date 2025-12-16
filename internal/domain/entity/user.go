package entity

import "github.com/tian841224/stock-bot/internal/domain/valueobject"

// User 使用者
type User struct {
	ID        uint
	AccountID string
	UserType  valueobject.UserType
	Status    bool
}

func (u *User) IsActive() bool {
	return u.Status
}

func (u *User) Enable() { u.Status = true }

func (u *User) Disable() { u.Status = false }
