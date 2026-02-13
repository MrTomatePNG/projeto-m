package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) (*JWTManager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is empty")
	}
	return &JWTManager{secret: []byte(secret), ttl: ttl}, nil
}

func (m *JWTManager) Generate(userID int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(m.ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Verify(tokenString string) (int64, error) {
	parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return m.secret, nil
	})
	if err != nil {
		return 0, err
	}
	if !parsed.Valid {
		return 0, jwt.ErrTokenInvalidClaims
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	switch v := sub.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	default:
		return 0, jwt.ErrTokenInvalidClaims
	}
}
