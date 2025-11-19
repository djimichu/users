package main

import (
	"JWTproject/internal/repository/postgreSQL"
	"log"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal("err:", err)
	}
	defer db.Close()
}
