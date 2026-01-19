# Bot 命令處理架構說明

## 架構概覽

本專案採用 Clean Architecture 設計，Bot 命令處理流程遵循以下層次：

```
interfaces/bot/          # 入口層 (Webhook Handler)
    ↓
application/usecase/bot/ # 應用層 (Message Processor + Command UseCase)
    ↓
infrastructure/external/ # 基礎設施層 (Bot Client)
```

## 架構層次說明

### 1. Interfaces Layer (介面層)

**位置**：`internal/interfaces/bot/`

**職責**：
- 接收來自外部的 webhook 請求
- 驗證請求的合法性（Secret Token）
- 解析請求 JSON
- 回應 HTTP 200
- 將請求轉發到 Application Layer

**實作檔案**：
- `telegram/handler.go` - Telegram webhook handler
- `line/handler.go` - LINE webhook handler


### 2. Application Layer (應用層)

#### 2.1 Message Processor (訊息處理器)

**位置**：`internal/application/usecase/bot/`

**職責**：
- 解析訊息內容，提取命令和參數
- 驗證使用者（確保使用者存在於資料庫）
- 命令路由（根據命令類型分派到對應的處理器）
- 統一的錯誤處理和 panic recovery
- 參數驗證（如：檢查是否為空、日期格式驗證）

**實作檔案**：
- `telegram_message_processor.go` - Telegram 訊息處理器
- `line_message_processor.go` - LINE 訊息處理器


#### 2.2 Command UseCase (命令用例)

**位置**：`internal/application/usecase/bot/`

**職責**：
- 呼叫 domain/application 層的業務邏輯
- 格式化訊息內容
- 透過 Bot Client 發送回應

**實作檔案**：
- `tg_bot_command.go` - Telegram 命令用例
- `line_bot_command.go` - LINE 命令用例
- `bot_command.go` - 共用的業務邏輯

**設計模式**：

```
TgBotCommandUsecase (介面)
    ↓ 實作
tgBotCommandUsecase (結構)
    ↓ 依賴
BotCommandUsecase (共用業務邏輯)
```

### 3. Infrastructure Layer (基礎設施層)

**位置**：`internal/infrastructure/external/bot/`

**職責**：
- 封裝第三方 Bot SDK
- 提供訊息發送功能（文字、圖片、按鈕等）
- 不包含業務邏輯

**實作檔案**：
- `telegram/client.go` - Telegram Bot Client
- `line/client.go` - LINE Bot Client


## 命令處理流程

### Telegram 流程圖

```
Telegram Server
    ↓ POST webhook
TgHandler.Webhook()
    ↓ 驗證 & 解析 & 回應 200
TelegramMessageProcessor.ProcessUpdate()
    ↓ 解析命令
TelegramMessageProcessor.routeCommand()
    ↓ 根據命令路由
TgBotCommandUsecase.GetXXX()
    ↓ 呼叫業務邏輯
BotCommandUsecase.GetXXX()
    ↓ 格式化 & 發送
TgBotClient.SendMessage() / SendPhoto()
    ↓
Telegram Server
```

### LINE 流程圖

```
LINE Server
    ↓ POST webhook
LineBotHandler.Webhook()
    ↓ 解析 & 回應 200
LineMessageProcessor.ProcessTextMessage()
    ↓ 解析命令
LineMessageProcessor.routeCommand()
    ↓ 根據命令路由
LineBotCommandUsecase.GetXXX()
    ↓ 呼叫業務邏輯
BotCommandUsecase.GetXXX()
    ↓ 格式化 & 發送
LineBotClient.ReplyMessage() / ReplyPhoto()
    ↓
LINE Server
```

## 新增命令指南

### 1. 在 Message Processor 新增命令路由

**檔案**：`telegram_message_processor.go` 或 `line_message_processor.go`

```go
func (p *TelegramMessageProcessor) routeCommand(ctx context.Context, command, arg1, arg2 string, chatID int64) error {
    switch command {
    // ... 其他命令
    case "/mynewcommand":
        return p.handleMyNewCommand(ctx, chatID, arg1)
    default:
        return p.handleUnknownCommand(chatID)
    }
}

// 實作命令處理函數
func (p *TelegramMessageProcessor) handleMyNewCommand(ctx context.Context, chatID int64, arg string) error {
    if arg == "" {
        return p.sendError(chatID, "請輸入參數")
    }
    return p.tgCommandUsecase.GetMyNewData(ctx, arg, chatID)
}
```

### 2. 在 Command UseCase 實作業務邏輯

**檔案**：`tg_bot_command.go` 或 `line_bot_command.go`

```go
// 1. 在介面定義新方法
type TgBotCommandUsecase interface {
    // ... 其他方法
    GetMyNewData(ctx context.Context, arg string, chatID int64) error
}

// 2. 實作方法
func (u *tgBotCommandUsecase) GetMyNewData(ctx context.Context, arg string, chatID int64) error {
    // 呼叫共用業務邏輯
    data, err := u.botCommandUsecase.GetMyNewData(ctx, UserTypeTelegram, arg)
    if err != nil {
        return err
    }
    
    // 發送訊息
    return u.client.SendMessage(chatID, data)
}
```

### 3. (可選) 在共用業務邏輯實作

**檔案**：`bot_command.go`

```go
// 1. 在介面定義新方法
type BotCommandUsecase interface {
    // ... 其他方法
    GetMyNewData(ctx context.Context, userType valueobject.UserType, arg string) (string, error)
}

// 2. 實作方法
func (u *botCommandUsecase) GetMyNewData(ctx context.Context, userType valueobject.UserType, arg string) (string, error) {
    // 呼叫 domain 層的業務邏輯
    result, err := u.someUsecase.GetData(ctx, arg)
    if err != nil {
        return "", err
    }
    
    // 格式化訊息
    return u.formatterPort.FormatMyData(result, userType), nil
}
```

## 集中管理命令的優勢

### 1. 單一職責原則 (SRP)

- **Handler**：只負責 HTTP 層面
- **Message Processor**：只負責命令路由和參數驗證
- **Command UseCase**：只負責業務邏輯編排
- **Bot Client**：只負責發送訊息

### 2. 易於維護

- 所有命令都在 `routeCommand()` 方法中統一管理
- 新增命令時只需要修改路由表
- 命令處理邏輯獨立，互不影響

### 3. 易於測試

- 每一層都可以獨立測試
- Message Processor 可以 mock Command UseCase
- Command UseCase 可以 mock Bot Client

### 4. 易於擴展

- 新增 Bot 平台（如 Discord）只需要實作新的 Handler、Message Processor 和 Command UseCase
- 共用的業務邏輯不需要重複實作

### 5. 符合 Clean Architecture

- 依賴方向正確：外層依賴內層
- 業務邏輯與技術實現分離
- 易於替換 Bot SDK 或底層實現

## 測試建議

### 1. Unit Test

```go
// Message Processor 測試
func TestTelegramMessageProcessor_ProcessUpdate(t *testing.T) {
    // Mock dependencies
    mockUsecase := &mockTgBotCommandUsecase{}
    mockUserPort := &mockUserAccountPort{}
    mockClient := &mockTgBotClient{}
    
    processor := NewTelegramMessageProcessor(
        mockUsecase,
        mockUserPort,
        mockClient,
        logger,
    )
    
    // Test case
    update := &tgbot.Update{
        Message: &tgbot.Message{
            Chat: tgbot.Chat{ID: 123},
            Text: "/start",
        },
    }
    
    err := processor.ProcessUpdate(context.Background(), update)
    assert.NoError(t, err)
    assert.True(t, mockUsecase.GetUseGuideMessageCalled)
}
```

### 2. Integration Test

```go
// Command UseCase 整合測試
func TestTgBotCommandUsecase_GetStockPrice(t *testing.T) {
    // Setup real dependencies
    cfg := config.LoadTestConfig()
    client, _ := tgbotapi.NewBot(cfg, logger)
    usecase := NewTgBotCommandUsecase(
        botCommandUsecase,
        client,
    )
    
    // Test with real data
    err := usecase.GetStockPrice(context.Background(), "2330", nil, 123)
    assert.NoError(t, err)
}
```
