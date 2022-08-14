package main

import (
	"context"
	"fmt"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/adapters/db/psqlStorage"
	"github.com/VrMolodyakov/vote-service/internal/config"
	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/pkg/client/postgresql"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

func main() {
	fmt.Println("start")
	logger := logging.GetLogger("info")
	cfg := config.GetConfig()
	pgConfig := postgresql.NewPgConfig(
		cfg.PostgreSql.Username,
		cfg.PostgreSql.Password,
		cfg.PostgreSql.Host,
		cfg.PostgreSql.Port,
		cfg.PostgreSql.Dbname)
	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}
	repo := psqlStorage.NewChoiceStorage(pgClient, logger)
	v1 := entity.Choice{"a", 1, 1}
	repo.Insert(context.Background(), v1)
	if err != nil {
		logger.Info("inside save")
		logger.Fatal(err)
	}
	v2, err := repo.FindChoicesByVoteId(context.Background(), 1)
	if err != nil {
		logger.Info("inside save")
		logger.Fatal(err)
	}
	logger.Info(v2)
	logger.Info(cfg)
	logger.Info("end")
}

// rdCfg := redis.NewRdConfig(cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DbNumber)
// 	client, err := redis.NewClient(context.Background(), &rdCfg)
// 	if err != nil {
// 		logger.Fatal(err)
// 	}
// 	cch := choiceCache.NewChoiceCache(client, logger)
// 	choice := entity.Choice{"a", 1, 1}
// 	err = cch.Save("title", choice, 1*time.Second	)
// 	if err != nil {
// 		logger.Info("inside save")
// 		logger.Fatal(err)
// 	}
// 	c, err := cch.Get("title")
// 	if err != nil {
// 		logger.Info("inside get")
// 		logger.Fatal(err)
// 	}
// 	logger.Info(c)

// pgConfig := postgresql.NewPgConfig(
// 	cfg.PostgreSql.Username,
// 	cfg.PostgreSql.Password,
// 	cfg.PostgreSql.Host,
// 	cfg.PostgreSql.Port,
// 	cfg.PostgreSql.Dbname)
// pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
// if err != nil {
// 	logger.Fatal(err)
// }
// repo := psqlStorage.NewVoteStorage(pgClient, logger)
// v := entity.Vote{"title", 1}
// repo.Insert(context.Background(), v)
// v1, err := repo.FindIdByTitle(context.Background(), "title")
// if err != nil {
// 	logger.Info("inside save")
// 	logger.Fatal(err)
// }
// logger.Info(v1)
