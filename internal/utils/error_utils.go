package utils

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/unify-z/go-surl/internal/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return true
	}
	return false
}

func IsInternalServerError(err error) bool {
	if err == nil {
		return false
	}
	var hashVersionTooNewError bcrypt.HashVersionTooNewError
	var invalidCostError bcrypt.InvalidCostError
	var invalidHashPrefixError bcrypt.InvalidHashPrefixError
	if errors.As(err, &hashVersionTooNewError) || errors.As(err, &invalidCostError) || errors.As(err, &invalidHashPrefixError) {
		logger.Error(err.Error())
		return true
	}
	switch {
	case errors.Is(err, gorm.ErrInvalidTransaction),
		errors.Is(err, gorm.ErrNotImplemented),
		errors.Is(err, gorm.ErrMissingWhereClause),
		errors.Is(err, gorm.ErrRegistered),
		errors.Is(err, gorm.ErrModelValueRequired),
		errors.Is(err, gorm.ErrPrimaryKeyRequired),
		errors.Is(err, gorm.ErrInvalidField),
		errors.Is(err, gorm.ErrInvalidData),
		errors.Is(err, gorm.ErrUnsupportedDriver),
		errors.Is(err, gorm.ErrUnsupportedRelation),
		errors.Is(err, gorm.ErrDryRunModeUnsupported):
		logger.Error(err.Error())
		return true
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return false
	}
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return false
	}
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return false
	}
	if errors.Is(err, io.EOF) {

		return false
	}
	msg := err.Error()
	if strings.Contains(msg, "unsupported destination, should be a struct") ||
		strings.Contains(msg, "unsupported unmarshal") ||
		strings.Contains(msg, "should be a struct") ||
		strings.Contains(msg, "unsupported destination") ||
		strings.Contains(msg, "reflect:") || // reflect errors
		strings.Contains(msg, "gomail:") || // gomail errors
		strings.Contains(msg, "invalid bind") ||
		strings.Contains(msg, "invalid destination") ||
		strings.Contains(msg, "cannot unmarshal") && !strings.Contains(msg, "cannot unmarshal string into Go value of type") {
		logger.Error(err.Error())
		return true
	}
	return false
}
