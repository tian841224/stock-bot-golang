package repository

import (
	"context"
	"time"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

type postgresTradeDateRepository struct {
	db *gorm.DB
}

func NewPostgresTradeDateRepository(db *gorm.DB) *postgresTradeDateRepository {
	return &postgresTradeDateRepository{db: db}
}

var _ repo.TradeDateReader = (*postgresTradeDateRepository)(nil)
var _ repo.TradeDateWriter = (*postgresTradeDateRepository)(nil)

func (r *postgresTradeDateRepository) toEntity(model *models.TradeDate) *entity.TradeDate {
	return &entity.TradeDate{
		ID:       model.ID,
		Date:     model.Date,
		Exchange: model.Exchange,
	}
}

func (r *postgresTradeDateRepository) toModel(entity *entity.TradeDate) *models.TradeDate {
	return &models.TradeDate{
		Model: models.Model{
			ID: entity.ID,
		},
		Date:     entity.Date,
		Exchange: entity.Exchange,
	}
}

// GetByID 根據 ID 取得交易日資料
func (r *postgresTradeDateRepository) GetByID(ctx context.Context, id uint) (*entity.TradeDate, error) {
	var tradeDate models.TradeDate
	err := r.db.WithContext(ctx).First(&tradeDate, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&tradeDate), nil
}

// GetByDate 根據日期取得交易日資料
func (r *postgresTradeDateRepository) GetByDate(ctx context.Context, date time.Time) (*entity.TradeDate, error) {
	var tradeDate models.TradeDate
	err := r.db.WithContext(ctx).Where("date = ?", date).First(&tradeDate).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&tradeDate), nil
}

// GetByDateRange 根據日期範圍取得交易日資料
func (r *postgresTradeDateRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.TradeDate, error) {
	var tradeDates []*models.TradeDate
	err := r.db.WithContext(ctx).Where("date BETWEEN ? AND ?", startDate, endDate).Find(&tradeDates).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	var entities []*entity.TradeDate
	for _, tradeDate := range tradeDates {
		entities = append(entities, r.toEntity(tradeDate))
	}
	return entities, nil
}

// Create 建立新交易日資料
func (r *postgresTradeDateRepository) Create(ctx context.Context, tradeDate *entity.TradeDate) error {
	dbModel := r.toModel(tradeDate)
	return r.db.WithContext(ctx).Create(dbModel).Error
}

// BatchCreateTradeDates 批次建立交易日資料
func (r *postgresTradeDateRepository) BatchCreateTradeDates(ctx context.Context, tradeDates []*entity.TradeDate) error {
	dbModels := make([]*models.TradeDate, 0, len(tradeDates))
	for _, tradeDate := range tradeDates {
		dbModels = append(dbModels, r.toModel(tradeDate))
	}
	return r.db.WithContext(ctx).CreateInBatches(dbModels, 100).Error
}
