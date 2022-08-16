package service

import (
	"errors"
	"testing"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/service/mocks"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	//var mockedRedis *mocks.MockRedisCache

	ctrl := gomock.NewController(t)
	mockedRedis := mocks.NewMockRedisCache(ctrl)
	defer ctrl.Finish()
	type mock func() *cacheService
	type args struct {
		voteTitle   string
		choiceTitle string
		count       int
		expireAt    time.Duration
	}
	testCases := []struct {
		title    string
		mockCall mock
		input    args
		isError  bool
	}{
		{
			title: "Success save and return nil",
			mockCall: func() *cacheService {
				mockedRedis.EXPECT().Set(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(nil)
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"voteTitle", "choiceTitle", 1, time.Minute},
			isError: false,
		},
		{
			title: "Unsuccessful save and should return error",
			mockCall: func() *cacheService {
				mockedRedis.EXPECT().Set(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(errors.New("cache internal error"))
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"voteTitle", "choiceTitle", 1, time.Minute},
			isError: true,
		},
		{
			title: "wrong vote titlle and save should return error",
			mockCall: func() *cacheService {
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"", "choiceTitle", 1, time.Minute},
			isError: true,
		},
		{
			title: "wrong choice titlle and save should return error",
			mockCall: func() *cacheService {
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"vote title", "", 1, time.Minute},
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			cacheService := test.mockCall()
			err := cacheService.Save(test.input.voteTitle,
				test.input.choiceTitle,
				test.input.count,
				test.input.expireAt)
			if !test.isError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockedRedis := mocks.NewMockRedisCache(ctrl)
	defer ctrl.Finish()
	type mock func() *cacheService
	type args struct {
		voteTitle   string
		choiceTitle string
	}
	testCases := []struct {
		title    string
		mockCall mock
		input    args
		want     int
		isError  bool
	}{
		{
			title: "Success save and return nil",
			mockCall: func() *cacheService {
				mockedRedis.EXPECT().Get(
					gomock.Any(),
					gomock.Any()).Return(1, nil)
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"voteTitle", "choiceTitle"},
			want:    1,
			isError: false,
		},
		{
			title: "Unsuccessful save and should return error",
			mockCall: func() *cacheService {
				mockedRedis.EXPECT().Get(
					gomock.Any(),
					gomock.Any()).Return(-1, errors.New("cache internal error"))
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"voteTitle", "choiceTitle"},
			want:    -1,
			isError: true,
		},
		{
			title: "wrong vote titlle and Get should return error",
			mockCall: func() *cacheService {
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"", "choiceTitle"},
			want:    -1,
			isError: true,
		},
		{
			title: "wrong choice titlle and Get should return error",
			mockCall: func() *cacheService {
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"vote title", ""},
			want:    -1,
			isError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			cacheService := test.mockCall()
			got, err := cacheService.Get(test.input.voteTitle,
				test.input.choiceTitle)
			if !test.isError {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			} else {
				assert.Equal(t, test.want, got)
				assert.Error(t, err)
			}

		})
	}
}
