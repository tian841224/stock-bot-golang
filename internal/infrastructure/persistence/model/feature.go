package models

// 功能模型
type Feature struct {
	Model
	Name        string `gorm:"column:name;type:varchar(255);uniqueIndex;not null" json:"name"`
	Code        string `gorm:"column:code;type:varchar(255);uniqueIndex;not null" json:"code"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
}

func (Feature) TableName() string {
	return "features"
}

func init() {
	RegisterModel(&Feature{})
}

