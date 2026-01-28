package repository

import (
	"context"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type featureRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ repo.FeatureReader = (*featureRepository)(nil)
var _ repo.FeatureWriter = (*featureRepository)(nil)

func NewFeatureRepository(db *gorm.DB, log logger.Logger) (repo.FeatureReader, repo.FeatureWriter) {
	repository := &featureRepository{
		db:     db,
		logger: log,
	}

	if err := repository.createDefaultFeatures(); err != nil {
		repository.logger.Error("Failed to create default features", logger.Error(err))
	}

	return repository, repository
}

// GetByID 根據 ID 取得功能
func (r *featureRepository) GetByID(ctx context.Context, id uint) (*entity.Feature, error) {
	var feature models.Feature
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&feature).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Feature not found", logger.Any("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get feature by ID", logger.Error(err), logger.Any("id", id))
		return nil, err
	}

	return &entity.Feature{
		ID:          feature.ID,
		Name:        feature.Name,
		Code:        feature.Code,
		Description: feature.Description,
	}, nil
}

// GetByCode 根據代碼取得功能
func (r *featureRepository) GetByCode(ctx context.Context, code string) (*entity.Feature, error) {
	var feature models.Feature
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&feature).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity.Feature{
		ID:          feature.ID,
		Name:        feature.Name,
		Code:        feature.Code,
		Description: feature.Description,
	}, nil
}

// GetByName 根據名稱取得功能
func (r *featureRepository) GetByName(ctx context.Context, name string) (*entity.Feature, error) {
	var feature models.Feature
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&feature).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity.Feature{
		ID:          feature.ID,
		Name:        feature.Name,
		Code:        feature.Code,
		Description: feature.Description,
	}, nil
}

// Create 建立新功能
func (r *featureRepository) Create(ctx context.Context, feature *entity.Feature) error {
	r.logger.Info("Creating feature", logger.String("name", feature.Name), logger.String("code", feature.Code))

	err := r.db.WithContext(ctx).Create(&models.Feature{
		Name:        feature.Name,
		Code:        feature.Code,
		Description: feature.Description,
	}).Error
	if err != nil {
		r.logger.Error("Failed to create feature", logger.Error(err), logger.String("name", feature.Name))
		return err
	}

	r.logger.Info("Feature created successfully", logger.String("name", feature.Name))
	return nil
}

// Update 更新功能資料
func (r *featureRepository) Update(ctx context.Context, feature *entity.Feature) error {
	r.logger.Info("Updating feature", logger.Any("id", feature.ID), logger.String("name", feature.Name))

	err := r.db.WithContext(ctx).Model(&models.Feature{}).
		Where("id = ?", feature.ID).
		Updates(map[string]interface{}{
			"name":        feature.Name,
			"code":        feature.Code,
			"description": feature.Description,
		}).Error

	if err != nil {
		r.logger.Error("Failed to update feature", logger.Error(err), logger.Any("id", feature.ID))
		return err
	}

	r.logger.Info("Feature updated successfully", logger.Any("id", feature.ID))
	return nil
}

// Delete 刪除功能
func (r *featureRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting feature", logger.Any("id", id))

	result := r.db.WithContext(ctx).Delete(&models.Feature{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete feature", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Feature not found for deletion", logger.Any("id", id))
		return nil
	}

	r.logger.Info("Feature deleted successfully", logger.Any("id", id))
	return nil
}

// List 取得功能列表
func (r *featureRepository) List(ctx context.Context, offset, limit int) ([]*entity.Feature, error) {
	var features []*models.Feature
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&features).Error
	if err != nil {
		return nil, err
	}

	var entities []*entity.Feature
	for _, feature := range features {
		entities = append(entities, &entity.Feature{
			ID:          feature.ID,
			Name:        feature.Name,
			Code:        feature.Code,
			Description: feature.Description,
		})
	}

	return entities, nil
}

// GetAll 取得所有功能
func (r *featureRepository) GetAll(ctx context.Context) ([]*entity.Feature, error) {
	var features []*models.Feature
	err := r.db.WithContext(ctx).Find(&features).Error
	if err != nil {
		return nil, err
	}

	var entities []*entity.Feature
	for _, feature := range features {
		entities = append(entities, &entity.Feature{
			ID:          feature.ID,
			Name:        feature.Name,
			Code:        feature.Code,
			Description: feature.Description,
		})
	}

	return entities, nil
}

// createDefaultFeatures 建立預設功能資料
func (r *featureRepository) createDefaultFeatures() error {
	defaultFeatures := []*models.Feature{
		{
			Name:        "Stock Info",
			Code:        "1",
			Description: models.SubscriptionItemStockInfo.GetName(),
		},
		{
			Name:        "Stock News",
			Code:        "2",
			Description: models.SubscriptionItemStockNews.GetName(),
		},
		{
			Name:        "Daily Market Info",
			Code:        "3",
			Description: models.SubscriptionItemDailyMarketInfo.GetName(),
		},
		{
			Name:        "Top Volume Items",
			Code:        "4",
			Description: models.SubscriptionItemTopVolumeItems.GetName(),
		},
	}

	for _, feature := range defaultFeatures {
		// 檢查是否已存在
		var existingFeature models.Feature
		err := r.db.Where("code = ?", feature.Code).First(&existingFeature).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 不存在則建立
				if err := r.db.Create(feature).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
