package parser

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser-new/internal/commands/model"
	"github.com/satont/twir/apps/parser-new/internal/entity"
	commandrepo "github.com/satont/twir/apps/parser-new/internal/repositories/command"
	dbmodel "github.com/satont/twir/libs/gomodels"
	genericcacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	CommandPrefix = '!'
)

type CommandParser struct {
	commandRepo  commandrepo.RepositoryContract
	commandCache *genericcacher.GenericCacher[[]dbmodel.ChannelsCommands]
}

type Params struct {
	fx.In

	CommandRepo  commandrepo.RepositoryContract
	CommandCache *genericcacher.GenericCacher[[]dbmodel.ChannelsCommands]
}

func NewCommandParser(params Params) CommandParser {
	return CommandParser{
		commandRepo:  params.CommandRepo,
		commandCache: params.CommandCache,
	}
}

// TODO: work with domain models instead of db models.

func (cp *CommandParser) Execute(
	ctx context.Context,
	input, broadcasterID string,
) (*entity.FullCommandResponse, error) {
	// We should cast input to the slice of runes because we want to get number
	// of characters in input independently of the length of each character.
	if len([]rune(input)) < 2 {
		return nil, nil
	}

	if input[0] != CommandPrefix {
		return nil, nil
	}

	candidate := input[1:]

	commands, err := cp.commandCache.Get(ctx, broadcasterID)
	if err != nil {
		return nil, fmt.Errorf(
			"get user commands from cache: %w", err,
		)
	}

	command, found := cp.Find(candidate, commands)
	if !found {
		return nil, nil
	}

	if pass, err := cp.filter(ctx, command); !pass {
		if err != nil {
			return nil, fmt.Errorf("filter command: %w", err)
		}

		return nil, nil
	}

	go func() {
		_ = cp.sendPostEvents()
	}()

	// ...

	return nil, nil
}

func (cp *CommandParser) filter(
	ctx context.Context,
	command dbmodel.ChannelsCommands,
) (bool, error) {
	if command.ExpiresAt.Valid && command.ExpiresType != nil && command.ExpiresAt.Time.Before(time.Now().UTC()) {
		switch *command.ExpiresType {
		case dbmodel.ChannelCommandExpiresTypeDisable:
			if command.Enabled {
				break
			}

			cp.commandRepo.UpdateByID()

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

			if err := cp.commandCache.Invalidate(ctx, data.BroadcasterUserId); err != nil {
				c.services.Logger.Sugar().Error(err)
				return nil, err
			}

		case dbmodel.ChannelCommandExpiresTypeDelete:
			if command.Default {
				break
			}
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

	if command.OnlineOnly {
		stream := &dbmodel.ChannelsStreams{}
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

	if len(command.EnabledCategories) != 0 {
		stream := &dbmodel.ChannelsStreams{}
		err = c.services.Gorm.
			WithContext(ctx).
			Where(`"userId" = ?`, data.BroadcasterUserId).
			Find(stream).Error
		if err != nil {
			return nil, err
		}

		if stream.ID != "" {
			if !lo.ContainsBy(
				command.EnabledCategories,
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

	if command.CooldownType == "PER_USER" && command.Cooldown.Int64 > 0 && shouldCheckCooldown {
		key := fmt.Sprintf("commands:%s:cooldowns:user:%s", command.ID, data.ChatterUserId)
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

	permitted := c.isUserHasPermissionToCommand(
		data.ChatterUserId,
		data.BroadcasterUserId,
		command,
		dbUser,
		userRoles,
		commandRoles,
	)

	return permitted, nil
}

func (cp *CommandParser) sendPostEvents() error {
	c.services.GrpcClients.Events.CommandUsed(
		context.Background(),
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
		context.Background(),
		&websockets.TriggerAlertRequest{
			ChannelId: data.BroadcasterUserId,
			AlertId:   alert.ID,
		},
	)

	return nil
}

func (cp *CommandParser) Find(
	input string,
	candidates []dbmodel.ChannelsCommands,
) (dbmodel.ChannelsCommands, bool) {
	fragments := strings.Split(
		strings.ToLower(
			strings.TrimSpace(input),
		), " ",
	)

	var elected *dbmodel.ChannelsCommands

	for index := range fragments {
		command := strings.Join(fragments[:index+1], " ")

		for _, candidate := range candidates {
			if candidate.Name == command {
				elected = &candidate
				break
			}

			if slices.Contains(candidate.Aliases, command) {
				elected = &candidate
				break
			}
		}
	}

	if elected == nil {
		return dbmodel.ChannelsCommands{}, false
	}

	// Sort responses in the order specified by the user.
	sort.Slice(
		elected.Responses, func(a, b int) bool {
			return elected.Responses[a].Order < elected.Responses[b].Order
		},
	)

	return *elected, false
}
