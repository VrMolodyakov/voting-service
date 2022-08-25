package service

import (
	"time"

	"github.com/VrMolodyakov/vote-service/internal/errs"
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
	if voteTitle == "" {
		return errs.ErrEmptyVoteTitle
	}
	if choiceTitle == "" {
		return errs.ErrEmptyChoiceTitle
	}
	return c.cache.Set(voteTitle, choiceTitle, count, expireAt)
}

func (c *cacheService) Get(voteTitle string, choiceTitle string) (int, error) {
	if voteTitle == "" {
		return -1, errs.ErrEmptyVoteTitle
	}
	if choiceTitle == "" {
		return -1, errs.ErrEmptyChoiceTitle
	}
	return c.cache.Get(voteTitle, choiceTitle)
}
