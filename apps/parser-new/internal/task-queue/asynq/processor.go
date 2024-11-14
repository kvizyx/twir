package asynq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/satont/twir/apps/parser-new/internal/task-queue"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
)

const (
	queueDefault         = "default"
	queueDefaultPriority = 5
)

type TaskProcessor struct {
	server *asynq.Server

	config       config.Config
	logger       logger.Logger
	tokensClient tokens.TokensClient
}

var _ taskqueue.TaskProcessor = (*TaskProcessor)(nil)

type TaskProcessorParams struct {
	fx.In

	Config       config.Config
	Logger       logger.Logger
	ClientOpts   asynq.RedisClientOpt
	TokensClient tokens.TokensClient
}

func NewTaskProcessor(params TaskProcessorParams) *TaskProcessor {
	taskProcessor := &TaskProcessor{
		config:       params.Config,
		logger:       params.Logger,
		tokensClient: params.TokensClient,
	}

	server := asynq.NewServer(
		params.ClientOpts,
		asynq.Config{
			Queues: map[string]int{
				queueDefault: queueDefaultPriority,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(taskProcessor.handleError),
			LogLevel:     asynq.ErrorLevel,
		},
	)

	taskProcessor.server = server

	return taskProcessor
}

func (tp *TaskProcessor) Start() error {
	router := asynq.NewServeMux()

	router.HandleFunc(TaskModUser, tp.fromHandler(tp.ProcessDistributeModUser))

	tp.logger.Info("task processor is listening")

	if err := tp.server.Start(router); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	return nil
}

func (tp *TaskProcessor) Stop() error {
	tp.server.Stop()
	tp.server.Shutdown()

	return nil
}

func (tp *TaskProcessor) handleError(_ context.Context, task *asynq.Task, err error) {
	taskID := task.ResultWriter().TaskID()

	tp.logger.Error(
		"failed to process task from task queue",
		slog.String("task_id", taskID),
		slog.Any("error", err),
	)
}

func (tp *TaskProcessor) fromHandler(handler taskqueue.HandlerFunc) asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		return handler(
			ctx, NewTask(
				TaskParams{
					Type:    task.Type(),
					Payload: task.Payload(),
				},
			),
		)
	}
}
