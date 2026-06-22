package service

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrAppNotFound     = errors.New("app not found")
)
