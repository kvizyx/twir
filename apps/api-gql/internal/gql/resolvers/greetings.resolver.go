package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/satont/twir/libs/utils"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/mappers"
)

// TwitchProfile is the resolver for the twitchProfile field.
func (r *greetingResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.Greeting) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.UserID)
}

// GreetingsCreate is the resolver for the greetingsCreate field.
func (r *mutationResolver) GreetingsCreate(ctx context.Context, opts gqlmodel.GreetingsCreateInput) (*gqlmodel.Greeting, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	entity := &model.ChannelsGreetings{
		ID:        uuid.NewString(),
		ChannelID: dashboardId,
		UserID:    opts.UserID,
		Enabled:   opts.Enabled,
		Text:      opts.Text,
		IsReply:   opts.IsReply,
		Processed: false,
	}

	if err := r.gorm.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, fmt.Errorf("cannot create greeting: %w", err)
	}

	r.logger.Audit(
		"New greeting",
		audit.Fields{
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
			OperationType: audit.OperationCreate,
			ObjectID:      &entity.ID,
		},
	)

	return &gqlmodel.Greeting{
		ID:      entity.ID,
		UserID:  entity.UserID,
		Enabled: entity.Enabled,
		IsReply: entity.IsReply,
		Text:    entity.Text,
	}, nil
}

// GreetingsUpdate is the resolver for the greetingsUpdate field.
func (r *mutationResolver) GreetingsUpdate(ctx context.Context, id string, opts gqlmodel.GreetingsUpdateInput) (*gqlmodel.Greeting, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	entity := &model.ChannelsGreetings{}
	if err := r.gorm.WithContext(ctx).Where(
		`"channelId" = ? AND id = ?`,
		dashboardId,
		id,
	).First(entity).Error; err != nil {
		return nil, fmt.Errorf("cannot find greeting: %w", err)
	}

	var entityCopy model.ChannelsGreetings
	if err := utils.DeepCopy(entity, &entityCopy); err != nil {
		return nil, fmt.Errorf("cannot copy greeting: %w", err)
	}

	if opts.IsReply.IsSet() {
		entity.IsReply = *opts.IsReply.Value()
	}

	if opts.Enabled.IsSet() {
		entity.Enabled = *opts.Enabled.Value()
	}

	if opts.Text.IsSet() {
		entity.Text = *opts.Text.Value()
	}

	if opts.UserID.IsSet() {
		entity.UserID = *opts.UserID.Value()
	}

	if err := r.gorm.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, fmt.Errorf("cannot update greeting: %w", err)
	}

	r.logger.Audit(
		"Update greeting",
		audit.Fields{
			OldValue:      entityCopy,
			NewValue:      entity,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
			OperationType: audit.OperationUpdate,
			ObjectID:      &entity.ID,
		},
	)

	return &gqlmodel.Greeting{
		ID:      entity.ID,
		UserID:  entity.UserID,
		Enabled: entity.Enabled,
		IsReply: entity.IsReply,
		Text:    entity.Text,
	}, nil
}

// GreetingsRemove is the resolver for the greetingsRemove field.
func (r *mutationResolver) GreetingsRemove(ctx context.Context, id string) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	greeting := &model.ChannelsGreetings{}
	if err := r.gorm.WithContext(ctx).Where(
		`"channelId" = ? AND id = ?`,
		dashboardId,
		id,
	).
		First(greeting).Error; err != nil {
		return false, fmt.Errorf("cannot find greeting: %w", err)
	}

	if err := r.gorm.WithContext(ctx).Delete(greeting).Error; err != nil {
		return false, fmt.Errorf("cannot remove greeting: %w", err)
	}

	r.logger.Audit(
		"Remove greeting",
		audit.Fields{
			OldValue:      greeting,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardId),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelGreeting),
			OperationType: audit.OperationDelete,
			ObjectID:      &greeting.ID,
		},
	)

	return true, nil
}

// Greetings is the resolver for the greetings field.
func (r *queryResolver) Greetings(ctx context.Context) ([]gqlmodel.Greeting, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.ChannelsGreetings
	if err := r.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Order(`"userId" ASC`).
		Find(&entities).
		Error; err != nil {
		return nil, fmt.Errorf("cannot find greetings: %w", err)
	}

	var greetings []gqlmodel.Greeting
	for _, entity := range entities {
		greetings = append(
			greetings,
			gqlmodel.Greeting{
				ID:      entity.ID,
				UserID:  entity.UserID,
				Enabled: entity.Enabled,
				IsReply: entity.IsReply,
				Text:    entity.Text,
			},
		)
	}

	return greetings, nil
}

// Greeting returns graph.GreetingResolver implementation.
func (r *Resolver) Greeting() graph.GreetingResolver { return &greetingResolver{r} }

type greetingResolver struct{ *Resolver }
