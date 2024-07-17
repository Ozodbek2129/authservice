package main

import (
	"auth/api"
	"auth/config"
	"auth/storage/postgres"
	"log"

	"github.com/go-redis/redis/v8"
)

func main() {
	db, err := postgres.ConnectDB()
	if err != nil {		
		log.Fatal(err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	router := api.NewRouter(db, redisClient)
	router.Run(config.Load().USER_ROUTER)
}
