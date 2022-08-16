package service

import (
	"testing"
	"time"

	"github.com/VrMolodyakov/vote-service/internal/service/mocks"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	ctrl := gomock.NewController(t)
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
				mockedRedis := mocks.NewMockRedisCache(ctrl)
				mockedRedis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				logger := logging.GetLogger("debug")
				return NewCahceService(mockedRedis, logger)
			},
			input:   args{"voteTitle", "choiceTitle", 1, time.Minute},
			isError: false,
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
			}

		})
	}
}
