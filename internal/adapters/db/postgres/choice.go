package postgres

import (
	"context"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	psql "github.com/VrMolodyakov/vote-service/pkg/client/postgresql"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/jackc/pgx/v4/pgxpool"
)

type choiceRepository struct {
	client *pgxpool.Pool
	logger *logging.Logger
}

func NewChoiceStorage(pool *pgxpool.Pool, logger *logging.Logger) *voteRepository {
	return &voteRepository{client: pool, logger: logger}
}

func (c *choiceRepository) Insert(ctx context.Context, vote entity.Vote) error {
	sql := `INSERT INTO vote_choice(choice_title,count,vote_id) VALUES($1,$2,$3)`
	_, err := c.client.Exec(ctx, sql, vote.Title)
	if err != nil {
		err = psql.ErrExecuteQuery(err)
		c.logger.Error(err)
		return err
	}
	return nil
}

func (c *choiceRepository) FindIdByVoteId(ctx context.Context, id int) ([]entity.Choice, error) {
	sql := `SELECT * FROM vote_choice WHERE vote_id = $1`
	rows, err := c.client.Query(ctx, sql, id)
	if err != nil {
		err = psql.ErrExecuteQuery(err)
		c.logger.Error(err)
		return nil, err
	}
	defer rows.Close()
	choices := make([]entity.Choice, 0)
	for rows.Next() {
		var choice entity.Choice
		if err = rows.Scan(&choice.Title, &choice.Count, &choice.VoteId); err != nil {
			c.logger.Error(err)
			return nil, err
		}
		choices = append(choices, choice)
	}
	return choices, nil
}
