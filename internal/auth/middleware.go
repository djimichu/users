package auth

import (
	"JWTproject/internal/repository"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey = contextKey("userID")

func JWTAuthMiddleware(jwtManager *JWTManager, userRepo *repository.UserRepo) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// заголовок запроса
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing auth header request", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid auth header", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			//проверяем токен
			token, err := jwtManager.Verify(tokenStr)
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			//достаем userID из claims(payload)
			claims, ok := token.Claims.(jwt.MapClaims) //claims это map, где ключ к примеру user_id, а значение ID
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
			}

			idStr, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "userID not found in token", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(idStr)
			if err != nil {
				http.Error(w, "invalid userID in token", http.StatusUnauthorized)
				return
			}
			//проверяем наличие пользователя в БД
			user, err := userRepo.GetUserByID(userID)
			if err != nil || user == "" {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
