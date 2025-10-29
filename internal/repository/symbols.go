package repository

import (
	"github.com/tian841224/stock-bot/internal/db/models"

	"gorm.io/gorm"
)

type SymbolRepository interface {
	Create(symbol *models.Symbol) error
	GetByID(id uint) (*models.Symbol, error)
	GetBySymbolAndMarket(symbol, market string) (*models.Symbol, error)
	Update(symbol *models.Symbol) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.Symbol, error)
	GetByMarket(market string) ([]*models.Symbol, error)
	BatchCreate(symbols []*models.Symbol) error
	BatchUpsert(symbols []*models.Symbol) (successCount, errorCount int, err error)
	GetMarketStats() (map[string]int, error)
}

type symbolRepository struct {
	db *gorm.DB
}

func NewSymbolRepository(db *gorm.DB) SymbolRepository {
	return &symbolRepository{db: db}
}

// Create 建立新股票代號
func (r *symbolRepository) Create(symbol *models.Symbol) error {
	return r.db.Create(symbol).Error
}

// GetByID 根據 ID 取得股票代號
func (r *symbolRepository) GetByID(id uint) (*models.Symbol, error) {
	var symbol models.Symbol
	err := r.db.First(&symbol, id).Error
	if err != nil {
		return nil, err
	}
	return &symbol, nil
}

// GetBySymbolAndMarket 根據股票代號和市場取得資料
func (r *symbolRepository) GetBySymbolAndMarket(symbol, market string) (*models.Symbol, error) {
	var symbolData models.Symbol
	err := r.db.Where("symbol = ? AND market = ?", symbol, market).First(&symbolData).Error
	if err != nil {
		return nil, err
	}
	return &symbolData, nil
}

// Update 更新股票代號資料
func (r *symbolRepository) Update(symbol *models.Symbol) error {
	return r.db.Save(symbol).Error
}

// Delete 刪除股票代號
func (r *symbolRepository) Delete(id uint) error {
	return r.db.Delete(&models.Symbol{}, id).Error
}

// List 取得股票代號列表
func (r *symbolRepository) List(offset, limit int) ([]*models.Symbol, error) {
	var symbols []*models.Symbol
	err := r.db.Offset(offset).Limit(limit).Find(&symbols).Error
	return symbols, err
}

// GetByMarket 根據市場取得股票代號列表
func (r *symbolRepository) GetByMarket(market string) ([]*models.Symbol, error) {
	var symbols []*models.Symbol
	err := r.db.Where("market = ?", market).Find(&symbols).Error
	return symbols, err
}

// BatchCreate 批次建立股票代號
func (r *symbolRepository) BatchCreate(symbols []*models.Symbol) error {
	return r.db.CreateInBatches(symbols, 100).Error
}

// BatchUpsert 批次更新或建立股票代號
func (r *symbolRepository) BatchUpsert(symbols []*models.Symbol) (successCount, errorCount int, err error) {
	tx := r.db.Begin()
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
			if r.shouldUpdate(existingSymbol, symbol) {
				existingSymbol.Name = symbol.Name
				existingSymbol.Market = symbol.Market
				if err := tx.Save(existingSymbol).Error; err != nil {
					errorCount++
					continue
				}
			}
		} else {
			// 建立新記錄
			if err := tx.Create(symbol).Error; err != nil {
				errorCount++
				continue
			}
		}
		successCount++
	}

	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}

	return successCount, errorCount, nil
}

// GetMarketStats 取得各市場統計資訊
func (r *symbolRepository) GetMarketStats() (map[string]int, error) {
	stats := make(map[string]int)

	var results []struct {
		Market string
		Count  int64
	}

	err := r.db.Model(&models.Symbol{}).
		Select("market, count(*) as count").
		Group("market").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	total := 0
	for _, result := range results {
		stats[result.Market] = int(result.Count)
		total += int(result.Count)
	}
	stats["total"] = total

	return stats, nil
}

// getBySymbolAndMarketTx 在交易中根據股票代號和市場取得資料
func (r *symbolRepository) getBySymbolAndMarketTx(tx *gorm.DB, symbol, market string) (*models.Symbol, error) {
	var symbolData models.Symbol
	err := tx.Where("symbol = ? AND market = ?", symbol, market).First(&symbolData).Error
	if err != nil {
		return nil, err
	}
	return &symbolData, nil
}

// shouldUpdate 判斷是否需要更新股票資料
func (r *symbolRepository) shouldUpdate(existing, new *models.Symbol) bool {
	return existing.Name != new.Name || existing.Market != new.Market
}
