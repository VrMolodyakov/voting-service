package service

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/VrMolodyakov/vote-service/internal/domain/service/mocks"
// 	"github.com/VrMolodyakov/vote-service/pkg/logging"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestCreate(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	mockRepo := mocks.NewMockVoteRepository(ctrl)
// 	defer ctrl.Finish()
// 	type mock func() *voteService
// 	testCases := []struct {
// 		title    string
// 		mockCall mock
// 		input    string
// 		want     int
// 		isError  bool
// 	}{
// 		{
// 			title: "Success Create and return nil",
// 			mockCall: func() *voteService {
// 				mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(1, nil)
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "voteTitle",
// 			want:    1,
// 			isError: false,
// 		},
// 		{
// 			title: "Error in repo Insert and return nil",
// 			mockCall: func() *voteService {
// 				mockRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(-1, errors.New("repo internal error"))
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "voteTitle",
// 			want:    -1,
// 			isError: true,
// 		},
// 		{
// 			title: "wrong vote titlle and Get should return error",
// 			mockCall: func() *voteService {
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "",
// 			want:    -1,
// 			isError: true,
// 		},
// 	}
// 	for _, test := range testCases {
// 		t.Run(test.title, func(t *testing.T) {
// 			voteService := test.mockCall()
// 			got, err := voteService.Create(context.Background(), test.input)
// 			if !test.isError {
// 				assert.NoError(t, err)
// 				assert.Equal(t, test.want, got)
// 			} else {
// 				assert.Equal(t, test.want, got)
// 				assert.Error(t, err)
// 			}
// 		})
// 	}
// }

// func TestGetByTitle(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	mockRepo := mocks.NewMockVoteRepository(ctrl)
// 	defer ctrl.Finish()
// 	type mock func() *voteService
// 	testCases := []struct {
// 		title    string
// 		mockCall mock
// 		input    string
// 		want     int
// 		isError  bool
// 	}{
// 		{
// 			title: "Success Create and return nil",
// 			mockCall: func() *voteService {
// 				mockRepo.EXPECT().FindIdByTitle(gomock.Any(), gomock.Any()).Return(1, nil)
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "vote title",
// 			want:    1,
// 			isError: false,
// 		},
// 		{
// 			title: "Error in repo Insert and return nil",
// 			mockCall: func() *voteService {
// 				mockRepo.EXPECT().FindIdByTitle(gomock.Any(), gomock.Any()).Return(-1, errors.New("repo internal error"))
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "vote title",
// 			want:    -1,
// 			isError: true,
// 		},
// 		{
// 			title: "wrong vote titlle and Get should return error",
// 			mockCall: func() *voteService {
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "",
// 			want:    -1,
// 			isError: true,
// 		},
// 	}
// 	for _, test := range testCases {
// 		t.Run(test.title, func(t *testing.T) {
// 			voteService := test.mockCall()
// 			got, err := voteService.GetByTitle(context.Background(), test.input)
// 			if !test.isError {
// 				assert.NoError(t, err)
// 				assert.Equal(t, test.want, got)
// 			} else {
// 				assert.Equal(t, test.want, got)
// 				assert.Error(t, err)
// 			}
// 		})
// 	}
// }

// func TestDeleteByTitle(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	mockRepo := mocks.NewMockVoteRepository(ctrl)
// 	defer ctrl.Finish()
// 	type mock func() *voteService
// 	testCases := []struct {
// 		title    string
// 		mockCall mock
// 		input    string
// 		want     int
// 		isError  bool
// 	}{
// 		{
// 			title: "Success Create and return nil",
// 			mockCall: func() *voteService {
// 				mockRepo.EXPECT().DeleteVote(gomock.Any(), gomock.Any()).Return(nil)
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "some id",
// 			want:    1,
// 			isError: false,
// 		},
// 		{
// 			title: "Error in repo Insert and return nil",
// 			mockCall: func() *voteService {
// 				mockRepo.EXPECT().DeleteVote(gomock.Any(), gomock.Any()).Return(errors.New("repo internal error"))
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "some id",
// 			want:    -1,
// 			isError: true,
// 		},
// 		{
// 			title: "wrong vote titlle and Get should return error",
// 			mockCall: func() *voteService {
// 				logger := logging.GetLogger("debug")
// 				return NewVoteService(mockRepo, logger)
// 			},
// 			input:   "",
// 			want:    -1,
// 			isError: true,
// 		},
// 	}
// 	for _, test := range testCases {
// 		t.Run(test.title, func(t *testing.T) {
// 			voteService := test.mockCall()
// 			err := voteService.DeleteVoteById(context.Background(), test.input)
// 			if !test.isError {
// 				assert.NoError(t, err)
// 			} else {
// 				assert.Error(t, err)
// 			}
// 		})
// 	}
// }
