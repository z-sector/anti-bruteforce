package internal

import "errors"

var (
	ErrInvalidArgs       = errors.New("invalid arguments")
	ErrBlackListExists   = errors.New("value for blacklist already exists")
	ErrBlackListNotFound = errors.New("value for blacklist not found")
	ErrWhiteListExists   = errors.New("value for whitelist already exists")
	ErrWhiteListNotFound = errors.New("value for whitelist not found")
	ErrInvalidIP         = errors.New("invalid ip")
)
