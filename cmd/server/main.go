package main

import (
	"JWTproject/internal/auth"
	"JWTproject/internal/config"
	"JWTproject/internal/httpx"
	"JWTproject/internal/repository"
	"JWTproject/internal/repository/postgres"
	"JWTproject/internal/service/user"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("config initialised")

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpHours)
	fmt.Println("jwtManager initialised")

	connectDB, err := postgres.ConnectDB(cfg.DSN())
	if err != nil {
		panic(err)
	}
	fmt.Println("connect DB initialised")

	userRepo := repository.NewUserRepo(connectDB)
	fmt.Println("connect DB initialised")

	userService := user.NewUserService(userRepo, jwtManager)
	fmt.Println("userService initialised")

	httpHandlers := httpx.NewHTTPHandlers(userService)
	fmt.Println("httpHandlers initialised")

	server := httpx.NewHTTPServer(httpHandlers)
	fmt.Println("server initialised")

	fmt.Println("server starting...")
	if err := server.Start(cfg.HTTPPort, jwtManager, userRepo); err != nil {
		panic(err)
	}
}
