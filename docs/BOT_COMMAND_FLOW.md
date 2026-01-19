# Bot 命令處理流程圖

## 整體架構

```
┌─────────────────────────────────────────────────────────────────┐
│                        外部服務層                                  │
│                    (Telegram/LINE Server)                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ Webhook POST
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Interfaces Layer (介面層)                      │
│                                                                   │
│  ┌─────────────────────┐         ┌─────────────────────┐       │
│  │ TgHandler.Webhook() │         │LineBotHandler       │       │
│  │                     │         │.Webhook()           │       │
│  │ - 驗證 Secret Token │         │ - 解析 LINE 事件    │       │
│  │ - 解析 JSON         │         │ - 回應 200          │       │
│  │ - 回應 200          │         │                     │       │
│  └─────────────────────┘         └─────────────────────┘       │
│           │                                │                     │
└───────────┼────────────────────────────────┼─────────────────────┘
            │                                │
            │ ProcessUpdate()                │ ProcessTextMessage()
            ▼                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                 Application Layer (應用層)                        │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │         Message Processor (訊息處理器)                      │  │
│  │                                                            │  │
│  │  ┌────────────────────┐    ┌────────────────────┐       │  │
│  │  │TelegramMessage     │    │LineMessage         │       │  │
│  │  │Processor           │    │Processor           │       │  │
│  │  │                    │    │                    │       │  │
│  │  │1. 確保使用者存在   │    │1. 確保使用者存在   │       │  │
│  │  │2. 解析命令和參數   │    │2. 解析命令和參數   │       │  │
│  │  │3. 命令路由         │    │3. 命令路由         │       │  │
│  │  │4. 參數驗證         │    │4. 參數驗證         │       │  │
│  │  └────────────────────┘    └────────────────────┘       │  │
│  │           │                          │                    │  │
│  └───────────┼──────────────────────────┼────────────────────┘  │
│              │                          │                        │
│              │ GetXXX()                 │ GetXXX()              │
│              ▼                          ▼                        │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │         Command UseCase (命令用例)                         │  │
│  │                                                            │  │
│  │  ┌────────────────────┐    ┌────────────────────┐       │  │
│  │  │TgBotCommand        │    │LineBotCommand      │       │  │
│  │  │Usecase             │    │Usecase             │       │  │
│  │  │                    │    │                    │       │  │
│  │  │- GetStockPrice()   │    │- GetStockPrice()   │       │  │
│  │  │- GetTopVolume()    │    │- GetTopVolume()    │       │  │
│  │  │- GetChart()        │    │- GetChart()        │       │  │
│  │  └────────────────────┘    └────────────────────┘       │  │
│  │           │                          │                    │  │
│  │           └──────────┬───────────────┘                    │  │
│  │                      │                                     │  │
│  │                      ▼                                     │  │
│  │         ┌────────────────────────┐                        │  │
│  │         │BotCommandUsecase       │                        │  │
│  │         │(共用業務邏輯)           │                        │  │
│  │         │                        │                        │  │
│  │         │- 呼叫 Domain 服務      │                        │  │
│  │         │- 格式化訊息            │                        │  │
│  │         └────────────────────────┘                        │  │
│  │                      │                                     │  │
│  └──────────────────────┼─────────────────────────────────────┘  │
│                         │                                        │
└─────────────────────────┼────────────────────────────────────────┘
                          │
                          │ SendMessage() / SendPhoto()
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│              Infrastructure Layer (基礎設施層)                    │
│                                                                   │
│  ┌────────────────────┐         ┌────────────────────┐         │
│  │TgBotClient         │         │LineBotClient       │         │
│  │                    │         │                    │         │
│  │- SendMessage()     │         │- ReplyMessage()    │         │
│  │- SendPhoto()       │         │- ReplyPhoto()      │         │
│  │- SendKeyboard()    │         │- ReplyButtons()    │         │
│  └────────────────────┘         └────────────────────┘         │
│           │                              │                       │
└───────────┼──────────────────────────────┼───────────────────────┘
            │                              │
            │ API Call                     │ API Call
            ▼                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        外部服務層                                  │
│                    (Telegram/LINE Server)                        │
└─────────────────────────────────────────────────────────────────┘
```

## 命令處理詳細流程

### 以 `/d 2330` (查詢股價) 為例

```
1. Telegram Server 發送 webhook
   POST /telegram/webhook
   Body: {
     "message": {
       "chat": {"id": 123456},
       "text": "/d 2330"
     }
   }

2. TgHandler.Webhook()
   ├─ 驗證 X-Telegram-Bot-Api-Secret-Token
   ├─ 解析 JSON 為 tgbotapi.Update
   ├─ 回應 HTTP 200
   └─ 啟動 goroutine 處理訊息

3. TelegramMessageProcessor.ProcessUpdate()
   ├─ 提取 chatID = 123456
   ├─ 提取 messageText = "/d 2330"
   ├─ 確保使用者存在 (GetOrCreate)
   ├─ 解析命令: command="/d", arg1="2330", arg2=""
   └─ 路由命令

4. TelegramMessageProcessor.routeCommand()
   ├─ 匹配到 case "/d"
   ├─ 呼叫 handleStockPrice()
   ├─ 驗證 symbol 不為空
   ├─ 解析 date (如果有提供)
   └─ 呼叫 usecase

5. TgBotCommandUsecase.GetStockPrice()
   ├─ 呼叫 BotCommandUsecase.GetStockPrice()
   │   ├─ 呼叫 marketDataUsecase.GetStockPrice()
   │   ├─ 取得股價資料
   │   └─ 使用 formatterPort.FormatStockPrice() 格式化
   ├─ 取得格式化後的訊息
   └─ 呼叫 client.SendMessage()

6. TgBotClient.SendMessage()
   ├─ 建立 tgbotapi.NewMessage()
   ├─ 設定 ParseMode = HTML
   └─ 呼叫 Telegram Bot API

7. Telegram Server 將訊息發送給使用者
```

## 錯誤處理流程

```
任何層級發生錯誤
    │
    ▼
向上拋出 error
    │
    ▼
Message Processor 捕捉
    │
    ├─ 記錄 log
    │
    └─ 呼叫 sendError()
        │
        ▼
    發送錯誤訊息給使用者
```

## Panic Recovery

每個 webhook handler 都有 panic recovery 機制：

```go
go func(u tgbotapi.Update) {
    defer func() {
        if r := recover(); r != nil {
            h.logger.Error("處理 Telegram 更新發生 panic", logger.Any("recover", r))
        }
    }()

    ctx := context.Background()
    if err := h.messageProcessor.ProcessUpdate(ctx, &u); err != nil {
        h.logger.Error("處理 Telegram 更新失敗", logger.Error(err))
    }
}(update)
```

## 關鍵設計決策

### 1. 為什麼使用 goroutine？

```go
// 先回應 200，避免 webhook 超時
c.Status(http.StatusOK)

// 背景處理，避免 Telegram/LINE 重送
go func(u tgbotapi.Update) {
    // 處理邏輯...
}(update)
```

**原因**：
- Telegram/LINE 對 webhook 有超時限制（通常 30 秒）
- 如果處理時間過長，會導致 webhook 超時
- 超時後平台會重送請求，造成重複處理
- 先回應 200 表示已收到，再背景處理

### 2. 為什麼分離 Message Processor 和 Command UseCase？

**Message Processor**：
- 負責「路由」：決定該執行哪個命令
- 負責「驗證」：參數格式、使用者存在等
- 平台相關的邏輯（Telegram 用 chatID, LINE 用 replyToken）

**Command UseCase**：
- 負責「業務邏輯」：呼叫 domain 服務
- 負責「格式化」：將資料格式化為訊息
- 負責「發送」：透過 client 發送訊息

**好處**：
- 單一職責，易於測試
- 可以獨立修改路由規則或業務邏輯
- 容易新增新平台（如 Discord）

### 3. 為什麼有 BotCommandUsecase？

**目的**：共用業務邏輯

**範例**：
```go
// Telegram 版本
func (u *tgBotCommandUsecase) GetStockPrice(...) error {
    message, err := u.botCommandUsecase.GetStockPrice(...) // 共用邏輯
    return u.client.SendMessage(chatID, message)           // Telegram 特定
}

// LINE 版本
func (u *lineBotCommandUsecase) GetStockPrice(...) error {
    message, err := u.botCommandUsecase.GetStockPrice(...) // 共用邏輯
    return u.client.ReplyMessage(replyToken, message)      // LINE 特定
}
```

**好處**：
- 避免重複程式碼
- 業務邏輯修改只需要改一處
- Telegram 和 LINE 的差異只在發送方式
