package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/internal/domain/service/mocks"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateChoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	choiceRepo := mocks.NewMockСhoiceRepository(ctrl)
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
			got, err := choiceService.Create(ctx, test.input)
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

func TestGetVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	choiceRepo := mocks.NewMockСhoiceRepository(ctrl)
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
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choices := []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}}
				choiceRepo.EXPECT().FindChoices(gomock.Any(), gomock.Any()).Return(choices, nil)
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
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(-1, errors.New("title now found"))
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
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().FindChoices(gomock.Any(), gomock.Any()).Return(nil, errors.New("cannot find choices"))
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			choiceService := test.mock()
			ctx := context.Background()
			got, err := choiceService.Get(ctx, test.input)
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
	choiceRepo := mocks.NewMockСhoiceRepository(ctrl)
	voteService := mocks.NewMockVoteService(ctrl)
	cacheService := mocks.NewMockCacheService(ctrl)
	type args struct {
		voteTitle   string
		choiceTitle string
		count       int
	}
	type mockCall func(wg *sync.WaitGroup) *choiceService
	testCases := []struct {
		title   string
		input   args
		mock    mockCall
		want    []entity.Choice
		isError bool
		wait    bool
	}{
		{
			title: "success UpdateChoice() method vote result",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			mock: func(wg *sync.WaitGroup) *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(-1, errors.New("empty cache"))
				cacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(
					func(voteTitle interface{}, choiceTitle interface{}, count interface{}, expireAt interface{}) {
						wg.Done()
					})
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
			wait:    true,
		},
		{
			title: "success UpdateChoice() method vote result and error in goroutine",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			mock: func(wg *sync.WaitGroup) *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(-1, errors.New("empty cache"))
				cacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(
					func(voteTitle interface{}, choiceTitle interface{}, count interface{}, expireAt interface{}) {
						wg.Done()
					}).Return(errors.New("cannot save in cache"))
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
			wait:    true,
		},
		{
			title: "UpdateChoice() find choice in cache and succes method execute",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func(wg *sync.WaitGroup) *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				cacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(
					func(ctx interface{}, count interface{}, voteId interface{}, title interface{}) {
						wg.Done()
					})
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
			wait:    true,
		},
		{
			title: "UpdateChoice() find choice in cache and succes method execute but error in goroutine",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func(wg *sync.WaitGroup) *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				cacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(
					func(ctx interface{}, count interface{}, voteId interface{}, title interface{}) {
						wg.Done()
					}).Return(1, errors.New("cannot update "))
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
			wait:    true,
		},
		{
			title: " vote title not found and service.Update() should return error",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func(wg *sync.WaitGroup) *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(-1, errors.New("empty cache"))
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(-1, errors.New("title not found"))
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: true,
			wait:    false,
		},
		{
			title: " cannot execute repo.Update() and service should return error",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func(wg *sync.WaitGroup) *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(-1, errors.New("empty cache"))
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New("internal db error"))
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: true,
			wait:    false,
		},
		{
			title: "found in cache but couldn't save and then succes db update",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func(wg *sync.WaitGroup) *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				cacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("reddis internal error"))
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: false,
			wait:    false,
		},
		{
			title: "found in cache but couldn't save and then couldn't execute db.update and service.Update() should return error",
			input: args{voteTitle: "vote title", choiceTitle: "choice title", count: 1},
			want:  []entity.Choice{{Title: "title1", VoteId: 1, Count: 1}, {Title: "title2", VoteId: 1, Count: 1}},
			mock: func(wg *sync.WaitGroup) *choiceService {
				logger := logging.GetLogger("debug")
				cacheService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				cacheService.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("reddis internal error"))
				voteService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New("internal service error"))
				return NewChoiceService(cacheService, voteService, choiceRepo, logger)
			},
			isError: true,
			wait:    false,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			var wg sync.WaitGroup
			if test.wait {
				wg.Add(1)
			}
			choiceService := test.mock(&wg)
			ctx := context.Background()
			err := choiceService.Update(ctx, test.input.voteTitle, test.input.choiceTitle, test.input.count)
			wg.Wait()
			if !test.isError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
