package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type pgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	PoolSize string
}

func NewPgConfig(username string, password string, host string, port string, database string, poolSize string) *pgConfig {
	return &pgConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
		PoolSize: poolSize,
	}
}

func NewClient(ctx context.Context, maxAttempts int, delay time.Duration, cfg *pgConfig) (pool *pgxpool.Pool, err error) {
	connectUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.PoolSize)
	err = DoWithAttempts(func() error {
		config, err := pgxpool.ParseConfig(connectUrl)
		if err != nil {
			log.Fatalf("Failed while parsing config: %v\n", err)
		}
		pool, err = pgxpool.ConnectConfig(ctx, config)
		if err != nil {
			log.Println("Connection failed... Going to do the next attempt")

			return err
		}
		return nil
	}, maxAttempts, delay)
	if err != nil {
		log.Fatal("cannot connect to Postgresql")
	}
	log.Println("Connection has been established")
	return pool, nil

}

func DoWithAttempts(f func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = f(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil

	}
	return
}
