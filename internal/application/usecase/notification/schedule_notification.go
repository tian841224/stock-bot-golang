package notification

import (
	"context"
	"sync"

	"go.uber.org/multierr"

	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// ScheduleHandlerUsecase 提供給 scheduler 呼叫的入口（佔位）。
type ScheduleHandlerUsecase interface {
	RunScheduledTasks(ctx context.Context) error
}

type scheduleHandlerUsecase struct {
	notification SendNotificationUsecase
	log          logger.Logger
}

func NewScheduleHandlerUsecase(notification SendNotificationUsecase, log logger.Logger) ScheduleHandlerUsecase {
	return &scheduleHandlerUsecase{
		notification: notification,
		log:          log,
	}
}

func (u *scheduleHandlerUsecase) RunScheduledTasks(ctx context.Context) error {
	if u.notification == nil {
		return nil
	}

	tasks := []struct {
		Name string
		Func func(context.Context) error
	}{
		{"SendStockPriceNotification", u.notification.SendStockPriceNotification},
		{"SendStockNewsNotification", u.notification.SendStockNewsNotification},
		{"SendMarketInfoNotification", u.notification.SendMarketInfoNotification},
		{"SendTopVolumeNotification", u.notification.SendTopVolumeNotification},
	}

	errChan := make(chan error, len(tasks))
	var wg sync.WaitGroup

	for _, task := range tasks {
		// 紀錄顯示當前執行的服務
		u.log.Info("正在執行排程任務...", logger.String("task", task.Name))

		wg.Add(1)
		go func(t func(context.Context) error) {
			defer wg.Done()
			if err := t(ctx); err != nil {
				errChan <- err
			}
		}(task.Func)
	}

	wg.Wait()
	close(errChan)

	var errs error
	for err := range errChan {
		errs = multierr.Append(errs, err)
	}
	return errs
}
