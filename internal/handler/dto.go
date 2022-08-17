package handler

type VoteRequest struct {
	VoteTitle string   `json:"vote"`
	Choices   []string `json:"choices"`
}

type ChoiceRequest struct {
	ChoiceTitle string
}

type VoteResponse struct {
	VoteTitle string           `json:"vote"`
	Choices   []ChoiceResponse `json:"choices"`
}

type ChoiceResponse struct {
	ChoiceTitle string `json:"choice"`
	Count       int    `json:"vote_count"`
}
