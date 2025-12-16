package dto

type ChartAsset struct {
	// 圖表標題
	Caption string
	// 檔案名稱
	FileName string
	// 圖表資料
	Data []byte
}
