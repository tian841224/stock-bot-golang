package repository

import "gorm.io/gorm"

// QueryOption 查詢選項函數
type QueryOption func(*gorm.DB) *gorm.DB

// WithPreload 增加預加載關聯
func WithPreload(associations ...string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		for _, association := range associations {
			db = db.Preload(association)
		}
		return db
	}
}

// WithOrder 增加排序
func WithOrder(order string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

// WithLimit 增加限制
func WithLimit(limit int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// WithOffset 增加偏移
func WithOffset(offset int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

// WithPagination 增加分頁
func WithPagination(offset, limit int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}
