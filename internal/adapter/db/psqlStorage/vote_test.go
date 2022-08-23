package psqlStorage

import (
	"context"
	"errors"
	"testing"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

type voteMockRow struct {
	Id  int
	Err error
}

func (this voteMockRow) Scan(dest ...interface{}) error {
	if this.Id == 0 {
		return pgx.ErrNoRows
	}
	id := dest[0].(*int)
	*id = this.Id
	return nil
}

func TestInsertVote(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	voteRepo := voteRepository{client: mockPool, logger: logger}

	type mockCall func()
	tests := []struct {
		title   string
		mock    mockCall
		input   string
		want    int
		isError bool
	}{
		{
			title: "should insert successfully",
			input: "test title",
			mock: func() {
				row := voteMockRow{1, nil}
				mockPool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			want:    1,
			isError: false,
		},
		{
			title: "",
			input: "",
			mock: func() {
				row := voteMockRow{0, errors.New("psql error")}
				mockPool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			want:    1,
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			got, err := voteRepo.InsertVote(context.Background(), test.input)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestFindByTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	voteRepo := voteRepository{client: mockPool, logger: logger}

	type args struct {
		v entity.Vote
	}

	type mockCall func()
	tests := []struct {
		title   string
		mock    mockCall
		input   string
		want    int
		isError bool
	}{
		{
			title: "should insert successfully",
			input: "title to find",
			mock: func() {
				row := voteMockRow{1, nil}
				mockPool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			want:    1,
			isError: false,
		},
		{
			title: "psql not found and should return error ",
			input: "",
			mock: func() {
				row := voteMockRow{0, errors.New("psql error")}
				mockPool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			want:    1,
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			got, err := voteRepo.FindVote(context.Background(), test.input)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	voteRepo := voteRepository{client: mockPool, logger: logger}

	type mockCall func()
	tests := []struct {
		title   string
		mock    mockCall
		input   string
		isError bool
	}{
		{
			title: "should update successfully",
			input: "title to delete",
			mock: func() {
				mockPool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			isError: false,
		},
		{
			title: "update should return error",
			input: "title to delete",
			mock: func() {
				mockPool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
			},
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			err := voteRepo.DeleteVote(context.Background(), test.input)
			if !test.isError {
				assert.Equal(t, err, nil)
			} else {
				assert.Error(t, err)
			}

		})
	}
}
