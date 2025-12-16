package notification

import "context"

// SendNotificationUsecase 負責推播通知相關的應用邏輯（佔位，用以符合分層結構）。
type SendNotificationUsecase interface {
	SendDailyDigest(ctx context.Context) error
}

type sendNotificationUsecase struct{}

func NewSendNotificationUsecase() SendNotificationUsecase {
	return &sendNotificationUsecase{}
}

func (u *sendNotificationUsecase) SendDailyDigest(ctx context.Context) error {
	// TODO: implement notification delivery
	return nil
}
