package port

import (
	"context"

	"github.com/tian841224/stock-bot/internal/domain/entity"
	"github.com/tian841224/stock-bot/internal/domain/valueobject"
)

// UserAccountPort 提供 bot usecase 查詢或建立使用者的能力。
type UserAccountPort interface {
	GetOrCreate(ctx context.Context, accountID string, userType valueobject.UserType) (*entity.User, error)
}
