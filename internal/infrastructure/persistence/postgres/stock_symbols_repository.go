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

func (r *postgresStockSymbolsRepository) toEntity(model *models.StockSymbol) *entity.StockSymbol {
	return &entity.StockSymbol{
		ID:     model.ID,
		Symbol: model.Symbol,
		Market: model.Market,
		Name:   model.Name,
	}
}

func (r *postgresStockSymbolsRepository) toModel(entity *entity.StockSymbol) *models.StockSymbol {
	return &models.StockSymbol{
		Model: models.Model{
			ID: entity.ID,
		},
		Symbol: entity.Symbol,
		Market: entity.Market,
		Name:   entity.Name,
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

	return r.toEntity(&symbol), nil
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

	return r.toEntity(&symbolData), nil
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

	return r.toEntity(&symbolData), nil
}

// GetBySubscriptionID 根據訂閱 ID 取得股票代號列表
func (r *postgresStockSymbolsRepository) GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*entity.StockSymbol, error) {
	var symbols []*models.StockSymbol
	err := r.db.WithContext(ctx).
		Table("stock_symbols").
		Joins("JOIN subscription_symbols ON subscription_symbols.symbol_id = stock_symbols.id").
		Joins("JOIN subscriptions ON subscriptions.user_id = subscription_symbols.user_id").
		Where("subscriptions.id = ?", subscriptionID).
		Find(&symbols).Error

	if err != nil {
		return nil, err
	}

	var entities []*entity.StockSymbol
	for _, symbol := range symbols {
		entities = append(entities, r.toEntity(symbol))
	}
	return entities, nil
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
		entities = append(entities, r.toEntity(symbol))
	}
	return entities, nil
}

// GetBySubscriptionAndSymbol 根據訂閱 ID 和股票 ID 取得股票代號
func (r *postgresStockSymbolsRepository) GetBySubscriptionAndSymbol(ctx context.Context, subscriptionID, symbolID uint) (*entity.StockSymbol, error) {
	var symbol models.StockSymbol
	err := r.db.WithContext(ctx).
		Table("stock_symbols").
		Joins("JOIN subscription_symbols ON subscription_symbols.symbol_id = stock_symbols.id").
		Joins("JOIN subscriptions ON subscriptions.user_id = subscription_symbols.user_id").
		Where("subscriptions.id = ? AND stock_symbols.id = ?", subscriptionID, symbolID).
		First(&symbol).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toEntity(&symbol), nil
}

// Create 建立新股票代號
func (r *postgresStockSymbolsRepository) Create(ctx context.Context, symbol *entity.StockSymbol) error {
	r.logger.Info("Creating stock symbol", logger.String("symbol", symbol.Symbol), logger.String("market", symbol.Market))

	model := r.toModel(symbol)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		r.logger.Error("Failed to create stock symbol", logger.Error(err), logger.String("symbol", symbol.Symbol))
		return err
	}

	r.logger.Info("Stock symbol created successfully", logger.String("symbol", symbol.Symbol))
	return nil
}

func (r *postgresStockSymbolsRepository) Update(ctx context.Context, symbol *entity.StockSymbol) error {
	r.logger.Info("Updating stock symbol", logger.Any("id", symbol.ID), logger.String("symbol", symbol.Symbol))

	model := r.toModel(symbol)
	err := r.db.WithContext(ctx).Model(&models.StockSymbol{}).
		Where("id = ?", symbol.ID).
		Updates(model).Error

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
		entities = append(entities, r.toEntity(symbol))
	}

	return entities, nil
}

// BatchCreate 批次建立股票代號
func (r *postgresStockSymbolsRepository) BatchCreate(ctx context.Context, symbols []*entity.StockSymbol) error {
	r.logger.Info("Batch creating stock symbols", logger.Int("count", len(symbols)))

	modelSymbols := make([]*models.StockSymbol, len(symbols))
	for i, symbol := range symbols {
		modelSymbols[i] = r.toModel(symbol)
	}

	err := r.db.WithContext(ctx).CreateInBatches(modelSymbols, 100).Error
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

	// Deduplication map to prevent duplicates in the same batch
	deduplicatedMap := make(map[string]bool)
	modelSymbols := make([]models.StockSymbol, 0, len(symbols))

	for _, symbol := range symbols {
		key := fmt.Sprintf("%s|%s", symbol.Symbol, symbol.Market)
		if _, exists := deduplicatedMap[key]; exists {
			continue // Skip duplicate
		}
		deduplicatedMap[key] = true

		modelSymbols = append(modelSymbols, *r.toModel(symbol))
	}

	if len(modelSymbols) == 0 {
		return 0, 0, nil
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
