# Stock Bot 系統架構圖

## 整體系統架構

```mermaid
graph TB
    subgraph "使用者介面層"
        TG[Telegram Bot]
        LINE[LINE Bot]
    end

    subgraph "API Gateway"
        GIN[Gin Web Framework<br/>Port 8080]
    end

    subgraph "Clean Architecture 分層"
        subgraph "Interfaces Layer 介面層"
            HTTP[HTTP Handlers]
            BOT[Bot Handlers]
            PRES[Presenters]
        end

        subgraph "Application Layer 應用層"
            UC[Use Cases<br/>業務邏輯]
            PORTS[Ports<br/>介面定義]
            DTO[DTOs<br/>資料傳輸物件]
        end

        subgraph "Domain Layer 領域層"
            ENT[Entities<br/>實體]
            VO[Value Objects<br/>值物件]
            ERR[Domain Errors<br/>領域錯誤]
        end

        subgraph "Infrastructure Layer 基礎設施層"
            REPO[Repositories<br/>資料持久化]
            EXT[External APIs<br/>外部服務]
            LOG[Logger<br/>日誌系統]
            CFG[Config<br/>配置管理]
        end
    end

    subgraph "Docker 服務"
        BOT_SVC[Stock Bot Service<br/>Port 8080]
        SYNC_SVC[Sync Service<br/>Port 8081]
        SCHED_SVC[Scheduler Service<br/>Port 8082]
    end

    subgraph "資料庫"
        PG[(PostgreSQL 16)]
    end

    subgraph "外部服務"
        TWSE[TWSE API<br/>台灣證交所]
        FUGLE[Fugle API<br/>富果]
        FINMIND[FinMind API<br/>金融資料]
        IMGBB[ImgBB API<br/>圖片儲存]
    end

    %% 使用者到 API Gateway
    TG --> GIN
    LINE --> GIN

    %% API Gateway 到介面層
    GIN --> HTTP
    GIN --> BOT

    %% 介面層到應用層
    HTTP --> UC
    BOT --> UC
    UC --> PRES

    %% 應用層到領域層
    UC --> PORTS
    PORTS --> ENT
    PORTS --> VO

    %% 應用層到基礎設施層
    PORTS --> REPO
    PORTS --> EXT
    UC --> LOG

    %% 基礎設施層到資料庫
    REPO --> PG

    %% 基礎設施層到外部服務
    EXT --> TWSE
    EXT --> FUGLE
    EXT --> FINMIND
    EXT --> IMGBB

    %% 服務層連接
    BOT_SVC --> PG
    SYNC_SVC --> PG
    SCHED_SVC --> PG

    SYNC_SVC --> TWSE
    SYNC_SVC --> FINMIND
    
    SCHED_SVC --> FUGLE
    SCHED_SVC --> IMGBB

    style TG fill:#0088cc,stroke:#006699,color:#fff
    style LINE fill:#00b900,stroke:#009900,color:#fff
    style GIN fill:#00add8,stroke:#0099cc,color:#fff
    style HTTP fill:#e3f2fd,stroke:#90caf9
    style BOT fill:#e3f2fd,stroke:#90caf9
    style PRES fill:#e3f2fd,stroke:#90caf9
    style UC fill:#e8f5e9,stroke:#81c784
    style PORTS fill:#e8f5e9,stroke:#81c784
    style DTO fill:#e8f5e9,stroke:#81c784
    style ENT fill:#fff9c4,stroke:#fff176
    style VO fill:#fff9c4,stroke:#fff176
    style ERR fill:#fff9c4,stroke:#fff176
    style REPO fill:#ffe0b2,stroke:#ffb74d
    style EXT fill:#ffe0b2,stroke:#ffb74d
    style LOG fill:#ffe0b2,stroke:#ffb74d
    style CFG fill:#ffe0b2,stroke:#ffb74d
    style BOT_SVC fill:#bbdefb,stroke:#64b5f6
    style SYNC_SVC fill:#bbdefb,stroke:#64b5f6
    style SCHED_SVC fill:#bbdefb,stroke:#64b5f6
    style PG fill:#c8e6c9,stroke:#66bb6a
    style TWSE fill:#f8bbd0,stroke:#f06292
    style FUGLE fill:#f8bbd0,stroke:#f06292
    style FINMIND fill:#f8bbd0,stroke:#f06292
    style IMGBB fill:#f8bbd0,stroke:#f06292
```

## Clean Architecture 依賴方向

```mermaid
graph LR
    A[Interfaces Layer] --> B[Application Layer]
    B --> C[Domain Layer]
    A --> D[Infrastructure Layer]
    B --> D
    D --> C

    style A fill:#e3f2fd,stroke:#90caf9
    style B fill:#e8f5e9,stroke:#81c784
    style C fill:#fff9c4,stroke:#fff176
    style D fill:#ffe0b2,stroke:#ffb74d
```

## Docker 服務架構

```mermaid
graph TB
    subgraph "Docker Compose"
        subgraph "應用服務"
            BOT[stock-bot<br/>主要 Bot 服務<br/>Port: 8080]
            SYNC[sync-stock-info<br/>資料同步服務<br/>Port: 8081]
            SCHED[scheduler<br/>定時通知服務<br/>Port: 8082]
        end

        DB[(postgres<br/>PostgreSQL 16<br/>Port: 5432)]
    end

    BOT --> DB
    SYNC --> DB
    SCHED --> DB

    BOT -.->|Health Check| BOT
    SYNC -.->|Health Check| SYNC
    SCHED -.->|Health Check| SCHED
    DB -.->|Health Check| DB

    style BOT fill:#bbdefb,stroke:#64b5f6
    style SYNC fill:#bbdefb,stroke:#64b5f6
    style SCHED fill:#bbdefb,stroke:#64b5f6
    style DB fill:#c8e6c9,stroke:#66bb6a
```

## CI/CD 部署流程

```mermaid
graph LR
    A[GitHub Push] --> B[GitHub Actions]
    
    subgraph "Build Job"
        B --> C[Go Build]
        C --> D[Go Test]
        D --> E[Go Vet]
    end

    subgraph "Push Job"
        E --> F[Build Docker Images]
        F --> G[Push to Docker Hub]
    end

    subgraph "Deploy Job"
        G --> H[SSH to EC2]
        H --> I[Pull Images]
        I --> J[Docker Compose Up]
    end

    J --> K[AWS EC2]

    style A fill:#f8bbd0,stroke:#f06292
    style B fill:#e1bee7,stroke:#ba68c8
    style C fill:#c5cae9,stroke:#7986cb
    style D fill:#c5cae9,stroke:#7986cb
    style E fill:#c5cae9,stroke:#7986cb
    style F fill:#b2dfdb,stroke:#4db6ac
    style G fill:#b2dfdb,stroke:#4db6ac
    style H fill:#fff9c4,stroke:#fff176
    style I fill:#fff9c4,stroke:#fff176
    style J fill:#fff9c4,stroke:#fff176
    style K fill:#ffccbc,stroke:#ff8a65
```

## 控制流向 (Control Flow)

展示使用者指令如何在系統中被處理和執行:

```mermaid
sequenceDiagram
    participant User as 使用者
    participant Bot as Telegram/LINE Bot
    participant Handler as Bot Handler
    participant UseCase as Use Case
    participant Validator as 驗證器
    participant Formatter as 格式化器

    User->>Bot: 發送指令 (/k 2330)
    Bot->>Handler: 解析訊息
    Handler->>Handler: 識別指令類型
    Handler->>Validator: 驗證股票代碼
    
    alt 驗證失敗
        Validator-->>Handler: 返回錯誤
        Handler-->>Bot: 錯誤訊息
        Bot-->>User: 顯示錯誤提示
    else 驗證成功
        Validator-->>Handler: 驗證通過
        Handler->>UseCase: 執行對應 Use Case
        UseCase->>UseCase: 執行業務邏輯
        UseCase-->>Handler: 返回處理結果
        Handler->>Formatter: 格式化輸出
        Formatter-->>Handler: 格式化完成
        Handler-->>Bot: 返回訊息
        Bot-->>User: 顯示結果
    end
```

## 資料流向 (Data Flow)

展示資料在系統各層之間的流動和儲存:

```mermaid
graph LR
    subgraph "外部資料源"
        TWSE[TWSE API<br/>即時股價]
        FUGLE[Fugle API<br/>K線資料]
        FINMIND[FinMind API<br/>歷史資料]
    end

    subgraph "資料同步層"
        SYNC[Sync Service<br/>定時同步]
    end

    subgraph "資料儲存層"
        DB[(PostgreSQL<br/>資料庫)]
        CACHE[快取層<br/>可選]
    end

    subgraph "業務邏輯層"
        UC[Use Cases<br/>業務處理]
        REPO[Repositories<br/>資料存取]
    end

    subgraph "展示層"
        CHART[圖表生成器]
        FORMAT[資料格式化]
    end

    subgraph "使用者介面"
        USER[使用者]
    end

    %% 資料同步流程
    TWSE -->|股價資料| SYNC
    FUGLE -->|K線資料| SYNC
    FINMIND -->|歷史資料| SYNC
    SYNC -->|寫入| DB

    %% 資料查詢流程
    USER -->|查詢請求| UC
    UC -->|讀取| REPO
    REPO -->|查詢| DB
    DB -->|返回資料| REPO
    
    %% 快取流程(可選)
    REPO -.->|快取| CACHE
    CACHE -.->|命中| REPO

    %% 資料不存在時的流程
    REPO -->|資料缺失| TWSE
    REPO -->|資料缺失| FUGLE
    TWSE -->|即時資料| REPO
    FUGLE -->|即時資料| REPO
    REPO -->|更新| DB

    %% 資料展示流程
    REPO -->|原始資料| UC
    UC -->|處理後資料| CHART
    UC -->|處理後資料| FORMAT
    CHART -->|圖表| USER
    FORMAT -->|文字訊息| USER

    style TWSE fill:#f8bbd0,stroke:#f06292
    style FUGLE fill:#f8bbd0,stroke:#f06292
    style FINMIND fill:#f8bbd0,stroke:#f06292
    style SYNC fill:#bbdefb,stroke:#64b5f6
    style DB fill:#c8e6c9,stroke:#66bb6a
    style CACHE fill:#fff9c4,stroke:#fff176
    style UC fill:#e8f5e9,stroke:#81c784
    style REPO fill:#ffe0b2,stroke:#ffb74d
    style CHART fill:#e1bee7,stroke:#ba68c8
    style FORMAT fill:#e1bee7,stroke:#ba68c8
    style USER fill:#b2dfdb,stroke:#4db6ac
```

## 資料流向詳細說明

### 1. 資料同步流程
- **Sync Service** 定時從外部 API 獲取最新資料
- 資料經過驗證和轉換後儲存到 PostgreSQL
- 確保資料庫中的資料保持最新

### 2. 資料查詢流程
- 使用者發送查詢請求
- Use Case 透過 Repository 從資料庫讀取資料
- 如果資料不存在或過期,直接呼叫外部 API 獲取
- 新資料會被儲存到資料庫供後續使用

### 3. 資料展示流程
- 原始資料經過業務邏輯處理
- 根據需求生成圖表或格式化為文字
- 最終呈現給使用者
