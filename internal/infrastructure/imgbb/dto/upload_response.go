// Package dto 提供 ImgBB API 的資料傳輸物件
package dto

import (
	"strconv"
)

// ImgBBUploadResponse 代表 ImgBB API 上傳回應
type ImgBBUploadResponse struct {
	Data    ImgBBData `json:"data"`
	Success bool      `json:"success"`
	Status  int       `json:"status"`
}

// ImgBBData 代表上傳圖片的主要資料
type ImgBBData struct {
	ID               string      `json:"id"`
	Title            string      `json:"title"`
	URLViewer        string      `json:"url_viewer"`
	URL              string      `json:"url"`
	DisplayURL       string      `json:"display_url"`
	Width            interface{} `json:"width"`
	Height           interface{} `json:"height"`
	Size             interface{} `json:"size"`
	Time             interface{} `json:"time"`
	Expiration       interface{} `json:"expiration"`
	OriginalFilename string      `json:"original_filename"`
	Image            ImgBBImage  `json:"image"`
	Thumb            ImgBBImage  `json:"thumb"`
	Medium           ImgBBImage  `json:"medium"`
	DeleteURL        string      `json:"delete_url"`
}

// ImgBBImage 代表圖片的詳細資訊
type ImgBBImage struct {
	Filename  string `json:"filename"`
	Name      string `json:"name"`
	Mime      string `json:"mime"`
	Extension string `json:"extension"`
	URL       string `json:"url"`
}

// GetWidth 安全地取得寬度
func (d *ImgBBData) GetWidth() int {
	return d.getIntValue(d.Width)
}

// GetHeight 安全地取得高度
func (d *ImgBBData) GetHeight() int {
	return d.getIntValue(d.Height)
}

// GetSize 安全地取得檔案大小
func (d *ImgBBData) GetSize() int {
	return d.getIntValue(d.Size)
}

// GetTime 安全地取得時間戳
func (d *ImgBBData) GetTime() int64 {
	return d.getInt64Value(d.Time)
}

// GetExpiration 安全地取得過期時間
func (d *ImgBBData) GetExpiration() int {
	return d.getIntValue(d.Expiration)
}

// getIntValue 將 interface{} 轉換為 int
func (d *ImgBBData) getIntValue(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal
		}
	}
	return 0
}

// getInt64Value 將 interface{} 轉換為 int64
func (d *ImgBBData) getInt64Value(value interface{}) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
			return intVal
		}
	}
	return 0
}
