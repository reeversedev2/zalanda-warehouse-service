package cache

import (
	"errors"
	"os"

	"github.com/go-redis/redis"
)

type Redis struct {
	RedisClient redis.Client
}

func NewRedis() Redis {

	var client = redis.NewClient(&redis.Options{
		// Container name + port since we are using docker
		Addr:     "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if client == nil {
		errors.New("Failed to connect to Redis")
	}

	return Redis{
		RedisClient: *client,
	}
}
