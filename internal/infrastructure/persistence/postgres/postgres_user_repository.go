// Package repository 提供資料存取層實作
package repository

import (
	"context"
	"fmt"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

var _ repo.UserReader = (*postgresUserRepository)(nil)
var _ repo.UserWriter = (*postgresUserRepository)(nil)
var _ repo.UserAccountPort = (*postgresUserRepository)(nil)

func NewPostgresUserRepository(db *gorm.DB) *postgresUserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) toEntity(model *models.User) *entity.User {
	return &entity.User{
		ID:        model.ID,
		AccountID: model.AccountID,
		UserType:  model.UserType,
		Status:    model.Status,
	}
}

func (r *postgresUserRepository) toModel(entity *entity.User) *models.User {
	return &models.User{
		Model: models.Model{
			ID: entity.ID,
		},
		AccountID: entity.AccountID,
		UserType:  entity.UserType,
		Status:    entity.Status,
	}
}

// GetByID 根據 ID 取得使用者
func (r *postgresUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&user), nil
}

// GetByAccountID 根據帳號 ID 和使用者類型取得使用者
func (r *postgresUserRepository) GetByAccountID(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("account_id = ? AND user_type = ?", accountID, userType).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&user), nil
}

func (r *postgresUserRepository) GetOrCreate(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error) {
	user, err := r.GetByAccountID(ctx, accountID, userType)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &entity.User{
			AccountID: accountID,
			UserType:  userType,
			Status:    true,
		}
		err = r.Create(ctx, user)
	}
	return user, nil
}

// List 取得使用者列表
func (r *postgresUserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.User, 0, len(users))
	for _, user := range users {
		entities = append(entities, r.toEntity(user))
	}

	return entities, nil
}

// Create 建立新使用者
func (r *postgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	dbModel := r.toModel(user)
	return r.db.WithContext(ctx).Create(dbModel).Error
}

// Update 更新使用者資料
func (r *postgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	dbModel := r.toModel(user)

	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(dbModel)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found with id: %d", user.ID)
	}

	return nil
}

// Delete 刪除使用者
func (r *postgresUserRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found with id: %d", id)
	}

	return nil
}
