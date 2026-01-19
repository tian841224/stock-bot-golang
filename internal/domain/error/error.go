package error

import "errors"

var (
	// ErrNotFound 代表資源不存在
	ErrNotFound = errors.New("resource not found")
	// ErrInvalidArgument 代表參數不合法
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrAlreadyExists 代表資源已存在
	ErrAlreadyExists = errors.New("resource already exists")
)
