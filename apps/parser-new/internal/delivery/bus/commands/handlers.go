package commands

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func (bl *BusListener) GetCommandResponse(
	ctx context.Context,
	message twitch.TwitchChatMessage,
) parser.CommandParseResponse {
	// TODO: implement me
	return parser.CommandParseResponse{}
}

func (bl *BusListener) ParseVariablesInText(
	ctx context.Context,
	request parser.ParseVariablesInTextRequest,
) parser.ParseVariablesInTextResponse {
	// TODO: implement me
	return parser.ParseVariablesInTextResponse{}
}

func (bl *BusListener) ProcessMessageAsCommand(
	ctx context.Context,
	message twitch.TwitchChatMessage,
) struct{} {
	responses, err := bl.commandParser.Execute(
		ctx,
		message.Message.Text,
		message.BroadcasterUserId,
	)

	if responses == nil {
		if err != nil {
			bl.logger.Error("failed to execute parser", slog.Any("error", err))
		}

		return struct{}{}
	}

	var replyTo string

	if responses.IsReply {
		replyTo = message.MessageId
	}

	for _, response := range responses.Responses {
		request := bots.SendMessageRequest{
			ChannelId:   message.BroadcasterUserId,
			ChannelName: &message.BroadcasterUserLogin,
			Message:     response.Text.String,
			ReplyTo:     replyTo,
		}

		if err = bl.bus.Bots.SendMessage.Publish(request); err != nil {
			bl.logger.Error(
				"failed to publish send message request for bots",
				slog.Any("error", err),
			)
		}
	}

	return struct{}{}
}
