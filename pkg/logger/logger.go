// Package logger 提供日誌記錄功能
package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger 初始化全域日誌器
// - 預設依據 GIN_MODE 選擇對應設定：
//   - debug: 使用 Development 設定（易讀、含呼叫來源）
//   - release: 使用 Production 設定（精簡 JSON）
//
// - 無論模式，輸出改為 stdout/stderr，避免在不同平台被吞掉
// - 支援 LOG_LEVEL 覆蓋層級：debug, info, warn, error
func InitLogger() {
	mode := os.Getenv("GIN_MODE")

	var cfg zap.Config
	if strings.EqualFold(mode, "release") {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// 確保輸出到 stdout/stderr
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	// 允許用 LOG_LEVEL 覆蓋層級
	if lvl := strings.TrimSpace(os.Getenv("LOG_LEVEL")); lvl != "" {
		var level zapcore.Level
		if err := level.Set(strings.ToLower(lvl)); err == nil {
			cfg.Level = zap.NewAtomicLevelAt(level)
		}
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Log = l
}
