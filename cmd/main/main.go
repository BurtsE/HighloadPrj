package main

import (
	"highload/internal/analytics"
	"highload/internal/cache"
	"highload/internal/handlers"
	"log"
	"os"
)

func main() {
	redisClient := cache.NewRedisClient(os.Getenv("REDIS_URL"))
	analyzer := analytics.NewAnalyzer(50, 2.0) // window=50, threshold=2.0

	server := handlers.NewServer(redisClient, analyzer)

	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
