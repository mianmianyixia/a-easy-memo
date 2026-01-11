package database

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func ConnectRedis() (*redis.Client, context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	ctx := context.Background()
	return client, ctx
}
