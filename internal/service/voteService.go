package service

import (
	"context"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/internal/errors"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

type VoteRepository interface {
	Insert(ctx context.Context, vote entity.Vote) (int, error)
	FindIdByTitle(ctx context.Context, title string) (int, error)
}

type voteService struct {
	repo   VoteRepository
	logger *logging.Logger
}

func NewVoteService(repo VoteRepository, logger *logging.Logger) *voteService {
	return &voteService{repo: repo, logger: logger}
}

func (v *voteService) Create(ctx context.Context, vote entity.Vote) (int, error) {
	if vote.Title == "" {
		return -1, errors.ErrEmptyTitle
	}
	return v.repo.Insert(ctx, vote)
}

func (v *voteService) GetByTitle(ctx context.Context, title string) (int, error) {
	if title == "" {
		return -1, errors.ErrEmptyTitle
	}
	return v.repo.FindIdByTitle(ctx, title)
}
