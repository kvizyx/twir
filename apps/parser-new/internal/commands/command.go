package commands

import (
	"context"

	"github.com/satont/twir/apps/parser-new/internal/commands/argument"
	"github.com/satont/twir/apps/parser-new/internal/commands/model"
)

type DefaultCommand interface {
	// Args returns list of arguments which might be specified (if not optional) along
	// with the command invocation.
	Args() []argument.Argument

	// Handle handles command invocation with provided state and arguments.
	Handle(context.Context, model.HandleState, argument.Provider) (*model.HandleResult, error)

	// Settings returns default settings for default command.
	// They can be changed in the future by user.
	Settings() model.DefaultCommandSettings
}
