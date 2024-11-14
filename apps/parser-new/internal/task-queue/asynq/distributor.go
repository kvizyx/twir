package asynq

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/satont/twir/apps/parser-new/internal/task-queue"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
)

type TaskDistributor struct {
	client *asynq.Client

	logger logger.Logger
}

var _ taskqueue.TaskDistributor = (*TaskDistributor)(nil)

type TaskDistributorParams struct {
	fx.In

	Logger     logger.Logger
	ClientOpts *asynq.RedisClientOpt
}

func NewTaskDistributor(params TaskDistributorParams) *TaskDistributor {
	client := asynq.NewClient(params.ClientOpts)

	return &TaskDistributor{
		client: client,
		logger: params.Logger,
	}
}

// fromOptions takes generic task queue options and returns equivalents for the asynq options.
func (td *TaskDistributor) fromOptions(opts ...taskqueue.Option) []asynq.Option {
	implOpts := make([]asynq.Option, len(opts))

	for index, opt := range opts {
		var implOpt asynq.Option

		switch opt.Type() {
		case taskqueue.OptionTypeProcessIn:
			value := opt.Value().(time.Duration)
			implOpt = asynq.ProcessIn(value)
		}

		implOpts[index] = implOpt
	}

	return implOpts
}
