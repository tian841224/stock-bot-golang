# Bot 命令處理重構總結

## 重構目標

✅ **依照 Clean Architecture 設計，集中管理 Bot 命令處理邏輯**

## 變更內容

### 📁 新增檔案

| 檔案路徑 | 說明 |
|---------|------|
| `internal/application/usecase/bot/telegram_message_processor.go` | Telegram 訊息處理器，負責命令路由和參數驗證 |
| `internal/application/usecase/bot/line_message_processor.go` | LINE 訊息處理器，負責命令路由和參數驗證 |
| `internal/application/usecase/bot/line_bot_command.go` | LINE Bot 命令用例實作 |
| `docs/BOT_COMMAND_ARCHITECTURE.md` | Bot 命令架構詳細說明文檔 |
| `docs/BOT_COMMAND_FLOW.md` | Bot 命令處理流程圖和範例 |
| `docs/SUMMARY_BOT_REFACTORING.md` | 本文檔 |

### 📝 修改檔案

| 檔案路徑 | 修改內容 |
|---------|----------|
| `internal/interfaces/bot/telegram/handler.go` | 移除 `tgClient` 依賴，改為注入 `messageProcessor`，將命令處理邏輯委派給 processor |
| `internal/interfaces/bot/line/handler.go` | 移除 `botClient.HandleTextMessage()` 呼叫，改為使用 `messageProcessor.ProcessTextMessage()` |
| `internal/infrastructure/external/bot/line/client.go` | 移除對 `application/usecase/bot` 的循環依賴，刪除不必要的業務邏輯 |

## 架構設計

### 三層架構

```
┌─────────────────────────────────────┐
│   Interfaces Layer (介面層)          │
│   - telegram/handler.go             │
│   - line/handler.go                 │
│   職責：HTTP/Webhook 處理            │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│   Application Layer (應用層)         │
│   - telegram_message_processor.go   │
│   - line_message_processor.go       │
│   - tg_bot_command.go               │
│   - line_bot_command.go             │
│   - bot_command.go                  │
│   職責：命令路由、業務邏輯編排        │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│   Infrastructure Layer (基礎設施層)   │
│   - telegram/client.go              │
│   - line/client.go                  │
│   職責：Bot SDK 封裝、發送訊息        │
└─────────────────────────────────────┘
```

### 責任分離

#### 1. Handler (Interfaces Layer)
- ✅ 驗證 webhook secret token
- ✅ 解析 webhook JSON
- ✅ 回應 HTTP 200
- ✅ 啟動 goroutine 背景處理
- ✅ Panic recovery
- ❌ 不包含命令處理邏輯

#### 2. Message Processor (Application Layer)
- ✅ 提取訊息內容
- ✅ 確保使用者存在
- ✅ 解析命令和參數
- ✅ 命令路由（集中管理）
- ✅ 參數驗證
- ✅ 錯誤訊息發送
- ❌ 不直接呼叫 Domain 服務

#### 3. Command UseCase (Application Layer)
- ✅ 呼叫 Domain/Application 服務
- ✅ 格式化訊息內容
- ✅ 透過 Bot Client 發送訊息
- ✅ 共用業務邏輯（BotCommandUsecase）
- ❌ 不包含路由邏輯

#### 4. Bot Client (Infrastructure Layer)
- ✅ 封裝 Bot SDK
- ✅ 提供發送訊息功能
- ✅ 處理 API 錯誤
- ❌ 不包含業務邏輯

## 命令集中管理

所有命令都在 **Message Processor** 的 `routeCommand()` 方法中集中管理：

### Telegram

```go
// telegram_message_processor.go
func (p *TelegramMessageProcessor) routeCommand(ctx context.Context, command, arg1, arg2 string, chatID int64) error {
    switch command {
    case "/start":
        return p.tgCommandUsecase.GetUseGuideMessage(chatID)
    case "/k":
        return p.handleHistoricalCandles(ctx, chatID, arg1)
    case "/p":
        return p.handlePerformanceChart(ctx, chatID, arg1)
    case "/d":
        return p.handleStockPrice(ctx, chatID, arg1, arg2)
    case "/t":
        return p.tgCommandUsecase.GetTopVolumeStock(ctx, chatID)
    case "/i":
        return p.handleStockProfile(ctx, chatID, arg1)
    case "/r":
        return p.handleRevenueChart(ctx, chatID, arg1)
    case "/m":
        return p.handleDailyMarket(ctx, chatID, arg1)
    default:
        return p.handleUnknownCommand(chatID)
    }
}
```

### LINE

```go
// line_message_processor.go
func (p *LineMessageProcessor) routeCommand(ctx context.Context, command, arg1, arg2, replyToken string) error {
    switch command {
    case "/start":
        return p.lineCommandUsecase.GetUseGuideMessage(replyToken)
    case "/k":
        return p.handleHistoricalCandles(ctx, replyToken, arg1)
    case "/p":
        return p.handlePerformanceChart(ctx, replyToken, arg1)
    case "/d":
        return p.handleStockPrice(ctx, replyToken, arg1, arg2)
    case "/t":
        return p.lineCommandUsecase.GetTopVolumeStock(ctx, replyToken)
    case "/i":
        return p.handleStockProfile(ctx, replyToken, arg1)
    case "/r":
        return p.handleRevenueChart(ctx, replyToken, arg1)
    case "/m":
        return p.lineCommandUsecase.GetDailyMarketInfo(ctx, replyToken)
    default:
        return p.handleUnknownCommand(replyToken)
    }
}
```

## 修復的問題

### 1. 循環依賴問題

**問題**：
```
internal/infrastructure/external/bot/line/client.go
    ↓ import
internal/application/usecase/bot
    ↓ import
internal/infrastructure/external/bot/line
```

**解決方案**：
- 移除 `LineBotClient` 對 `application/usecase/bot` 的依賴
- 業務邏輯由 `LineBotCommandUsecase` 和 `LineMessageProcessor` 處理
- `LineBotClient` 只保留純粹的發送功能

### 2. 命令處理邏輯分散

**問題**：
- 原本的 `h.tgClient.ProcessUpdate(&u)` 方法不存在
- 沒有統一的命令路由機制
- 命令處理邏輯散落在各處

**解決方案**：
- 建立 `TelegramMessageProcessor` 和 `LineMessageProcessor`
- 集中所有命令路由邏輯在 `routeCommand()` 方法
- 清晰的命令 → 處理器映射

### 3. 職責不清晰

**問題**：
- Handler 混合了 HTTP 處理和命令路由
- Client 混合了 SDK 呼叫和業務邏輯

**解決方案**：
- Handler 只負責 HTTP 層面
- Message Processor 負責命令路由
- Command UseCase 負責業務邏輯
- Client 只負責 SDK 呼叫

## 優勢

### 1. 符合 Clean Architecture 原則 ✅

- **依賴方向正確**：外層依賴內層
- **關注點分離**：每層職責清晰
- **依賴反轉**：透過介面依賴抽象

### 2. 易於維護 ✅

- **命令集中管理**：所有命令都在 `routeCommand()` 方法中
- **單一修改點**：新增命令只需修改一處
- **清晰的結構**：容易找到對應的程式碼

### 3. 易於測試 ✅

- **獨立測試**：每一層都可以獨立測試
- **Mock 容易**：所有依賴都透過建構函式注入
- **單元測試**：可以針對每個方法進行測試

### 4. 易於擴展 ✅

- **新增平台**：只需實作新的 Handler、Processor 和 UseCase
- **共用邏輯**：BotCommandUsecase 避免重複程式碼
- **新增命令**：在 `routeCommand()` 中新增一個 case 即可

### 5. 錯誤處理完善 ✅

- **統一的錯誤處理**：所有錯誤都會記錄 log 並發送訊息給使用者
- **Panic recovery**：goroutine 中的 panic 會被捕捉
- **友善的錯誤訊息**：參數錯誤會告知正確的使用方式

## 測試建議

### Unit Test

```go
// Message Processor 測試
func TestTelegramMessageProcessor_RouteCommand(t *testing.T) {
    tests := []struct{
        name    string
        command string
        arg1    string
        arg2    string
        expect  string
    }{
        {"start command", "/start", "", "", "呼叫 GetUseGuideMessage"},
        {"stock price", "/d", "2330", "", "呼叫 GetStockPrice"},
        {"unknown command", "/xyz", "", "", "呼叫 handleUnknownCommand"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Test

```go
// Command UseCase 測試
func TestTgBotCommandUsecase_Integration(t *testing.T) {
    // Setup real dependencies
    cfg := config.LoadTestConfig()
    client, _ := tgbotapi.NewBot(cfg, logger)
    usecase := NewTgBotCommandUsecase(botCommandUsecase, client)
    
    // Test with real data
    err := usecase.GetStockPrice(context.Background(), "2330", nil, testChatID)
    assert.NoError(t, err)
}
```

## 下一步建議

### 1. 實作缺少的功能

- [ ] `GetStockProfile()` - 查詢公司資訊
- [ ] 訂閱功能 (`/sub`, `/unsub`, `/add`, `/del`, `/list`)
- [ ] 新聞功能 (`/n`)

### 2. 增加測試

- [ ] Message Processor 單元測試
- [ ] Command UseCase 單元測試
- [ ] Integration 測試
- [ ] E2E 測試

### 3. 錯誤處理優化

- [ ] 更詳細的錯誤訊息
- [ ] 錯誤分類（系統錯誤 vs 使用者錯誤）
- [ ] 錯誤監控和告警

### 4. 效能優化

- [ ] 增加 cache 機制
- [ ] 批次處理訊息
- [ ] 並發控制

### 5. 文檔補充

- [ ] API 文檔
- [ ] 部署文檔
- [ ] 故障排查指南

## 驗證清單

- [x] 移除循環依賴
- [x] 實作 Telegram Message Processor
- [x] 實作 LINE Message Processor
- [x] 實作 LINE Bot Command UseCase
- [x] 修改 Telegram Handler
- [x] 修改 LINE Handler
- [x] 清理 LINE Bot Client
- [x] 通過 linter 檢查
- [x] 建立架構文檔
- [x] 建立流程圖文檔
- [ ] 撰寫單元測試
- [ ] 撰寫整合測試
- [ ] 手動測試驗證

## 相關文檔

- [Bot 命令架構說明](./BOT_COMMAND_ARCHITECTURE.md) - 詳細的架構設計和實作說明
- [Bot 命令處理流程](./BOT_COMMAND_FLOW.md) - 流程圖和範例
- [Clean Architecture 遷移報告](../CLEAN_ARCHITECTURE_MIGRATION.md) - 整體架構重構說明

---

**重構日期**：2025-12-09  
**重構者**：Claude Code Assistant  
**架構風格**：Clean Architecture (Robert C. Martin)

