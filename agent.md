# Agent Context & Guidelines

此檔案為 AI Agent 協助開發時的專案指南，包含架構說明、開發規範與注意事項。閱讀此檔案可快速了解專案全貌與開發準則。

## 1. 專案概觀
本專案為「台股查詢機器人 (Stock Bot)」，支援 Telegram 與 Line 平台。
主要功能包含：即時股價查詢、K線圖繪製、新聞追蹤、個人化訂閱通知、大盤資訊等。

## 2. 技術堆疊
- **語言**: Golang 1.24+
- **Web 框架**: Gin
- **資料庫**: PostgreSQL 16 (使用 GORM)
- **容器化**: Docker, Docker Compose
- **架構設計**: Clean Architecture (整潔架構)
- **CI/CD**: GitHub Actions, AWS EC2
- **外部 API**: 
  - TWSE (台灣證交所)
  - Fugle (富果)
  - FinMind
  - ImgBB (圖片圖床)

## 3. 系統架構 (Clean Architecture)
專案嚴格遵循 Clean Architecture 分層原則，依賴方向僅能由外向內。

### 分層詳細說明
1.  **Domain Layer (領域層)**
    - **路徑**: `internal/domain`
    - **內容**: Entities (實體), Value Objects (值物件), Domain Errors (領域錯誤)。
    - **原則**: 最內層，不依賴任何外部庫或外層程式碼，包含核心業務規則。

2.  **Application Layer (應用層)**
    - **路徑**: `internal/application`
    - **內容**: Use Cases (業務邏輯), Ports (介面定義), DTOs (資料傳輸物件)。
    - **原則**: 定義應用程式的具體行為，協調 Domain Object，定義與外層溝通的介面 (Ports)。

3.  **Interfaces Layer (介面層)**
    - **路徑**: `internal/interfaces`
    - **內容**: HTTP Handlers (Gin), Bot Handlers (TG/Line), Presenters。
    - **原則**: 負責接收外部請求 (HTTP/Bot Webhook)，解析參數，呼叫 Use Case，並將結果格式化回傳。

4.  **Infrastructure Layer (基礎設施層)**
    - **路徑**: `internal/infrastructure`
    - **內容**: Persistence (Repository 實作), External APIs (HTTP Client), Logger, Config。
    - **原則**: 實作 Application 層定義的 Port 介面，處理具體的技術細節 (DB 連線、API 呼叫)。

## 4. 專案目錄結構
```
stock-bot/
├── cmd/                          # 應用程式入口
│   ├── bot/                      # 主服務 (Web/Bot)
│   ├── sync_stock_info/          # 同步服務 (Sync Service)
│   └── notification_stock_info/  # 排程通知服務 (Scheduler)
├── internal/                     # 核心程式碼 (Clean Architecture)
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   └── interfaces/
├── pkg/                          # 共用工具庫 (Utils, Helpers)
├── docs/                         # 系統文件
├── docker-compose.yml            # Docker 編排
└── .github/workflows/            # CI/CD 設定
```

## 5. 開發規範 (User Rules)
**⚠️ AI Agent 必須嚴格遵守以下 User 定義的規則：**

1.  **語言與用語**: 
    - 一律使用 **正體中文 (台灣用語)**。
2.  **分析與建議**:
    - 在修改代碼前，必須**分析修改前後的優缺點**。
    - 若發現現有寫法難以維護、結構不佳或有效能疑慮，**必須主動提出改善建議**。
3.  **溝通確認**: 
    - 當 User 的指令不夠精確、模糊或名詞使用錯誤時，**優先反問確認需求**，不要盲目執行。
4.  **Git Commit 規範**:
    - **僅針對暫存 (staged) 的變更**生成訊息。
    - **第一行**: 簡單扼要描述修改重點。
    - **後續內容**: 若為多項修改，使用條列式詳細說明。
    - 範例格式：
      ```text
      Feat: 新增個股到價通知功能

      - 實作 PriceAlertUseCase 邏輯
      - 新增 Telegram 通知推送適配器
      - 修正 UserEntity 缺少 AlertSettings 的問題
      ```

## 6. 開發與測試指令
- **啟動所有服務 (Docker)**: 
  `docker compose up -d`
- **本地執行 (Bot)**: 
  `go run ./cmd/bot`
- **執行單元測試**: 
  `go test -v ./...`
- **執行測試並查看覆蓋率**: 
  `go test -v -cover ./...`

## 7. 重要注意事項
- **依賴注入 (DI)**: 專案手動管理依賴注入，新增 Service 或 Repository 時，記得在 `cmd/` 下的 `main.go` 中進行初始化與注入。
- **環境變數**: 本地開發需複製 `.env.example` 為 `.env` 並填入對應金鑰。
- **錯誤處理**: 
  - 盡量使用 `internal/domain/error` 定義的錯誤類型。
  - 在 Interface 層統一攔截錯誤並轉換為使用者友善的訊息。
- **Git Flow**: 主要分支為 `master` (或 `main`)，開發建議開立 feature branch。

---
*Created by AI Agent based on project analysis.*
