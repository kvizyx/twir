package parser

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/cacher"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/commands"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/types/services"
	"github.com/satont/twir/apps/parser/internal/variables"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	splitNameRegexp = regexp.MustCompile(`[^\s]+`)
)

type Parser struct {
	services        *services.Services
	defaultCommands map[string]*types.DefaultCommand
	variables       *variables.Variables
}

type Params struct {
	Services        *services.Services
	DefaultCommands map[string]*types.DefaultCommand
	Variables       *variables.Variables
}

func NewParser(params Params) Parser {
	return Parser{
		services:        params.Services,
		defaultCommands: params.DefaultCommands,
		variables:       params.Variables,
	}
}

// ParseCommandResponses parses responses of the provided command with context of the chat message associated
// with this command and returns CommandResponses ready to be sent to the user's chat.
func (p *Parser) ParseCommandResponses(
	ctx context.Context,
	command *commands.FindByMessageResult,
	message twitch.TwitchChatMessage,
) *commands.CommandResponses {
	result := &commands.CommandResponses{
		KeepOrder: command.Command.KeepResponsesOrder,
		IsReply:   command.Command.IsReply,
	}

	var cmdParams *string

	params := strings.TrimSpace(message.Message.Text[1:][len(command.FoundBy):])
	// this shit comes from 7tv for bypass message duplicate
	params = strings.ReplaceAll(params, "\U000e0000", "")
	params = strings.TrimSpace(params)

	if len(params) > 0 {
		cmdParams = &params
	}

	var defaultCommand *types.DefaultCommand

	if command.Command.Default {
		if cmd, ok := p.defaultCommands[command.Command.DefaultName.String]; ok {
			defaultCommand = cmd
		}
	}

	go func() {
		err := p.services.Gorm.WithContext(context.Background()).Create(
			&model.ChannelsCommandsUsages{
				ID:        uuid.New().String(),
				UserID:    message.ChatterUserId,
				ChannelID: message.BroadcasterUserId,
				CommandID: command.Command.ID,
			},
		).Error
		if err != nil {
			p.services.Logger.Sugar().Error(
				"failed to create command usage",
				zap.Error(err),
			)
		}
	}()

	badges := make([]string, 0, len(message.Badges))
	for _, badge := range message.Badges {
		badges = append(badges, strings.ToUpper(badge.SetId))
	}

	mentions := make([]twitch.ChatMessageMessageFragmentMention, 0, len(message.Message.Fragments))
	for _, fragment := range message.Message.Fragments {
		if fragment.Type != twitch.FragmentType_MENTION {
			continue
		}

		mentions = append(mentions, *fragment.Mention)
	}

	emotes := make([]*types.ParseContextEmote, 0, len(message.Message.Fragments))
	for _, fragment := range message.Message.Fragments {
		if fragment.Type != twitch.FragmentType_EMOTE {
			continue
		}

		emotes = append(
			emotes, &types.ParseContextEmote{
				Name:  fragment.Text,
				ID:    fragment.Emote.Id,
				Count: 1,
				Positions: []*types.ParseContextEmotePosition{
					{
						Start: int64(fragment.Position.Start),
						End:   int64(fragment.Position.End),
					},
				},
			},
		)
	}

	var (
		parseCtxChannel = &types.ParseContextChannel{
			ID:   message.BroadcasterUserId,
			Name: message.BroadcasterUserLogin,
		}

		parseCtxSender = &types.ParseContextSender{
			ID:          message.ChatterUserId,
			Name:        message.ChatterUserLogin,
			DisplayName: message.ChatterUserName,
			Badges:      badges,
			Color:       message.Color,
		}

		parseCtx = &types.ParseContext{
			Message:   message,
			MessageId: message.MessageId,
			Channel:   parseCtxChannel,
			Sender:    parseCtxSender,
			Text:      cmdParams,
			RawText:   message.Message.Text[1:],
			IsCommand: true,
			Services:  p.services,
			Cacher: cacher.NewCacher(
				&cacher.CacherOpts{
					Services:        p.services,
					ParseCtxText:    cmdParams,
					ParseCtxChannel: parseCtxChannel,
					ParseCtxSender:  parseCtxSender,
				},
			),
			Emotes:   emotes,
			Mentions: mentions,
			Command:  command.Command,
		}
	)

	var channelStream model.ChannelsStreams

	if err := p.services.Gorm.WithContext(ctx).First(&channelStream).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			p.services.Logger.Sugar().Error(
				"failed to get channel stream",
				zap.String("channel_id", message.BroadcasterUserId),
				zap.Error(err),
			)

			return nil
		}
	}

	if command.Command.Default && defaultCommand != nil {
		argsParser, err := command_arguments.NewParser(defaultCommand.Args, params)
		if err != nil {
			if argsParser == nil {
				return nil
			}

			usage := argsParser.BuildUsageString(defaultCommand.Args, defaultCommand.Name)

			return &commands.CommandResponses{
				Responses: []string{fmt.Sprintf("[Usage]: %s", usage)},
				IsReply:   command.Command.IsReply,
			}
		}

		parseCtx.ArgsParser = argsParser

		results, err := defaultCommand.Handler(ctx, parseCtx)
		if err != nil {
			p.services.Logger.Sugar().Error(
				"error happened on default command execution",
				zap.Error(err),
				zap.Dict(
					"channel",
					zap.String("id", message.BroadcasterUserId),
					zap.String("name", message.BroadcasterUserLogin),
				),
				zap.Dict(
					"sender",
					zap.String("id", message.ChatterUserId),
					zap.String("name", message.ChatterUserLogin),
				),
				zap.String("message", message.Message.Text),
				zap.Dict(
					"command",
					zap.String("id", command.Command.ID),
					zap.String("name", command.Command.Name),
				),
			)

			var commandErr *types.CommandHandlerError

			if errors.As(err, &commandErr) {
				results = &types.CommandsHandlerResult{
					Result: []string{
						fmt.Sprintf("[Twir error]: %s", commandErr.Message),
					},
				}
			} else {
				results = &types.CommandsHandlerResult{
					Result: []string{"[Twir error]: unknown error happened. Please contact developers."},
				}
			}
		}

		result.Responses = lo.
			IfF(results == nil, func() []string { return []string{} }).
			ElseF(
				func() []string {
					return results.Result
				},
			)
	} else {
		responsesForCategory := make(
			[]model.ChannelsCommandsResponses,
			0,
			len(command.Command.Responses),
		)

		for _, response := range command.Command.Responses {
			if len(response.TwitchCategoryIDs) > 0 && channelStream.ID != "" {
				if !lo.ContainsBy(
					response.TwitchCategoryIDs,
					func(categoryId string) bool {
						return categoryId == channelStream.GameId
					},
				) {
					continue
				}
			}

			responsesForCategory = append(responsesForCategory, *response)
		}

		result.Responses = lo.Map(
			responsesForCategory,
			func(response model.ChannelsCommandsResponses, _ int) string {
				return response.Text.String
			},
		)
	}

	var wg sync.WaitGroup

	for index, response := range result.Responses {
		wg.Add(1)

		go func() {
			defer wg.Done()
			result.Responses[index] = p.variables.ParseVariablesInText(ctx, parseCtx, response)
		}()
	}

	wg.Wait()

	return result
}

func (p *Parser) FindChannelCommandInInput(
	input string,
	cmds []model.ChannelsCommands,
) *commands.FindByMessageResult {
	var command commands.FindByMessageResult

	var (
		msg       = strings.ToLower(input)
		splitName = splitNameRegexp.FindAllString(msg, -1)
		length    = len(splitName)
	)

	for i := 0; i < length; i++ {
		query := strings.Join(splitName, " ")

		for _, cmd := range cmds {
			if cmd.Name == query {
				command.FoundBy = query
				command.Command = &cmd
				break
			}

			if lo.SomeBy(
				cmd.Aliases, func(item string) bool {
					return item == query
				},
			) {
				command.FoundBy = query
				command.Command = &cmd
				break
			}
		}

		if command.Command != nil {
			break
		}

		splitName = splitName[:len(splitName)-1]
	}

	// Sort command responses in right order, which set from dashboard ui.
	if command.Command != nil {
		sort.Slice(
			command.Command.Responses, func(a, b int) bool {
				return command.Command.Responses[a].Order < command.Command.Responses[b].Order
			},
		)
	}

	return &command
}
