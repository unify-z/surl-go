package helpers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	intlerror "github.com/unify-z/go-surl/internal/errors"
)

type JWTHelper struct {
	secret   string
	duration int
}

func NewJWTHelper(secret string, duration int) *JWTHelper {
	return &JWTHelper{
		secret:   secret,
		duration: duration,
	}
}

func (h *JWTHelper) GenerateToken(customClaims map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range customClaims {
		claims[k] = v
	}
	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(time.Second * time.Duration(h.duration)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.secret))
}

func (h *JWTHelper) GenerateUserToken(userID uint, username string, isAdmin bool) (string, error) {
	claims := map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"is_admin": isAdmin,
	}
	return h.GenerateToken(claims)
}

func (h *JWTHelper) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.secret), nil
	})

	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, intlerror.ErrTokenExpired
			}
		}
		return nil, intlerror.ErrInvalidToken
	}

	if !token.Valid {
		return nil, intlerror.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, intlerror.ErrInvalidToken
}
