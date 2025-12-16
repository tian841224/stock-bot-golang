// Package models 提供資料庫模型定義
package models

import (
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

// User 使用者模型
type User struct {
	Model
	AccountID string `gorm:"column:account_id;type:varchar(255);uniqueIndex;not null" json:"account_id"`
	// 使用者類型 TG or LINE
	UserType valueobject.UserType `gorm:"column:user_type;type:SMALLINT;not null;check:user_type IN (1,2)" json:"user_type"`
	Status   bool                 `gorm:"column:status;type:boolean" json:"status"`
}

func (u *User) GetUserType() valueobject.UserType {
	return valueobject.UserType(u.UserType)
}

func (User) TableName() string {
	return "users"
}

func init() {
	RegisterModel(&User{})
}
