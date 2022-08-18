package handler

type VoteRequest struct {
	VoteTitle string   `json:"vote"`
	Choices   []string `json:"choices"`
}

type VoteInfoRequest struct {
	VoteTitle string `json:"vote"`
}

type ChoiceRequest struct {
	VoteTitle   string `json:"vote"`
	ChoiceTitle string `json:"choice"`
}

type VoteResponse struct {
	VoteTitle string           `json:"vote"`
	Choices   []ChoiceResponse `json:"choices"`
}

type ChoiceResponse struct {
	ChoiceTitle string `json:"choice"`
	Count       int    `json:"vote_count"`
}
