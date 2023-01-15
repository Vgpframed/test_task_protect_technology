package vote_service

import (
	"beta/internal/adapters/db/postgresql"
	"beta/internal/domain/models"
	"context"
)

var GService *Service

type Service struct {
	interfaceDB DB
}

type DB interface {
	AddVote(ctx context.Context, vote postgresql.RequestVote) (err error)
	GetVote(ctx context.Context, vote postgresql.RequestVote) (err error)
}

func NewServer(storage DB) *Service {
	GService = &Service{
		interfaceDB: storage,
	}
	return GService
}

func AddVote(ctx context.Context, vote models.RequestVote) {
	_ = GService.interfaceDB.AddVote(ctx, models.ConvertVoteToDB(vote))
}
