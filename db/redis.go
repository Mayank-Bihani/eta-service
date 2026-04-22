package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedis(redisUrl string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return client
}