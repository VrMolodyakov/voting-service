package psqlStorage

import (
	"context"
	"errors"

	"github.com/VrMolodyakov/vote-service/internal/errs"
	psql "github.com/VrMolodyakov/vote-service/pkg/client/postgresql"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type voteRepository struct {
	client PostgresClient
	logger *logging.Logger
}

func NewVoteStorage(pool *pgxpool.Pool, logger *logging.Logger) *voteRepository {
	return &voteRepository{client: pool, logger: logger}
}

func (v *voteRepository) Insert(ctx context.Context, vote string) (int, error) {
	sql := `INSERT INTO vote(vote_title)
			SELECT $1
			WHERE NOT EXISTS (SELECT vote_id FROM vote WHERE vote_title=$2) RETURNING vote_id`
	var id int
	err := v.client.QueryRow(ctx, sql, vote, vote).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, errs.ErrTitleAlreadyExist
		}
		return -1, err
	}
	return id, nil
}

func (v *voteRepository) Find(ctx context.Context, title string) (int, error) {
	sql := `SELECT vote_id FROM vote WHERE vote_title = $1`
	var id int
	err := v.client.QueryRow(ctx, sql, title).Scan(&id)
	if err != nil {
		err = psql.ErrExecuteQuery(err)
		v.logger.Error(err)
		return -1, err
	}
	return id, nil
}

func (v *voteRepository) Delete(ctx context.Context, id string) error {
	sql := `DELETE FROM VOTE
			WHERE vote_id = $1`
	_, err := v.client.Exec(ctx, sql, id)
	if err != nil {
		err = psql.ErrExecuteQuery(err)
		v.logger.Error(err)
		return err
	}
	return nil

}
