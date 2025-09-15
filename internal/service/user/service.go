package user

import (
	"errors"
	"stock-bot/internal/db/models"
	"stock-bot/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser 建立新使用者
func (s *UserService) CreateUser(accountID string, userType models.UserType) (*models.User, error) {
	// 檢查使用者是否已存在
	existingUser, err := s.userRepo.GetByAccountID(accountID, userType)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUser != nil {
		return existingUser, nil // 使用者已存在，回傳現有使用者
	}

	// 建立新使用者
	user := &models.User{
		AccountID: accountID,
		UserType:  userType,
		Status:    true,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByAccountID 根據帳號 ID 取得使用者
func (s *UserService) GetUserByAccountID(accountID string, userType models.UserType) (*models.User, error) {
	return s.userRepo.GetByAccountID(accountID, userType)
}

// UpdateUserStatus 更新使用者狀態
func (s *UserService) UpdateUserStatus(userID uint, status bool) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Status = status
	return s.userRepo.Update(user)
}

// GetUserList 取得使用者列表
func (s *UserService) GetUserList(page, pageSize int) ([]*models.User, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}
