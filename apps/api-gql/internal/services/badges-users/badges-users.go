package badges_users

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/badges_users"
	"github.com/twirapp/twir/libs/repositories/badges_users/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	BadgesUsersRepository badges_users.Repository
}

func New(opts Opts) *Service {
	return &Service{
		badgesUsersRepository: opts.BadgesUsersRepository,
	}
}

type Service struct {
	badgesUsersRepository badges_users.Repository
}

type GetManyInput struct {
	BadgeID uuid.UUID
}

func modelToEntity(b model.BadgeUser) entity.BadgeUser {
	return entity.BadgeUser{
		ID:        b.ID,
		BadgeID:   b.BadgeID,
		UserID:    b.UserID,
		CreatedAt: b.CreatedAt,
	}
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) ([]entity.BadgeUser, error) {
	selectedBadgeUsers, err := c.badgesUsersRepository.GetMany(
		ctx,
		badges_users.GetManyInput{
			BadgeID: input.BadgeID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]entity.BadgeUser, 0, len(selectedBadgeUsers))
	for _, b := range selectedBadgeUsers {
		result = append(result, modelToEntity(b))
	}

	return result, nil
}

type CreateInput struct {
	BadgeID uuid.UUID
	UserID  string
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.BadgeUser, error) {
	badgeUser, err := c.badgesUsersRepository.Create(
		ctx,
		badges_users.CreateInput{
			BadgeID: input.BadgeID,
			UserID:  input.UserID,
		},
	)
	if err != nil {
		return entity.BadgeUserNil, err
	}

	return modelToEntity(badgeUser), nil
}

type DeleteInput struct {
	BadgeID uuid.UUID
	UserID  string
}

func (c *Service) Delete(ctx context.Context, input DeleteInput) error {
	return c.badgesUsersRepository.Delete(
		ctx,
		badges_users.DeleteInput{
			BadgeID: input.BadgeID,
			UserID:  input.UserID,
		},
	)
}
