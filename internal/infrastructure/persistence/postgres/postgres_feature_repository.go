package repository

import (
	"context"

	"github.com/tian841224/stock-bot/internal/domain/entity"
	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type featureRepository struct {
	db *gorm.DB
}

var _ repo.FeatureReader = (*featureRepository)(nil)
var _ repo.FeatureWriter = (*featureRepository)(nil)

func NewFeatureRepository(db *gorm.DB) (repo.FeatureReader, repo.FeatureWriter) {
	repository := &featureRepository{db: db}

	// 自動建立預設功能資料
	_ = repository.createDefaultFeatures()

	return repository, repository
}

// GetByID 根據 ID 取得功能
func (r *featureRepository) GetByID(ctx context.Context, id uint) (*entity.Feature, error) {
	var feature models.Feature
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&feature).Error
	if err != nil {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
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
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
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
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
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
	err := r.db.WithContext(ctx).Create(&models.Feature{
		Name:        feature.Name,
		Code:        feature.Code,
		Description: feature.Description,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

// Update 更新功能資料
func (r *featureRepository) Update(ctx context.Context, feature *entity.Feature) error {
	err := r.db.WithContext(ctx).Model(&models.Feature{}).
		Where("id = ?", feature.ID).
		Updates(map[string]interface{}{
			"name":        feature.Name,
			"code":        feature.Code,
			"description": feature.Description,
		}).Error
	return err
}

// Delete 刪除功能
func (r *featureRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&models.Feature{}, id).Error
	if err != nil {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return r.db.WithContext(ctx).Delete(&models.Feature{}, id).Error
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

// createDefaultFeatures 建立預設功能資料（私有方法）
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
		if err == gorm.ErrRecordNotFound {
			// 不存在則建立
			if err := r.db.Create(feature).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
