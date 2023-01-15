package models

import "beta/internal/adapters/db/postgresql"

func ConvertVoteToDB(vote RequestVote) postgresql.RequestVote {

	return postgresql.RequestVote{
		VoteId:   vote.VoteId,
		VotingId: vote.VotingId,
		OptionId: vote.OptionId,
	}
}
