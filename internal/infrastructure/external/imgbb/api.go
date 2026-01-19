package imgbb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tian841224/stock-bot/internal/infrastructure/external/imgbb/dto"
)

const (
	ImgBBBaseURL = "https://api.imgbb.com/1/upload"
)

// ImgBBClient ImgBB API 客戶端
type ImgBBClient struct {
	apiKey string
	client *http.Client
}

// NewImgBBClient 建立新的 ImgBB 客戶端
func NewImgBBClient(apiKey string) *ImgBBClient {
	return &ImgBBClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// UploadOptions 上傳選項
type UploadOptions struct {
	Name       string // 檔案名稱
	Expiration int    // 過期時間（秒），60-15552000
}

// UploadFromBase64 從 base64 字串上傳圖片
func (c *ImgBBClient) UploadFromBase64(base64Data string, options *UploadOptions) (*dto.ImgBBUploadResponse, error) {
	// 建立表單資料
	formData := url.Values{}
	formData.Set("key", c.apiKey)
	formData.Set("image", base64Data)

	if options != nil {
		if options.Name != "" {
			formData.Set("name", options.Name)
		}
		if options.Expiration > 0 {
			formData.Set("expiration", strconv.Itoa(options.Expiration))
		}
	}

	// 建立 POST 請求
	req, err := http.NewRequest("POST", ImgBBBaseURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("建立請求失敗: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.sendRequest(req)
}

// UploadFromFile 從檔案上傳圖片
func (c *ImgBBClient) UploadFromFile(file io.Reader, filename string, options *UploadOptions) (*dto.ImgBBUploadResponse, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 加入 API key
	if err := writer.WriteField("key", c.apiKey); err != nil {
		return nil, fmt.Errorf("寫入 API key 失敗: %w", err)
	}

	// 加入檔案
	fileWriter, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, fmt.Errorf("建立檔案欄位失敗: %w", err)
	}

	if _, err := io.Copy(fileWriter, file); err != nil {
		return nil, fmt.Errorf("複製檔案內容失敗: %w", err)
	}

	// 加入選項
	if options != nil {
		if options.Name != "" {
			if err := writer.WriteField("name", options.Name); err != nil {
				return nil, fmt.Errorf("寫入檔案名稱失敗: %w", err)
			}
		}
		if options.Expiration > 0 {
			if err := writer.WriteField("expiration", strconv.Itoa(options.Expiration)); err != nil {
				return nil, fmt.Errorf("寫入過期時間失敗: %w", err)
			}
		}
	}

	writer.Close()

	// 建立 POST 請求
	req, err := http.NewRequest("POST", ImgBBBaseURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("建立請求失敗: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.sendRequest(req)
}

// UploadFromURL 從 URL 上傳圖片
func (c *ImgBBClient) UploadFromURL(imageURL string, options *UploadOptions) (*dto.ImgBBUploadResponse, error) {
	// 建立表單資料
	formData := url.Values{}
	formData.Set("key", c.apiKey)
	formData.Set("image", imageURL)

	if options != nil {
		if options.Name != "" {
			formData.Set("name", options.Name)
		}
		if options.Expiration > 0 {
			formData.Set("expiration", strconv.Itoa(options.Expiration))
		}
	}

	// 建立 POST 請求
	req, err := http.NewRequest("POST", ImgBBBaseURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("建立請求失敗: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.sendRequest(req)
}

// sendRequest 發送請求並解析回應
func (c *ImgBBClient) sendRequest(req *http.Request) (*dto.ImgBBUploadResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("發送請求失敗: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("讀取回應失敗: %w", err)
	}

	var uploadResp dto.ImgBBUploadResponse
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		return nil, fmt.Errorf("解析 JSON 回應失敗: %w", err)
	}

	if !uploadResp.Success {
		return nil, fmt.Errorf("上傳失敗，狀態碼: %d", uploadResp.Status)
	}

	return &uploadResp, nil
}
