package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTManager struct {
	secret   string
	expHours time.Duration
}

func NewJWTManager(secret string, exp time.Duration) *JWTManager {
	return &JWTManager{
		secret:   secret,
		expHours: exp,
	}
}
func (m *JWTManager) Generate(userID uuid.UUID) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(m.expHours).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secret))
}
func (m *JWTManager) Verify(tokenStr string) (*jwt.Token, error) {

	return jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(m.secret), nil
	})
}
