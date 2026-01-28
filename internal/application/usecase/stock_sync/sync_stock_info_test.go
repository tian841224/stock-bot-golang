package stock_sync

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tian841224/stock-bot/internal/domain/entity"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type mockStockSymbolRepo struct {
	batchUpsertFunc    func(ctx context.Context, symbols []*entity.StockSymbol) (int, int, error)
	getMarketStatsFunc func(ctx context.Context) (map[string]int, error)
}

func (m *mockStockSymbolRepo) GetByID(ctx context.Context, id uint) (*entity.StockSymbol, error) {
	return nil, nil
}

func (m *mockStockSymbolRepo) GetBySymbol(ctx context.Context, symbol string) (*entity.StockSymbol, error) {
	return nil, nil
}

func (m *mockStockSymbolRepo) GetBySymbolID(ctx context.Context, symbolID uint) ([]*entity.StockSymbol, error) {
	return nil, nil
}

func (m *mockStockSymbolRepo) GetBySubscriptionID(ctx context.Context, subscriptionID uint) ([]*entity.StockSymbol, error) {
	return nil, nil
}

func (m *mockStockSymbolRepo) GetBySubscriptionAndSymbol(ctx context.Context, subscriptionID, symbolID uint) (*entity.StockSymbol, error) {
	return nil, nil
}

func (m *mockStockSymbolRepo) GetMarketStats(ctx context.Context) (map[string]int, error) {
	if m.getMarketStatsFunc != nil {
		return m.getMarketStatsFunc(ctx)
	}
	return nil, nil
}

func (m *mockStockSymbolRepo) Create(ctx context.Context, stockSymbol *entity.StockSymbol) error {
	return nil
}

func (m *mockStockSymbolRepo) Update(ctx context.Context, stockSymbol *entity.StockSymbol) error {
	return nil
}

func (m *mockStockSymbolRepo) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *mockStockSymbolRepo) BatchUpsert(ctx context.Context, symbols []*entity.StockSymbol) (int, int, error) {
	if m.batchUpsertFunc != nil {
		return m.batchUpsertFunc(ctx, symbols)
	}
	return len(symbols), 0, nil
}

type mockStockInfoProvider struct {
	getTaiwanStockInfoFunc        func(ctx context.Context) ([]*entity.StockSymbol, error)
	getUSStockInfoFunc            func(ctx context.Context) ([]*entity.StockSymbol, error)
	getTaiwanStockTradingDateFunc func(ctx context.Context) ([]*entity.TradeDate, error)
}

func (m *mockStockInfoProvider) GetTaiwanStockInfo(ctx context.Context) ([]*entity.StockSymbol, error) {
	if m.getTaiwanStockInfoFunc != nil {
		return m.getTaiwanStockInfoFunc(ctx)
	}
	return nil, nil
}

func (m *mockStockInfoProvider) GetUSStockInfo(ctx context.Context) ([]*entity.StockSymbol, error) {
	if m.getUSStockInfoFunc != nil {
		return m.getUSStockInfoFunc(ctx)
	}
	return nil, nil
}

func (m *mockStockInfoProvider) GetTaiwanStockTradingDate(ctx context.Context) ([]*entity.TradeDate, error) {
	if m.getTaiwanStockTradingDateFunc != nil {
		return m.getTaiwanStockTradingDateFunc(ctx)
	}
	return nil, nil
}

type mockSyncMetadataRepo struct {
	getByMarketFunc func(ctx context.Context, market string) (*entity.SyncMetadata, error)
	upsertFunc      func(ctx context.Context, metadata *entity.SyncMetadata) error
}

func (m *mockSyncMetadataRepo) GetByMarket(ctx context.Context, market string) (*entity.SyncMetadata, error) {
	if m.getByMarketFunc != nil {
		return m.getByMarketFunc(ctx, market)
	}
	return nil, nil
}

func (m *mockSyncMetadataRepo) Upsert(ctx context.Context, metadata *entity.SyncMetadata) error {
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, metadata)
	}
	return nil
}

type mockTradeDateRepository struct {
	getByIDFunc               func(ctx context.Context, id uint) (*entity.TradeDate, error)
	getByDateFunc             func(ctx context.Context, date time.Time) (*entity.TradeDate, error)
	getByDateRangeFunc        func(ctx context.Context, startDate, endDate time.Time) ([]*entity.TradeDate, error)
	createFunc                func(ctx context.Context, tradeDate *entity.TradeDate) error
	batchCreateTradeDatesFunc func(ctx context.Context, tradeDates []*entity.TradeDate) error
}

func (m *mockTradeDateRepository) GetByID(ctx context.Context, id uint) (*entity.TradeDate, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTradeDateRepository) GetByDate(ctx context.Context, date time.Time) (*entity.TradeDate, error) {
	if m.getByDateFunc != nil {
		return m.getByDateFunc(ctx, date)
	}
	return nil, nil
}

func (m *mockTradeDateRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.TradeDate, error) {
	if m.getByDateRangeFunc != nil {
		return m.getByDateRangeFunc(ctx, startDate, endDate)
	}
	return nil, nil
}

func (m *mockTradeDateRepository) Create(ctx context.Context, tradeDate *entity.TradeDate) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, tradeDate)
	}
	return nil
}

func (m *mockTradeDateRepository) BatchCreateTradeDates(ctx context.Context, tradeDates []*entity.TradeDate) error {
	if m.batchCreateTradeDatesFunc != nil {
		return m.batchCreateTradeDatesFunc(ctx, tradeDates)
	}
	return nil
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Error(msg string, fields ...logger.Field) {}
func (m *mockLogger) Warn(msg string, fields ...logger.Field)  {}
func (m *mockLogger) Debug(msg string, fields ...logger.Field) {}
func (m *mockLogger) Panic(msg string, fields ...logger.Field) {}
func (m *mockLogger) Fatal(msg string, fields ...logger.Field) {}
func (m *mockLogger) Sync() error                              { return nil }

func TestStockSyncUsecase_SyncTaiwanStockInfo_Success(t *testing.T) {
	mockRepo := &mockStockSymbolRepo{
		batchUpsertFunc: func(ctx context.Context, symbols []*entity.StockSymbol) (int, int, error) {
			return len(symbols), 0, nil
		},
	}

	mockProvider := &mockStockInfoProvider{
		getTaiwanStockInfoFunc: func(ctx context.Context) ([]*entity.StockSymbol, error) {
			return []*entity.StockSymbol{
				{Symbol: "2330", Name: "台積電", Market: "TW"},
				{Symbol: "2317", Name: "鴻海", Market: "TW"},
			}, nil
		},
	}

	mockSyncRepo := &mockSyncMetadataRepo{}
	mockTradeDateRepo := &mockTradeDateRepository{}
	usecase := NewStockSyncUsecase(mockRepo, mockProvider, mockSyncRepo, mockTradeDateRepo, &mockLogger{})

	err := usecase.SyncTaiwanStockInfo(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestStockSyncUsecase_SyncTaiwanStockInfo_ProviderError(t *testing.T) {
	mockRepo := &mockStockSymbolRepo{}
	mockProvider := &mockStockInfoProvider{
		getTaiwanStockInfoFunc: func(ctx context.Context) ([]*entity.StockSymbol, error) {
			return nil, errors.New("API error")
		},
	}

	mockSyncRepo := &mockSyncMetadataRepo{}
	mockTradeDateRepo := &mockTradeDateRepository{}
	usecase := NewStockSyncUsecase(mockRepo, mockProvider, mockSyncRepo, mockTradeDateRepo, &mockLogger{})

	err := usecase.SyncTaiwanStockInfo(context.Background())
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStockSyncUsecase_SyncUSStockInfo_Success(t *testing.T) {
	mockRepo := &mockStockSymbolRepo{
		batchUpsertFunc: func(ctx context.Context, symbols []*entity.StockSymbol) (int, int, error) {
			return len(symbols), 0, nil
		},
	}

	mockProvider := &mockStockInfoProvider{
		getUSStockInfoFunc: func(ctx context.Context) ([]*entity.StockSymbol, error) {
			return []*entity.StockSymbol{
				{Symbol: "AAPL", Name: "Apple Inc.", Market: "US"},
				{Symbol: "GOOGL", Name: "Alphabet Inc.", Market: "US"},
			}, nil
		},
	}

	mockSyncRepo := &mockSyncMetadataRepo{}
	mockTradeDateRepo := &mockTradeDateRepository{}
	usecase := NewStockSyncUsecase(mockRepo, mockProvider, mockSyncRepo, mockTradeDateRepo, &mockLogger{})

	err := usecase.SyncUSStockInfo(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestStockSyncUsecase_SyncUSStockInfo_ProviderError(t *testing.T) {
	mockRepo := &mockStockSymbolRepo{}
	mockProvider := &mockStockInfoProvider{
		getUSStockInfoFunc: func(ctx context.Context) ([]*entity.StockSymbol, error) {
			return nil, errors.New("API error")
		},
	}

	mockSyncRepo := &mockSyncMetadataRepo{}
	mockTradeDateRepo := &mockTradeDateRepository{}
	usecase := NewStockSyncUsecase(mockRepo, mockProvider, mockSyncRepo, mockTradeDateRepo, &mockLogger{})

	err := usecase.SyncUSStockInfo(context.Background())
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestStockSyncUsecase_GetSyncStats_Success(t *testing.T) {
	expectedStats := map[string]int{
		"TW": 1800,
		"US": 8000,
	}

	mockRepo := &mockStockSymbolRepo{
		getMarketStatsFunc: func(ctx context.Context) (map[string]int, error) {
			return expectedStats, nil
		},
	}
	mockProvider := &mockStockInfoProvider{}

	mockSyncRepo := &mockSyncMetadataRepo{}
	mockTradeDateRepo := &mockTradeDateRepository{}
	usecase := NewStockSyncUsecase(mockRepo, mockProvider, mockSyncRepo, mockTradeDateRepo, &mockLogger{})

	stats, err := usecase.GetSyncStats(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if stats["TW"] != expectedStats["TW"] || stats["US"] != expectedStats["US"] {
		t.Errorf("Expected stats %v, got %v", expectedStats, stats)
	}
}

func TestStockSyncUsecase_AsyncBatchUpsert_PartialFailure(t *testing.T) {
	callCount := 0
	mockRepo := &mockStockSymbolRepo{
		batchUpsertFunc: func(ctx context.Context, symbols []*entity.StockSymbol) (int, int, error) {
			callCount++
			if callCount == 2 {
				return 0, len(symbols), errors.New("batch error")
			}
			return len(symbols), 0, nil
		},
	}

	mockProvider := &mockStockInfoProvider{
		getTaiwanStockInfoFunc: func(ctx context.Context) ([]*entity.StockSymbol, error) {
			symbols := make([]*entity.StockSymbol, 250)
			for i := 0; i < 250; i++ {
				symbols[i] = &entity.StockSymbol{
					Symbol: "TEST",
					Name:   "Test Stock",
					Market: "TW",
				}
			}
			return symbols, nil
		},
	}

	mockSyncRepo := &mockSyncMetadataRepo{}
	mockTradeDateRepo := &mockTradeDateRepository{}
	usecase := NewStockSyncUsecase(mockRepo, mockProvider, mockSyncRepo, mockTradeDateRepo, &mockLogger{})

	err := usecase.SyncTaiwanStockInfo(context.Background())
	if err != nil {
		t.Errorf("Expected no error (partial failure allowed), got %v", err)
	}
}

func TestStockSyncUsecase_AsyncBatchUpsert_ContextCancellation(t *testing.T) {
	mockRepo := &mockStockSymbolRepo{
		batchUpsertFunc: func(ctx context.Context, symbols []*entity.StockSymbol) (int, int, error) {
			return len(symbols), 0, nil
		},
	}

	mockProvider := &mockStockInfoProvider{
		getTaiwanStockInfoFunc: func(ctx context.Context) ([]*entity.StockSymbol, error) {
			symbols := make([]*entity.StockSymbol, 500)
			for i := 0; i < 500; i++ {
				symbols[i] = &entity.StockSymbol{
					Symbol: "TEST",
					Name:   "Test Stock",
					Market: "TW",
				}
			}
			return symbols, nil
		},
	}

	mockSyncRepo := &mockSyncMetadataRepo{}
	mockTradeDateRepo := &mockTradeDateRepository{}
	usecase := NewStockSyncUsecase(mockRepo, mockProvider, mockSyncRepo, mockTradeDateRepo, &mockLogger{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := usecase.SyncTaiwanStockInfo(ctx)
	if err != nil {
		t.Errorf("Expected no error (context cancellation handled by worker), got %v", err)
	}
}

func TestStockSyncUsecase_SplitIntoBatches(t *testing.T) {
	mockRepo := &mockStockSymbolRepo{}
	mockProvider := &mockStockInfoProvider{}
	mockSyncRepo := &mockSyncMetadataRepo{}
	mockTradeDateRepo := &mockTradeDateRepository{}
	usecase := &stockSyncUsecase{
		stockSymbolRepo:   mockRepo,
		stockInfoProvider: mockProvider,
		syncMetadataRepo:  mockSyncRepo,
		tradeDateRepo:     mockTradeDateRepo,
		logger:            &mockLogger{},
	}

	tests := []struct {
		name            string
		symbolsCount    int
		batchSize       int
		expectedBatches int
	}{
		{
			name:            "exactly divisible",
			symbolsCount:    200,
			batchSize:       100,
			expectedBatches: 2,
		},
		{
			name:            "with remainder",
			symbolsCount:    250,
			batchSize:       100,
			expectedBatches: 3,
		},
		{
			name:            "less than batch size",
			symbolsCount:    50,
			batchSize:       100,
			expectedBatches: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			symbols := make([]*entity.StockSymbol, tt.symbolsCount)
			for i := 0; i < tt.symbolsCount; i++ {
				symbols[i] = &entity.StockSymbol{
					Symbol: "TEST",
					Name:   "Test Stock",
					Market: "TW",
				}
			}

			batches := usecase.splitIntoBatches(symbols, tt.batchSize)
			if len(batches) != tt.expectedBatches {
				t.Errorf("Expected %d batches, got %d", tt.expectedBatches, len(batches))
			}

			totalSymbols := 0
			for _, batch := range batches {
				totalSymbols += len(batch)
			}
			if totalSymbols != tt.symbolsCount {
				t.Errorf("Expected total symbols %d, got %d", tt.symbolsCount, totalSymbols)
			}
		})
	}
}
