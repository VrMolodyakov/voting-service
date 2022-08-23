package service

import (
	"context"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/internal/errs"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

const (
	expire        time.Duration = 5 * time.Minute
	updateTimeout               = 1 * time.Minute
)

type CacheService interface {
	Save(voteTitle string, choiceTitle string, count int, expireAt time.Duration) error
	Get(voteTitle string, choiceTitle string) (int, error)
}

type VoteService interface {
	Create(ctx context.Context, title string) (int, error)
	GetByTitle(ctx context.Context, title string) (int, error)
}

type СhoiceRepository interface {
	Insert(ctx context.Context, choice entity.Choice) (string, error)
	FindChoices(ctx context.Context, id int) ([]entity.Choice, error)
	FindChoice(ctx context.Context, id int, choiceTitle string) (entity.Choice, error)
	IncrementUpdate(ctx context.Context, count int, voteId int, title string) (int, error)
	SetUpdate(ctx context.Context, count int, voteId int, title string) (int, error)
}

type choiceService struct {
	cache  CacheService
	vote   VoteService
	repo   СhoiceRepository
	logger *logging.Logger
}

func NewChoiceService(cache CacheService, vote VoteService, repo СhoiceRepository, logger *logging.Logger) *choiceService {
	return &choiceService{vote: vote, cache: cache, repo: repo, logger: logger}
}

func (c *choiceService) UpdateChoice(ctx context.Context, voteTitle string, choiceTitle string, count int) error {
	c.logger.Debugf("try to update choice with vote title = %v, choice title = %v,count = %v", voteTitle, choiceTitle, count)
	lastCount, err := c.cache.Get(voteTitle, choiceTitle)
	if err != nil {
		updCount, err := c.incrUpdate(ctx, voteTitle, choiceTitle, count)
		if err != nil {
			c.logger.Errorf("cannot update for vote title = %v , choice title = %v", voteTitle, choiceTitle)
			return err
		}
		go func() {
			err := c.cache.Save(voteTitle, choiceTitle, updCount, expire)
			if err != nil {
				c.logger.Errorf("cache.Save() error due to %v", err)
			}
		}()
		return nil

	} else {
		newCount := lastCount + count
		c.logger.Debugf("last choice count = %v", lastCount)
		err := c.cache.Save(voteTitle, choiceTitle, newCount, expire)
		if err != nil {
			_, err := c.incrUpdate(ctx, voteTitle, choiceTitle, count)
			if err != nil {
				c.logger.Errorf("couldn't update count due to %v", err)
				return err
			}
			return nil
		}
		go func() {
			updCtx, cancel := context.WithTimeout(context.Background(), updateTimeout)
			defer cancel()
			_, err := c.incrUpdate(updCtx, voteTitle, choiceTitle, count)
			c.logger.Errorf("couldn't update count due to %v", err)
		}()
		return nil
	}
}

func (c *choiceService) incrUpdate(ctx context.Context, voteTitle string, choiceTitle string, count int) (int, error) {
	id, err := c.vote.GetByTitle(ctx, voteTitle)
	if err != nil {
		return -1, errs.ErrTitleNotExist
	}
	return c.repo.IncrementUpdate(ctx, count, id, choiceTitle)
}

func (c *choiceService) GetVoteResult(ctx context.Context, voteTitle string) ([]entity.Choice, error) {
	c.logger.Debugf("try to find with choices title %v", voteTitle)
	id, err := c.vote.GetByTitle(ctx, voteTitle)
	if err != nil {
		c.logger.Errorf("GetVoteResult() error due to %v", err)
		return nil, errs.ErrTitleNotExist
	}
	choices, err := c.repo.FindChoices(ctx, id)
	if err != nil {
		c.logger.Errorf("GetVoteResult() error due to %v", err)
		return nil, err
	}
	return choices, nil

}

func (c *choiceService) CreateChoice(ctx context.Context, choice entity.Choice) (string, error) {
	if choice.Title == "" {
		return "", errs.ErrEmptyChoiceTitle
	}
	return c.repo.Insert(ctx, choice)
}
