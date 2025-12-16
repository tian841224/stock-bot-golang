// Package logger provides structured logging for the application.
package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field 定義日誌欄位，不依賴具體的實作
type Field interface{}

// Logger 定義日誌記錄介面，完全抽象化，不依賴 zap
type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Sync() error
}

// 便利函數：建立各種類型的日誌欄位

// String 建立字串欄位
func String(key, value string) Field {
	return zap.String(key, value)
}

// Int 建立整數欄位
func Int(key string, value int) Field {
	return zap.Int(key, value)
}

// Int64 建立 64 位整數欄位
func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

// Float64 建立浮點數欄位
func Float64(key string, value float64) Field {
	return zap.Float64(key, value)
}

// Bool 建立布林欄位
func Bool(key string, value bool) Field {
	return zap.Bool(key, value)
}

// Error 建立錯誤欄位
func Error(err error) Field {
	return zap.Error(err)
}

// Time 建立時間欄位
func Time(key string, value time.Time) Field {
	return zap.Time(key, value)
}

// Any 建立任意類型欄位
func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

// zapLogger 實作 Logger 介面
type zapLogger struct {
	logger *zap.Logger
}

// convertFields 將 Field 轉換為 zap.Field
func convertFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		if zapField, ok := field.(zap.Field); ok {
			zapFields = append(zapFields, zapField)
		} else {
			zapFields = append(zapFields, zap.Any("field", field))
		}
	}
	return zapFields
}

// Info 記錄資訊級別日誌
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, convertFields(fields...)...)
}

// Error 記錄錯誤級別日誌
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, convertFields(fields...)...)
}

// Warn 記錄警告級別日誌
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, convertFields(fields...)...)
}

// Debug 記錄除錯級別日誌
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, convertFields(fields...)...)
}

// Panic 記錄並觸發 panic
func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, convertFields(fields...)...)
}

// Fatal 記錄並終止程式
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, convertFields(fields...)...)
}

// Sync 同步日誌緩衝區
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// NewLogger 建立新的 Logger 實例
func NewLogger() (Logger, error) {
	mode := os.Getenv("GIN_MODE")

	var cfg zap.Config
	if strings.EqualFold(mode, "release") {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	if lvl := strings.TrimSpace(os.Getenv("LOG_LEVEL")); lvl != "" {
		var level zapcore.Level
		if err := level.Set(strings.ToLower(lvl)); err != nil {
			return nil, err
		}
		cfg.Level = zap.NewAtomicLevelAt(level)
	}

	zapLog, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: zapLog}, nil
}
