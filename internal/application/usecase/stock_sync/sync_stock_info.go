package stock_sync

import (
	"context"
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
	s.logger.Info("開始同步台灣股票資訊...")
	now := time.Now()

	symbols, err := s.stockInfoProvider.GetTaiwanStockInfo(ctx)
	if err != nil {
		s.logger.Error("取得台灣股票資訊失敗", logger.Error(err))
		s.updateSyncMetadata(ctx, "TW", &now, nil, 0, err.Error())
		return err
	}

	s.logger.Info("成功取得股票資訊", logger.Int("count", len(symbols)))

	successCount, errorCount, err := s.asyncBatchUpsert(ctx, symbols)
	if err != nil {
		s.logger.Error("批次更新股票資訊失敗", logger.Error(err))
		s.updateSyncMetadata(ctx, "TW", &now, nil, successCount, err.Error())
		return err
	}

	s.logger.Info("股票資訊同步完成",
		logger.Int("成功", successCount),
		logger.Int("失敗", errorCount),
		logger.Int("總計", len(symbols)))

	s.updateSyncMetadata(ctx, "TW", &now, &now, successCount, "")

	return nil
}

func (s *stockSyncUsecase) SyncUSStockInfo(ctx context.Context) error {
	s.logger.Info("開始同步美股股票資訊...")
	now := time.Now()

	symbols, err := s.stockInfoProvider.GetUSStockInfo(ctx)
	if err != nil {
		s.logger.Error("取得美股股票資訊失敗", logger.Error(err))
		s.updateSyncMetadata(ctx, "US", &now, nil, 0, err.Error())
		return err
	}

	s.logger.Info("成功取得股票資訊", logger.Int("count", len(symbols)))

	successCount, errorCount, err := s.asyncBatchUpsert(ctx, symbols)
	if err != nil {
		s.logger.Error("批次更新股票資訊失敗", logger.Error(err))
		s.updateSyncMetadata(ctx, "US", &now, nil, successCount, err.Error())
		return err
	}

	s.logger.Info("股票資訊同步完成",
		logger.Int("成功", successCount),
		logger.Int("失敗", errorCount),
		logger.Int("總計", len(symbols)))

	s.updateSyncMetadata(ctx, "US", &now, &now, successCount, "")

	return nil
}

func (s *stockSyncUsecase) SyncTaiwanStockTradingDate(ctx context.Context) error {
	s.logger.Info("開始同步台股交易日...")

	tradeDates, err := s.stockInfoProvider.GetTaiwanStockTradingDate(ctx)
	if err != nil {
		s.logger.Error("取得台股交易日失敗", logger.Error(err))
		return err
	}

	s.logger.Info("成功取得交易日", logger.Int("count", len(tradeDates)))

	err = s.tradeDateRepo.BatchCreateTradeDates(ctx, tradeDates)
	if err != nil {
		s.logger.Error("批次更新交易日失敗", logger.Error(err))
		return err
	}

	s.logger.Info("交易日同步完成",
		logger.Int("總計", len(tradeDates)))

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

	batches := s.splitIntoBatches(symbols, batchSize)
	s.logger.Info("開始非同步批次處理",
		logger.Int("總數量", len(symbols)),
		logger.Int("批次數", len(batches)),
		logger.Int("工作者數", maxWorkers))

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
			s.logger.Warn("批次處理失敗",
				logger.Int("批次ID", result.batchID),
				logger.Error(result.err))
		}
		totalSuccess += result.successCount
		totalError += result.errorCount
	}

	s.logger.Info("非同步批次處理完成",
		logger.Int("成功", totalSuccess),
		logger.Int("失敗", totalError))

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
		s.logger.Debug("工作者開始處理批次",
			logger.Int("工作者ID", workerID),
			logger.Int("批次ID", batchID),
			logger.Int("批次大小", len(batch)))

		successCount, errorCount, err := s.stockSymbolRepo.BatchUpsert(ctx, batch)

		resultChan <- batchResult{
			batchID:      batchID,
			successCount: successCount,
			errorCount:   errorCount,
			err:          err,
		}

		s.logger.Debug("工作者完成批次處理",
			logger.Int("工作者ID", workerID),
			logger.Int("批次ID", batchID),
			logger.Int("成功", successCount),
			logger.Int("失敗", errorCount))
	}
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
		s.logger.Error("更新同步元資料失敗", logger.Error(err), logger.String("market", market))
	}
}
