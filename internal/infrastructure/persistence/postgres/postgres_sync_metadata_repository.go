package repository

import (
	"context"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type postgresSyncMetadataRepository struct {
	db *gorm.DB
}

var _ repo.SyncMetadataRepository = (*postgresSyncMetadataRepository)(nil)

func NewSyncMetadataRepository(db *gorm.DB) *postgresSyncMetadataRepository {
	return &postgresSyncMetadataRepository{db: db}
}

func (r *postgresSyncMetadataRepository) GetByMarket(ctx context.Context, market string) (*entity.SyncMetadata, error) {
	var model models.SyncMetadata
	err := r.db.WithContext(ctx).Where("market = ?", market).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity.SyncMetadata{
		ID:            model.ID,
		Market:        model.Market,
		LastSyncAt:    model.LastSyncAt,
		LastSuccessAt: model.LastSuccessAt,
		LastError:     model.LastError,
		TotalCount:    model.TotalCount,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}, nil
}

func (r *postgresSyncMetadataRepository) Upsert(ctx context.Context, metadata *entity.SyncMetadata) error {
	model := models.SyncMetadata{
		Market:        metadata.Market,
		LastSyncAt:    metadata.LastSyncAt,
		LastSuccessAt: metadata.LastSuccessAt,
		LastError:     metadata.LastError,
		TotalCount:    metadata.TotalCount,
	}

	var existing models.SyncMetadata
	err := r.db.WithContext(ctx).Where("market = ?", metadata.Market).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(&model).Error
	} else if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Model(&existing).Updates(model).Error
}
