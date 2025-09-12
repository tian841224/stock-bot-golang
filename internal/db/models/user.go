package models

// 使用者模型
type User struct {
	Model
	AccountID string `gorm:"column:account_id;type:varchar(255);uniqueIndex;not null" json:"account_id"`
	// 使用者類型 TG or LINE
	UserType UserType `gorm:"column:user_type;type:SMALLINT;not null;check:user_type IN (1,2)" json:"user_type"`
	Status   bool     `gorm:"column:status;type:boolean" json:"status"`
}

type UserType int

const (
	UserTypeTelegram UserType = 1
	UserTypeLine     UserType = 2
)

func (u *User) GetUserType() UserType {
	return UserType(u.UserType)
}

func (User) TableName() string {
	return "users"
}

func init() {
	RegisterModel(&User{})
}
