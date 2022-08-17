package service

import (
	"context"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/internal/errors"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

const (
	expire time.Duration = 5
)

type CacheService interface {
	Save(voteTitle string, choiceTitle string, count int, expireAt time.Duration) error
	Get(voteTitle string, choiceTitle string) (int, error)
}

type VoteService interface {
	Create(ctx context.Context, vote entity.Vote) (int, error)
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
	lastCount, err := c.cache.Get(voteTitle, choiceTitle)
	if err != nil {
		id, err := c.vote.GetByTitle(ctx, voteTitle)
		if err != nil {
			return errors.ErrTitleNotExist
		}
		choice, err := c.repo.FindChoicesByVoteIdAndTitle(ctx, id, choiceTitle)
		if err != nil {
			return errors.ErrChoiceTitleNotExist
		}
		updateCount := choice.Count + count
		go func() {
			err := c.cache.Save(voteTitle, choice.Title, updateCount, time.Minute*expire)
			if err != nil {
				c.logger.Errorf("cache.Save() error due to %v", err)
			}
		}()
		return c.repo.UpdateByTitleAndId(ctx, updateCount, id, choiceTitle)

	} else {
		newCount := lastCount + count
		duration := time.Minute * expire
		err := c.cache.Save(voteTitle, choiceTitle, newCount, duration)
		go func() {
			id, err := c.vote.GetByTitle(ctx, voteTitle)
			if err != nil {
				c.logger.Errorf("vote.GetByTitle error due to %v", err)
				return
			}
			err = c.repo.UpdateByTitleAndId(ctx, newCount, id, choiceTitle)
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
		return nil, errors.ErrTitleNotExist
	}
	choices, err := c.repo.FindChoicesByVoteId(ctx, id)
	if err != nil {
		c.logger.Errorf("GetVoteResult() error due to %v", err)
		return nil, err
	}
	return choices, nil

}

/*
	request -> (
		vot title
		choice title

	)->update Choice (
		voteTitle string
		choiceTitle string
		count int
	)->if in cache (
		update cache
		update psql
		return nil (200 response)
	)->not in cache(
		check title is present
		choice is present
		update count
		save to cache
		update psql
		return nil (200 response)
	)











for i := 0; i < len(choices); i++ {
			if choices[i].Title == choiceTitle {
				choices[i].Count += count

				return c.repo.UpdateByTitleAndId(ctx, choices[i].Count, id, choiceTitle)
			}
		}


*/
