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
