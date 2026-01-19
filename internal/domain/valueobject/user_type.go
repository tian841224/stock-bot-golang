package valueobject

// UserType types
type UserType int

const (
	UserTypeTelegram UserType = iota + 1
	UserTypeLine
)

// UserTypeMap mapping table for user types
var UserTypeMap = map[string]UserType{
	"1": UserTypeTelegram,
	"2": UserTypeLine,
}

// GetName 回傳使用者類型名稱
func (u UserType) GetName() string {
	switch u {
	case UserTypeTelegram:
		return "Telegram"
	case UserTypeLine:
		return "Line"
	default:
		return "Unknown"
	}
}
