package handler

import (
	"context"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
)

type VoteService interface {
	Create(ctx context.Context, vote string) (int, error)
	Get(ctx context.Context, title string) (int, error)
	Delete(ctx context.Context, id string) error
}

type ChoiceService interface {
	Create(ctx context.Context, choice entity.Choice) (string, error)
	Get(ctx context.Context, voteTitle string) ([]entity.Choice, error)
	Update(ctx context.Context, voteTitle string, choiceTitle string, count int) error
}
