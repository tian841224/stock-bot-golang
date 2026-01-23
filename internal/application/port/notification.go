package port

import "context"

type NotificationPort interface {
	// 推送股票股價
	SendStockPriceNotification(ctx context.Context) error
	// 推送股票新聞
	SendStockNewsNotification(ctx context.Context) error
	// 推送大盤資訊
	SendMarketInfoNotification(ctx context.Context) error
	//推送交易量排行
	SendTopVolumeNotification(ctx context.Context) error
}
