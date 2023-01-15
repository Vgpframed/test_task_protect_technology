package postgresql

// RequestVote - входная структура голосования
type RequestVote struct {
	VoteId    string
	VotingId  string
	OptionId  string
	CreatedAt int64
}

func (r RequestVote) getTableName() string {
	return "vote.votes"
}
