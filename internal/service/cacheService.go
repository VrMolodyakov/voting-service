package service

import (
	"time"

	"github.com/VrMolodyakov/vote-service/internal/errors"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

type RedisCache interface {
	Set(voteTitle string, choiceTitle string, count int, expireAt time.Duration) error
	Get(voteTitle string, choiceTitle string) (int, error)
}

type cacheService struct {
	cache  RedisCache
	logger *logging.Logger
}

func NewCahceService(cache RedisCache, logger *logging.Logger) *cacheService {
	return &cacheService{logger: logger, cache: cache}
}

func (c *cacheService) Save(voteTitle string, choiceTitle string, count int, expireAt time.Duration) error {
	if voteTitle == "" || choiceTitle == "" {
		return errors.ErrEmptyTitle
	}
	return c.cache.Set(voteTitle, choiceTitle, count, expireAt)
}

func (c *cacheService) Get(voteTitle string, choiceTitle string) (int, error) {
	if voteTitle == "" || choiceTitle == "" {
		return -1, errors.ErrEmptyTitle
	}
	return c.cache.Get(voteTitle, choiceTitle)
}
