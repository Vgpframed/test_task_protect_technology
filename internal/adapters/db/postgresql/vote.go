package postgresql

import (
	"beta/internal/config"
	lg "beta/internal/domain/logger"
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type VoteStorage struct {
	DBClient *gorm.DB
	Tracer   opentracing.Tracer
	Logger   lg.Factory
}

// NewVoteStorage - Инициализация клиента.
func NewVoteStorage(newTracer opentracing.Tracer, newLogger lg.Factory, cfg *config.Config) (client VoteStorage, err error) {

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable`,
		cfg.Db.Xhost,
		cfg.Db.Xuser,
		cfg.Db.Xpassword,
		cfg.Db.Xdbname,
		cfg.Db.Xport,
	)

	client.DBClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	client.Tracer = newTracer
	client.Logger = newLogger

	return
}

// AddVote -
func (v VoteStorage) AddVote(ctx context.Context, vote RequestVote) (err error) {
	tx := v.DBClient.WithContext(ctx)
	defer ctx.Done()

	vote.CreatedAt = time.Now().UTC().Unix()

	res := tx.Table(vote.getTableName()).Create(&vote)
	if res.Error != nil {
		err = res.Error
		v.Logger.For(ctx).Error("db request AddVote", zap.Error(res.Error))
		return
	}

	return
}

// GetVote -
func (v VoteStorage) GetVote(ctx context.Context, vote RequestVote) (newvote RequestVote, err error) {
	tx := v.DBClient.WithContext(ctx)
	defer ctx.Done()

	res := tx.Table(vote.getTableName()).Where("vote_id = ? and voting_id = ? ", vote.VoteId, vote.VotingId).Scan(&newvote)
	if res.Error != nil {
		err = res.Error
		v.Logger.For(ctx).Error("db request GetVote", zap.Error(res.Error))
		return
	}
	return
}

// UpdateVote -
func (v VoteStorage) UpdateVote(ctx context.Context, vote RequestVote) (err error) {
	tx := v.DBClient.WithContext(ctx)
	defer ctx.Done()

	res := tx.Table(vote.getTableName()).Where("vote_id = ? and voting_id = ?", vote.VoteId, vote.VotingId).Update("option_id", vote.OptionId)
	if res.Error != nil {
		err = res.Error
		v.Logger.For(ctx).Error("db request UpdateVote", zap.Error(res.Error))
		return
	}

	return
}

// GetAllVotes -
func (v VoteStorage) GetAllVotes(ctx context.Context) (votes []SendVote, err error) {
	tx := v.DBClient.WithContext(ctx)
	defer ctx.Done()

	var vote RequestVote

	res := tx.Table(vote.getTableName()).Select("option_id, count(vote_id) cnt").Group("option_id").Scan(&votes)
	if res.Error != nil {
		err = res.Error
		v.Logger.For(ctx).Error("db request UpdateVote", zap.Error(res.Error))
		return
	}

	return
}
