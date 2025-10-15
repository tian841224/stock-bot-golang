# 使用官方 Go 映像作為建置階段
FROM golang:1.24.1-alpine AS builder

# 安裝必要的套件
RUN apk add --no-cache git ca-certificates tzdata

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 檔案
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製原始碼
COPY . .

# 建置應用程式
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -ldflags='-w -s' \
    -o main ./cmd/bot

# 使用輕量級的 Alpine 映像作為執行階段
FROM alpine:latest

# 安裝 ca-certificates、tzdata 和下載工具
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    fontconfig \
    && rm -rf /var/cache/apk/*

# 下載繁體中文字型 (Variable Font TTF 格式)
RUN mkdir -p /usr/share/fonts/custom \
    && cd /usr/share/fonts/custom \
    && wget -O NotoSansTC-VariableFont.ttf \
       "https://github.com/google/fonts/raw/main/ofl/notosanstc/NotoSansTC%5Bwght%5D.ttf" \
    && fc-cache -f -v \
    && rm -rf /tmp/*

# 設定時區為台北時間
ENV TZ=Asia/Taipei

# 建立非 root 使用者
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 設定工作目錄
WORKDIR /root/

# 從建置階段複製執行檔
COPY --from=builder /app/main .

# 變更檔案擁有者
RUN chown -R appuser:appgroup /root/

# 切換到非 root 使用者
USER appuser

# 暴露埠號
EXPOSE 8080

# 執行應用程式
CMD ["./main"]