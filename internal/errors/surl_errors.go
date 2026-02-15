package errors

import "errors"

var (
	ErrShortURLNotFound = errors.New("short URL not found")
	ErrInvalidShortCode = errors.New("invalid short code")
)
