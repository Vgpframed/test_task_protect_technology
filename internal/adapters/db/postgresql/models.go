package postgresql

// RequestVote - входная структура голосования
type RequestVote struct {
	VoteId    string
	VotingId  string
	OptionId  string
	CreatedAt int64
}

// getTableName - получение имени таблицы
func (r RequestVote) getTableName() string {
	return "vote.votes"
}

// SendVote - входная структура для сервиса gamma
type SendVote struct {
	OptionId string `json:"optionId"`
	Cnt      int64  `json:"count"`
}
