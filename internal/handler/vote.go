package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/internal/errs"
	"github.com/VrMolodyakov/vote-service/pkg/logging"
)

const (
	prefix string = ""
	indent string = "   "
)

type handler struct {
	logger        *logging.Logger
	voteService   VoteService
	choiceService ChoiceService
}

func NewVoteHandler(logger *logging.Logger, voteService VoteService, choiceService ChoiceService) *handler {
	return &handler{logger: logger, voteService: voteService, choiceService: choiceService}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var vote VoteRequest
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	id, err := h.voteService.Create(ctx, vote.VoteTitle)
	if err != nil {
		errorResponse(w, err)
		return
	}
	choices := dtoToChoices(id, vote.Choices)
	var respose VoteResponse
	respose.VoteTitle = vote.VoteTitle
	for _, choice := range choices {
		choice, err := h.choiceService.CreateChoice(ctx, choice)
		if err != nil {
			errorResponse(w, err)
			return
		}
		respose.Choices = append(respose.Choices, ChoiceResponse{ChoiceTitle: choice, Count: 0})

	}
	jsonReponce, err := json.MarshalIndent(respose, prefix, indent)
	if err != nil {
		errorResponse(w, err)
	}
	h.logger.Debug(vote)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonReponce)

}

func GetChoices(w http.ResponseWriter, r *http.Request) {

}

func errorResponse(w http.ResponseWriter, err error) {
	if errors.Is(err, errs.ErrEmptyChoiceTitle) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if errors.Is(err, errs.ErrEmptyVoteTitle) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func dtoToChoices(voteId int, choices []string) []entity.Choice {
	chs := make([]entity.Choice, 0, len(choices))
	for _, choice := range choices {
		choice := entity.Choice{Title: choice, VoteId: voteId, Count: 0}
		chs = append(chs, choice)
	}
	return chs
}

/*

choices := make([]ChoiceResponse, 0)
	for _, v := range vote.Choices {
		choices = append(choices, ChoiceResponse{v, 0})
	}
	resp := VoteResponse{vote.VoteTitle, choices}
	w.WriteHeader(http.StatusOK)
	jsonReponce, _ := json.MarshalIndent(resp, prefix, indent)
	w.Write(jsonReponce)
	h.logger.Debug(vote)

*/
