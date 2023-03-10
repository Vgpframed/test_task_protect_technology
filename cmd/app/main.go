package main

import (
	"beta/internal/adapters/db/postgresql"
	"beta/internal/config"
	http_v1 "beta/internal/controller/http/v1"
	lg "beta/internal/domain/logger"
	vote_service "beta/internal/domain/service"
	"beta/internal/domain/tracer"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	jexpvar "github.com/uber/jaeger-lib/metrics/expvar"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	serviceName = "beta"
)

func main() {

	metricsFactory := jexpvar.NewFactory(100)

	logger, _ := zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel), zap.AddCallerSkip(1))
	zapLogger := logger.With(zap.String("service", serviceName))
	newLogger := lg.NewFactory(zapLogger)
	cfg := config.GetConfig(newLogger)

	newLogger.Bg().Info("Server configuration", zap.Any("Config", cfg))

	tracer := tracing.Init(serviceName, metricsFactory, newLogger, cfg.BaseConfig.JaegerEndpoint)

	start(tracer, newLogger, *cfg)
}

func start(tracer opentracing.Tracer, logger lg.Factory, cfg config.Config) {

	pg, err := postgresql.NewVoteStorage(tracer, logger, &cfg)

	if err != nil {
		logger.Bg().Fatal("shutdown", zap.Error(err))
	}

	serv := vote_service.NewServer(pg, logger, cfg)
	if serv == nil {
		logger.Bg().Fatal("shutdown", zap.Error(errors.New("Nil deference pointer in the serv")))
	}
	logger.Bg().Info("Start setting http server")
	router := gin.Default()
	router.POST("/voting", http_v1.PostVote)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", cfg.BaseConfig.HttpPort),
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		logger.Bg().Info("Start listening gateway of server")

		if err := router.Run(server.Addr); err != http.ErrServerClosed {
			logger.Bg().Fatal("shutdown", zap.Error(err))
		}
	}()
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Bg().Info("Closing HTTP Server")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Bg().Fatal("shutdown", zap.Error(err))
		}
	}()
	wg.Wait()
	logger.Bg().Info("I'm leaving, bye!")
}
