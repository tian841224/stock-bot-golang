package repository

import (
	"context"
	"fmt"

	"github.com/tian841224/stock-bot/internal/domain/entity"
	repo "github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"

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
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, symbol := range symbols {
		// 嘗試查找現有記錄
		existingSymbol, err := r.getBySymbolAndMarketTx(tx, symbol.Symbol, symbol.Market)
		if err != nil && err != gorm.ErrRecordNotFound {
			errorCount++
			continue
		}

		if existingSymbol != nil {
			// 更新現有記錄
			if r.shouldUpdate(existingSymbol, &models.StockSymbol{
				Symbol: symbol.Symbol,
				Market: symbol.Market,
				Name:   symbol.Name,
			}) {
				existingSymbol.Name = symbol.Name
				existingSymbol.Market = symbol.Market
				if err := tx.Save(&existingSymbol).Error; err != nil {
					errorCount++
					continue
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}

	return successCount, errorCount, nil
}

// getBySymbolAndMarketTx 在交易中根據股票代號和市場取得資料
func (r *postgresStockSymbolsRepository) getBySymbolAndMarketTx(tx *gorm.DB, symbol, market string) (*models.StockSymbol, error) {
	var symbolData models.StockSymbol
	err := tx.Where("symbol = ? AND market = ?", symbol, market).First(&symbolData).Error
	if err != nil {
		return nil, err
	}
	return &symbolData, nil
}

// shouldUpdate 判斷是否需要更新股票資料
func (r *postgresStockSymbolsRepository) shouldUpdate(existing, new *models.StockSymbol) bool {
	return existing.Name != new.Name || existing.Market != new.Market
}
