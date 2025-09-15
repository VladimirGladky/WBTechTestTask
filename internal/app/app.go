package app

import (
	"WBTechTestTask/internal/config"
	kafka "WBTechTestTask/internal/kafka/consumer"
	"WBTechTestTask/internal/repository"
	"WBTechTestTask/internal/service"
	"WBTechTestTask/internal/transport"
	"WBTechTestTask/pkg/logger"
	"WBTechTestTask/pkg/postgres"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type App struct {
	SubscriptionServer *transport.OderServer
	cfg                *config.Config
	ctx                context.Context
	wg                 sync.WaitGroup
	cancel             context.CancelFunc
	consumer           *kafka.Consumer
}

func New(cfg *config.Config, ctx context.Context) *App {
	db, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		panic(err)
	}
	repo := repository.NewOrderRepository(ctx, db.Pool)
	srv := service.NewOrderService(repo, ctx)
	if err := srv.InitCache(ctx); err != nil {
		panic(err)
	}
	server := transport.New(cfg, ctx, srv)
	consumer := kafka.NewConsumer(cfg.Kafka.KafkaBrokers, cfg.Kafka.KafkaTopic, srv, ctx)
	return &App{
		SubscriptionServer: server,
		cfg:                cfg,
		ctx:                ctx,
		consumer:           consumer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 2)
	a.wg.Add(1)
	go func() {
		logger.GetLoggerFromCtx(a.ctx).Info("Server started on address", zap.Any("address", a.cfg.Host+":"+a.cfg.Port))
		defer a.wg.Done()
		if err := a.SubscriptionServer.Start(); err != nil {
			errCh <- err
		}
	}()

	a.wg.Add(1)
	go func() {
		logger.GetLoggerFromCtx(a.ctx).Info("Kafka consumer starting",
			zap.Strings("brokers", a.cfg.Kafka.KafkaBrokers),
			zap.String("topic", a.cfg.Kafka.KafkaTopic))
		defer a.wg.Done()
		if err := a.consumer.ConsumeMessages(); err != nil {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		logger.GetLoggerFromCtx(a.ctx).Error("error running app", zap.Error(err))
		a.Stop()
	case sig := <-sigCh:
		logger.GetLoggerFromCtx(a.ctx).Info("signal received", zap.Any("signal", sig))
		a.Stop()
	case <-a.ctx.Done():
		logger.GetLoggerFromCtx(a.ctx).Info("context done")
	}

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.GetLoggerFromCtx(a.ctx).Info("all components stopped gracefully")
	case <-time.After(35 * time.Second):
		logger.GetLoggerFromCtx(a.ctx).Warn("graceful shutdown timed out")
	}

	return nil
}

func (a *App) Stop() {
	logger.GetLoggerFromCtx(a.ctx).Info("Shutting down app")

	var errs []error

	if a.cancel != nil {
		a.cancel()
	}

	if a.SubscriptionServer != nil {
		if err := a.SubscriptionServer.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("error stopping server: %w", err))
		}
	}

	if a.consumer != nil {
		if err := a.consumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing consumer: %w", err))
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			logger.GetLoggerFromCtx(a.ctx).Error("error stopping app", zap.Error(err))
		}
	}
	logger.GetLoggerFromCtx(a.ctx).Info("shutdown completed")
}
