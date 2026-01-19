package notification

import "context"

// ScheduleHandlerUsecase 提供給 scheduler 呼叫的入口（佔位）。
type ScheduleHandlerUsecase interface {
	RunScheduledTasks(ctx context.Context) error
}

type scheduleHandlerUsecase struct {
	notification SendNotificationUsecase
}

func NewScheduleHandlerUsecase(notification SendNotificationUsecase) ScheduleHandlerUsecase {
	return &scheduleHandlerUsecase{notification: notification}
}

func (u *scheduleHandlerUsecase) RunScheduledTasks(ctx context.Context) error {
	if u.notification == nil {
		return nil
	}
	return u.notification.SendDailyDigest(ctx)
}
