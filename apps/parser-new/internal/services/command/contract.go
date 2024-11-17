package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/parser-new/internal/repositories/command"
	model "github.com/satont/twir/libs/gomodels"
)

type ServiceContract interface {
	DeleteByID(ctx context.Context, commandID uuid.UUID) error
}

type Service struct {
	commandRepo command.RepositoryContract
}

var _ ServiceContract = (*Service)(nil)

func NewService() Service {
	return Service{}
}

func (s *Service) DeleteByID(ctx context.Context, commandID uuid.UUID) error {
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

	return nil
}
