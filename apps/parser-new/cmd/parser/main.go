package main

import (
	"github.com/satont/twir/apps/parser-new/internal/delivery/bus/commands"
	"github.com/satont/twir/apps/parser-new/internal/delivery/bus/variables"
	"github.com/satont/twir/apps/parser-new/internal/delivery/grpc"
	"github.com/satont/twir/apps/parser-new/internal/delivery/http"
	"github.com/satont/twir/apps/parser-new/internal/task-queue"
	"github.com/satont/twir/apps/parser-new/internal/task-queue/asynq"
	"github.com/twirapp/twir/libs/baseapp"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp(
			baseapp.Opts{
				AppName: "parser",
			},
		),
		fx.Provide(
			asynq.NewRedisClientOpts,
			fx.Annotate(
				asynq.NewTaskDistributor,
				fx.As(new(taskqueue.TaskDistributor)),
			),
			fx.Annotate(
				asynq.NewTaskProcessorFx,
				fx.As(new(taskqueue.TaskProcessorListener)),
			),
			grpc.NewServerFx,
			http.NewServerFx,
			commands.NewBusListenerFx,
			variables.NewBusListenerFx,
		),
		fx.Invoke(
			func(grpc.Server) {},
			func(http.Server) {},
			func(taskqueue.TaskProcessorListener) {},
			func(commands.BusListener) {},
			func(variables.BusListener) {},
		),
	).Run()
}
