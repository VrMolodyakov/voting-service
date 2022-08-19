package psqlStorage

import (
	"context"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	psql "github.com/VrMolodyakov/vote-service/pkg/client/postgresql"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/jackc/pgx/v4/pgxpool"
)

type choiceRepository struct {
	client PostgresClient
	logger *logging.Logger
}

func NewChoiceStorage(pool *pgxpool.Pool, logger *logging.Logger) *choiceRepository {
	return &choiceRepository{client: pool, logger: logger}
}

func (c *choiceRepository) Insert(ctx context.Context, choice entity.Choice) (string, error) {
	sql := `INSERT INTO choice(choice_title,count,vote_id) VALUES($1,$2,$3) RETURNING choice_title`
	var title string
	err := c.client.QueryRow(ctx, sql, choice.Title, choice.Count, choice.VoteId).Scan(&title)
	if err != nil {
		err = psql.ErrExecuteQuery(err)
		c.logger.Error(err)
		return "", err
	}
	return title, nil
}

func (c *choiceRepository) FindChoicesByVoteId(ctx context.Context, id int) ([]entity.Choice, error) {
	sql := `SELECT * FROM choice WHERE vote_id = $1`
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

func (c *choiceRepository) FindChoicesByVoteIdAndTitle(ctx context.Context, id int, choiceTitle string) (entity.Choice, error) {
	sql := `SELECT choice_title,vote_id,count
			FROM choice 
			WHERE vote_id = $1 AND choice_title = $2`
	var choice entity.Choice
	err := c.client.QueryRow(ctx, sql, id, choiceTitle).Scan(&choice.Title, &choice.VoteId, &choice.Count)
	if err != nil {
		err = psql.ErrExecuteQuery(err)
		c.logger.Error(err)
		return entity.Choice{}, err
	}
	c.logger.Debugf("Find choice.count = %v , choice.voteId = %v , choice.Title = %v ,", choice.Count, choice.VoteId, choice.VoteId)
	return choice, nil
}

func (c *choiceRepository) UpdateByTitleAndId(ctx context.Context, count int, voteId int, title string) error {
	sql := `UPDATE choice
			SET count = $1
			WHERE choice_title = $2 AND vote_id = $3`
	_, err := c.client.Exec(ctx, sql, count, title, voteId)
	if err != nil {
		err = psql.ErrExecuteQuery(err)
		c.logger.Error(err)
		return err
	}
	return nil
}
