package httpx

import (
	"JWTproject/internal/auth"
	"JWTproject/internal/logger"
	"JWTproject/internal/repository"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userIDKey = contextKey("userID")

func JWTAuthMiddleware(jwtManager *auth.JWTManager, userRepo *repository.UserRepo) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			// заголовок запроса
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Logger.Warn("auth header is missing",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.String("ip", r.RemoteAddr),
				)
				http.Error(w, "missing auth header request", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Logger.Warn("invalid auth header format",
					zap.String("path", r.URL.Path),
					zap.String("ip", r.RemoteAddr),
					zap.String("method", r.Method),
				)
				http.Error(w, "invalid auth header", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			//проверяем токен
			token, err := jwtManager.Verify(tokenStr)
			if err != nil || !token.Valid {
				logger.Logger.Warn("invalid token",
					zap.Error(err),
					zap.String("path", r.URL.Path),
					zap.String("ip", r.RemoteAddr),
					zap.String("method", r.Method),
				)
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			//достаем userID из claims(payload)
			claims, ok := token.Claims.(jwt.MapClaims) //claims это map, где ключ к примеру user_id, а значение ID
			if !ok {
				logger.Logger.Warn("invalid token claims",
					zap.String("path", r.URL.Path),
					zap.String("ip", r.RemoteAddr),
					zap.String("method", r.Method),
				)
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			idStr, ok := claims["user_id"].(string)
			if !ok {
				logger.Logger.Warn("userID not found in token",
					zap.String("path", r.URL.Path),
					zap.String("ip", r.RemoteAddr),
					zap.String("method", r.Method),
				)
				http.Error(w, "userID not found in token", http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(idStr)
			if err != nil {
				logger.Logger.Warn("invalid userID in token",
					zap.Error(err),
					zap.String("path", r.URL.Path),
					zap.String("id", idStr),
					zap.String("ip", r.RemoteAddr),
					zap.String("method", r.Method),
				)
				http.Error(w, "invalid userID in token", http.StatusUnauthorized)
				return
			}
			//проверяем наличие пользователя в БД
			user, err := userRepo.GetUserByID(userID)
			if err != nil || user == "" {
				logger.Logger.Error("user not found",
					zap.Error(err),
					zap.String("path", r.URL.Path),
					zap.String("user_id", userID.String()),
				)
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)

			logger.Logger.Info("auth successful",
				zap.String("path", r.URL.Path),
				zap.String("user_id", userID.String()),
				zap.String("ip", r.RemoteAddr),
				zap.String("method", r.Method),
				zap.Duration("duration", time.Since(start)),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
