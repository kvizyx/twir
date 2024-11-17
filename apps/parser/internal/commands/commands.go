package commands

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/types/services"
	"github.com/satont/twir/apps/parser/internal/variables"
	model "github.com/satont/twir/libs/gomodels"
	busparser "github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/zap"
)

type Commands struct {
	DefaultCommands  map[string]*types.DefaultCommand
	services         *services.Services
	variablesService *variables.Variables
}

type Opts struct {
	Services         *services.Services
	VariablesService *variables.Variables
}

func New(opts *Opts) *Commands {
	return &Commands{
		DefaultCommands:  DefaultCommands(),
		services:         opts.Services,
		variablesService: opts.VariablesService,
	}
}

func (c *Commands) GetChannelCommands(
	ctx context.Context,
	channelId string,
) ([]model.ChannelsCommands, error) {
	return c.services.CommandsCache.Get(ctx, channelId)
}

var splittedNameRegexp = regexp.MustCompile(`[^\s]+`)

// FindChannelCommandInInput splits chat message by spaces, then read message from end to start
// and delete one word from end while message gets empty, or we found a command in message.
func (c *Commands) FindChannelCommandInInput(
	input string,
	cmds []model.ChannelsCommands,
) *FindByMessageResult {
	msg := strings.ToLower(input)
	splittedName := splittedNameRegexp.FindAllString(msg, -1)

	res := FindByMessageResult{}

	length := len(splittedName)

	for i := 0; i < length; i++ {
		query := strings.Join(splittedName, " ")
		for _, cmd := range cmds {
			if cmd.Name == query {
				res.FoundBy = query
				res.Command = &cmd
				break
			}

			if lo.SomeBy(
				cmd.Aliases, func(item string) bool {
					return item == query
				},
			) {
				res.FoundBy = query
				res.Command = &cmd
				break
			}
		}

		if res.Command != nil {
			break
		} else {
			splittedName = splittedName[:len(splittedName)-1]
			continue
		}
	}

	// sort command responses in right order, which set from dashboard ui
	if res.Command != nil {
		sort.Slice(
			res.Command.Responses, func(a, b int) bool {
				return res.Command.Responses[a].Order < res.Command.Responses[b].Order
			},
		)
	}

	return &res
}

func (c *Commands) ProcessChatMessage(ctx context.Context, data twitch.TwitchChatMessage) (
	*busparser.CommandParseResponse,
	error,
) {
	if data.Message.Text[0] != '!' {
		return nil, nil
	}

	cmds, err := c.GetChannelCommands(ctx, data.BroadcasterUserId)
	if err != nil {
		return nil, err
	}

	cmd := c.FindChannelCommandInInput(data.Message.Text[1:], cmds)
	if cmd.Command == nil {
		return nil, nil
	}

	if cmd.Command.ExpiresAt.Valid && cmd.Command.ExpiresType != nil && cmd.Command.ExpiresAt.Time.Before(time.Now().UTC()) {
		if *cmd.Command.ExpiresType == model.ChannelCommandExpiresTypeDisable && cmd.Command.Enabled {
			err = c.services.Gorm.
				WithContext(ctx).
				Where(`"id" = ?`, cmd.Command.ID).
				Model(&model.ChannelsCommands{}).
				Updates(
					map[string]interface{}{
						"enabled": false,
					},
				).Error
			if err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}

			if err := c.services.CommandsCache.Invalidate(ctx, data.BroadcasterUserId); err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}
		} else if *cmd.Command.ExpiresType == model.ChannelCommandExpiresTypeDelete && !cmd.Command.Default {
			err = c.services.Gorm.
				WithContext(ctx).
				Where(`"id" = ?`, cmd.Command.ID).
				Delete(&model.ChannelsCommands{}).Error
			if err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}

			if err := c.services.CommandsCache.Invalidate(ctx, data.BroadcasterUserId); err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}
		}

		return nil, nil
	}

	if cmd.Command.OnlineOnly {
		stream := &model.ChannelsStreams{}
		err = c.services.Gorm.
			WithContext(ctx).
			Where(`"userId" = ?`, data.BroadcasterUserId).
			Find(stream).Error
		if err != nil {
			return nil, err
		}
		if stream == nil || stream.ID == "" {
			return nil, nil
		}
	}

	if len(cmd.Command.EnabledCategories) != 0 {
		stream := &model.ChannelsStreams{}
		err = c.services.Gorm.
			WithContext(ctx).
			Where(`"userId" = ?`, data.BroadcasterUserId).
			Find(stream).Error
		if err != nil {
			return nil, err
		}

		if stream.ID != "" {
			if !lo.ContainsBy(
				cmd.Command.EnabledCategories,
				func(category string) bool {
					return category == stream.GameId
				},
			) {
				return nil, nil
			}
		}
	}

	convertedBadges := make([]string, 0, len(data.Badges))
	for _, badge := range data.Badges {
		convertedBadges = append(convertedBadges, strings.ToUpper(badge.SetId))
	}

	dbUser, _, userRoles, commandRoles, err := c.prepareCooldownAndPermissionsCheck(
		ctx,
		data.ChatterUserId,
		data.BroadcasterUserId,
		convertedBadges,
		cmd.Command,
	)
	if err != nil {
		return nil, err
	}

	shouldCheckCooldown := c.shouldCheckCooldown(convertedBadges, cmd.Command, userRoles)
	if cmd.Command.CooldownType == "GLOBAL" && cmd.Command.Cooldown.Int64 > 0 && shouldCheckCooldown {
		key := fmt.Sprintf("commands:%s:cooldowns:global", cmd.Command.ID)
		rErr := c.services.Redis.Get(ctx, key).Err()

		if errors.Is(rErr, redis.Nil) {
			c.services.Redis.Set(ctx, key, "", time.Duration(cmd.Command.Cooldown.Int64)*time.Second)
		} else if rErr != nil {
			c.services.Logger.Sugar().Error(rErr)
			return nil, errors.New("error while setting redis cooldown for command")
		} else {
			return nil, nil
		}
	}

	if cmd.Command.CooldownType == "PER_USER" && cmd.Command.Cooldown.Int64 > 0 && shouldCheckCooldown {
		key := fmt.Sprintf("commands:%s:cooldowns:user:%s", cmd.Command.ID, data.ChatterUserId)
		rErr := c.services.Redis.Get(ctx, key).Err()

		if rErr == redis.Nil {
			c.services.Redis.Set(ctx, key, "", time.Duration(cmd.Command.Cooldown.Int64)*time.Second)
		} else if rErr != nil {
			zap.S().Error(rErr)
			return nil, errors.New("error while setting redis cooldown for command")
		} else {
			return nil, nil
		}
	}

	hasPerm := c.isUserHasPermissionToCommand(
		data.ChatterUserId,
		data.BroadcasterUserId,
		cmd.Command,
		dbUser,
		userRoles,
		commandRoles,
	)

	if !hasPerm {
		return nil, nil
	}

	go func() {
		gCtx := context.Background()

		c.services.GrpcClients.Events.CommandUsed(
			// this should be background, because we don't want to wait for response
			gCtx,
			&events.CommandUsedMessage{
				BaseInfo:           &events.BaseInfo{ChannelId: data.BroadcasterUserId},
				CommandId:          cmd.Command.ID,
				CommandName:        cmd.Command.Name,
				CommandInput:       strings.TrimSpace(data.Message.Text[len(cmd.FoundBy)+1:]),
				UserName:           data.ChatterUserLogin,
				UserDisplayName:    data.ChatterUserName,
				UserId:             data.ChatterUserId,
				IsDefault:          cmd.Command.Default,
				DefaultCommandName: cmd.Command.DefaultName.String,
				MessageId:          data.MessageId,
			},
		)

		alert := model.ChannelAlert{}
		if err := c.services.Gorm.Where(
			"channel_id = ? AND command_ids && ?",
			data.BroadcasterUserId,
			pq.StringArray{cmd.Command.ID},
		).Find(&alert).Error; err != nil {
			zap.S().Error(err)
			return
		}

		if alert.ID == "" {
			return
		}
		c.services.GrpcClients.WebSockets.TriggerAlert(
			gCtx,
			&websockets.TriggerAlertRequest{
				ChannelId: data.BroadcasterUserId,
				AlertId:   alert.ID,
			},
		)
	}()

	// TODO: refactor parsectx to new chat message struct
	result := c.services.CommandsParser.ParseCommandResponses(ctx, cmd, data)

	return &busparser.CommandParseResponse{
		Responses: result.Responses,
		IsReply:   result.IsReply,
		KeepOrder: result.KeepOrder,
	}, nil
}
