package handler

import (
	"context"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
)

type VoteService interface {
	Create(ctx context.Context, vote string) (int, error)
	GetByTitle(ctx context.Context, title string) (int, error)
}

type ChoiceService interface {
	CreateChoice(ctx context.Context, choice entity.Choice) (string, error)
	GetVoteResult(ctx context.Context, voteTitle string) ([]entity.Choice, error)
	UpdateChoice(ctx context.Context, voteTitle string, choiceTitle string, count int) error
}
