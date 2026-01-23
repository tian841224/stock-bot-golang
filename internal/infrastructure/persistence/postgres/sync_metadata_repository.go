package repository

import (
	"context"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type postgresSyncMetadataRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ repo.SyncMetadataRepository = (*postgresSyncMetadataRepository)(nil)

func NewSyncMetadataRepository(db *gorm.DB, log logger.Logger) *postgresSyncMetadataRepository {
	return &postgresSyncMetadataRepository{
		db:     db,
		logger: log,
	}
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
	r.logger.Info("Upserting sync metadata", logger.String("market", metadata.Market), logger.Int("total_count", metadata.TotalCount))

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
		r.logger.Info("Creating new sync metadata", logger.String("market", metadata.Market))
		err := r.db.WithContext(ctx).Create(&model).Error
		if err != nil {
			r.logger.Error("Failed to create sync metadata", logger.Error(err), logger.String("market", metadata.Market))
			return err
		}
		r.logger.Info("Sync metadata created successfully", logger.String("market", metadata.Market))
		return nil
	} else if err != nil {
		r.logger.Error("Failed to query sync metadata", logger.Error(err), logger.String("market", metadata.Market))
		return err
	}

	err = r.db.WithContext(ctx).Model(&existing).Updates(model).Error
	if err != nil {
		r.logger.Error("Failed to update sync metadata", logger.Error(err), logger.String("market", metadata.Market))
		return err
	}

	r.logger.Info("Sync metadata updated successfully", logger.String("market", metadata.Market))
	return nil
}
