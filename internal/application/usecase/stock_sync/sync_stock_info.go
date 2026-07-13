package stock_sync

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type StockSyncUsecase interface {
	SyncTaiwanStockInfo(ctx context.Context) error
	SyncUSStockInfo(ctx context.Context) error
	SyncTaiwanStockTradingDate(ctx context.Context) error
	GetSyncStats(ctx context.Context) (map[string]int, error)
}

type stockSyncUsecase struct {
	stockSymbolRepo   port.StockSymbolRepository
	stockInfoProvider port.StockInfoProvider
	syncMetadataRepo  port.SyncMetadataRepository
	tradeDateRepo     port.TradeDateRepository
	logger            logger.Logger
}

func NewStockSyncUsecase(
	stockSymbolRepo port.StockSymbolRepository,
	stockInfoProvider port.StockInfoProvider,
	syncMetadataRepo port.SyncMetadataRepository,
	tradeDateRepo port.TradeDateRepository,
	log logger.Logger,
) StockSyncUsecase {
	return &stockSyncUsecase{
		stockSymbolRepo:   stockSymbolRepo,
		stockInfoProvider: stockInfoProvider,
		syncMetadataRepo:  syncMetadataRepo,
		tradeDateRepo:     tradeDateRepo,
		logger:            log,
	}
}

func (s *stockSyncUsecase) SyncTaiwanStockInfo(ctx context.Context) error {
	syncID := fmt.Sprintf("tw-sync-%d", time.Now().UnixNano())
	log := logger.FromContext(ctx, s.logger).With(logger.String("sync_id", syncID))
	log.Info("start syncing TW stock symbols")
	now := time.Now()

	symbols, err := s.stockInfoProvider.GetTaiwanStockInfo(ctx)
	if err != nil {
		log.Error("failed to get TW stock symbols", logger.Error(err))
		s.updateSyncMetadata(ctx, "TW", &now, nil, 0, err.Error())
		return err
	}

	log.Info("fetched TW stock symbols", logger.Int("count", len(symbols)))

	successCount, errorCount, err := s.asyncBatchUpsert(ctx, symbols)
	if err != nil {
		log.Error("batch upsert TW stock symbols failed", logger.Error(err))
		s.updateSyncMetadata(ctx, "TW", &now, nil, successCount, err.Error())
		return err
	}

	log.Info("TW stock symbol sync completed",
		logger.Int("success", successCount),
		logger.Int("failed", errorCount),
		logger.Int("total", len(symbols)))

	s.updateSyncMetadata(ctx, "TW", &now, &now, successCount, "")

	return nil
}

func (s *stockSyncUsecase) SyncUSStockInfo(ctx context.Context) error {
	syncID := fmt.Sprintf("us-sync-%d", time.Now().UnixNano())
	log := logger.FromContext(ctx, s.logger).With(logger.String("sync_id", syncID))
	log.Info("start syncing US stock symbols")
	now := time.Now()

	symbols, err := s.stockInfoProvider.GetUSStockInfo(ctx)
	if err != nil {
		log.Error("failed to get US stock symbols", logger.Error(err))
		s.updateSyncMetadata(ctx, "US", &now, nil, 0, err.Error())
		return err
	}

	log.Info("fetched US stock symbols", logger.Int("count", len(symbols)))

	successCount, errorCount, err := s.asyncBatchUpsert(ctx, symbols)
	if err != nil {
		log.Error("batch upsert US stock symbols failed", logger.Error(err))
		s.updateSyncMetadata(ctx, "US", &now, nil, successCount, err.Error())
		return err
	}

	log.Info("US stock symbol sync completed",
		logger.Int("success", successCount),
		logger.Int("failed", errorCount),
		logger.Int("total", len(symbols)))

	s.updateSyncMetadata(ctx, "US", &now, &now, successCount, "")

	return nil
}

func (s *stockSyncUsecase) SyncTaiwanStockTradingDate(ctx context.Context) error {
	syncID := fmt.Sprintf("tw-date-sync-%d", time.Now().UnixNano())
	log := logger.FromContext(ctx, s.logger).With(logger.String("sync_id", syncID))
	log.Info("start syncing TW trade dates")

	tradeDates, err := s.stockInfoProvider.GetTaiwanStockTradingDate(ctx)
	if err != nil {
		log.Error("failed to get TW trade dates", logger.Error(err))
		return err
	}

	log.Info("fetched TW trade dates", logger.Int("count", len(tradeDates)))

	err = s.tradeDateRepo.BatchCreateTradeDates(ctx, tradeDates)
	if err != nil {
		log.Error("failed to batch create trade dates", logger.Error(err))
		return err
	}

	log.Info("TW trade date sync completed",
		logger.Int("total", len(tradeDates)))

	return nil
}

func (s *stockSyncUsecase) GetSyncStats(ctx context.Context) (map[string]int, error) {
	return s.stockSymbolRepo.GetMarketStats(ctx)
}

func (s *stockSyncUsecase) asyncBatchUpsert(ctx context.Context, symbols []*entity.StockSymbol) (totalSuccess, totalError int, err error) {
	const (
		batchSize  = 100
		maxWorkers = 5
	)

	// 全域去重：確保不同批次間也不會有重複的資料
	uniqueSymbolsMap := make(map[string]*entity.StockSymbol)
	for _, s := range symbols {
		key := s.Symbol + "|" + s.Market
		if _, exists := uniqueSymbolsMap[key]; !exists {
			uniqueSymbolsMap[key] = s
		}
	}

	// 重建去重後的切片
	var uniqueSymbols []*entity.StockSymbol
	for _, s := range uniqueSymbolsMap {
		uniqueSymbols = append(uniqueSymbols, s)
	}

	// 使用去重後的資料進行批次切分
	batches := s.splitIntoBatches(uniqueSymbols, batchSize)
	s.logger.Info("start async batch upsert",
		logger.Int("original_count", len(symbols)),
		logger.Int("deduplicated_count", len(uniqueSymbols)),
		logger.Int("batch_count", len(batches)),
		logger.Int("worker_count", maxWorkers))

	batchChan := make(chan []*entity.StockSymbol, len(batches))
	resultChan := make(chan batchResult, len(batches))

	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go s.worker(ctx, i+1, batchChan, resultChan, &wg)
	}

	go func() {
		for _, batch := range batches {
			batchChan <- batch
		}
		close(batchChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.err != nil {
			s.logger.Warn("batch processing failed",
				logger.Int("batch_id", result.batchID),
				logger.Error(result.err))
		}
		totalSuccess += result.successCount
		totalError += result.errorCount
	}

	s.logger.Info("async batch upsert completed",
		logger.Int("success", totalSuccess),
		logger.Int("failed", totalError))

	return totalSuccess, totalError, nil
}

type batchResult struct {
	batchID      int
	successCount int
	errorCount   int
	err          error
}

func (s *stockSyncUsecase) worker(ctx context.Context, workerID int, batchChan <-chan []*entity.StockSymbol, resultChan chan<- batchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	batchID := 0
	for batch := range batchChan {
		select {
		case <-ctx.Done():
			resultChan <- batchResult{
				batchID: batchID,
				err:     ctx.Err(),
			}
			return
		default:
		}

		batchID++
		s.logger.Debug("worker processing batch",
			logger.Int("worker_id", workerID),
			logger.Int("batch_id", batchID),
			logger.Int("batch_size", len(batch)))

		successCount, errorCount, err := s.processBatch(ctx, batch)

		resultChan <- batchResult{
			batchID:      batchID,
			successCount: successCount,
			errorCount:   errorCount,
			err:          err,
		}

		s.logger.Debug("worker finished batch",
			logger.Int("worker_id", workerID),
			logger.Int("batch_id", batchID),
			logger.Int("success", successCount),
			logger.Int("failed", errorCount))
	}
}

// processBatch 執行單一批次的 upsert，並攔截 panic 避免整個 worker 崩潰。
// 發生 panic 時將該批次視為全部失敗並回傳錯誤，讓 worker 能繼續處理後續批次。
func (s *stockSyncUsecase) processBatch(ctx context.Context, batch []*entity.StockSymbol) (successCount, errorCount int, err error) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic recovered in batch upsert", logger.Any("recover", r))
			successCount = 0
			errorCount = len(batch)
			err = fmt.Errorf("批次處理發生 panic: %v", r)
		}
	}()
	return s.stockSymbolRepo.BatchUpsert(ctx, batch)
}

func (s *stockSyncUsecase) splitIntoBatches(symbols []*entity.StockSymbol, batchSize int) [][]*entity.StockSymbol {
	var batches [][]*entity.StockSymbol

	for i := 0; i < len(symbols); i += batchSize {
		end := i + batchSize
		if end > len(symbols) {
			end = len(symbols)
		}
		batches = append(batches, symbols[i:end])
	}

	return batches
}

func (s *stockSyncUsecase) updateSyncMetadata(ctx context.Context, market string, lastSyncAt, lastSuccessAt *time.Time, totalCount int, errorMsg string) {
	var lastError *string
	if errorMsg != "" {
		lastError = &errorMsg
	}

	metadata := &entity.SyncMetadata{
		Market:        market,
		LastSyncAt:    lastSyncAt,
		LastSuccessAt: lastSuccessAt,
		LastError:     lastError,
		TotalCount:    totalCount,
	}

	if err := s.syncMetadataRepo.Upsert(ctx, metadata); err != nil {
		s.logger.Error("failed to update sync metadata", logger.Error(err), logger.String("market", market))
	}
}
