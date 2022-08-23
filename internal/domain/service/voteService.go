package service

import (
	"context"
	"errors"

	"github.com/VrMolodyakov/vote-service/internal/errs"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

type VoteRepository interface {
	DeleteVote(ctx context.Context, id string) error
	FindVote(ctx context.Context, title string) (int, error)
	InsertVote(ctx context.Context, vote string) (int, error)
}

type voteService struct {
	repo   VoteRepository
	logger *logging.Logger
}

func NewVoteService(repo VoteRepository, logger *logging.Logger) *voteService {
	return &voteService{repo: repo, logger: logger}
}

func (v *voteService) Create(ctx context.Context, title string) (int, error) {
	v.logger.Debugf("try to create vote with title %v", title)
	if title == "" {
		return -1, errs.ErrEmptyVoteTitle
	}

	vote, err := v.repo.InsertVote(ctx, title)
	if err != nil {
		if errors.Is(err, errs.ErrTitleAlreadyExist) {
			v.logger.Errorf("create error due to %v", err)
			return -1, err
		}
	}
	return vote, nil
}

func (v *voteService) GetByTitle(ctx context.Context, title string) (int, error) {
	v.logger.Debugf("try to get vote with title %v", title)
	if title == "" {
		return -1, errs.ErrEmptyVoteTitle
	}
	return v.repo.FindVote(ctx, title)
}

func (v *voteService) DeleteVoteById(ctx context.Context, id string) error {
	v.logger.Debugf("try to get vote with title %v", id)
	if id == "" {
		return errs.ErrEmptyVoteTitle
	}
	return v.repo.DeleteVote(ctx, id)
}
