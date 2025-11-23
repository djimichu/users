package main

import (
	"JWTproject/internal/auth"
	"JWTproject/internal/config"
	"JWTproject/internal/httpx"
	"JWTproject/internal/logger"
	"JWTproject/internal/repository"
	"JWTproject/internal/repository/postgres"
	"JWTproject/internal/service/user"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {

	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "local"
	}

	if err := logger.Init(env); err != nil {
		log.Fatal("failed logger:", err)
	}
	defer logger.Close()

	// либо структурированный zap Logger, либо sugar для printf-подобных вызовов
	zl := logger.Logger

	zl.Info("application starting",
		zap.String("env", env),
	)

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpHours)

	connectDB, err := postgres.ConnectDB(cfg.DSN())
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepo(connectDB)

	userService := user.NewUserService(userRepo, jwtManager)

	httpHandlers := httpx.NewHTTPHandlers(userService)

	server := httpx.NewHTTPServer(httpHandlers)

	fmt.Println("server starting...")
	if err := server.Start(cfg.HTTPPort, jwtManager, userRepo); err != nil {
		panic(err)
	}
}
