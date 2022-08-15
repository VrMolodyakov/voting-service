package service

import (
	"context"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
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
	FindChoicesByVoteId(ctx context.Context, id int) ([]entity.Choice, error)
	Insert(ctx context.Context, choice entity.Choice) (string, error)
}

type choiceService struct {
	cache cacheService
	vote  VoteService
	repo  choiceRepository
}

func NewChoiceService(cache cacheService, vote VoteService, repo choiceRepository) *choiceService {
	return &choiceService{vote: vote, cache: cache, repo: repo}
}

func (c *choiceService) UpdateChoice(ctx context.Context, voteTitle string, choiceTitle string, count int) error {

	return nil
}

/*
lastCount, err := c.cache.Get(voteTitle, choiceTitle)
	if err != nil {
		id, err := c.vote.GetByTitle(ctx, voteTitle)
		if err != nil {
			return errors.ErrTitleNotExist
		}
		choices, err := c.repo.FindChoicesByVoteId(ctx, id)
		if err != nil {
			return err
		}
		for i := 0; i < len(choices); i++ {
			if choices[i].Title == choiceTitle {
				choices[i].Count += count

				return c.repo.UpdateByTitleAndId(ctx, choices[i].Count, id, choiceTitle)
			}
		}
	} else {
		duration := time.Minute * 5
		err := c.cache.Save(voteTitle, choiceTitle, lastCount+count, duration)
		if err != nil {
			return err
		}
	}

*/
