package models

// RequestVote -
type RequestVote struct {
	VoteId   string `json:"voteId"`
	VotingId string `json:"votingId"`
	OptionId string `json:"optionId"`
}

// Response -
type Response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}
