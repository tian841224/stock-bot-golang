package repository

import (
	"context"
	"fmt"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgresStockSymbolsRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

var _ repo.StockSymbolReader = (*postgresStockSymbolsRepository)(nil)
var _ repo.StockSymbolWriter = (*postgresStockSymbolsRepository)(nil)

func NewSymbolRepository(db *gorm.DB, log logger.Logger) *postgresStockSymbolsRepository {
	return &postgresStockSymbolsRepository{
		db:     db,
		logger: log,
	}
}

// GetByID 根據 ID 取得股票代號
func (r *postgresStockSymbolsRepository) GetByID(ctx context.Context, id uint) (*entity.StockSymbol, error) {
	var symbol models.StockSymbol
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&symbol).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity.StockSymbol{
		ID:     symbol.ID,
		Symbol: symbol.Symbol,
		Market: symbol.Market,
		Name:   symbol.Name,
	}, nil
}

// GetBySymbolAndMarket 根據股票代號和市場取得資料
func (r *postgresStockSymbolsRepository) GetBySymbolAndMarket(ctx context.Context, symbol, market string) (*entity.StockSymbol, error) {
	var symbolData models.StockSymbol
	err := r.db.WithContext(ctx).Where("symbol = ? AND market = ?", symbol, market).First(&symbolData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity.StockSymbol{
		ID:     symbolData.ID,
		Symbol: symbolData.Symbol,
		Market: symbolData.Market,
		Name:   symbolData.Name,
	}, nil
}

// GetBySymbol 根據股票代號取得資料
func (r *postgresStockSymbolsRepository) GetBySymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
	var symbolData models.StockSymbol
	err := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&symbolData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity.StockSymbol{
		ID:     symbolData.ID,
		Symbol: symbolData.Symbol,
		Market: symbolData.Market,
		Name:   symbolData.Name,
	}, nil
}

// GetBySymbolID 根據股票 ID 取得股票代號列表
func (r *postgresStockSymbolsRepository) GetBySymbolID(ctx context.Context, symbolID uint) ([]*entity.StockSymbol, error) {
	var symbols []*models.StockSymbol
	err := r.db.WithContext(ctx).Where("symbol_id = ?", symbolID).Find(&symbols).Error
	if err != nil {
		return nil, err
	}

	var entities []*entity.StockSymbol
	for _, symbol := range symbols {
		entities = append(entities, &entity.StockSymbol{
			ID:     symbol.ID,
			Symbol: symbol.Symbol,
			Market: symbol.Market,
			Name:   symbol.Name,
		})
	}
	return entities, nil
}

// Create 建立新股票代號
func (r *postgresStockSymbolsRepository) Create(ctx context.Context, symbol *entity.StockSymbol) error {
	r.logger.Info("Creating stock symbol", logger.String("symbol", symbol.Symbol), logger.String("market", symbol.Market))

	err := r.db.WithContext(ctx).Create(&models.StockSymbol{
		Symbol: symbol.Symbol,
		Market: symbol.Market,
		Name:   symbol.Name,
	}).Error
	if err != nil {
		r.logger.Error("Failed to create stock symbol", logger.Error(err), logger.String("symbol", symbol.Symbol))
		return err
	}

	r.logger.Info("Stock symbol created successfully", logger.String("symbol", symbol.Symbol))
	return nil
}

func (r *postgresStockSymbolsRepository) Update(ctx context.Context, symbol *entity.StockSymbol) error {
	r.logger.Info("Updating stock symbol", logger.Any("id", symbol.ID), logger.String("symbol", symbol.Symbol))

	err := r.db.WithContext(ctx).Model(&models.StockSymbol{}).
		Where("id = ?", symbol.ID).
		Updates(map[string]interface{}{
			"symbol": symbol.Symbol,
			"market": symbol.Market,
			"name":   symbol.Name,
		}).Error
	if err != nil {
		r.logger.Error("Failed to update stock symbol", logger.Error(err), logger.Any("id", symbol.ID))
		return err
	}

	r.logger.Info("Stock symbol updated successfully", logger.Any("id", symbol.ID))
	return nil
}

// Delete 刪除股票代號
func (r *postgresStockSymbolsRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting stock symbol", logger.Any("id", id))

	result := r.db.WithContext(ctx).Delete(&models.StockSymbol{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete stock symbol", logger.Error(result.Error), logger.Any("id", id))
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Stock symbol not found for deletion", logger.Any("id", id))
		return fmt.Errorf("symbol not found with id: %d", id)
	}

	r.logger.Info("Stock symbol deleted successfully", logger.Any("id", id))
	return nil
}

// List 取得股票代號列表
func (r *postgresStockSymbolsRepository) List(ctx context.Context, offset, limit int) ([]*entity.StockSymbol, error) {
	var symbols []*models.StockSymbol
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&symbols).Error
	if err != nil {
		return nil, err
	}

	var entities []*entity.StockSymbol
	for _, symbol := range symbols {
		entities = append(entities, &entity.StockSymbol{
			ID:     symbol.ID,
			Symbol: symbol.Symbol,
			Market: symbol.Market,
			Name:   symbol.Name,
		})
	}

	return entities, nil
}

// BatchCreate 批次建立股票代號
func (r *postgresStockSymbolsRepository) BatchCreate(ctx context.Context, symbols []*entity.StockSymbol) error {
	r.logger.Info("Batch creating stock symbols", logger.Int("count", len(symbols)))

	err := r.db.WithContext(ctx).CreateInBatches(symbols, 100).Error
	if err != nil {
		r.logger.Error("Failed to batch create stock symbols", logger.Error(err), logger.Int("count", len(symbols)))
		return err
	}

	r.logger.Info("Batch create completed successfully", logger.Int("count", len(symbols)))
	return nil
}

// BatchUpsert 批次更新或建立股票代號
func (r *postgresStockSymbolsRepository) BatchUpsert(ctx context.Context, symbols []*entity.StockSymbol) (successCount, errorCount int, err error) {
	if len(symbols) == 0 {
		r.logger.Debug("BatchUpsert called with empty symbols list")
		return 0, 0, nil
	}

	r.logger.Info("Starting batch upsert", logger.Int("count", len(symbols)))

	modelSymbols := make([]models.StockSymbol, 0, len(symbols))
	for _, symbol := range symbols {
		modelSymbols = append(modelSymbols, models.StockSymbol{
			Symbol: symbol.Symbol,
			Market: symbol.Market,
			Name:   symbol.Name,
		})
	}

	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "symbol"}, {Name: "market"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).CreateInBatches(modelSymbols, 100)

	if result.Error != nil {
		r.logger.Error("Failed to batch upsert stock symbols", logger.Error(result.Error), logger.Int("count", len(symbols)))
		return 0, len(symbols), result.Error
	}

	r.logger.Info("Batch upsert completed successfully", logger.Int("affected_rows", int(result.RowsAffected)))
	return int(result.RowsAffected), 0, nil
}

// GetMarketStats 取得各市場的股票統計資訊
func (r *postgresStockSymbolsRepository) GetMarketStats(ctx context.Context) (map[string]int, error) {
	type Result struct {
		Market string
		Count  int
	}

	var results []Result
	err := r.db.WithContext(ctx).
		Model(&models.StockSymbol{}).
		Select("market, COUNT(*) as count").
		Group("market").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	for _, result := range results {
		stats[result.Market] = result.Count
	}

	return stats, nil
}
