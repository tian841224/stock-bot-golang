package port

import "context"

// NotificationSender 定義通知發送介面
type NotificationSender interface {
	// SendMessage 發送文字訊息
	SendMessage(ctx context.Context, userID string, message string) error
	// SendImage 發送圖片
	SendImage(ctx context.Context, userID string, imageURL string, caption string) error
}
