package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

type rdConfig struct {
	Password string
	Host     string
	Port     string
	DbNumber int
}

func NewRdConfig(password string, host string, port string, dbNumber int) rdConfig {
	return rdConfig{
		Password: password,
		Host:     host,
		Port:     port,
		DbNumber: dbNumber,
	}
}

func NewClient(ctx context.Context, cfg *rdConfig) (*redis.Client, error) {
	address := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	client := redis.NewClient(
		&redis.Options{
			Addr:     address,
			Password: cfg.Password,
			DB:       cfg.DbNumber,
		})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}
