package seventv

import (
	"context"

	"github.com/satont/twir/apps/parser-new/internal/commands"
	"github.com/satont/twir/apps/parser-new/internal/commands/argument"
	"github.com/satont/twir/apps/parser-new/internal/commands/model"
)

type CommandEmote struct{}

var _ commands.DefaultCommand = (*CommandEmote)(nil)

func NewCommandEmote() *CommandEmote {
	return &CommandEmote{}
}

func (c *CommandEmote) Args() []argument.Argument {
	return []argument.Argument{
		argument.NewString("emote-name", false),
	}
}

func (c *CommandEmote) Settings() model.DefaultCommandSettings {
	return model.DefaultCommandSettings{
		Name:        "7tv emote",
		Description: "Search emote by name in current emote set",
		Module:      "7tv",
		IsReply:     true,
		IsVisible:   true,
		IsEnabled:   false,
	}
}

func (c *CommandEmote) Handle(
	ctx context.Context,
	state model.HandleState,
	argument argument.Provider,
) (*model.HandleResult, error) {
	// TODO: implement me
	return nil, nil
}
