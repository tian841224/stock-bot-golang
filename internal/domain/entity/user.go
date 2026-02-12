package entity

import (
	"github.com/tian841224/stock-bot/internal/domain/valueobject"

	domainerror "github.com/tian841224/stock-bot/internal/domain/error"
)

// User 使用者
type User struct {
	ID        uint
	AccountID string
	UserType  valueobject.UserType
	Status    bool
}

// Validate 驗證使用者資料的合法性
func (u *User) Validate() error {
	if u.AccountID == "" {
		return domainerror.ErrInvalidArgument
	}
	if !u.IsValidUserType() {
		return domainerror.NewInvalidUserTypeError(string(rune(u.UserType)))
	}
	return nil
}

// IsValidUserType 檢查使用者類型是否有效
func (u *User) IsValidUserType() bool {
	return u.UserType == valueobject.UserTypeTelegram ||
		u.UserType == valueobject.UserTypeLine
}

func (u *User) IsActive() bool {
	return u.Status
}

func (u *User) Enable() { u.Status = true }

func (u *User) Disable() { u.Status = false }
