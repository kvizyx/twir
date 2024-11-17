package spam

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/commands"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"
)

const (
	spamCountArgName   = "count"
	spamMessageArgName = "message"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "spam",
		Description: null.StringFrom("Spam into chat. Example usage: <b>!spam 5 Follow me on twitter"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
	},
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name: spamCountArgName,
			Min:  lo.ToPtr(1),
			Max:  lo.ToPtr(10),
		},
		command_arguments.VariadicString{
			Name: spamMessageArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		count := parseCtx.ArgsParser.Get(spamCountArgName).Int()
		text := parseCtx.ArgsParser.Get(spamMessageArgName).String()

		cmd, isCmd := strings.CutPrefix(
			strings.TrimSpace(text),
			"!",
		)
		if !isCmd {
			return buildSpamResult(count, text), nil
		}

		cmd = strings.ToLower(cmd)

		cmdFragments := strings.Split(cmd, " ")
		if len(cmdFragments) == 0 {
			return buildSpamResult(count, text), nil
		}

		cmdName := cmdFragments[0]

		responses, err := getCommandResponses(ctx, cmdName, parseCtx)
		if err != nil {
			return buildSpamResult(count, text), nil
		}

		return buildSpamResult(count, responses...), nil
	},
}

func getCommandResponses(
	ctx context.Context,
	cmdName string,
	parseCtx *types.ParseContext,
) ([]string, error) {
	var command model.ChannelsCommands

	err := parseCtx.Services.Gorm.
		WithContext(ctx).
		Model(&model.ChannelsCommands{}).
		Joins("Responses").
		Where("channelId = ?", parseCtx.Channel.ID).
		Where("name = ?", cmdName).Or("? = ANY(aliases)", cmdName).
		Find(&command).Error
	if err != nil {
		parseCtx.Services.Logger.Error(
			"failed to get channel command",
			zap.Error(err),
		)

		return nil, fmt.Errorf("get channel command: %w", err)
	}

	if len(command.Responses) == 0 {
		return nil, errors.New("command has no responses")
	}

	var (
		message = parseCtx.Message.Message.Text
		cmdAt   = strings.Index(message, cmdName)
	)

	if cmdAt < 0 {
		return nil, errors.New("command not found in message")
	}

	cmdAt += len(cmdName) + 1

	responses := parseCtx.Services.CommandsParser.ParseCommandResponses(
		ctx,
		&commands.FindByMessageResult{
			Command: &command,
			FoundBy: message[1:cmdAt],
		},
		parseCtx.Message,
	)

	return responses.Responses, nil
}

func buildSpamResult(count int, responses ...string) *types.CommandsHandlerResult {
	if len(responses) == 0 || count == 0 {
		return nil
	}

	result := &types.CommandsHandlerResult{
		Result: make([]string, count*len(responses)),
	}

	var (
		offset         int
		responsesCount = len(responses)
	)

	for countIndex := 1; countIndex < count+1; countIndex++ {
		for responseIndex, response := range responses {
			result.Result[countIndex*responseIndex+offset] = response
		}

		offset += responsesCount
	}

	return result
}
