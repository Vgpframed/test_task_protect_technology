package models

import "beta/internal/adapters/db/postgresql"

func ConvertVoteToDB(vote RequestVote) postgresql.RequestVote {

	return postgresql.RequestVote{
		VoteId:   vote.VoteId,
		VotingId: vote.VotingId,
		OptionId: vote.OptionId,
	}
}

func ExistOptionVote(slice []string, element string) (exist bool) {
	exist = false
	for _, el := range slice {
		if el == element {
			exist = true
			break
		}
	}
	return
}
