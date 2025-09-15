package stock_sync

import (
	"stock-bot/internal/db/models"
	"stock-bot/internal/infrastructure/finmindtrade"
	"stock-bot/internal/infrastructure/finmindtrade/dto"
	"stock-bot/internal/repository"
	"stock-bot/pkg/logger"
	"sync"

	"go.uber.org/zap"
)

type StockSyncService struct {
	symbolsRepo   repository.SymbolsRepository
	finmindClient finmindtrade.FinmindTradeAPIInterface
}

func NewStockSyncService(symbolsRepo repository.SymbolsRepository, finmindClient finmindtrade.FinmindTradeAPIInterface) *StockSyncService {
	return &StockSyncService{
		symbolsRepo:   symbolsRepo,
		finmindClient: finmindClient,
	}
}

// SyncTaiwanStockInfo 同步台灣股票資訊
func (s *StockSyncService) SyncTaiwanStockInfo() error {
	logger.Log.Info("開始同步台灣股票資訊...")

	// 呼叫 FinMind API
	requestDto := dto.FinmindtradeRequestDto{}
	response, err := s.finmindClient.GetTaiwanStockInfo(requestDto)
	if err != nil {
		logger.Log.Error("呼叫 FinMind API 失敗", zap.Error(err))
		return err
	}

	if response.Status != 200 {
		logger.Log.Error("FinMind API 回應錯誤",
			zap.Int("status", response.Status),
			zap.String("message", response.Msg))
		return nil // 不返回錯誤，避免程式中斷
	}

	logger.Log.Info("成功取得股票資訊", zap.Int("count", len(response.Data)))

	// 轉換為 models.Symbols
	symbols := make([]*models.Symbols, 0, len(response.Data))
	for _, stockInfo := range response.Data {
		symbol := &models.Symbols{
			Symbol: stockInfo.StockID,
			Name:   stockInfo.StockName,
			Market: stockInfo.Type,
		}
		symbols = append(symbols, symbol)
	}

	// 非同步批次處理
	successCount, errorCount, err := s.asyncBatchUpsert(symbols)
	if err != nil {
		logger.Log.Error("批次更新股票資訊失敗", zap.Error(err))
		return err
	}

	logger.Log.Info("股票資訊同步完成",
		zap.Int("成功", successCount),
		zap.Int("失敗", errorCount),
		zap.Int("總計", len(response.Data)))

	return nil
}

// asyncBatchUpsert 非同步批次更新或建立股票資訊
func (s *StockSyncService) asyncBatchUpsert(symbols []*models.Symbols) (totalSuccess, totalError int, err error) {
	const (
		batchSize  = 100 // 每個批次的大小
		maxWorkers = 5   // 最大工作者數量
	)

	// 將資料分割成批次
	batches := s.splitIntoBatches(symbols, batchSize)
	logger.Log.Info("開始非同步批次處理",
		zap.Int("總數量", len(symbols)),
		zap.Int("批次數", len(batches)),
		zap.Int("工作者數", maxWorkers))

	// 建立通道
	batchChan := make(chan []*models.Symbols, len(batches))
	resultChan := make(chan batchResult, len(batches))

	// 啟動工作者
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go s.worker(i+1, batchChan, resultChan, &wg)
	}

	// 發送批次到通道
	go func() {
		for _, batch := range batches {
			batchChan <- batch
		}
		close(batchChan)
	}()

	// 等待所有工作者完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集結果
	for result := range resultChan {
		if result.err != nil {
			logger.Log.Warn("批次處理失敗",
				zap.Int("批次ID", result.batchID),
				zap.Error(result.err))
			// 不直接返回錯誤，繼續處理其他批次
		}
		totalSuccess += result.successCount
		totalError += result.errorCount
	}

	logger.Log.Info("非同步批次處理完成",
		zap.Int("成功", totalSuccess),
		zap.Int("失敗", totalError))

	return totalSuccess, totalError, nil
}

// batchResult 批次處理結果
type batchResult struct {
	batchID      int
	successCount int
	errorCount   int
	err          error
}

// worker 工作者函式
func (s *StockSyncService) worker(workerID int, batchChan <-chan []*models.Symbols, resultChan chan<- batchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	batchID := 0
	for batch := range batchChan {
		batchID++
		logger.Log.Debug("工作者開始處理批次",
			zap.Int("工作者ID", workerID),
			zap.Int("批次ID", batchID),
			zap.Int("批次大小", len(batch)))

		successCount, errorCount, err := s.symbolsRepo.BatchUpsert(batch)

		resultChan <- batchResult{
			batchID:      batchID,
			successCount: successCount,
			errorCount:   errorCount,
			err:          err,
		}

		logger.Log.Debug("工作者完成批次處理",
			zap.Int("工作者ID", workerID),
			zap.Int("批次ID", batchID),
			zap.Int("成功", successCount),
			zap.Int("失敗", errorCount))
	}
}

// splitIntoBatches 將資料分割成批次
func (s *StockSyncService) splitIntoBatches(symbols []*models.Symbols, batchSize int) [][]*models.Symbols {
	var batches [][]*models.Symbols

	for i := 0; i < len(symbols); i += batchSize {
		end := i + batchSize
		if end > len(symbols) {
			end = len(symbols)
		}
		batches = append(batches, symbols[i:end])
	}

	return batches
}

// GetSyncStats 取得同步統計資訊
func (s *StockSyncService) GetSyncStats() (map[string]int, error) {
	return s.symbolsRepo.GetMarketStats()
}
