package repository

import (
	"github.com/tian841224/stock-bot/internal/db/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByAccountID(accountID string, userType models.UserType) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 建立新使用者
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID 根據 ID 取得使用者
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByAccountID 根據帳號 ID 和使用者類型取得使用者
func (r *userRepository) GetByAccountID(accountID string, userType models.UserType) (*models.User, error) {
	var user models.User
	err := r.db.Where("account_id = ? AND user_type = ?", accountID, userType).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新使用者資料
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete 刪除使用者
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List 取得使用者列表
func (r *userRepository) List(offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}
