# AGENTS.md - Stock Bot 開發指南

本文件提供 Stock Bot 專案的編碼代理程式開發的全面指南。它涵蓋建置命令、測試程序、程式碼風格慣例以及開發最佳實務。

## 建置/檢查/測試命令

### 建置應用程式
```bash
# 建置主要應用程式（Bot 服務）
go build ./cmd/bot

# 建置股票同步應用程式
go build ./cmd/sync_stock_info

# 建置通知排程應用程式
go build ./cmd/notification_stock_info

# 建置所有套件
go build ./...
```

### 程式碼格式化和檢查
```bash
# 格式化所有 Go 程式碼（標準 Go 格式化）
go fmt ./...

# 檢查潛在問題（檢查）
go vet ./...

# 格式化和檢查匯入
goimports -w .
```

### 測試
```bash
# 執行所有測試
go test ./...

# 以詳細輸出執行測試
go test -v ./...

# 為特定套件執行測試
go test ./internal/application/usecase/bot
go test ./internal/application/usecase/stock
go test ./internal/application/usecase/health
go test ./internal/application/usecase/notification
go test ./internal/application/usecase/stock_sync

# 執行單一測試函式
go test -v -run TestBotCommandUsecase_ProcessCommand ./internal/application/usecase/bot
go test -v -run TestValidateStock ./internal/application/usecase/stock
go test -v -run TestSendNotificationUsecase ./internal/application/usecase/notification

# 以涵蓋率執行測試
go test -cover ./...

# 產生涵蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 相依性
```bash
# 下載相依性
go mod download

# 整理 go.mod 和 go.sum
go mod tidy

# 驗證相依性
go mod verify
```

### Docker 相關命令
```bash
# 建置 Docker 映像檔
docker build -t stock-bot-go:latest -f Dockerfile .
docker build -t stock-bot-sync-go:latest -f Dockerfile.sync .
docker build -t stock-bot-scheduler-go:latest -f Dockerfile.scheduler .

# 使用 docker-compose 啟動服務（本地開發）
docker-compose up -d

# 使用 docker-compose 啟動服務（CI/CD 部署）
docker-compose -f docker-compose_cicd.yml up -d

# 查看服務狀態
docker-compose ps

# 查看服務日誌
docker-compose logs -f stock-bot-go
docker-compose logs -f sync-stock-info-go
docker-compose logs -f scheduler-go

# 停止服務
docker-compose down

# 重新建置並啟動服務
docker-compose up -d --build
```

## 程式碼風格指南

### 專案架構
此專案遵循 **Clean Architecture** 原則，並結合 **Hexagonal Architecture** 的適配器模式：
- **領域層**：`internal/domain/` - 核心業務實體、值物件、領域服務和領域錯誤
- **應用層**：`internal/application/` - 使用案例、DTO 和連接埠（介面）
- **基礎設施層**：`internal/infrastructure/` - 外部服務、資料庫、API適配器、配置和記錄
- **介面層**：`internal/interfaces/` - HTTP 處理器、訊息處理器和展示器

### 命名慣例

#### 套件
- 使用小寫、單字名稱（例如：`entity`、`usecase`、`port`）
- 遵循 Go 命名慣例
- 為私有套件使用 `internal/`

#### 型別和結構
```go
// 實體
type User struct {
    ID        uint
    AccountID string
    UserType  valueobject.UserType
    Status    bool
}

// DTO（資料傳輸物件）
type StockPrice struct {
    Symbol    string
    Price     float64
    Timestamp time.Time
}

// 介面（連接埠）
type BotCommandPort interface {
    GetDailyMarketInfo(ctx context.Context, userType valueobject.UserType, count int) (string, error)
    GetStockPerformance(ctx context.Context, userType valueobject.UserType, symbol string) (string, error)
}
```

#### 函式和方法
- 為匯出函式使用 PascalCase
- 為未匯出函式使用 camelCase
- 接收器名稱應簡短（1-2 個字元）
```go
func (u *User) IsActive() bool {
    return u.Status
}

func (uc *botCommandUsecase) processCommand(ctx context.Context, command string) (string, error) {
    // 實作邏輯
}
```

#### 變數和常數
```go
// 常數
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// 變數
var (
    defaultLogger Logger
)

// 區域變數
userID := "12345"
stockSymbol := "AAPL"
```

### 匯入
- 依標準函式庫、第三方和內部套件分組匯入
- 在群組之間使用空白行
- 移除未使用的匯入

```go
import (
    "context"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "github.com/tian841224/stock-bot/internal/application/dto"
    "github.com/tian841224/stock-bot/internal/domain/valueobject"
)
```

### 錯誤處理
- 明確傳回錯誤，不要恐慌
- 為業務邏輯使用自訂領域錯誤
- 以適當的上下文記錄錯誤

```go
// 領域錯誤
var (
    ErrNotFound       = errors.New("resource not found")
    ErrInvalidArgument = errors.New("invalid argument")
    ErrAlreadyExists   = errors.New("resource already exists")
)

// 錯誤處理模式
func (uc *userUsecase) GetUser(ctx context.Context, userID uint) (*entity.User, error) {
    user, err := uc.userRepo.GetByID(ctx, userID)
    if err != nil {
        uc.logger.Error("Failed to get user", zap.Error(err), zap.Uint("userID", userID))
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    if user == nil {
        return nil, domain.ErrNotFound
    }

    return user, nil
}
```

### 記錄
- 使用 Zap 進行結構化記錄
- 包含相關的上下文欄位
- 使用適當的記錄層級

```go
// 記錄器介面使用方式
logger.Info("User login successful",
    zap.String("accountID", accountID),
    zap.String("userType", string(userType)))

logger.Error("Database connection failed",
    zap.Error(err),
    zap.String("database", dbName))
```

### 測試
- 為多個測試案例使用表格驅動測試
- 模擬外部相依性
- 測試成功和錯誤情境
- 使用描述性的測試名稱

```go
func TestBotCommandUsecase_ProcessCommand(t *testing.T) {
    tests := []struct {
        name     string
        command  string
        userType valueobject.UserType
        want     string
        wantErr  bool
    }{
        {
            name:     "valid stock price command",
            command:  "/price AAPL",
            userType: valueobject.PremiumUser,
            want:     "AAPL: $150.00",
            wantErr:  false,
        },
        {
            name:     "invalid command",
            command:  "/invalid",
            userType: valueobject.FreeUser,
            want:     "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 測試實作
        })
    }
}
```

### 註解
- 使用中文撰寫註解（遵循專案慣例）
- 使用 `//` 作為單行註解
- 記錄匯出的函式和型別
- 新增套件層級文件

```go
// User 使用者實體
type User struct {
    ID        uint
    AccountID string // 帳號 ID
    UserType  valueobject.UserType // 使用者類型
    Status    bool // 啟用狀態
}

// GetUser 根據 ID 取得使用者
func (uc *userUsecase) GetUser(ctx context.Context, userID uint) (*entity.User, error) {
    // 實作邏輯
}
```

### 檔案組織
- 盡可能每個檔案一個型別
- 將相關型別分組在一起
- 使用描述性的檔案名稱

```
internal/
├── domain/
│   ├── entity/
│   │   ├── user.go
│   │   ├── stock.go
│   │   ├── subscription.go
│   │   ├── subscription_symbol.go
│   │   ├── stock_symbol.go
│   │   ├── trade_date.go
│   │   └── feature.go
│   ├── valueobject/
│   │   ├── user_type.go
│   │   ├── subscription_type.go
│   │   └── market_type.go
│   ├── service/
│   │   └── domain_service.go
│   └── error/
│       └── error.go
├── application/
│   ├── usecase/
│   │   ├── bot/
│   │   │   ├── bot_command.go
│   │   │   └── bot_command_test.go
│   │   ├── stock/
│   │   │   ├── generate_chart.go
│   │   │   ├── market_data.go
│   │   │   └── generate_chart_test.go
│   │   ├── health/
│   │   │   ├── health_check.go
│   │   │   └── health_check_test.go
│   │   ├── notification/
│   │   │   ├── send_notification.go
│   │   │   └── schedule_handler.go
│   │   ├── stock_sync/
│   │   │   └── sync_stock.go
│   │   └── user/
│   │       └── user_management.go
│   ├── dto/
│   │   ├── stock_price.go
│   │   ├── user_subscription.go
│   │   └── notification.go
│   └── port/
│       ├── bot_command_port.go
│       ├── repository.go
│       ├── notification_port.go
│       └── market_data_port.go
├── infrastructure/
│   ├── adapter/
│   │   ├── formatter/
│   │   │   ├── telegram_formatter.go
│   │   │   ├── line_formatter.go
│   │   │   └── formatter_adapter.go
│   │   ├── market/
│   │   │   ├── market_data_gateway.go
│   │   │   └── market_chart_gateway.go
│   │   └── presenter/
│   │       └── validation_gateway.go
│   ├── config/
│   │   └── config.go
│   ├── external/
│   │   ├── bot/
│   │   │   ├── telegram/
│   │   │   └── line/
│   │   └── stock/
│   │       ├── twse/
│   │       ├── cnyes/
│   │       ├── fugle/
│   │       └── finmindtrade/
│   ├── logging/
│   │   └── logger.go
│   └── persistence/
│       ├── model/
│       │   └── gorm_models.go
│       └── postgres/
│           ├── user_repository.go
│           ├── subscription_repository.go
│           ├── stock_symbols_repository.go
│           ├── subscription_symbol_repository.go
│           └── trade_date_repository.go
└── interfaces/
    ├── bot/
    │   ├── telegram/
    │   │   └── handler.go
    │   └── line/
    │       └── handler.go
    ├── health/
    │   └── handler.go
    └── presenter/
        └── response.go

cmd/
├── bot/
│   └── main.go                    # 主要 Bot 服務
├── sync_stock_info/
│   └── main.go                    # 股票資料同步服務
└── notification_stock_info/
    └── main.go                    # 排程通知服務
```

### 額外的開發工具
```bash
# 檢查程式碼品質
goimports -w .

# 執行效能測試
go test -bench=. ./...

# 生成測試覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 執行 golangci-lint（如果已安裝）
golangci-lint run ./...
```

## Cursor 規則
- 遵循 Clean Architecture 原則
- 維持目前的專案結構設計

## 開發工作流程

### 提交之前
1. 執行測試：`go test ./...`
2. 格式化程式碼：`go fmt ./...`
3. 檢查問題：`go vet ./...`
4. 確保沒有編譯錯誤：`go build ./...`
5. 整理相依性：`go mod tidy`

### 新增功能
1. 從領域層開始（實體、值物件）
2. 定義應用程式連接埠（介面）
3. 實作使用案例
4. 新增基礎設施適配器
5. 建立介面處理器
6. 撰寫全面的測試
7. 更新相關文件

### 新增排程任務
1. 在 `internal/application/usecase/notification/` 中實作業務邏輯
2. 使用 `ScheduleHandlerUsecase` 包裝排程任務
3. 在 `cmd/notification_stock_info/main.go` 中配置排程時間
4. 確保任務支援 context 取消機制
5. 添加適當的錯誤處理和日誌記錄
6. 測試排程任務的執行邏輯

### Docker 部署流程
1. 本地測試：`docker-compose up -d`
2. 檢查服務健康狀態：`docker-compose ps`
3. 查看日誌：`docker-compose logs -f [service-name]`
4. CI/CD 部署：使用 `docker-compose_cicd.yml`
5. 確保環境變數正確配置在 `.env` 檔案中

### 程式碼審查檢查清單
- [ ] Clean Architecture 原則已遵循
- [ ] 適當的錯誤處理
- [ ] 全面的測試涵蓋率
- [ ] 已套用程式碼格式化
- [ ] 沒有檢查問題
- [ ] 已新增適當的記錄
- [ ] 中文註解用於文件
- [ ] 介面相容性已維持
- [ ] Docker 配置已更新（如適用）
- [ ] 環境變數已記錄
- [ ] 資料庫遷移已處理（如適用）</content>
<parameter name="filePath">C:\Users\Tian\source\repos\stock-bot-clean\AGENTS.md