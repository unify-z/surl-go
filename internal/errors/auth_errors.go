package errors

import "errors"

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrTokenExpired         = errors.New("token has expired")
	ErrInvalidToken         = errors.New("invalid token")
)
