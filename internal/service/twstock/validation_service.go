package twstock

// ========== 驗證相關方法 ==========

// ValidateStockID 驗證股票代號是否存在
func (s *StockService) ValidateStockID(stockID string) (bool, string, error) {
	// 先從資料庫查詢
	symbol, err := s.symbolsRepo.GetBySymbolAndMarket(stockID, "TW")
	if err == nil && symbol != nil {
		return true, symbol.Name, nil
	}

	return false, "", nil
}
