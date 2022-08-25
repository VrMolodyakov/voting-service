package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VrMolodyakov/vote-service/internal/domain/entity"
	"github.com/VrMolodyakov/vote-service/internal/errs"
	"github.com/gorilla/mux"
)

const (
	prefix string = ""
	indent string = "   "
)

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("inside Create handler")
	var vote FullVoteRequest
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
	var response VoteResponse
	response.VoteTitle = vote.VoteTitle
	for _, choice := range choices {
		choice, err := h.choiceService.Create(ctx, choice)
		if err != nil {
			errorResponse(w, err)
			return
		}
		response.Choices = append(response.Choices, ChoiceResponse{ChoiceTitle: choice, Count: 0})

	}
	jsonReponce, err := json.MarshalIndent(response, prefix, indent)
	if err != nil {
		errorResponse(w, err)
		return
	}
	h.logger.Debug(vote)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonReponce)

}

func (h *handler) GetChoices(w http.ResponseWriter, r *http.Request) {
	var vote VoteTitleRequest
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		errorResponse(w, err)
		return
	}
	h.logger.Debugf("try to get choices for %v", vote)
	ctx := r.Context()
	choices, err := h.choiceService.Get(ctx, vote.VoteTitle)
	if err != nil {
		errorResponse(w, err)
		return
	}
	choiceDto := make([]ChoiceResponse, 0, len(choices))
	for _, choice := range choices {
		choiceDto = append(choiceDto, choiceToDto(choice))
	}
	jsonReponce, err := json.MarshalIndent(choiceDto, prefix, indent)
	if err != nil {
		errorResponse(w, err)
		return
	}
	h.logger.Debug(vote)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonReponce)
}

func (h *handler) UpdateChoice(w http.ResponseWriter, r *http.Request) {
	var updateReq UpdateChoiceRequest
	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		errorResponse(w, err)
		return
	}
	h.logger.Debugf("try tot update choice %v", updateReq)
	ctx := r.Context()
	err = h.choiceService.Update(ctx, updateReq.VoteTitle, updateReq.ChoiceTitle, 1)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	h.logger.Debugf("try tot delete vote %v", id)
	ctx := r.Context()
	err := h.voteService.Delete(ctx, id)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func errorResponse(w http.ResponseWriter, err error) {
	if errors.Is(err, errs.ErrEmptyChoiceTitle) ||
		errors.Is(err, errs.ErrEmptyVoteTitle) ||
		errors.Is(err, errs.ErrTitleNotExist) ||
		errors.Is(err, errs.ErrChoiceTitleNotExist) ||
		errors.Is(err, errs.ErrTitleAlreadyExist) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func dtoToChoices(voteId int, choices []string) []entity.Choice {
	chs := make([]entity.Choice, 0, len(choices))
	for _, choice := range choices {
		choice := entity.Choice{Title: choice, VoteId: voteId, Count: 0}
		chs = append(chs, choice)
	}
	return chs
}

func choiceToDto(choice entity.Choice) ChoiceResponse {
	return ChoiceResponse{ChoiceTitle: choice.Title, Count: choice.Count}
}
