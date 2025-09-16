package user

import (
	"errors"
	"stock-bot/internal/db/models"
	"stock-bot/internal/repository"

	"gorm.io/gorm"
)

// UserService 使用者服務介面
type UserService interface {
	CreateUser(accountID string, userType models.UserType) (*models.User, error)
	GetUserByAccountID(accountID string, userType models.UserType) (*models.User, error)
	UpdateUserStatus(userID uint, status bool) error
	GetUserList(page, pageSize int) ([]*models.User, error)
	GetOrCreate(accountID string, userType models.UserType) (*models.User, error)
}

// userService 使用者服務實作
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 建立新的使用者服務
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// CreateUser 建立新使用者
func (s *userService) CreateUser(accountID string, userType models.UserType) (*models.User, error) {
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
func (s *userService) GetUserByAccountID(accountID string, userType models.UserType) (*models.User, error) {
	return s.userRepo.GetByAccountID(accountID, userType)
}

// UpdateUserStatus 更新使用者狀態
func (s *userService) UpdateUserStatus(userID uint, status bool) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Status = status
	return s.userRepo.Update(user)
}

// GetUserList 取得使用者列表
func (s *userService) GetUserList(page, pageSize int) ([]*models.User, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

// GetOrCreate 取得或建立使用者
func (s *userService) GetOrCreate(accountID string, userType models.UserType) (*models.User, error) {
	// 先嘗試取得現有使用者
	user, err := s.userRepo.GetByAccountID(accountID, userType)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 如果使用者不存在，建立新使用者
	if user == nil {
		user = &models.User{
			AccountID: accountID,
			UserType:  userType,
			Status:    true,
		}

		err = s.userRepo.Create(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
