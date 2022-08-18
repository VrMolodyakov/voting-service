package main

import (
	_ "context"
	_ "fmt"
	_ "time"

	"github.com/VrMolodyakov/vote-service/internal"
	_ "github.com/VrMolodyakov/vote-service/internal/adapter/db/choiceCache"
	"github.com/VrMolodyakov/vote-service/internal/config"
	_ "github.com/VrMolodyakov/vote-service/internal/config"
	_ "github.com/VrMolodyakov/vote-service/pkg/client/redis"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger(cfg.LogLvl)
	app := internal.NewApp(logger, cfg)
	app.Run()
}

/*

	pgConfig := postgresql.NewPgConfig(
	cfg.PostgreSql.Username,
	cfg.PostgreSql.Password,
	cfg.PostgreSql.Host,
	cfg.PostgreSql.Port,
	cfg.PostgreSql.Dbname)











	logger := logging.GetLogger("info")
	cfg := config.GetConfig()
	rdCfg := redis.NewRdConfig(cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DbNumber)
	client, err := redis.NewClient(context.Background(), &rdCfg)
	check(logger, err)
	cache := choiceCache.NewChoiceCache(client, logger)
	err = cache.Save("vote", "choice", 1, 1*time.Second)
	check(logger, err)
	time.Sleep(7 * time.Second)
	logger.Info("after time")
	count, err := cache.Get("vote", "choice")
	check(logger, err)
	logger.Info("count ", count)
	logger.Info("end")





	fmt.Println("start")
	logger := logging.GetLogger("info")
	cfg := config.GetConfig()
	rdCfg := redis.NewRdConfig(cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DbNumber)
	client, err := redis.NewClient(context.Background(), &rdCfg)
	check(logger, err)
	cache := choiceCache.NewChoiceCache(client, logger)
	err = cache.Set("vote", "choice", 1, 1*time.Second)
	check(logger, err)
	time.Sleep(7 * time.Second)
	logger.Info("after time")
	count, err := cache.Get("vote", "choice")
	check(logger, err)
	logger.Info("count ", count)
	logger.Info("end")
*/
