package service

import (
	"context"
	"errors"
	"testing"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/internal/domain/service/mocks"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateChoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	choiceRepo := mocks.NewMockChoiceRepository(ctrl)
	voteRepo := mocks.NewMockVoteRepository(ctrl)
	mockedRedis := mocks.NewMockRedisCache(ctrl)
	type mockCall func() *choiceService
	testCases := []struct {
		title   string
		input   entity.Choice
		mock    mockCall
		want    string
		isError bool
	}{
		{
			title: "success choice creation",
			input: entity.Choice{Title: "choice title", VoteId: 1, Count: 1},
			want:  "choice title",
			mock: func() *choiceService {
				logger := logging.GetLogger("debug")
				cacheService := NewCahceService(mockedRedis, logger)
				voteService := NewVoteService(voteRepo, logger)
				choiceRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return("choice title", nil)
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
		},
		{
			title: "title is empty and create should return error",
			input: entity.Choice{Title: "", VoteId: 1, Count: 1},
			want:  "",
			mock: func() *choiceService {
				logger := logging.GetLogger("debug")
				cacheService := NewCahceService(mockedRedis, logger)
				voteService := NewVoteService(voteRepo, logger)
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: true,
		},
		{
			title: "internal db error and create should return error",
			input: entity.Choice{Title: "choice title", VoteId: 1, Count: 1},
			want:  "",
			mock: func() *choiceService {
				logger := logging.GetLogger("debug")
				cacheService := NewCahceService(mockedRedis, logger)
				voteService := NewVoteService(voteRepo, logger)
				choiceRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return("", errors.New("internal db error"))
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			choiceService := test.mock()
			ctx := context.Background()
			got, err := choiceService.CreateChoice(ctx, test.input)
			if !test.isError {
				assert.NoError(t, err)
				assert.Equal(t, got, test.want)
			} else {
				assert.Equal(t, got, test.want)
				assert.Error(t, err)
			}
		})
	}
}

func TestGetVoteResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	choiceRepo := mocks.NewMockChoiceRepository(ctrl)
	voteService := mocks.NewMockVoteService(ctrl)
	mockedRedis := mocks.NewMockRedisCache(ctrl)
	type mockCall func() *choiceService
	testCases := []struct {
		title   string
		input   string
		mock    mockCall
		want    []entity.Choice
		isError bool
	}{
		{
			title: "success GetVoteResult() method vote result",
			input: "vote title",
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func() *choiceService {
				logger := logging.GetLogger("debug")
				cacheService := NewCahceService(mockedRedis, logger)
				voteService.EXPECT().GetByTitle(gomock.Any(), gomock.Any()).Return(1, nil)
				choices := []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}}
				choiceRepo.EXPECT().FindChoicesByVoteId(gomock.Any(), gomock.Any()).Return(choices, nil)
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
		},
		{
			title: "cannot find by title and GetVoteResult() method should return error",
			input: "vote title",
			want:  nil,
			mock: func() *choiceService {
				logger := logging.GetLogger("debug")
				cacheService := NewCahceService(mockedRedis, logger)
				voteService.EXPECT().GetByTitle(gomock.Any(), gomock.Any()).Return(-1, errors.New("title now found"))
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: true,
		},
		{
			title: "cannot find by choice by vote id and GetVoteResult() method should return error",
			input: "vote title",
			want:  nil,
			mock: func() *choiceService {
				logger := logging.GetLogger("debug")
				cacheService := NewCahceService(mockedRedis, logger)
				voteService.EXPECT().GetByTitle(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().FindChoicesByVoteId(gomock.Any(), gomock.Any()).Return(nil, errors.New("cannot find choices"))
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			choiceService := test.mock()
			ctx := context.Background()
			got, err := choiceService.GetVoteResult(ctx, test.input)
			if !test.isError {
				assert.NoError(t, err)
				assert.Equal(t, got, test.want)
			} else {
				assert.Equal(t, got, test.want)
				assert.Error(t, err)
			}
		})
	}
}

func TestUpdateChoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	choiceRepo := mocks.NewMockChoiceRepository(ctrl)
	voteService := mocks.NewMockVoteService(ctrl)
	cacheService := mocks.NewMockCacheService(ctrl)
	type args struct {
		voteTitle   string
		choiceTitle string
		count       int
	}
	type mockCall func() *choiceService
	testCases := []struct {
		title   string
		input   args
		mock    mockCall
		want    []entity.Choice
		isError bool
	}{
		{
			title: "success UpdateChoice() method vote result",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func() *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(-1, errors.New("empty cache"))
				cacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				voteService.EXPECT().GetByTitle(gomock.Any(), gomock.Any()).Return(1, nil)
				choice := entity.Choice{Title: "title1", VoteId: 1, Count: 1}
				choiceRepo.EXPECT().FindChoicesByVoteIdAndTitle(gomock.Any(), gomock.Any(), gomock.Any()).Return(choice, nil)
				choiceRepo.EXPECT().UpdateByTitleAndId(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
		},
		{
			title: "success UpdateChoice() method vote result",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func() *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				cacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				voteService.EXPECT().GetByTitle(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().UpdateByTitleAndId(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			choiceService := test.mock()
			ctx := context.Background()
			err := choiceService.UpdateChoice(ctx, test.input.voteTitle, test.input.choiceTitle, test.input.count)
			if !test.isError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
