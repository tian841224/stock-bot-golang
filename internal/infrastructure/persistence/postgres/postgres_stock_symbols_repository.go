package repository

import (
	"context"
	"fmt"

	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

	"gorm.io/gorm"
)

type postgresStockSymbolsRepository struct {
	db *gorm.DB
}

var _ repo.StockSymbolReader = (*postgresStockSymbolsRepository)(nil)
var _ repo.StockSymbolWriter = (*postgresStockSymbolsRepository)(nil)

func NewSymbolRepository(db *gorm.DB) *postgresStockSymbolsRepository {
	return &postgresStockSymbolsRepository{db: db}
}

// GetByID 根據 ID 取得股票代號
func (r *postgresStockSymbolsRepository) GetByID(ctx context.Context, id uint) (*entity.StockSymbol, error) {
	var symbol models.StockSymbol
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&symbol).Error
	if err != nil {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
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
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
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
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &entity.StockSymbol{
		ID:     symbolData.ID,
		Symbol: symbolData.Symbol,
		Market: symbolData.Market,
		Name:   symbolData.Name,
	}, nil
}

// GetBySubscriptionID 根據訂閱 ID 取得股票代號列表
func (r *postgresStockSymbolsRepository) GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*entity.StockSymbol, error) {
	var symbols []*models.StockSymbol
	err := r.db.WithContext(ctx).Where("subscription_id = ?", subscriptionID).Find(&symbols).Error
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

// GetBySubscriptionAndSymbol 根據訂閱和股票取得股票代號列表
func (r *postgresStockSymbolsRepository) GetBySubscriptionAndSymbol(ctx context.Context, subscriptionID, symbolID uint) (*entity.StockSymbol, error) {
	var symbol models.StockSymbol
	err := r.db.WithContext(ctx).Where("subscription_id = ? AND symbol_id = ?", subscriptionID, symbolID).First(&symbol).Error
	if err != nil {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &entity.StockSymbol{
		ID:     symbol.ID,
		Symbol: symbol.Symbol,
		Market: symbol.Market,
		Name:   symbol.Name,
	}, nil
}

// Create 建立新股票代號
func (r *postgresStockSymbolsRepository) Create(ctx context.Context, symbol *entity.StockSymbol) error {
	err := r.db.WithContext(ctx).Create(&models.StockSymbol{
		Symbol: symbol.Symbol,
		Market: symbol.Market,
		Name:   symbol.Name,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresStockSymbolsRepository) Update(ctx context.Context, symbol *entity.StockSymbol) error {
	err := r.db.WithContext(ctx).Model(&models.StockSymbol{}).
		Where("id = ?", symbol.ID).
		Updates(map[string]interface{}{
			"symbol": symbol.Symbol,
			"market": symbol.Market,
			"name":   symbol.Name,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

// Delete 刪除股票代號
func (r *postgresStockSymbolsRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&models.StockSymbol{}, id).Error
	if err != nil {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("symbol not found: %w", err)
	}

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
	err := r.db.WithContext(ctx).CreateInBatches(symbols, 100).Error
	if err != nil {
		return err
	}
	return nil
}

// BatchUpsert 批次更新或建立股票代號
func (r *postgresStockSymbolsRepository) BatchUpsert(ctx context.Context, symbols []*entity.StockSymbol) (successCount, errorCount int, err error) {
	for _, symbol := range symbols {
		var existingSymbol models.StockSymbol
		err := r.db.WithContext(ctx).Where("symbol = ? AND market = ?", symbol.Symbol, symbol.Market).First(&existingSymbol).Error

		if err == gorm.ErrRecordNotFound {
			newSymbol := models.StockSymbol{
				Symbol: symbol.Symbol,
				Market: symbol.Market,
				Name:   symbol.Name,
			}
			if err := r.db.WithContext(ctx).Create(&newSymbol).Error; err != nil {
				errorCount++
				continue
			}
			successCount++
		} else if err != nil {
			errorCount++
			continue
		} else {
			// 記錄已存在，跳過但計入成功
			successCount++
		}
	}

	return successCount, errorCount, nil
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
