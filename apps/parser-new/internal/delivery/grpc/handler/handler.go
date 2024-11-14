package handler

import (
	"context"

	defaultcommands "github.com/satont/twir/apps/parser-new/internal/commands/defaults"
	defaultvariables "github.com/satont/twir/apps/parser-new/internal/variables/defaults"
	"github.com/twirapp/twir/libs/grpc/parser"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerHandler struct {
	parser.UnimplementedParserServer

	defaultCommands  defaultcommands.DefaultCommands
	defaultVariables defaultvariables.DefaultVariables
}

type Params struct {
	DefaultCommands  defaultcommands.DefaultCommands
	DefaultVariables defaultvariables.DefaultVariables
}

func NewServerHandler(params Params) *ServerHandler {
	return &ServerHandler{
		defaultCommands:  params.DefaultCommands,
		defaultVariables: params.DefaultVariables,
	}
}

func (s *ServerHandler) GetDefaultCommands(_ context.Context, _ *emptypb.Empty) (
	*parser.GetDefaultCommandsResponse,
	error,
) {
	defaultCommands := s.defaultCommands.DefaultCommands()

	return &parser.GetDefaultCommandsResponse{
		List: fromDefaultCommands(defaultCommands),
	}, nil
}

func (s *ServerHandler) GetDefaultVariables(_ context.Context, _ *emptypb.Empty) (
	*parser.GetVariablesResponse,
	error,
) {
	defaultVariables := s.defaultVariables.DefaultVariables()

	return &parser.GetVariablesResponse{
		List: fromDefaultVariables(defaultVariables),
	}, nil
}
