package defaults

import (
	"github.com/satont/twir/apps/parser-new/internal/commands"
	seventv "github.com/satont/twir/apps/parser-new/internal/commands/defaults/7tv"
	"go.uber.org/fx"
)

type DefaultCommands struct {
	defaultCommands []commands.DefaultCommand
}

type Params struct {
	fx.In
}

func NewDefaultCommands(params Params) DefaultCommands {
	defaultCommands := []commands.DefaultCommand{
		seventv.NewCommandEmote(),
	}

	return DefaultCommands{
		defaultCommands: defaultCommands,
	}
}

// DefaultCommands is a getter for default commands.
func (dc *DefaultCommands) DefaultCommands() []commands.DefaultCommand {
	return dc.defaultCommands
}
