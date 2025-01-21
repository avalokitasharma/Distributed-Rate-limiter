package main

import (
	"log"
	"net/http"
	"rate-limiter/server/api"
	"rate-limiter/server/config"
	"rate-limiter/server/redis"
)

func main() {
	cfg := config.NewConfig()

	redisClient, err := redis.NewClient(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Fatal(err)
	}

	handler := api.NewHandler(redisClient)
	router := api.SetupRoutes(handler)

	log.Printf("server is starting onport %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))

}
