package model

import (
	"github.com/satont/twir/apps/parser-new/internal/commands/argument"
	"github.com/satont/twir/apps/parser-new/internal/entity"
)

type CommandWithArgs struct {
	Args    []argument.Argument
	Command entity.Command
}

type DefaultCommandSettings struct {
	Name                string
	Description         string
	Roles               []string
	Module              string
	IsReply             bool
	IsVisible           bool
	IsEnabled           bool
	IsKeepResponseOrder bool
	Aliases             []string
}
