package handler

import (
	"github.com/satont/twir/apps/parser-new/internal/commands"
	"github.com/satont/twir/apps/parser-new/internal/variables"
	"github.com/twirapp/twir/libs/grpc/parser"
)

type (
	grpcDefaultCommand  = parser.GetDefaultCommandsResponse_DefaultCommand
	grpcDefaultCommands = []*grpcDefaultCommand
)

func fromDefaultCommands(commands []commands.DefaultCommand) grpcDefaultCommands {
	grpcCommands := make(grpcDefaultCommands, len(commands))

	for index, command := range commands {
		settings := command.Settings()

		grpcCommands[index] = &grpcDefaultCommand{
			Name:               settings.Name,
			Description:        settings.Description,
			Visible:            settings.IsVisible,
			RolesNames:         settings.Roles,
			Module:             settings.Module,
			IsReply:            settings.IsReply,
			KeepResponsesOrder: settings.IsKeepResponseOrder,
			Aliases:            settings.Aliases,
		}
	}

	return grpcCommands
}

type (
	grpcDefaultVariable  = parser.GetVariablesResponse_Variable
	grpcDefaultVariables = []*grpcDefaultVariable
)

func fromDefaultVariables(variables []variables.DefaultVariable) grpcDefaultVariables {
	grpcVariables := make(grpcDefaultVariables, len(variables))

	for index, variable := range variables {
		// TODO: implement me
		_ = variable
		grpcVariables[index] = nil
	}

	return grpcVariables
}
