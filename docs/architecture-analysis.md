# 專案結構分析與重構評估

本文件記錄一次針對 `stock-bot` 的全面結構盤點:架構分層、死碼、重複程式碼與品質問題。
結論分為兩部分——**本次已執行的低風險清理**(見第 3 節),以及**建議後續處理的大型重構**(見第 7 節)。

分析方法:以 `deadcode`(golang.org/x/tools)、`go vet` 搭配全庫 `grep` 交叉驗證每一項發現,
避免誤報。所有「已執行」項目皆通過 `go build ./...`、`go vet ./...`、`go test ./...`。

---

## 1. 概述

| 項目 | 內容 |
|------|------|
| 模組 | `github.com/tian841224/stock-bot` |
| Go 版本 | 1.24 |
| 進入點 | `cmd/bot`(webhook 伺服器)、`cmd/sync_stock_info`(同步)、`cmd/notification_stock_info`(排程通知) |
| 架構 | Clean / Hexagonal,四層:`domain` → `application` → `infrastructure` / `interfaces` |

整體評價:**洋蔥核心健全,問題集中在 bot/notification 這一片與 wiring/命名層。**

---

## 2. 架構評估

### 2.1 健全之處

- `internal/domain` **純淨**:未 import 任何 application / infrastructure / interfaces,
  無 GORM / viper / gin 洩漏,實體只依賴標準函式庫與同層 `valueobject` / `error`。
- **Port 放對位置**:介面宣告於消費端 `internal/application/port/`,並採 Reader/Writer 分離(ISP)。
  基礎設施於生產端以 `var _ port.Xxx = (*impl)(nil)` 斷言實作,例如
  `adapter/stock/finmind_stock_info.go`。
- `infrastructure` 僅依賴 `application/port` + `domain`,方向正確。
- `adapter` 與 `external` 的分工對 **stock provider** 是清楚的:`external/*` 是原始 API client + 廠商 DTO;
  `adapter/*` 是防腐層(ACL),把廠商 DTO 轉成領域實體 / 應用 DTO。

### 2.2 分層違規(建議後續處理)

| 違規 | 位置 | 說明 |
|------|------|------|
| 應用層依賴具體基礎設施 | `usecase/bot/telegram_command.go`、`line_command.go`、`*_message_processor.go`、`usecase/notification/send_notification.go` | 直接 import `infrastructure/external/bot/{line,telegram}`、`external/imgbb`,而非透過 port 抽象 |
| 介面層依賴具體 client | `interfaces/bot/line/handler.go` | 直接 import `external/bot/line` |
| ACL 缺口 | bot 這一片沒有 `adapter/bot/…` | stock 有防腐層,bot 完全沒有,雙平台差異直接漏進 usecase |
| 命名混淆 | `adapter/presenter`(其實是 ValidationGateway)vs `interfaces/presenter` | 同名不同義,分處兩層 |

> 這些是「依賴規則反轉」,不影響目前執行,但讓 bot 這一片難以測試與抽換。修法見第 7 節①。

---

## 3. 本次已執行的清理(低風險)

### 3.1 死碼移除(`refactor(cleanup)`)

全部經 `deadcode` + `grep` 驗證零呼叫點:

- **整檔刪除**:未使用的 `watchlist` / `watchlist_item` / `notification_delivery` /
  `notification_event` repository、`postgres/query_options.go`、
  四個獨立 health checker(`api_checker` / `database_checker` / `resource_monitor` / `sync_status_checker`,
  已被 `health_checker.go` 內聯取代)、`pkg/errors`。
- `persistence/postgres.go` 移除向後相容全域 `db` 變數與 `GetDB` / `InitDB` / `Close`。
- `telegram_formatter.go` 只保留在用的 `FormatStockNews`。
- 跨 port / adapter / usecase / dto 四層皆無呼叫的 `GetUserSubscriptionDetail`。
- 零散 helper:`domain/error` 未用建構子與 `IsNotFound` / `IsInvalidArgument`、
  `ParseSubscriptionType`、`ParseSubscriptionItem`、`zap Float64`、`FormatFloatWithCommas`。
- 註解掉的殭屍程式碼(`handleUnknownCommand` 殘骸、font debug log)。
- 從未被讀取的 config 欄位:`TELEGRAM_ADMIN_CHAT_ID`、`SCHEDULER_STOCK_SPEC`、`SCHEDULER_TIMEZONE`。

> **保留** `model/` 下無 repository 引用的 model(如 `NotificationDelivery`):它們在 `init()` 中
> 透過 `RegisterModel` 註冊給 AutoMigrate,屬 DB schema 意圖,移除會刪表。

> **方法論限制與修正**:最初這次清理也移除了 12 個外部 API 方法(finmindtrade 8 個、
> fugle 3 個、twse 1 個)與其專屬 DTO,判斷依據是「`deadcode` 工具 + 全庫 grep 皆為零呼叫點」。
> 但零呼叫點只能證明「目前沒有東西呼叫」,無法分辨這是「舊功能殘骸」還是「已先串好、
> 尚未接上層邏輯的預備工作」——兩者在靜態分析下完全無法區分,只有作者知道意圖。
> 這批方法屬於後者(先把 API 串好、規劃中尚未串接),已於後續復原,詳見附註。
> 順帶一提:這批方法本來就不會被 `deadcode` CLI 或 `golangci-lint` 的 `unused` 檢查標記——
> 兩者對「已匯出、且所屬具體型別仍被建構並傳遞」的方法皆有已知盲點,只能靠 grep 零呼叫點
> 人工判斷,這也是為何需要作者確認意圖而非只信工具結果。

合計移除約 1,300 行(死碼實際淨額,已扣除復原的外部 API 方法)。

### 3.2 命名修正(`refactor(naming)`)

純檔名調整(Go 不以檔名參照,零行為風險):
`line_fommatter.go`→`line_formatter.go`、`notification_deliverie.go`→`notification_delivery.go`、
`user_subscription_Item.go`→`user_subscription_item.go`、`stock_company_Info.go`→`stock_company_info.go`、
`inmindtradeRequest.go`→`finmindtradeRequest.go`。

### 3.3 Bug 修復與加固(`fix(stock)`)

- **`GetStockPerformance` 資料遺失(實際 bug)**:原以 `make([]T, len(result))` 回傳
  長度正確但內容全為零值的空切片,查得的資料從未帶入。改為直接回傳 `result`,並新增回歸測試。
- **`GetStockPrice` 越界防護**:存取 `tradeDates` 最後兩筆前先檢查 `len >= 2`,避免交易日不足時 panic。
- **背景 goroutine 加固**:股票同步 worker 將單批 upsert 包入 `processBatch` 攔截 panic;
  排程 fan-out goroutine 加入 `recover()`,單一任務崩潰不再拖垮整體。

### 3.4 依賴整理(`chore(deps)`)

`go mod tidy`:移除未使用的舊版 `line-bot-sdk-go v7`(實際用 v8),並把誤標 `// indirect`
的直接依賴(`gin-contrib/zap`、`line-bot-sdk-go/v8`、`multierr`、`x/image`)歸位。

---

## 4. 重複程式碼熱點(建議後續處理)

### 4.1 Telegram / LINE 雙平台複製 — 影響最大

雙平台在**四層**幾乎鏡像複製(估計 85–90% 重疊):

| 層 | 檔案 | 差異點 |
|----|------|--------|
| Command usecase | `telegram_command.go` / `line_command.go` | 僅 `SendMessage(chatID,…)` vs `ReplyMessage(replyToken,…)` |
| Message processor | `telegram_message_processor.go` / `line_message_processor.go` | 路由與 `handleX` 幾乎相同,僅 `chatID int64` vs `replyToken string` |
| Formatter | `bot_text_formatter.go` | 每個方法內嵌 `if userType == Telegram {…HTML…} else {…plain…}`,同欄位寫兩次 |
| Handler | `telegram/handler.go` / `line/handler.go` | 相同的「先回 200 再背景處理 + recover + timeout」骨架 |

**缺少的抽象**:沒有 `MessageSender` / `Replier` 介面把 `Send(target, text)` / `SendPhoto` 抽掉。
`14:00 收盤前取前一日` 的邏輯甚至重複三處(兩個 processor + `query_market_data.go`)。

### 4.2 Repository CRUD 逐實體複製

10 個 `postgres/*_repository.go` 各自手寫 `toEntity` / `toModel` / CRUD,
`Create/Update/Delete` 主體結構一致,`err == gorm.ErrRecordNotFound { return nil, nil }`
在單一 `subscription_repository.go` 內就重複 8 次。無泛型基底。

### 4.3 外部 client HTTP 邏輯重造

各 client 各有一份 `do request → 檢查狀態 → JSON decode → 包錯` 的私有泛型 helper
(`finmindtrade.doRequest`、`fugle.getResponse`、`cnyes.getResponse`、`imgbb.sendRequest`);
`twse/api.go` 甚至逐方法內聯複製。無共用 `pkg/httpclient`。

### 4.4 應用層重複

- `send_notification.go` 四個方法(股價 / 新聞 / 大盤 / 量能通知)約 95% 相同:
  取訂閱者 → 迴圈 → `ParseInt` → 送出 → log-and-`continue`。
- `sync_stock_info.go` 的 `SyncTaiwanStockInfo` 與 `SyncUSStockInfo` 僅差 `"TW"` / `"US"`。

---

## 5. 品質問題(建議後續處理)

- **錯誤處理三套並行**:`fmt.Errorf` + 硬編碼中文(且中文同時當使用者 UI 文案,把展示耦合進應用層)、
  `pkg/errors`(已刪,原本即死)、`domain/error.DomainError`(僅少數 entity 使用)。
  `%w` 與 `%v` 混用,stock client 用 `%v` 丟失錯誤鏈。
  通知迴圈以 `continue` 吞掉逐筆錯誤,整批失敗仍回傳 `nil`。
- **硬編碼值**:所有外部 API base URL、HTTP timeout(`10s` 重複 4 處)、`14`(收盤線,3 處)、
  batch size / worker 數、取樣上限等散落程式碼,應集中於 config 或具名常數。
- **測試覆蓋薄且偏斜**:約 10 個測試檔對 ~135 個原始檔;`domain` 僅 `stock_math_test`,
  10 個 repository 與 5 個 API client **零測試**。關鍵未測邏輯:split 還原 / 取樣數學、通知廣播路徑、
  repository mapper。
- **巨型檔案**:`pkg/imageutil/chart.go`(約 1,155 行)遠超其他檔,建議依「座標系 / 繪圖元件 / 版面」拆分。

---

## 6. `mockUserSubscriptionPort` 過時(備註)

`usecase/stock/mock_helper_test.go` 的 `mockUserSubscriptionPort` 方法簽章與現行
`port.UserSubscriptionPort` 已不一致(例如 `AddUserSubscriptionStock` 回傳 `(bool,error)` vs 介面的 `error`),
且未被任何介面斷言。屬測試 scaffolding,非產品死碼,列為低優先,待雙平台重構時一併整理。

---

## 7. 後續大型重構建議與優先序

| # | 重構 | 效益 | 風險 |
|---|------|------|------|
| ① | 引入 `MessageSender` / `Replier` port,把 Telegram/LINE 的 command usecase 與 message processor 收斂為單一實作 + 兩個薄 adapter | 消除約 380 行鏡像複製;根除「每個新指令要寫兩次」的長期稅;順帶補齊 bot 的 ACL/port | 中(需完整回歸雙平台) |
| ② | 抽出 `internal/bootstrap`(或 container)共用 wiring | 消除三個 `main.go` 約 100–120 行重複的 config/logger/DB/健康檢查建構 | 低 |
| ③ | 新增 `pkg/httpclient.Get[T]`,五個 client(尤其 twse 內聯與 finmind 手刻方法)統一走它 | 統一狀態處理與 `%w` 包裝,消除重複 | 低 |
| ④ | 泛型 `Repository[E,M]` + mapper | 消除約 30 個近乎相同的 CRUD 與 8× 的 `ErrRecordNotFound` 判斷 | 中 |
| ⑤ | 錯誤策略收斂到 `DomainError`,使用者文案移出 error value | 錯誤可被呼叫端判別;展示與應用層解耦 | 中 |

建議順序:先做 ②③(低風險、獨立),再做 ①,最後 ④⑤。
