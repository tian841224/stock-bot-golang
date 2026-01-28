package notification

import (
	"context"
	"sync"

	"go.uber.org/multierr"
)

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

	tasks := []func(context.Context) error{
		u.notification.SendStockPriceNotification,
		u.notification.SendStockNewsNotification,
		u.notification.SendMarketInfoNotification,
		u.notification.SendTopVolumeNotification,
	}

	errChan := make(chan error, len(tasks))
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go func(t func(context.Context) error) {
			defer wg.Done()
			if err := t(ctx); err != nil {
				errChan <- err
			}
		}(task)
	}

	wg.Wait()
	close(errChan)

	var errs error
	for err := range errChan {
		errs = multierr.Append(errs, err)
	}
	return errs
}
