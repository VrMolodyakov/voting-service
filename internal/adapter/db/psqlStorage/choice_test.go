package psqlStorage

import (
	"context"
	"errors"
	"testing"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

type choiceMockRow struct {
	title string
	Err   error
}

type choiceEntityRow struct {
	title  string
	voteId int
	count  int
	Err    error
}

func (this choiceEntityRow) Scan(dest ...interface{}) error {
	if this.title == "" {
		return pgx.ErrNoRows
	}
	choice := dest[0].(*entity.Choice)

	choice.VoteId = this.voteId
	choice.Title = this.title
	choice.Count = this.count
	return nil
}

func (this choiceMockRow) Scan(dest ...interface{}) error {
	if this.title == "" {
		return pgx.ErrNoRows
	}
	title := dest[0].(*string)
	*title = this.title
	return nil
}

func TestInsertChoice(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	choiceRepo := choiceRepository{client: mockPool, logger: logger}

	type mockCall func()
	tests := []struct {
		title   string
		mock    mockCall
		input   entity.Choice
		want    string
		isError bool
	}{
		{
			title: "should insert successfully",
			input: entity.Choice{
				Title:  "test title",
				Count:  10,
				VoteId: 1,
			},
			mock: func() {
				row := choiceMockRow{"choice title", nil}
				mockPool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			want:    "choice title",
			isError: false,
		},
		{
			title: "should return error due to psql error",
			input: entity.Choice{
				Title:  "test title",
				Count:  10,
				VoteId: 1,
			},
			mock: func() {
				row := choiceMockRow{"", errors.New("psql error")}
				mockPool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			want:    "",
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			got, err := choiceRepo.Insert(context.Background(), test.input)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestFindIdByVoteId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	choiceRepo := choiceRepository{client: mockPool, logger: logger}

	type mockCall func()
	tests := []struct {
		title   string
		mock    mockCall
		input   int
		want    []entity.Choice
		isError bool
	}{
		{
			title: "should find successfully",
			input: 1,
			mock: func() {
				columns := []string{"title", "voteId", "count"}
				pgxRows := pgxpoolmock.NewRows(columns).
					AddRow("first", 10, 1).
					AddRow("second", 11, 1).
					AddRow("third", 12, 1).
					ToPgxRows()
				mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxRows, nil)
			},
			want:    []entity.Choice{{"first", 1, 10}, {"second", 1, 11}, {"third", 1, 12}},
			isError: false,
		},
		{
			title: "should return error inside psql Query request",
			input: 1,
			mock: func() {
				columns := []string{"title", "voteId", "count"}
				pgxRows := pgxpoolmock.NewRows(columns).
					AddRow("first", 10, 1).
					AddRow("second", 11, 1).
					AddRow("third", 12, 1).
					ToPgxRows()
				mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxRows, errors.New("not found"))
			},
			want:    []entity.Choice{},
			isError: true,
		},
		{
			title: "should return error while scanning row",
			input: 1,
			mock: func() {
				columns := []string{"title", "voteId", "count"}
				pgxRows := pgxpoolmock.NewRows(columns).
					AddRow("first", 10, 1).
					AddRow("second", 11, 1).
					AddRow("third", 12, 1).
					RowError(1, errors.New("internal psql erroe")).
					ToPgxRows()
				mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxRows, nil)
			},
			want:    []entity.Choice{},
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			got, err := choiceRepo.FindChoicesByVoteId(context.Background(), test.input)
			if !test.isError {
				for i, actual := range got {
					assert.Equal(t, test.want[i].Title, actual.Title)
					assert.Equal(t, test.want[i].Count, actual.Count)
					assert.Equal(t, test.want[i].VoteId, actual.VoteId)
				}
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestUpdateById(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	choiceRepo := choiceRepository{client: mockPool, logger: logger}

	type args struct {
		voteId int
		count  int
		title  string
	}

	type mockCall func()
	tests := []struct {
		title   string
		mock    mockCall
		input   args
		isError bool
	}{
		{
			title: "should update successfully",
			input: args{1, 1, "title"},
			mock: func() {
				mockPool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			isError: false,
		},
		{
			title: "update should return error",
			input: args{1, 1, "title"},
			mock: func() {
				mockPool.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("internal error"))
			},
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			err := choiceRepo.UpdateByTitleAndId(context.Background(), test.input.count, test.input.voteId, test.input.title)
			if !test.isError {
				assert.Equal(t, err, nil)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestFindChoicesByVoteIdAndTitle(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := logging.GetLogger("debug")
	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)
	choiceRepo := choiceRepository{client: mockPool, logger: logger}
	type mockCall func()
	type args struct {
		voteId int
		title  string
	}
	tests := []struct {
		title   string
		mock    mockCall
		input   args
		want    entity.Choice
		isError bool
	}{
		{
			title: "should find successfully",
			input: args{
				title:  "choice title",
				voteId: 1,
			},
			mock: func() {
				row := choiceEntityRow{"choice title", 1, 1, nil}
				mockPool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			want:    entity.Choice{"choice title", 1, 1},
			isError: false,
		},
		{
			title: "should return error due to psql error",
			input: args{
				title:  "choice title",
				voteId: 1,
			},
			mock: func() {
				row := choiceMockRow{"", errors.New("psql error")}
				mockPool.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
			},
			want:    entity.Choice{"choice title", 1, 1},
			isError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			got, err := choiceRepo.FindChoicesByVoteIdAndTitle(context.Background(), test.input.voteId, test.input.title)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}
