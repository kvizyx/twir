package grpc

import (
	"fmt"
	"log/slog"
	"net"

	defaultcommands "github.com/satont/twir/apps/parser-new/internal/commands/defaults"
	"github.com/satont/twir/apps/parser-new/internal/delivery/grpc/handler"
	defaultvariables "github.com/satont/twir/apps/parser-new/internal/variables/defaults"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/parser"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server

	logger           logger.Logger
	defaultCommands  defaultcommands.DefaultCommands
	defaultVariables defaultvariables.DefaultVariables
}

type Params struct {
	fx.In

	Logger           logger.Logger
	DefaultCommands  defaultcommands.DefaultCommands
	DefaultVariables defaultvariables.DefaultVariables
}

func NewServer(params Params) Server {
	return Server{
		logger:           params.Logger,
		defaultCommands:  params.DefaultCommands,
		defaultVariables: params.DefaultVariables,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%d", constants.ParserServerPort)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	s.server = grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
	)

	parser.RegisterParserServer(
		s.server, handler.NewServerHandler(
			handler.Params{
				DefaultCommands:  s.defaultCommands,
				DefaultVariables: s.defaultVariables,
			},
		),
	)

	s.logger.Info("grpc server is listening", slog.String("addr", addr))

	if err = s.server.Serve(listener); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
