package errors

import "errors"

var (
	// ErrNotFound 代表資源不存在。
	ErrNotFound = errors.New("not found")
	// ErrInvalidArgument 代表參數不合法。
	ErrInvalidArgument = errors.New("invalid argument")
)
