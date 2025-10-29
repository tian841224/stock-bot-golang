package imageutil

import (
	"fmt"
	"os"

	"github.com/tian841224/stock-bot/pkg/logger"

	"github.com/flopp/go-findfont"
	"github.com/golang/freetype/truetype"
	"go.uber.org/zap"
	"golang.org/x/image/font/gofont/goregular"
)

// FontLoader 字型載入器
type FontLoader struct {
	// 支援繁體中文的粗體字型路徑 (TTF 格式)
	FontPaths []string
	// 字型名稱列表
	FontNames []string
}

// NewFontLoader 建立新的字型載入器
func NewFontLoader() *FontLoader {
	return &FontLoader{
		FontPaths: []string{
			// Linux 字型路徑
			"/usr/share/fonts/custom/NotoSansTC-Bold.ttf",
			"/usr/share/fonts/custom/NotoSansTC-VariableFont.ttf",
			// Windows 字型路徑
			"C:/Windows/Fonts/msjh.ttc",                  // 微軟正黑體
			"C:/Windows/Fonts/msjhbd.ttc",                // 微軟正黑體粗體
			"C:/Windows/Fonts/NotoSansCJKtc-Regular.otf", // Noto Sans CJK TC
			"C:/Windows/Fonts/NotoSansCJKtc-Bold.otf",    // Noto Sans CJK TC 粗體
			"C:/Windows/Fonts/simsun.ttc",                // 宋體
			"C:/Windows/Fonts/simhei.ttf",                // 黑體
		},
		FontNames: []string{
			// Linux 字型名稱
			"Noto Sans TC Bold",
			"Noto Sans CJK TC Bold",
			"Noto Sans TC",
			"Noto Sans CJK TC",
			// Windows 字型名稱
			"Microsoft JhengHei",         // 微軟正黑體
			"Microsoft JhengHei Bold",    // 微軟正黑體粗體
			"Microsoft JhengHei UI",      // 微軟正黑體 UI
			"Microsoft JhengHei UI Bold", // 微軟正黑體 UI 粗體
			"Noto Sans CJK TC",           // Noto Sans CJK TC
			"Noto Sans CJK TC Bold",      // Noto Sans CJK TC 粗體
			"SimSun",                     // 宋體
			"SimHei",                     // 黑體
			"PMingLiU",                   // 新細明體
			"MingLiU",                    // 細明體
		},
	}
}

// LoadChineseFont 載入支援繁體中文的字型
func (fl *FontLoader) LoadChineseFont() (*truetype.Font, error) {
	// 先嘗試直接路徑載入
	for _, path := range fl.FontPaths {
		// logger.Log.Info("嘗試載入字型", zap.String("path", path))
		if font, err := fl.loadFontFromPath(path); err == nil {
			// logger.Log.Info("成功載入字型", zap.String("path", path))
			return font, nil
		}

		// 再嘗試使用字型名稱查找 (優先粗體)
		for _, name := range fl.FontNames {
			// logger.Log.Info("嘗試查找字型", zap.String("name", name))
			if fontPath, err := findfont.Find(name); err == nil {
				// logger.Log.Info("找到字型路徑", zap.String("name", name), zap.String("path", fontPath))
				if font, err := fl.loadFontFromPath(fontPath); err == nil {
					//logger.Log.Info("成功載入字型", zap.String("name", name), zap.String("path", fontPath))
					return font, nil
				}
			}
		}
	}

	// 最後使用內建字型
	logger.Log.Warn("無法載入繁體中文字型，使用內建字型")
	return truetype.Parse(goregular.TTF)
}

// loadFontFromPath 從指定路徑載入字型
func (fl *FontLoader) loadFontFromPath(path string) (*truetype.Font, error) {
	// 檢查檔案是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("字型檔案不存在: %s", path)
	}

	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("讀取字型檔案失敗: %v", err)
	}

	// 解析 TTF 格式字型
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("解析字型檔案失敗: %v", err)
	}

	logger.Log.Info("成功載入字型", zap.String("path", path))
	return font, nil
}

// 全域字型載入器實例
var defaultFontLoader = NewFontLoader()

// LoadChineseFont 全域函數，載入支援繁體中文的字型
func LoadChineseFont() (*truetype.Font, error) {
	return defaultFontLoader.LoadChineseFont()
}
