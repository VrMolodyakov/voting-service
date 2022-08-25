package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/internal/errs"
	"github.com/VrMolodyakov/vote-service/internal/handler/mocks"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	router := mux.NewRouter()
	ctrl := gomock.NewController(t)
	choiceServ := mocks.NewMockChoiceService(ctrl)
	voteServ := mocks.NewMockVoteService(ctrl)
	handler := NewVoteHandler(logging.GetLogger("debug"), voteServ, choiceServ)
	handler.InitRoutes(router)
	type mockCall func([]entity.Choice)
	type args struct {
		voteTitle string
		choices   []entity.Choice
	}
	testCases := []struct {
		title          string
		inputRequest   string
		inputBody      args
		want           string
		mock           mockCall
		expectedStatus int
	}{
		{
			title:        "create vote pool and 200 response",
			inputRequest: `{"vote":"Best pokemon","choices":["Pikachu","Mew","Noone"]}`,
			inputBody:    args{voteTitle: "Best pokemon", choices: []entity.Choice{{Title: "Pikachu", VoteId: 1}, {Title: "Noone", VoteId: 1}, {Title: "Mew", VoteId: 1}}},
			mock: func(choices []entity.Choice) {
				voteServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(1, nil)
				for _, c := range choices {
					choiceServ.EXPECT().Create(gomock.Any(), c).Return(c.Title, nil)
				}

			},
			want:           "{\"vote\": \"Best pokemon\",\"choices\": [{\"choice\": \"Pikachu\",\"vote_count\": 0},{\"choice\": \"Mew\",\"vote_count\": 0},{\"choice\": \"Noone\",\"vote_count\": 0}]}",
			expectedStatus: 201,
		},
		{
			title:        "empty request title and  400 code response",
			inputRequest: `{"vote":"","choices":["Pikachu","Mew","Noone"]}`,
			inputBody:    args{voteTitle: "", choices: []entity.Choice{{Title: "Pikachu", VoteId: 1}, {Title: "Noone", VoteId: 1}, {Title: "Mew", VoteId: 1}}},
			mock: func(choices []entity.Choice) {
				voteServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(-1, errs.ErrEmptyVoteTitle)

			},
			want:           "vote title is empty",
			expectedStatus: 400,
		},
		{
			title:        " service error and 500 code response",
			inputRequest: `{"vote":"Best pokemon","choices":["Pikachu","Mew","Noone"]}`,
			inputBody:    args{voteTitle: "Best pokemon", choices: []entity.Choice{{Title: "Pikachu", VoteId: 1}, {Title: "Noone", VoteId: 1}, {Title: "Mew", VoteId: 1}}},
			mock: func(choices []entity.Choice) {
				voteServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(-1, errors.New("internal service error"))

			},
			want:           "500 Internal Server Error",
			expectedStatus: 500,
		},
		{
			title:        "empty request choice title and  400 code response",
			inputRequest: `{"vote":"Best pokemon","choices":["","Mew","Noone"]}`,
			inputBody:    args{voteTitle: "Best pokemon", choices: []entity.Choice{{Title: "", VoteId: 1}, {Title: "Noone", VoteId: 1}, {Title: "Mew", VoteId: 1}}},
			mock: func(choices []entity.Choice) {
				voteServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return(1, nil)
				choiceServ.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", errs.ErrEmptyChoiceTitle)
			},
			want:           "choice title is empty",
			expectedStatus: 400,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock(test.inputBody.choices)
			req := httptest.NewRequest(
				"POST",
				"/api/vote",
				bytes.NewBufferString(test.inputRequest),
			)
			req.Header.Set("Content-type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, clearResponse(recorder.Body.String()), test.want)
			assert.Equal(t, test.expectedStatus, recorder.Code)

		})
	}
}

func TestGetChoiceHandler(t *testing.T) {
	router := mux.NewRouter()
	ctrl := gomock.NewController(t)
	choiceServ := mocks.NewMockChoiceService(ctrl)
	voteServ := mocks.NewMockVoteService(ctrl)
	handler := NewVoteHandler(logging.GetLogger("debug"), voteServ, choiceServ)
	handler.InitRoutes(router)
	type mockCall func()
	testCases := []struct {
		title          string
		inputRequest   string
		want           string
		mock           mockCall
		expectedStatus int
	}{
		{
			title:        "get vote result and 200 response",
			inputRequest: `{"vote":"Best pokemon"}`,
			mock: func() {
				choices := []entity.Choice{{Title: "Pikachu", VoteId: 1, Count: 1}, {Title: "Noone", VoteId: 1, Count: 3}, {Title: "Mew", VoteId: 1, Count: 2}}
				choiceServ.EXPECT().Get(gomock.Any(), gomock.Any()).Return(choices, nil)

			},
			want:           "[{\"choice\": \"Pikachu\",\"vote_count\": 1},{\"choice\": \"Noone\",\"vote_count\": 3},{\"choice\": \"Mew\",\"vote_count\": 2}]",
			expectedStatus: 200,
		},
		{
			title:        "title not found and 400 response",
			inputRequest: `{"vote":"wrong title"}`,
			mock: func() {
				choiceServ.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errs.ErrTitleNotExist)

			},
			want:           "the title doesn't exist",
			expectedStatus: 400,
		},
		{
			title:        "service internal error and 500 response",
			inputRequest: `{"vote":"Best pokemon"}`,
			mock: func() {
				choiceServ.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("internal service error"))

			},
			want:           "500 Internal Server Error",
			expectedStatus: 500,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			req := httptest.NewRequest(
				"POST",
				"/api/result",
				bytes.NewBufferString(test.inputRequest),
			)
			req.Header.Set("Content-type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, clearResponse(recorder.Body.String()), test.want)
			assert.Equal(t, test.expectedStatus, recorder.Code)

		})
	}
}

func TestUpdateChoiceHandler(t *testing.T) {
	router := mux.NewRouter()
	ctrl := gomock.NewController(t)
	choiceServ := mocks.NewMockChoiceService(ctrl)
	voteServ := mocks.NewMockVoteService(ctrl)
	handler := NewVoteHandler(logging.GetLogger("debug"), voteServ, choiceServ)
	handler.InitRoutes(router)
	type mockCall func()
	testCases := []struct {
		title          string
		inputRequest   string
		mock           mockCall
		expectedStatus int
	}{
		{
			title:        "success update and 204 response",
			inputRequest: `{"vote":"Best pokemon","choice":"Mew"}`,
			mock: func() {
				choiceServ.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			},
			expectedStatus: 204,
		},
		{
			title:        "vote title not found and 400 response",
			inputRequest: `{"vote":"Best pokemon","choice":"Mew"}`,
			mock: func() {
				choiceServ.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errs.ErrTitleNotExist)

			},
			expectedStatus: 400,
		},
		{
			title:        "choice title not found and 400 response",
			inputRequest: `{"vote":"Best pokemon","choice":"Mew"}`,
			mock: func() {
				choiceServ.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errs.ErrChoiceTitleNotExist)

			},
			expectedStatus: 400,
		},
		{
			title:        "internal service error and  500 response",
			inputRequest: `{"vote":"Best pokemon","choice":"Mew"}`,
			mock: func() {
				choiceServ.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("internal service error"))

			},
			expectedStatus: 500,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			req := httptest.NewRequest(
				"POST",
				"/api/choice",
				bytes.NewBufferString(test.inputRequest),
			)
			req.Header.Set("Content-type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.expectedStatus, recorder.Code)

		})
	}
}

func TestDeleteChoiceHandler(t *testing.T) {
	router := mux.NewRouter()
	ctrl := gomock.NewController(t)
	choiceServ := mocks.NewMockChoiceService(ctrl)
	voteServ := mocks.NewMockVoteService(ctrl)
	handler := NewVoteHandler(logging.GetLogger("debug"), voteServ, choiceServ)
	handler.InitRoutes(router)
	type mockCall func()
	testCases := []struct {
		title          string
		inputRequest   string
		mock           mockCall
		expectedStatus int
	}{
		{
			title:        "success update and 204 response",
			inputRequest: `{"vote":"Best pokemon","choice":"Mew"}`,
			mock: func() {
				voteServ.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

			},
			expectedStatus: 204,
		},
		{
			title:        "vote title is empty and 400 response",
			inputRequest: `{"vote":"Best pokemon","choice":"Mew"}`,
			mock: func() {
				voteServ.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errs.ErrEmptyVoteTitle)

			},
			expectedStatus: 400,
		},
		{
			title:        "choice title not found and 400 response",
			inputRequest: `{"vote":"Best pokemon","choice":"Mew"}`,
			mock: func() {
				voteServ.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errs.ErrTitleNotExist)

			},
			expectedStatus: 400,
		},
		{
			title:        "internal service error and  500 response",
			inputRequest: `{"vote":"Best pokemon","choice":"Mew"}`,
			mock: func() {
				voteServ.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("internal service error"))

			},
			expectedStatus: 500,
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			test.mock()
			req := httptest.NewRequest(
				"DELETE",
				"/api/vote/1",
				bytes.NewBufferString(test.inputRequest),
			)
			req.Header.Set("Content-type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			assert.Equal(t, test.expectedStatus, recorder.Code)

		})
	}
}

func clearResponse(s string) string {
	temp := strings.ReplaceAll(s, "   ", "")
	return strings.ReplaceAll(temp, "\n", "")
}

//			want: "{\n   \"vote\": \"Best pokemon\",\n   \"choices\": [\n      {\n         \"choice\": \"Pikachu\",\n         \"vote_count\": 0\n      },\n      {\n         \"choice\": \"Mew\",\n         \"vote_count\": 0\n      },\n      {\n         \"choice\": \"Noone\",\n         \"vote_count\": 0\n      }\n   ]\n}",
