package redis

import (
	"context"
	"fmt"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	goredis "github.com/redis/go-redis/v9"
)

func New(cfg *config.Config) (*goredis.Client, error) {

	client := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
