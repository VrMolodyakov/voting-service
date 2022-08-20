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

type choiceRepository interface {
	UpdateByTitleAndId(ctx context.Context, count int, voteId int, title string) error
	FindChoicesByVoteIdAndTitle(ctx context.Context, id int, choiceTitle string) (entity.Choice, error)
	FindChoicesByVoteId(ctx context.Context, id int) ([]entity.Choice, error)
	Insert(ctx context.Context, choice entity.Choice) (string, error)
}

type choiceService struct {
	cache  CacheService
	vote   VoteService
	repo   choiceRepository
	logger *logging.Logger
}

func NewChoiceService(cache CacheService, vote VoteService, repo choiceRepository, logger *logging.Logger) *choiceService {
	return &choiceService{vote: vote, cache: cache, repo: repo, logger: logger}
}

func (c *choiceService) UpdateChoice(ctx context.Context, voteTitle string, choiceTitle string, count int) error {
	c.logger.Debugf("try to update choice with vote title = %v, choice title = %v,count = %v", voteTitle, choiceTitle, count)
	lastCount, err := c.cache.Get(voteTitle, choiceTitle)
	if err != nil {
		c.logger.Debugf("vote title = %v and choice title = %v not found in cache", voteTitle, choiceTitle)
		id, err := c.vote.GetByTitle(ctx, voteTitle)
		if err != nil {
			return errs.ErrTitleNotExist
		}
		choice, err := c.repo.FindChoicesByVoteIdAndTitle(ctx, id, choiceTitle)
		if err != nil {
			return errs.ErrChoiceTitleNotExist
		}
		updateCount := choice.Count + count
		c.logger.Debugf("current choice = %v", choice)
		go func() {
			err := c.cache.Save(voteTitle, choice.Title, updateCount, expire)
			if err != nil {
				c.logger.Errorf("cache.Save() error due to %v", err)
			}
		}()
		return c.repo.UpdateByTitleAndId(ctx, updateCount, id, choiceTitle)

	} else {
		newCount := lastCount + count
		err := c.cache.Save(voteTitle, choiceTitle, newCount, expire)
		go func() {
			updCtx, cancel := context.WithTimeout(context.Background(), updateTimeout)
			defer cancel()
			id, err := c.vote.GetByTitle(updCtx, voteTitle)
			if err != nil {
				c.logger.Errorf("vote.GetByTitle error due to %v", err)
				return
			}
			err = c.repo.UpdateByTitleAndId(updCtx, newCount, id, choiceTitle)
			if err != nil {
				c.logger.Errorf("repo.UpdateByTitleAndId error due to %v", err)
			}
		}()
		if err != nil {
			c.logger.Errorf("cache.Save() error due to %v", err)
			return err
		}
		return nil
	}
}

func (c *choiceService) GetVoteResult(ctx context.Context, voteTitle string) ([]entity.Choice, error) {
	c.logger.Debugf("try to find with choices title %v", voteTitle)
	id, err := c.vote.GetByTitle(ctx, voteTitle)
	if err != nil {
		c.logger.Errorf("GetVoteResult() error due to %v", err)
		return nil, errs.ErrTitleNotExist
	}
	choices, err := c.repo.FindChoicesByVoteId(ctx, id)
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
