package bot

import (
	"context"
	"strings"
	"time"

	"github.com/tian841224/stock-bot/internal/application/port"
	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

// trackUserActivity 紀錄累計人次、確保使用者存在並更新其活躍時間；
// 供 LineMessageProcessor 與 TelegramMessageProcessor 共用，避免重複的
// GetOrCreate 呼叫與各自吞掉錯誤的問題。統計/活躍度失敗僅記錄，不中斷訊息處理。
func trackUserActivity(
	ctx context.Context,
	systemStatsRepo port.SystemStatisticsRepository,
	userAccountPort port.UserAccountPort,
	accountID string,
	userType valueobject.UserType,
	log logger.Logger,
) *entity.User {
	if err := systemStatsRepo.IncrementVisitCount(ctx); err != nil {
		log.Error("failed to increment visit count", logger.Error(err))
	}

	user, err := userAccountPort.GetOrCreate(ctx, accountID, userType)
	if err != nil {
		log.Error("failed to ensure user exists", logger.Error(err))
		return nil
	}

	if err := userAccountPort.UpdateActivity(ctx, user.ID); err != nil {
		log.Error("failed to update user activity", logger.Error(err))
	}

	return user
}

// parseMessageArgs 將訊息文字拆成命令與最多兩個參數
func parseMessageArgs(messageText string) (command, arg1, arg2 string) {
	parts := strings.Fields(messageText)
	if len(parts) == 0 {
		return "", "", ""
	}

	command = parts[0]
	if len(parts) > 1 {
		arg1 = parts[1]
	}
	if len(parts) > 2 {
		arg2 = parts[2]
	}
	return command, arg1, arg2
}

// parseDate 解析 YYYY-MM-DD 格式的日期字串
func parseDate(value string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", value, time.Local)
}
