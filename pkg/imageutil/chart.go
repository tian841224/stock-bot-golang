// Package imageutil 提供圖片處理與圖表生成功能
package imageutil

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
	"strings"
	"time"

	"github.com/golang/freetype"
)

// ChartColors 圖表顏色配置
type ChartColors struct {
	// 背景顏色
	BackgroundWhite     color.RGBA // 白色背景
	BackgroundLightGray color.RGBA // 淺灰色背景

	// 文字顏色
	TextDarkGray color.RGBA // 深灰色文字
	TextBlack    color.RGBA // 黑色文字
	TextRed      color.RGBA // 紅色文字
	TextGreen    color.RGBA // 綠色文字

	// 軸線顏色
	AxisDarkGray color.RGBA // 深灰色軸線
	AxisBlack    color.RGBA // 黑色軸線

	// 格線顏色
	GridLightGray  color.RGBA // 淺灰色格線
	GridMintGreen  color.RGBA // 薄荷綠格線
	GridDashedGray color.RGBA // 虛線灰色格線

	// 圖表元素顏色
	PositiveGreen   color.RGBA // 正值綠色
	NegativeRed     color.RGBA // 負值紅色
	KLineShadow     color.RGBA // K線影線
	KLineUpRed      color.RGBA // K線上漲紅色
	KLineDownGreen  color.RGBA // K線下跌綠色
	VolumeUpRed     color.RGBA // 成交量上漲紅色
	VolumeDownGreen color.RGBA // 成交量下跌綠色

	// 特殊標記顏色
	HighestPriceRed  color.RGBA // 最高價紅色
	LowestPriceGreen color.RGBA // 最低價綠色
	MonthlyAvgRed    color.RGBA // 月均價紅色
}

// DefaultChartColors 預設圖表顏色配置
func DefaultChartColors() ChartColors {
	return ChartColors{
		// 背景顏色
		BackgroundWhite:     color.RGBA{255, 255, 255, 255},
		BackgroundLightGray: color.RGBA{245, 245, 245, 255},

		// 文字顏色
		TextDarkGray: color.RGBA{15, 15, 15, 255},
		TextBlack:    color.RGBA{0, 0, 0, 255},
		TextRed:      color.RGBA{180, 30, 40, 255},
		TextGreen:    color.RGBA{20, 80, 40, 255},

		// 軸線顏色
		AxisDarkGray: color.RGBA{8, 8, 8, 255},
		AxisBlack:    color.RGBA{15, 15, 15, 255},

		// 格線顏色
		GridLightGray:  color.RGBA{200, 200, 200, 255},
		GridMintGreen:  color.RGBA{180, 200, 180, 180},
		GridDashedGray: color.RGBA{100, 100, 100, 255},

		// 圖表元素顏色
		PositiveGreen:   color.RGBA{100, 150, 120, 255},
		NegativeRed:     color.RGBA{180, 100, 100, 255},
		KLineShadow:     color.RGBA{60, 60, 60, 255},
		KLineUpRed:      color.RGBA{150, 50, 50, 255},
		KLineDownGreen:  color.RGBA{50, 120, 70, 255},
		VolumeUpRed:     color.RGBA{150, 50, 50, 255},
		VolumeDownGreen: color.RGBA{50, 120, 70, 255},

		// 特殊標記顏色
		HighestPriceRed:  color.RGBA{100, 20, 20, 255},
		LowestPriceGreen: color.RGBA{20, 80, 40, 255},
		MonthlyAvgRed:    color.RGBA{80, 15, 15, 255},
	}
}

// ChartTitle 圖表標題配置
type ChartTitle struct {
	// 標題字型大小
	FontSize int

	// 標題位置
	X int // X座標（相對於圖表寬度的比例，0-100）
	Y int // Y座標（相對於圖表高度的比例，0-100）

	// 標題顏色
	Color color.RGBA

	// 標題對齊方式
	Align string // "left", "center", "right"

	// 標題邊距
	MarginTop    int // 頂部邊距
	MarginBottom int // 底部邊距
}

// DefaultChartTitle 預設圖表標題配置
func DefaultChartTitle() ChartTitle {
	return ChartTitle{
		FontSize:     18,                          // 標題字型大小
		X:            50,                          // 水平居中
		Y:            8,                           // 距離頂部8%
		Color:        color.RGBA{15, 15, 15, 255}, // 深灰色
		Align:        "center",                    // 置中對齊
		MarginTop:    20,                          // 頂部邊距
		MarginBottom: 10,                          // 底部邊距
	}
}

// DrawTitle 繪製圖表標題
func (title ChartTitle) DrawTitle(c *freetype.Context, imgWidth, imgHeight int, titleText string) {
	// 設定字型大小
	c.SetFontSize(float64(title.FontSize))

	// 設定顏色
	c.SetSrc(image.NewUniform(title.Color))

	// 計算標題位置
	titleX := (imgWidth * title.X / 100)
	titleY := (imgHeight * title.Y / 100)

	// 根據對齊方式調整X座標
	switch title.Align {
	case "center":
		// 估算標題寬度並置中
		titleWidth := len(titleText) * (title.FontSize / 2)
		titleX = (imgWidth - titleWidth) / 2
	case "right":
		titleWidth := len(titleText) * (title.FontSize / 2)
		titleX = imgWidth - titleWidth - 20
	}

	// 繪製標題
	c.DrawString(titleText, freetype.Pt(titleX, titleY))
}

// PerformanceData 績效資料結構
type PerformanceData struct {
	Period      string
	PeriodName  string
	Performance string
}

// RevenueChartData 營收圖表資料結構
type RevenueChartData struct {
	Period        string  // 期間 (YYYY/MM)
	PeriodName    string  // 顯示名稱
	Revenue       int64   // 月營收 (千元)
	YoY           float64 // 年增率 (%)
	StockPrice    float64 // 股價
	LatestRevenue int64   // 最新月營收 (用於顯示)
	LatestYoY     float64 // 最新年增率 (用於顯示)
}

// ChartConfig 圖表設定
type ChartConfig struct {
	Title      string
	Width      int
	Height     int
	ShowGrid   bool
	ShowLegend bool
}

// DefaultChartConfig 預設圖表設定
func DefaultChartConfig() ChartConfig {
	return ChartConfig{
		Title:      "股票績效表現",
		Width:      1600, // 增加寬度
		Height:     800,  // 增加高度
		ShowGrid:   true,
		ShowLegend: true,
	}
}

// 生成績效圖表 (PNG格式)
func GeneratePerformanceChartPNG(data []PerformanceData, config ChartConfig) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("無績效資料可生成圖表")
	}

	// 取得顏色和標題配置
	colors := DefaultChartColors()
	titleConfig := DefaultChartTitle()

	// 建立圖片
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	// 填充白色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{colors.BackgroundWhite}, image.Point{}, draw.Src)

	// 載入字型 - 使用 go-findfont 動態查找支援中文的字型
	ttf, err := LoadChineseFont()
	if err != nil {
		return nil, fmt.Errorf("載入字型失敗: %v", err)
	}

	// 建立 freetype context
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ttf)
	c.SetFontSize(14)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(colors.TextDarkGray))
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

	// 圖表區域
	chartLeft := 120
	chartTop := 100
	chartWidth := config.Width - 300 // 增加右側空間
	chartHeight := config.Height - 200

	// 繪製標題
	titleConfig.DrawTitle(c, config.Width, config.Height, config.Title)

	// 重設字型大小
	c.SetFontSize(14)

	// 繪製座標軸
	drawLine(img, chartLeft, chartTop, chartLeft, chartTop+chartHeight, colors.AxisBlack)                        // Y軸
	drawLine(img, chartLeft, chartTop+chartHeight, chartLeft+chartWidth, chartTop+chartHeight, colors.AxisBlack) // X軸

	// 繪製格線和軸標籤
	if config.ShowGrid {
		// 繪製 Y 軸標籤（績效百分比）
		yGridLines := 5
		for i := 0; i <= yGridLines; i++ {
			y := chartTop + (chartHeight * i / yGridLines)
			value := maxVal - ((maxVal - minVal) * float64(i) / float64(yGridLines))

			// 水平格線
			if i > 0 && i < yGridLines {
				drawLine(img, chartLeft, y, chartLeft+chartWidth, y, colors.GridLightGray)
			}

			// Y軸標籤
			label := fmt.Sprintf("%.1f%%", value)
			pt := freetype.Pt(chartLeft-100, y+5)
			c.DrawString(label, pt)
		}

		// 繪製 X 軸標籤（時間）
		labelStep := 1
		if len(data) > 15 {
			labelStep = len(data) / 10
		}
		if labelStep < 2 {
			labelStep = 2 // 最小間隔2個資料點
		}

		for i, item := range data {
			x := chartLeft + (chartWidth * i / (len(data) - 1))

			// 垂直格線 - 減少格線數量
			if i > 0 && i < len(data)-1 && i%labelStep == 0 {
				drawLine(img, x, chartTop, x, chartTop+chartHeight, colors.GridLightGray)
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
				// 確保時間標籤使用黑色
				c.SetSrc(image.NewUniform(colors.TextBlack))
				c.SetFontSize(14)
				// 調整標籤位置，讓文字更居中，增加垂直間距
				pt := freetype.Pt(x-25, chartTop+chartHeight+25)
				c.DrawString(label, pt)

				// 在時間標籤下方顯示累積績效
				value := values[i]
				performanceText := fmt.Sprintf("%.2f%%", value)

				// 設定績效數值顏色 - 負值顯示綠色，正值顯示紅色
				c.SetFontSize(12)
				if value < 0 {
					c.SetSrc(image.NewUniform(colors.PositiveGreen)) // 負值用綠色
				} else {
					c.SetSrc(image.NewUniform(colors.NegativeRed)) // 正值用紅色
				}
				c.DrawString(performanceText, freetype.Pt(x-20, chartTop+chartHeight+45))

				// 重設為黑色，避免影響後續文字
				c.SetSrc(image.NewUniform(colors.TextBlack))
			}
		}
	}

	// 折線圖（X軸=時間，Y軸=績效）
	lineColor := colors.NegativeRed // 深一點的粉紅色

	// 找出最高績效和最低績效的索引
	var maxIndex, minIndex int
	maxValue := values[0]
	minValue := values[0]
	for i, value := range values {
		if value > maxValue {
			maxValue = value
			maxIndex = i
		}
		if value < minValue {
			minValue = value
			minIndex = i
		}
	}

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
			drawThickLine(img, prevX, prevY, x, y, 3, lineColor) // 使用3像素粗線
		}
	}

	// 標示最高績效
	if len(values) > 0 {
		maxX := chartLeft + (chartWidth * maxIndex / (len(data) - 1))
		maxY := chartTop + chartHeight - int((maxValue-minVal)/(maxVal-minVal)*float64(chartHeight))

		c.SetFontSize(16)
		c.SetSrc(image.NewUniform(colors.NegativeRed)) // 深粉紅色
		maxLabel := fmt.Sprintf("最高: %.2f%%", maxValue)
		c.DrawString(maxLabel, freetype.Pt(maxX-30, maxY-35)) // 調整位置避免與數值重疊

		// 繪製指向最高點的小圓圈
		drawCircle(img, maxX, maxY, 4, colors.NegativeRed)
	}

	// 標示最低績效
	if len(values) > 0 {
		minX := chartLeft + (chartWidth * minIndex / (len(data) - 1))
		minY := chartTop + chartHeight - int((minValue-minVal)/(maxVal-minVal)*float64(chartHeight))

		c.SetFontSize(16)
		c.SetSrc(image.NewUniform(colors.PositiveGreen)) // 深薄荷綠
		minLabel := fmt.Sprintf("最低: %.2f%%", minValue)
		c.DrawString(minLabel, freetype.Pt(minX-30, minY+35)) // 調整位置避免與數值重疊

		// 繪製指向最低點的小圓圈
		drawCircle(img, minX, minY, 4, colors.PositiveGreen)

		// 重設字型顏色
		c.SetSrc(image.NewUniform(colors.TextBlack))
		c.SetFontSize(14)
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
		drawDashedLine(img, chartLeft, zeroY, chartLeft+chartWidth, zeroY, colors.GridDashedGray)
	}

	// X軸標籤 - 顯示在圖表最右側，與日期文字同一水平線
	c.SetFontSize(14) // 使用與日期文字相同的字型大小
	c.SetSrc(image.NewUniform(colors.TextBlack))
	ptX := freetype.Pt(chartLeft+chartWidth+50, chartTop+chartHeight+25)
	c.DrawString("Time", ptX)

	// Y軸標籤
	pt2 := freetype.Pt(chartLeft-50, chartTop-10)
	c.DrawString("Performance (%)", pt2)

	// 右側標籤 - 累計績效，與績效數值同一水平線
	c.SetFontSize(12) // 使用與績效數值相同的字型大小
	c.SetSrc(image.NewUniform(colors.TextBlack))
	pt3 := freetype.Pt(chartLeft+chartWidth+50, chartTop+chartHeight+45)
	c.DrawString("累計績效", pt3)

	// 將圖片編碼為 PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("編碼 PNG 失敗: %v", err)
	}

	return buf.Bytes(), nil
}

// 生成折線圖 (PNG格式)
func GeneratePerformanceLineChart(data []PerformanceData, title string) ([]byte, error) {
	config := DefaultChartConfig()
	config.Title = title
	return GeneratePerformanceChartPNG(data, config)
}

// 生成營收圖表 (柱狀圖+折線圖組合)
func GenerateRevenueChart(data []RevenueChartData, stockName string, stockCode string) ([]byte, error) {

	if len(data) == 0 {
		return nil, fmt.Errorf("無營收資料可生成圖表")
	}

	// 取得顏色和標題配置
	colors := DefaultChartColors()
	titleConfig := DefaultChartTitle()

	// 圖表設定
	config := ChartConfig{
		Title:      fmt.Sprintf("%s (%s) 月營收", stockName, stockCode),
		Width:      1600, // 增加寬度以容納更多資訊
		Height:     800,  // 增加高度
		ShowGrid:   true,
		ShowLegend: true,
	}

	// 建立圖片
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	// 填充淺灰色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{colors.BackgroundLightGray}, image.Point{}, draw.Src)

	// 載入字型 - 使用 go-findfont 動態查找支援中文的字型
	ttf, err := LoadChineseFont()
	if err != nil {
		return nil, fmt.Errorf("載入字型失敗: %v", err)
	}

	// 建立 freetype context
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ttf)
	c.SetFontSize(16)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(colors.TextDarkGray))

	// 計算營收和年增率的範圍
	minRevenue := data[0].Revenue
	maxRevenue := data[0].Revenue
	minYoY := data[0].YoY
	maxYoY := data[0].YoY

	for _, item := range data {
		if item.Revenue < minRevenue {
			minRevenue = item.Revenue
		}
		if item.Revenue > maxRevenue {
			maxRevenue = item.Revenue
		}
		if item.YoY < minYoY {
			minYoY = item.YoY
		}
		if item.YoY > maxYoY {
			maxYoY = item.YoY
		}
	}

	// 營收範圍調整
	revenueMargin := (maxRevenue - minRevenue) * 1 / 5
	if revenueMargin == 0 {
		revenueMargin = maxRevenue / 5
	}
	minRevenue = 0 // 營收從0開始顯示
	maxRevenue += revenueMargin

	// 年增率範圍調整
	yoyMargin := (maxYoY - minYoY) * 1 / 10
	if yoyMargin == 0 {
		yoyMargin = 10
	}
	minYoY -= yoyMargin
	maxYoY += yoyMargin

	// 圖表區域
	chartLeft := 120
	chartTop := 120
	chartWidth := config.Width - 340
	chartHeight := config.Height - 250

	// 繪製標題
	titleConfig.DrawTitle(c, config.Width, config.Height, config.Title)

	// 在右上角顯示最新數據
	latestData := data[len(data)-1]
	c.SetFontSize(16) // 增加資訊字型大小

	// 顯示股票代碼和名稱
	stockInfoText := fmt.Sprintf("%s (%s)", stockName, stockCode)
	pt := freetype.Pt(config.Width-300, 30)
	c.DrawString(stockInfoText, pt)

	infoText := fmt.Sprintf("%s 營收: %.0f億", latestData.PeriodName, float64(latestData.LatestRevenue)/100000)
	pt = freetype.Pt(config.Width-300, 60)
	c.DrawString(infoText, pt)

	yoyText := fmt.Sprintf("YoY: %.2f%%", latestData.LatestYoY)
	pt = freetype.Pt(config.Width-300, 90) // 調整位置
	// YoY數據使用紅色
	c.SetSrc(image.NewUniform(colors.TextRed)) // 紅色
	c.DrawString(yoyText, pt)
	c.SetSrc(image.NewUniform(colors.TextDarkGray)) // 重設為黑色

	c.SetFontSize(16) // 調整基礎字型大小

	// 繪製座標軸
	drawLine(img, chartLeft, chartTop, chartLeft, chartTop+chartHeight, colors.AxisBlack)                        // Y軸
	drawLine(img, chartLeft, chartTop+chartHeight, chartLeft+chartWidth, chartTop+chartHeight, colors.AxisBlack) // X軸

	// 繪製右側Y軸（年增率）
	drawLine(img, chartLeft+chartWidth, chartTop, chartLeft+chartWidth, chartTop+chartHeight, colors.AxisBlack)

	// 繪製格線和軸標籤
	if config.ShowGrid {
		// 左側Y軸標籤（營收）
		yGridLines := 5
		for i := 0; i <= yGridLines; i++ {
			y := chartTop + (chartHeight * i / yGridLines)
			value := maxRevenue - ((maxRevenue - minRevenue) * int64(i) / int64(yGridLines))

			// 水平格線
			if i > 0 && i < yGridLines {
				drawLine(img, chartLeft, y, chartLeft+chartWidth, y, colors.TextBlack) // 使用黑色格線
			}

			// 左側Y軸標籤 (營收，單位：億)
			label := fmt.Sprintf("%.0f億", float64(value)/100000)
			pt := freetype.Pt(chartLeft-100, y+5)
			c.DrawString(label, pt)
		}

		// 右側Y軸標籤（年增率）
		for i := 0; i <= yGridLines; i++ {
			y := chartTop + (chartHeight * i / yGridLines)
			value := maxYoY - ((maxYoY - minYoY) * float64(i) / float64(yGridLines))

			// 右側Y軸標籤 (年增率)
			label := fmt.Sprintf("%.0f%%", value)
			pt := freetype.Pt(chartLeft+chartWidth+10, y+5)
			c.DrawString(label, pt)
		}

		// X軸標籤（時間）
		labelStep := 1
		if len(data) > 12 {
			labelStep = 2
		}

		for i, item := range data {
			x := chartLeft + (chartWidth * i / (len(data) - 1))

			// 垂直格線
			if i > 0 && i < len(data)-1 && i%labelStep == 0 {
				drawLine(img, x, chartTop, x, chartTop+chartHeight, colors.TextBlack) // 使用黑色格線
			}

			// X軸標籤
			if i%labelStep == 0 || i == len(data)-1 {
				label := item.PeriodName
				pt := freetype.Pt(x-15, chartTop+chartHeight+25)
				c.DrawString(label, pt)
			}
		}
	}

	// 繪製柱狀圖（營收）
	barWidth := chartWidth / len(data) * 6 / 10 // 60% 寬度
	barSpacing := chartWidth / len(data)
	barColor := colors.PositiveGreen // 深一點的薄荷綠

	for i, item := range data {
		// 計算柱狀圖位置
		x := chartLeft + (barSpacing * i) + (barSpacing-barWidth)/2
		barHeight := int((item.Revenue - minRevenue) * int64(chartHeight) / (maxRevenue - minRevenue))
		y := chartTop + chartHeight - barHeight

		// 繪製柱狀圖
		if barHeight > 0 {
			drawRect(img, x, y, barWidth, barHeight, barColor)
		}

		// 在柱狀圖上方顯示營收數字
		if item.Revenue > 0 {
			// 格式化營收數字（單位：億）
			revenueText := fmt.Sprintf("%.0f", float64(item.Revenue)/100000)

			// 計算文字位置（柱狀圖中心上方）
			textX := x + barWidth/2 - len(revenueText)*3 // 估算文字寬度的一半
			textY := y - 5                               // 柱狀圖上方5像素

			// 確保文字不會超出圖表頂部
			if textY < chartTop+15 {
				textY = chartTop + 15
			}

			// 設定較小的字型來顯示數值
			c.SetFontSize(14) // 增加營收數字字型大小
			pt := freetype.Pt(textX, textY)
			c.DrawString(revenueText, pt)
			c.SetFontSize(14) // 重設為原始字型大小
		}
	}

	// 繪製折線圖（年增率）
	lineColor := colors.NegativeRed

	for i := range data {
		value := data[i].YoY
		x := chartLeft + (chartWidth * i / (len(data) - 1))
		y := chartTop + chartHeight - int((value-minYoY)/(maxYoY-minYoY)*float64(chartHeight))

		// 繪製資料點
		drawCircle(img, x, y, 5, lineColor) // 增大資料點

		// 在資料點上方顯示YoY百分比
		yoyText := fmt.Sprintf("%.1f%%", value)

		// 計算文字位置（資料點上方）
		textX := x - len(yoyText)*3 // 估算文字寬度的一半
		textY := y - 15             // 資料點上方15像素

		// 確保文字不會超出圖表頂部
		if textY < chartTop+15 {
			textY = y + 25 // 如果上方空間不足，顯示在下方
		}

		// 設定字型來顯示數值
		c.SetFontSize(16) // 增加年增率字型大小
		// YoY數字統一使用紅色
		c.SetSrc(image.NewUniform(colors.TextRed))
		pt := freetype.Pt(textX, textY)
		c.DrawString(yoyText, pt)
		c.SetSrc(image.NewUniform(colors.TextBlack)) // 重設為黑色
		c.SetFontSize(16)                            // 重設為原始字型大小

		// 繪製粗線段 (除了第一個點)
		if i > 0 {
			prevValue := data[i-1].YoY
			prevX := chartLeft + (chartWidth * (i - 1) / (len(data) - 1))
			prevY := chartTop + chartHeight - int((prevValue-minYoY)/(maxYoY-minYoY)*float64(chartHeight))
			drawThickLine(img, prevX, prevY, x, y, 4, lineColor) // 使用4像素粗的線
		}
	}

	// 零線 (年增率)
	zeroY := chartTop + chartHeight - int((0-minYoY)/(maxYoY-minYoY)*float64(chartHeight))
	if minYoY < 0 && maxYoY > 0 {
		drawDashedLine(img, chartLeft, zeroY, chartLeft+chartWidth, zeroY, colors.GridDashedGray)
	}

	// 圖例
	if config.ShowLegend {
		legendY := config.Height - 80
		// 營收圖例
		drawRect(img, chartLeft, legendY, 15, 15, barColor)
		pt := freetype.Pt(chartLeft+25, legendY+12)
		c.DrawString("營收", pt)

		// 年增率圖例
		drawCircle(img, chartLeft+100, legendY+7, 5, lineColor)                             // 增大圓點
		drawThickLine(img, chartLeft+85, legendY+7, chartLeft+115, legendY+7, 4, lineColor) // 使用更粗線
		pt = freetype.Pt(chartLeft+125, legendY+12)
		c.DrawString("YoY", pt)
	}

	// X軸標籤 - 移到X軸右端，避免與時間標籤重疊
	ptX := freetype.Pt(chartLeft+chartWidth+80, chartTop+chartHeight+25)
	c.DrawString("Time", ptX)

	// 左側Y軸標籤 - 移到Y軸上端
	pt1 := freetype.Pt(chartLeft-50, chartTop-10)
	c.DrawString("營收 (億)", pt1)

	// 右側Y軸標籤 - 移到右側Y軸上端
	pt2 := freetype.Pt(chartLeft+chartWidth+50, chartTop-10)
	c.DrawString("YoY (%)", pt2)

	// 將圖片編碼為 PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("編碼 PNG 失敗: %v", err)
	}

	return buf.Bytes(), nil
}

// 生成K線圖 (PNG格式)
func GenerateCandlestickChart(data []CandlestickData, stockName string, symbol string) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("無K線資料可生成圖表")
	}

	// 取得顏色和標題配置
	colors := DefaultChartColors()
	titleConfig := DefaultChartTitle()

	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}

	var highestHigh, lowestLow float64
	var highestIndex, lowestIndex int
	if len(data) > 0 {
		highestHigh = data[0].High
		lowestLow = data[0].Low
		for i, d := range data {
			if d.High > highestHigh {
				highestHigh = d.High
				highestIndex = i
			}
			if d.Low < lowestLow {
				lowestLow = d.Low
				lowestIndex = i
			}
		}
	}

	config := DefaultChartConfig()
	config.Title = fmt.Sprintf("%s (%s) K線圖", stockName, symbol)
	config.Width = 1600
	config.Height = 900

	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))
	draw.Draw(img, img.Bounds(), &image.Uniform{colors.BackgroundWhite}, image.Point{}, draw.Src)

	// 載入字型
	ttf, err := LoadChineseFont()
	if err != nil {
		return nil, fmt.Errorf("載入字型失敗: %v", err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ttf)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(colors.TextDarkGray))

	// 找出價格和成交量的最大最小值
	minPrice := data[0].Low
	maxPrice := data[0].High
	maxVolume := data[0].Volume
	for _, d := range data {
		if d.High > maxPrice {
			maxPrice = d.High
		}
		if d.Low < minPrice {
			minPrice = d.Low
		}
		if d.Volume > maxVolume {
			maxVolume = d.Volume
		}
	}
	priceMargin := (maxPrice - minPrice) * 0.1
	maxPrice += priceMargin
	minPrice -= priceMargin
	if minPrice < 0 {
		minPrice = 0
	}

	chartLeft := 100
	chartTop := 100
	chartWidth := config.Width - 200
	chartHeight := config.Height - 300
	volumeHeight := 100
	priceChartHeight := chartHeight - volumeHeight - 50

	// 繪製標題
	titleConfig.DrawTitle(c, config.Width, config.Height, config.Title)

	// 繪製坐標軸
	drawLine(img, chartLeft, chartTop, chartLeft, chartTop+priceChartHeight, colors.AxisDarkGray)
	drawLine(img, chartLeft, chartTop+priceChartHeight, chartLeft+chartWidth, chartTop+priceChartHeight, colors.AxisDarkGray)

	// 繪製價格Y軸標籤
	yGridLines := 5
	c.SetFontSize(14)
	for i := 0; i <= yGridLines; i++ {
		y := chartTop + (priceChartHeight * i / yGridLines)
		price := maxPrice - (maxPrice-minPrice)*float64(i)/float64(yGridLines)
		label := fmt.Sprintf("%.2f", price)
		c.DrawString(label, freetype.Pt(chartLeft-60, y+5))
		if i > 0 && i < yGridLines {
			drawLine(img, chartLeft, y, chartLeft+chartWidth, y, colors.GridLightGray)
		}
	}

	// 繪製K線
	candleWidth := float64(chartWidth) / float64(len(data))
	for i, d := range data {
		x := chartLeft + int(candleWidth*float64(i)+candleWidth/2)
		highY := chartTop + int(float64(priceChartHeight)*(1-(d.High-minPrice)/(maxPrice-minPrice)))
		lowY := chartTop + int(float64(priceChartHeight)*(1-(d.Low-minPrice)/(maxPrice-minPrice)))
		openY := chartTop + int(float64(priceChartHeight)*(1-(d.Open-minPrice)/(maxPrice-minPrice)))
		closeY := chartTop + int(float64(priceChartHeight)*(1-(d.Close-minPrice)/(maxPrice-minPrice)))

		// 影線
		drawLine(img, x, highY, x, lowY, colors.KLineShadow)

		// 實體
		bodyWidth := int(candleWidth * 0.8)
		if bodyWidth < 1 {
			bodyWidth = 1
		}
		var candleColor color.RGBA
		if d.Close >= d.Open {
			candleColor = colors.KLineUpRed
			drawRect(img, x-bodyWidth/2, closeY, bodyWidth, openY-closeY, candleColor)
		} else {
			candleColor = colors.KLineDownGreen
			drawRect(img, x-bodyWidth/2, openY, bodyWidth, closeY-openY, candleColor)
		}

		// X軸標籤 - 顯示每月1號和該月均價
		if i == 0 || isFirstDayOfMonth(d.Date, data, i) {
			// 解析日期
			dateTime, err := time.Parse("2006-01-02", strings.Split(d.Date, "T")[0])
			if err == nil {
				monthLabel := dateTime.Format("1/2")

				// 繪製月份標籤
				c.SetFontSize(14)
				c.SetSrc(image.NewUniform(colors.TextBlack))
				c.DrawString(monthLabel, freetype.Pt(x-15, chartTop+priceChartHeight+20))

				// 只有不是第一個資料點時才顯示均價
				if i != 0 {
					// 計算該月均價
					monthAvg := calculateMonthlyAverage(data, i)
					avgLabel := fmt.Sprintf("%.2f", monthAvg)

					// 繪製均價標籤
					c.SetSrc(image.NewUniform(colors.MonthlyAvgRed))
					c.DrawString(avgLabel, freetype.Pt(x-20, chartTop+priceChartHeight+35))
				}

				// 繪製垂直虛線
				drawDashedVerticalLine(img, x, chartTop, chartTop+priceChartHeight, colors.GridLightGray)

				c.SetFontSize(14)
			}
		}
	}

	// 繪製成交量
	volumeChartTop := chartTop + priceChartHeight + 150
	drawLine(img, chartLeft, volumeChartTop, chartLeft, volumeChartTop+volumeHeight, colors.AxisDarkGray)
	drawLine(img, chartLeft, volumeChartTop+volumeHeight, chartLeft+chartWidth, volumeChartTop+volumeHeight, colors.AxisDarkGray)

	// 成交量Y軸標籤 (單位：千萬元)
	c.DrawString(fmt.Sprintf("%.1f千萬", maxVolume/10000000), freetype.Pt(chartLeft-80, volumeChartTop+5))
	c.DrawString("0", freetype.Pt(chartLeft-60, volumeChartTop+volumeHeight+5))

	for i, d := range data {
		x := chartLeft + int(candleWidth*float64(i)+candleWidth*0.1)
		barHeight := int(float64(volumeHeight) * (d.Volume / maxVolume))
		y := volumeChartTop + volumeHeight - barHeight
		barWidth := int(candleWidth * 0.8)
		if barWidth < 1 {
			barWidth = 1
		}

		var volColor color.RGBA
		if d.Close >= d.Open {
			volColor = colors.VolumeUpRed
		} else {
			volColor = colors.VolumeDownGreen
		}
		drawRect(img, x, y, barWidth, barHeight, volColor)
	}

	// 在最高價和最低價的K線上標示價格
	if len(data) > 0 {
		candleWidth := float64(chartWidth) / float64(len(data))

		// 標示最高價
		highX := chartLeft + int(candleWidth*float64(highestIndex)+candleWidth/2)
		highY := chartTop + int(float64(priceChartHeight)*(1-(highestHigh-minPrice)/(maxPrice-minPrice))) - 20
		c.SetFontSize(15)
		c.SetSrc(image.NewUniform(colors.HighestPriceRed))
		c.DrawString(fmt.Sprintf("最高: %.2f", highestHigh), freetype.Pt(highX-30, highY))

		// 標示最低價
		lowX := chartLeft + int(candleWidth*float64(lowestIndex)+candleWidth/2)
		lowY := chartTop + int(float64(priceChartHeight)*(1-(lowestLow-minPrice)/(maxPrice-minPrice))) + 30
		c.SetSrc(image.NewUniform(colors.LowestPriceGreen))
		c.DrawString(fmt.Sprintf("最低: %.2f", lowestLow), freetype.Pt(lowX-30, lowY))

		// 重設為黑色
		c.SetSrc(image.NewUniform(colors.TextBlack))
		c.SetFontSize(15)
	}

	c.SetFontSize(14)
	c.SetSrc(image.NewUniform(colors.MonthlyAvgRed))
	c.DrawString("月均價", freetype.Pt(chartLeft+chartWidth+10, chartTop+priceChartHeight+35))

	// 軸標籤
	c.SetSrc(image.NewUniform(colors.TextDarkGray))
	// X軸標籤 - 移到X軸右端
	c.DrawString("Time", freetype.Pt(chartLeft+chartWidth+10, chartTop+priceChartHeight+15))
	// Y軸標籤 - 移到Y軸上端
	c.DrawString("Price", freetype.Pt(chartLeft-30, chartTop-10))
	// 成交量Y軸標籤 - 移到成交量Y軸上端
	c.DrawString("Volume", freetype.Pt(chartLeft-40, volumeChartTop-10))

	buf := bytes.Buffer{}
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("編碼 PNG 失敗: %v", err)
	}
	return buf.Bytes(), nil
}

// CandlestickData K線資料結構
type CandlestickData struct {
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// isFirstDayOfMonth 檢查是否為月份的第一個交易日
func isFirstDayOfMonth(dateStr string, data []CandlestickData, currentIndex int) bool {
	if currentIndex == 0 {
		return true
	}

	currentDate, err := time.Parse("2006-01-02", strings.Split(dateStr, "T")[0])
	if err != nil {
		return false
	}

	prevDate, err := time.Parse("2006-01-02", strings.Split(data[currentIndex-1].Date, "T")[0])
	if err != nil {
		return false
	}

	return currentDate.Month() != prevDate.Month()
}

// calculateMonthlyAverage 計算該月的收盤價平均值
func calculateMonthlyAverage(data []CandlestickData, startIndex int) float64 {
	startDate, err := time.Parse("2006-01-02", strings.Split(data[startIndex].Date, "T")[0])
	if err != nil {
		return data[startIndex].Close
	}

	sum := 0.0
	count := 0

	for i := startIndex; i < len(data); i++ {
		currentDate, err := time.Parse("2006-01-02", strings.Split(data[i].Date, "T")[0])
		if err != nil {
			break
		}

		// 如果不是同一個月，停止計算
		if currentDate.Month() != startDate.Month() || currentDate.Year() != startDate.Year() {
			break
		}

		sum += data[i].Close
		count++
	}

	if count == 0 {
		return data[startIndex].Close
	}

	return sum / float64(count)
}

// drawDashedVerticalLine 繪製垂直虛線
func drawDashedVerticalLine(img *image.RGBA, x, y1, y2 int, col color.RGBA) {
	dash := 0
	for y := y1; y <= y2; y++ {
		if dash%8 < 4 { // 4像素實線，4像素空白
			img.Set(x, y, col)
		}
		dash++
	}
}

// drawThickLine 繪製粗線
func drawThickLine(img *image.RGBA, x1, y1, x2, y2, thickness int, col color.RGBA) {
	// 使用多條平行線來模擬粗線效果
	for t := -thickness / 2; t <= thickness/2; t++ {
		// 垂直和水平方向的偏移
		if abs(x2-x1) > abs(y2-y1) {
			// 主要是水平線，在垂直方向偏移
			drawLine(img, x1, y1+t, x2, y2+t, col)
		} else {
			// 主要是垂直線，在水平方向偏移
			drawLine(img, x1+t, y1, x2+t, y2, col)
		}
	}
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
