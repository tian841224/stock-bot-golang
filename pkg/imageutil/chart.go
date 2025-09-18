package imageutil

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// PerformanceData 績效資料結構
type PerformanceData struct {
	Period      string
	PeriodName  string
	Performance string
}

// ChartConfig 圖表設定
type ChartConfig struct {
	Title      string
	Width      int
	Height     int
	ShowGrid   bool
	ShowLegend bool
	ChartType  string // "line" 或 "bar"
}

// DefaultChartConfig 預設圖表設定
func DefaultChartConfig() ChartConfig {
	return ChartConfig{
		Title:      "股票績效表現",
		Width:      1200, // 增加寬度
		Height:     600,  // 增加高度
		ShowGrid:   true,
		ShowLegend: true,
		ChartType:  "line",
	}
}

// GeneratePerformanceChartPNG 生成績效圖表 (PNG格式)
func GeneratePerformanceChartPNG(data []PerformanceData, config ChartConfig) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("無績效資料可生成圖表")
	}

	// 建立圖片
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	// 填充白色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// 載入字型 - 嘗試載入支援中文的字型
	var ttf *truetype.Font
	var err error

	// 嘗試載入系統中支援中文的字型（Windows系統）
	chineseFontPaths := []string{
		"C:\\Windows\\Fonts\\msyh.ttc",   // 微軟雅黑
		"C:\\Windows\\Fonts\\simsun.ttc", // 宋體
		"C:\\Windows\\Fonts\\simhei.ttf", // 黑體
		"C:\\Windows\\Fonts\\simkai.ttf", // 楷體
	}

	fontLoaded := false
	for _, fontPath := range chineseFontPaths {
		if _, err := os.Stat(fontPath); err == nil {
			fontBytes, err := os.ReadFile(fontPath)
			if err == nil {
				// TTC 檔案需要特殊處理，先嘗試解析
				ttf, err = truetype.Parse(fontBytes)
				if err == nil {
					fontLoaded = true
					break
				}
			}
		}
	}

	// 如果找不到中文字型，使用預設字型
	if !fontLoaded {
		ttf, err = truetype.Parse(goregular.TTF)
		if err != nil {
			return nil, fmt.Errorf("載入字型失敗: %v", err)
		}
	}

	// 建立 freetype context
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ttf)
	c.SetFontSize(12)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 255})) // 黑色文字

	// 解析績效資料
	values := make([]float64, len(data))
	for i, item := range data {
		performanceStr := strings.TrimSuffix(item.Performance, "%")
		performance, err := strconv.ParseFloat(performanceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("解析績效數據失敗: %v", err)
		}
		values[i] = performance
	}

	// 計算圖表邊界
	minVal := values[0]
	maxVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// 加入一些邊距
	margin := (maxVal - minVal) * 0.1
	if margin == 0 {
		margin = 1
	}
	minVal -= margin
	maxVal += margin

	// 圖表區域 - 為更大圖片調整邊距
	chartLeft := 120
	chartTop := 100
	chartWidth := config.Width - 200   // 左右各留更多空間
	chartHeight := config.Height - 200 // 上下各留更多空間

	// 繪製標題
	c.SetFontSize(16)
	titleWidth := len(config.Title) * 8 // 估算標題寬度
	titleX := (config.Width - titleWidth) / 2
	pt := freetype.Pt(titleX, 50) // 增加頂部邊距，確保不被裁掉
	c.DrawString(config.Title, pt)

	// 重設字型大小
	c.SetFontSize(10)

	// 繪製座標軸
	drawLine(img, chartLeft, chartTop, chartLeft, chartTop+chartHeight, color.RGBA{0, 0, 0, 255})                        // Y軸
	drawLine(img, chartLeft, chartTop+chartHeight, chartLeft+chartWidth, chartTop+chartHeight, color.RGBA{0, 0, 0, 255}) // X軸

	// 繪製格線和軸標籤
	if config.ShowGrid {
		// 繪製 Y 軸標籤（績效百分比）
		yGridLines := 5
		for i := 0; i <= yGridLines; i++ {
			y := chartTop + (chartHeight * i / yGridLines)
			value := maxVal - ((maxVal - minVal) * float64(i) / float64(yGridLines))

			// 水平格線
			if i > 0 && i < yGridLines {
				drawLine(img, chartLeft, y, chartLeft+chartWidth, y, color.RGBA{200, 200, 200, 255})
			}

			// Y軸標籤
			label := fmt.Sprintf("%.1f%%", value)
			pt := freetype.Pt(chartLeft-100, y+5)
			c.DrawString(label, pt)
		}

		// 繪製 X 軸標籤（時間）- 智能顯示，避免標籤過密
		labelStep := 1
		if len(data) > 15 {
			labelStep = len(data) / 10 // 最多顯示10個標籤，利用更大圖片空間
		}
		if labelStep < 2 {
			labelStep = 2 // 最小間隔2個資料點
		}

		for i, item := range data {
			x := chartLeft + (chartWidth * i / (len(data) - 1))

			// 垂直格線 - 減少格線數量
			if i > 0 && i < len(data)-1 && i%labelStep == 0 {
				drawLine(img, x, chartTop, x, chartTop+chartHeight, color.RGBA{200, 200, 200, 255})
			}

			// X軸標籤 - 只顯示部分標籤避免擁擠
			showLabel := false
			if i%labelStep == 0 {
				showLabel = true
			} else if i == len(data)-1 {
				// 只有當最後一個標籤與前一個顯示的標籤距離足夠時才顯示
				lastLabelIndex := ((len(data) - 2) / labelStep) * labelStep
				if i-lastLabelIndex >= labelStep/2 {
					showLabel = true
				}
			}

			if showLabel {
				label := item.PeriodName
				// 調整標籤位置，讓文字更居中，增加垂直間距
				pt := freetype.Pt(x-25, chartTop+chartHeight+45)
				c.DrawString(label, pt)
			}
		}
	}

	if config.ChartType == "bar" {
		// 垂直柱狀圖（X軸=時間，Y軸=績效）
		barWidth := chartWidth / len(data) * 8 / 10 // 80% 寬度
		barSpacing := chartWidth / len(data)

		// 計算零點位置
		zeroY := chartTop + chartHeight - int((0-minVal)/(maxVal-minVal)*float64(chartHeight))

		for i, _ := range data {
			value := values[i]

			// 計算柱狀圖位置
			x := chartLeft + (barSpacing * i) + (barSpacing-barWidth)/2

			// 根據正負值選擇顏色和位置
			barColor := color.RGBA{76, 175, 80, 255} // 綠色 (正值)
			var y, barHeight int

			if value >= 0 {
				// 正值：從零點往上畫
				barHeight = int((value - 0) / (maxVal - minVal) * float64(chartHeight))
				y = zeroY - barHeight
			} else {
				// 負值：從零點往下畫
				barColor = color.RGBA{244, 67, 54, 255} // 紅色 (負值)
				barHeight = int((0 - value) / (maxVal - minVal) * float64(chartHeight))
				y = zeroY
			}

			// 繪製垂直柱狀圖
			if barHeight > 0 {
				drawRect(img, x, y, barWidth, barHeight, barColor)
			}
		}

	} else {
		// 折線圖（X軸=時間，Y軸=績效）
		lineColor := color.RGBA{33, 150, 243, 255} // 藍色

		// 繪製折線（只有連接線，不顯示資料點）
		for i := range data {
			value := values[i]
			x := chartLeft + (chartWidth * i / (len(data) - 1))
			y := chartTop + chartHeight - int((value-minVal)/(maxVal-minVal)*float64(chartHeight))

			// 繪製線段 (除了第一個點)
			if i > 0 {
				prevValue := values[i-1]
				prevX := chartLeft + (chartWidth * (i - 1) / (len(data) - 1))
				prevY := chartTop + chartHeight - int((prevValue-minVal)/(maxVal-minVal)*float64(chartHeight))
				drawLine(img, prevX, prevY, x, y, lineColor)
			}
		}

		// 零線 (如果有負值) - 水平線
		hasNegative := false
		for _, v := range values {
			if v < 0 {
				hasNegative = true
				break
			}
		}

		if hasNegative {
			zeroY := chartTop + chartHeight - int((0-minVal)/(maxVal-minVal)*float64(chartHeight))
			drawDashedLine(img, chartLeft, zeroY, chartLeft+chartWidth, zeroY, color.RGBA{128, 128, 128, 255})
		}
	}

	// 軸標籤
	ptX := freetype.Pt(chartLeft+chartWidth/2-20, config.Height-40) // 增加底部邊距
	c.DrawString("Time", ptX)

	// Y軸標籤 (垂直文字效果用簡化版本)
	pt2 := freetype.Pt(30, chartTop+chartHeight/2) // 增加左側邊距
	c.DrawString("Performance (%)", pt2)

	// 將圖片編碼為 PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("編碼 PNG 失敗: %v", err)
	}

	return buf.Bytes(), nil
}

// drawLine 繪製直線
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	sy := 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	x, y := x1, y1
	for {
		img.Set(x, y, col)
		if x == x2 && y == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// drawDashedLine 繪製虛線
func drawDashedLine(img *image.RGBA, x1, y1, x2, y2 int, col color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	sy := 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	x, y := x1, y1
	dash := 0
	for {
		if dash%10 < 5 { // 5像素實線，5像素空白
			img.Set(x, y, col)
		}
		dash++
		if x == x2 && y == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// drawRect 繪製矩形
func drawRect(img *image.RGBA, x, y, width, height int, col color.RGBA) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(x+i, y+j, col)
		}
	}
}

// drawCircle 繪製圓形
func drawCircle(img *image.RGBA, centerX, centerY, radius int, col color.RGBA) {
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			if x*x+y*y <= radius*radius {
				img.Set(centerX+x, centerY+y, col)
			}
		}
	}
}

// abs 取絕對值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GeneratePerformanceLineChart 生成折線圖 (PNG格式)
func GeneratePerformanceLineChart(data []PerformanceData, title string) ([]byte, error) {
	config := DefaultChartConfig()
	config.Title = title
	config.ChartType = "line"
	return GeneratePerformanceChartPNG(data, config)
}

// GeneratePerformanceBarChart 生成柱狀圖 (PNG格式)
func GeneratePerformanceBarChart(data []PerformanceData, title string) ([]byte, error) {
	config := DefaultChartConfig()
	config.Title = title
	config.ChartType = "bar"
	return GeneratePerformanceChartPNG(data, config)
}
