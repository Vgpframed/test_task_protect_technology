package vote_service

import (
	"beta/internal/adapters/db/postgresql"
	cfg "beta/internal/config"
	"beta/internal/domain/models"
	"beta/internal/domain/usecase/gamma"
	"context"
	lg "gitlab.satel.eyevox.ru/satel_vks/jaeger_tracer/log"
	"time"
)

var GService *Service
var timer = time.Now().UTC()

type Service struct {
	interfaceDB DB
	Logger      lg.Factory
	Config      *cfg.Config
}

type DB interface {
	AddVote(ctx context.Context, vote postgresql.RequestVote) (err error)
	UpdateVote(ctx context.Context, vote postgresql.RequestVote) (err error)
	GetVote(ctx context.Context, vote postgresql.RequestVote) (Vote postgresql.RequestVote, err error)
	GetAllVotes(ctx context.Context) (votes []postgresql.SendVote, err error)
}

func NewServer(storage DB, logger lg.Factory, config cfg.Config) *Service {
	GService = &Service{
		interfaceDB: storage,
		Logger:      logger,
		Config:      &config,
	}
	Votes, err := GService.interfaceDB.GetAllVotes(context.Background())
	if err != nil {
		return &Service{}
	}
	GService.countPercentVotes(Votes)
	return GService
}

func (s *Service) AddVote(ctx context.Context, vote models.RequestVote) {

	convertedVote := models.ConvertVoteToDB(vote)
	Vote, _ := s.interfaceDB.GetVote(ctx, convertedVote)

	if Vote.VoteId == vote.VoteId {
		s.interfaceDB.UpdateVote(ctx, convertedVote)
	}

	err := GService.interfaceDB.AddVote(ctx, convertedVote)
	if err != nil {
		return
	}

	s.addVoteArr(vote)
}

type Percents struct {
	OptionId string
	Cnt      int64
	Percent  int64
}

var CountPer []Percents

func (s *Service) countPercentVotes(votes []postgresql.SendVote) (CountPer []Percents) {

	for _, vote := range votes {

		CountPer = append(CountPer, Percents{
			OptionId: vote.OptionId,
			Cnt:      vote.Cnt,
			Percent:  vote.Cnt / int64(len(votes)) * 100,
		})
	}
	return
}

func (s *Service) addVoteArr(vote models.RequestVote) {
	var finish []models.SendVote

	exist := false
	changedPercents := false

	for _, v := range CountPer {

		current := v.Percent

		if v.OptionId == vote.OptionId {
			exist = true
			v.Cnt += 1
			v.Percent = v.Cnt / int64(len(CountPer)) * 100
		}

		if current < v.Percent {
			changedPercents = true
		}

		finish = append(finish, models.SendVote{
			OptionId: v.OptionId,
			Count:    v.Cnt,
		})
	}

	if !exist {

		current := 1 / int64(len(CountPer)) * 100

		CountPer = append(CountPer, Percents{
			OptionId: vote.OptionId,
			Cnt:      1,
			Percent:  current,
		})

		if current > 0 {
			changedPercents = true
		}

		finish = append(finish, models.SendVote{
			OptionId: vote.OptionId,
			Count:    1,
		})
	}

	if changedPercents {
		go func() {

			Timer := time.Since(timer)
			for {
				if Timer.Milliseconds() > 500 {
					timer = time.Now().UTC()

					gamma.SendVotesGamma(finish, s.Config.BaseConfig.GammaEndpoint, s.Logger)
					break
				}
			}

		}()

	}
}
