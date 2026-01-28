package repository

import (
	"github.com/tian841224/stock-bot/internal/domain/entity"
	models "github.com/tian841224/stock-bot/internal/infrastructure/persistence/model"
)

// ConversionHelpers 提供通用的資料轉換輔助函數

// convertUser 將 model User 轉換為 entity User
func convertUser(model *models.User) *entity.User {
	if model == nil {
		return nil
	}
	return &entity.User{
		ID:        model.ID,
		AccountID: model.AccountID,
		UserType:  model.UserType,
		Status:    model.Status,
	}
}

// convertFeature 將 model Feature 轉換為 entity Feature
func convertFeature(model *models.Feature) *entity.Feature {
	if model == nil {
		return nil
	}
	return &entity.Feature{
		ID:          model.ID,
		Name:        model.Name,
		Code:        model.Code,
		Description: model.Description,
	}
}

// convertStockSymbol 將 model StockSymbol 轉換為 entity StockSymbol
func convertStockSymbol(model *models.StockSymbol) *entity.StockSymbol {
	if model == nil {
		return nil
	}
	return &entity.StockSymbol{
		ID:     model.ID,
		Symbol: model.Symbol,
		Market: model.Market,
		Name:   model.Name,
	}
}

// convertStockSymbolSlice 批次轉換 StockSymbol 切片
func convertStockSymbolSlice(models []*models.StockSymbol) []*entity.StockSymbol {
	if models == nil {
		return nil
	}

	entities := make([]*entity.StockSymbol, 0, len(models))
	for _, model := range models {
		if model != nil {
			entities = append(entities, convertStockSymbol(model))
		}
	}
	return entities
}

// batchConvert 通用批次轉換函數
// 接受一個轉換函數，對切片中的每個元素應用該函數
func batchConvert[T any, R any](items []T, converter func(T) R) []R {
	if items == nil {
		return nil
	}

	results := make([]R, 0, len(items))
	for _, item := range items {
		results = append(results, converter(item))
	}
	return results
}
