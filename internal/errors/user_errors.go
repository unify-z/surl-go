package errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserBanned         = errors.New("user is banned")
	ErrInvalidVerifyCode  = errors.New("invalid verification code")
)
