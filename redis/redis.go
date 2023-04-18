package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"study02-chat-service/config"
)

func OpenNewRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:   0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CheckRedisConnection(cfg *config.RedisConfig) error {
	if c, err := OpenNewRedisClient(cfg); err != nil {
		return err
	} else {
		defer c.Close()
	}

	return nil
}
