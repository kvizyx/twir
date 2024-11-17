package commands

import (
	model "github.com/satont/twir/libs/gomodels"
)

type FindByMessageResult struct {
	Command *model.ChannelsCommands
	FoundBy string
}

type CommandResponses struct {
	Responses []string
	IsReply   bool
	KeepOrder bool
}
