package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/adapters/db/postgres"
	"github.com/VrMolodyakov/vote-service/internal/config"
	psqlClient "github.com/VrMolodyakov/vote-service/pkg/client/postgresql"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

func main() {
	fmt.Println("start")
	logger := logging.GetLogger("info")
	cfg := config.GetConfig()
	pgConfig := psqlClient.NewPgConfig(
		cfg.PostgreSql.Username,
		cfg.PostgreSql.Password,
		cfg.PostgreSql.Host,
		cfg.PostgreSql.Port,
		cfg.PostgreSql.Dbname)
	pgClient, err := psqlClient.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}
	repo := postgres.NewVoteStorage(pgClient, logger)
	fmt.Println(repo)

}
