package service

import (
	"context"

	"github.com/VrMolodyakov/vote-service/internal/errs"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

type VoteRepository interface {
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, title string) (int, error)
	Insert(ctx context.Context, vote string) (int, error)
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

	vote, err := v.repo.Insert(ctx, title)
	if err != nil {
		v.logger.Errorf("couldn't create for title = %v ", title)
		return -1, err
	}
	return vote, nil
}

func (v *voteService) Get(ctx context.Context, title string) (int, error) {
	v.logger.Debugf("try to get vote with title %v", title)
	if title == "" {
		return -1, errs.ErrEmptyVoteTitle
	}
	return v.repo.Find(ctx, title)
}

func (v *voteService) Delete(ctx context.Context, id string) error {
	v.logger.Debugf("try to get vote with title %v", id)
	if id == "" {
		return errs.ErrEmptyVoteTitle
	}
	return v.repo.Delete(ctx, id)
}
