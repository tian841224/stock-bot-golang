package repository

import (
	"context"
	"time"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresSystemStatisticsRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ repo.SystemStatisticsRepository = (*postgresSystemStatisticsRepository)(nil)

func NewPostgresSystemStatisticsRepository(db *gorm.DB, log logger.Logger) *postgresSystemStatisticsRepository {
	return &postgresSystemStatisticsRepository{
		db:     db,
		logger: log,
	}
}

func (r *postgresSystemStatisticsRepository) IncrementVisitCount(ctx context.Context) error {
	return r.db.WithContext(ctx).Model(&models.SystemStatistics{}).
		Where("key = ?", "total_interactions").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"value": gorm.Expr("value + 1")}),
		}).
		Create(&models.SystemStatistics{Key: "total_interactions", Value: 1}).Error
}

func (r *postgresSystemStatisticsRepository) GetTotalVisits(ctx context.Context) (int64, error) {
	var stats models.SystemStatistics
	err := r.db.WithContext(ctx).Where("key = ?", "total_interactions").First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return stats.Value, nil
}

func (r *postgresSystemStatisticsRepository) GetOnlineUsersCount(ctx context.Context, duration time.Duration) (int64, error) {
	var count int64
	threshold := time.Now().Add(-duration)
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("last_activity_at > ?", threshold).
		Count(&count).Error
	return count, err
}
