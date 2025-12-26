package database

import (
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

func ConnectRedis() (*redis.Client, context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 1*time.Minute)
	return client, ctx
}
