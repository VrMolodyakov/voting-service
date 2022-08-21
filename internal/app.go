package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/adapter/db/choiceCache"
	"github.com/VrMolodyakov/vote-service/internal/adapter/db/psqlStorage"
	"github.com/VrMolodyakov/vote-service/internal/config"
	"github.com/VrMolodyakov/vote-service/internal/domain/service"
	"github.com/VrMolodyakov/vote-service/internal/handler"
	"github.com/VrMolodyakov/vote-service/pkg/client/postgresql"
	"github.com/VrMolodyakov/vote-service/pkg/client/redis"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/VrMolodyakov/vote-service/pkg/shutdown"
	"github.com/gorilla/mux"
)

const (
	attemp       int = 5
	delay            = 5 * time.Second
	writeTimeout     = 15 * time.Second
	readTimeout      = 15 * time.Second
)

type app struct {
	logger *logging.Logger
	cfg    *config.Config
	router *mux.Router
}

func NewApp(logger *logging.Logger, cfg *config.Config) *app {
	return &app{logger: logger, cfg: cfg}
}

func (a *app) Run() {
	a.startHttp()
}

func (a *app) startHttp() {
	a.logger.Debug("start init handler")

	pgCfg := postgresql.NewPgConfig(
		a.cfg.PostgreSql.Username,
		a.cfg.PostgreSql.Password,
		a.cfg.PostgreSql.Host,
		a.cfg.PostgreSql.Port,
		a.cfg.PostgreSql.Dbname,
		a.cfg.PostgreSql.PoolSize)

	psqlClient, err := postgresql.NewClient(context.Background(), attemp, delay, pgCfg)
	a.checkErr(err)

	rdCfg := redis.NewRdConfig(a.cfg.Redis.Password, a.cfg.Redis.Host, a.cfg.Redis.Port, a.cfg.Redis.DbNumber)
	rdClient, err := redis.NewClient(context.Background(), &rdCfg)
	a.checkErr(err)

	voteRepo := psqlStorage.NewVoteStorage(psqlClient, a.logger)
	choiceRepo := psqlStorage.NewChoiceStorage(psqlClient, a.logger)
	voteService := service.NewVoteService(voteRepo, a.logger)
	choiceCache := choiceCache.NewChoiceCache(rdClient, a.logger)
	cacheService := service.NewCahceService(choiceCache, a.logger)
	choiceService := service.NewChoiceService(cacheService, voteService, choiceRepo, a.logger)

	a.router = mux.NewRouter()
	a.initializeRouters(choiceService, voteService)
	a.logger.Info("start listening...")
	port := fmt.Sprintf(":%s", a.cfg.Port)
	server := &http.Server{
		Addr:         port,
		Handler:      a.router,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
	a.checkErr(err)
	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, rdClient, server)
	defer psqlClient.Close()
	if err := server.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
}

func (a *app) initializeRouters(choiceService handler.ChoiceService, voteService handler.VoteService) {
	h := handler.NewVoteHandler(a.logger, voteService, choiceService)
	a.router.HandleFunc("/api/vote", h.Create).Methods("POST")
	a.router.HandleFunc("/api/result", h.GetChoices).Methods("POST")
	a.router.HandleFunc("/api/choice", h.UpdateChoice).Methods("POST")
	a.router.HandleFunc("/api/vote/{id:[0-9]+}", h.DeleteVote).Methods("DELETE")
}

func (a *app) checkErr(err error) {
	if err != nil {
		a.logger.Fatal(err)
	}
}
