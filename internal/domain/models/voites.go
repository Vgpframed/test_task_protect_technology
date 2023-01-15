package models

// RequestVote - входная структура голосования
type RequestVote struct {
	VoteId   string `json:"voteId"`
	VotingId string `json:"votingId"`
	OptionId string `json:"optionId"`
}

// Response - структура ответа на запрос
type Response struct {
	Result  string `json:"result" description:"результат выполнения операции, обязательное поле"`
	Message string `json:"message,omitempty" description:"подробное сообщение о результате, если он не положительный, не обязательное поле"`
}

// SendVote - входная структура для сервиса gamma
type SendVote struct {
	VoteId   string `json:"voteId"`
	VotingId string `json:"votingId"`
	OptionId string `json:"optionId"`
	Count    int64  `yaml:"count"`
}
