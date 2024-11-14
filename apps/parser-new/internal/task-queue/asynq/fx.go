package asynq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
)

func NewRedisClientOpts(config config.Config) (asynq.RedisClientOpt, error) {
	url, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		return asynq.RedisClientOpt{}, fmt.Errorf("parse redis uri: %w", err)
	}

	return asynq.RedisClientOpt{
		Addr:     url.Addr,
		Password: url.Password,
		DB:       url.DB,
		Username: url.Username,
		PoolSize: url.PoolSize,
	}, nil
}

func NewTaskProcessorFx(params TaskProcessorParams, lc fx.Lifecycle) *TaskProcessor {
	taskProcessor := NewTaskProcessor(params)

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := taskProcessor.Start(); err != nil {
						params.Logger.Error(
							"failed to start task processor",
							slog.Any("error", err),
						)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return taskProcessor.Stop()
			},
		},
	)

	return taskProcessor
}
