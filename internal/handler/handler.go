package handler

import (
	"github.com/VrMolodyakov/vote-service/pkg/logging"
	"github.com/gorilla/mux"
)

type handler struct {
	logger        *logging.Logger
	voteService   VoteService
	choiceService ChoiceService
}

func NewVoteHandler(logger *logging.Logger, voteService VoteService, choiceService ChoiceService) *handler {
	return &handler{logger: logger, voteService: voteService, choiceService: choiceService}
}

func (h *handler) InitRoutes(router *mux.Router) {
	router.HandleFunc("/api/vote", h.Create).Methods("POST")
	router.HandleFunc("/api/result", h.GetChoices).Methods("POST")
	router.HandleFunc("/api/choice", h.UpdateChoice).Methods("POST")
	router.HandleFunc("/api/vote/{id:[0-9]+}", h.DeleteVote).Methods("DELETE")
}
