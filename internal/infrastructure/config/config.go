package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
	var config Config

	// 設定 .env 檔案
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")     // 目前目錄
	viper.AddConfigPath("..")    // 上一層目錄
	viper.AddConfigPath("../..") // 專案根目錄（適用於 cmd 子目錄）

	// 嘗試讀取 .env 檔案（如果存在的話）
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("無法讀取 .env 檔案: %v，將使用環境變數\n", err)
	}

	// 啟用自動環境變數支援
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 將環境變數綁定到結構體
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析設定失敗: %w", err)
	}

	// 驗證配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置驗證失敗: %w", err)
	}

	return &config, nil
}
