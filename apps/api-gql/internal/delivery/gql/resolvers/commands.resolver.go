package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
)

// Responses is the resolver for the responses field.
func (r *commandResolver) Responses(ctx context.Context, obj *gqlmodel.Command) ([]gqlmodel.CommandResponse, error) {
	if obj == nil || obj.Default {
		return []gqlmodel.CommandResponse{}, nil
	}

	parsedUuid, err := uuid.Parse(obj.ID)
	if err != nil {
		return nil, err
	}

	return data_loader.GetCommandResponsesById(ctx, parsedUuid)
}

// Group is the resolver for the group field.
func (r *commandResolver) Group(ctx context.Context, obj *gqlmodel.Command) (*gqlmodel.CommandGroup, error) {
	if obj == nil || obj.GroupID == nil {
		return nil, nil
	}

	parsedUuid, err := uuid.Parse(*obj.GroupID)
	if err != nil {
		return nil, err
	}

	group, err := data_loader.GetCommandGroupById(ctx, parsedUuid)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// TwitchCategories is the resolver for the twitchCategories field.
func (r *commandResponseResolver) TwitchCategories(ctx context.Context, obj *gqlmodel.CommandResponse) ([]gqlmodel.TwitchCategory, error) {
	categories, err := data_loader.GetTwitchCategoriesByIDs(ctx, obj.TwitchCategoriesIds)
	if err != nil {
		return nil, err
	}

	resultedCategories := make([]gqlmodel.TwitchCategory, 0, len(categories))
	for _, category := range categories {
		resultedCategories = append(resultedCategories, *category)
	}

	return resultedCategories, nil
}

// CommandsCreate is the resolver for the commandsCreate field
func (r *mutationResolver) CommandsCreate(ctx context.Context, opts gqlmodel.CommandsCreateOpts) (*gqlmodel.CommandCreatePayload, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	var groupId *uuid.UUID
	if opts.GroupID.IsSet() && opts.GroupID.Value() != nil {
		parsedGroupId, err := uuid.Parse(*opts.GroupID.Value())
		if err != nil {
			return nil, err
		}

		groupId = &parsedGroupId
	}

	var expiresType *string
	if opts.ExpiresType.IsSet() && opts.ExpiresType.Value() != nil {
		expiresType = lo.ToPtr(opts.ExpiresType.Value().String())
	}

	responses := make([]commands.CreateInputResponse, len(opts.Responses))
	for idx, res := range opts.Responses {
		responses[idx] = commands.CreateInputResponse{
			Text:              &res.Text,
			Order:             idx,
			TwitchCategoryIDs: res.TwitchCategoriesIds,
		}
	}

	createInput := commands.CreateInput{
		ChannelID:                 dashboardId,
		ActorID:                   user.ID,
		Name:                      opts.Name,
		Cooldown:                  opts.Cooldown,
		CooldownType:              opts.CooldownType,
		Enabled:                   opts.Enabled,
		Aliases:                   opts.Aliases,
		Description:               opts.Description,
		Visible:                   opts.Visible,
		IsReply:                   opts.IsReply,
		KeepResponsesOrder:        opts.KeepResponsesOrder,
		DeniedUsersIDS:            opts.DeniedUsersIds,
		AllowedUsersIDS:           opts.AllowedUsersIds,
		RolesIDS:                  opts.RolesIds,
		OnlineOnly:                opts.OnlineOnly,
		CooldownRolesIDs:          opts.CooldownRolesIds,
		EnabledCategories:         opts.EnabledCategories,
		RequiredWatchTime:         opts.RequiredWatchTime,
		RequiredMessages:          opts.RequiredMessages,
		RequiredUsedChannelPoints: opts.RequiredUsedChannelPoints,
		GroupID:                   groupId,
		ExpiresAt:                 opts.ExpiresAt.Value(),
		ExpiresType:               expiresType,
		Responses:                 responses,
	}

	newCmd, err := r.commandsService.Create(ctx, createInput)
	if err != nil {
		return nil, err
	}

	return &gqlmodel.CommandCreatePayload{ID: newCmd.ID.String()}, nil
}

// CommandsUpdate is the resolver for the commandsUpdate field.
func (r *mutationResolver) CommandsUpdate(ctx context.Context, id string, opts gqlmodel.CommandsUpdateOpts) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("wrong command id: %w", err)
	}

	updateInput := commands_with_groups_and_responses.UpdateInput{
		ActorID:                   user.ID,
		ChannelID:                 dashboardId,
		Name:                      opts.Name.Value(),
		Cooldown:                  opts.Cooldown.Value(),
		CooldownType:              opts.CooldownType.Value(),
		Enabled:                   opts.Enabled.Value(),
		Aliases:                   opts.Aliases.Value(),
		Description:               opts.Description.Value(),
		Visible:                   opts.Visible.Value(),
		IsReply:                   opts.IsReply.Value(),
		KeepResponsesOrder:        opts.KeepResponsesOrder.Value(),
		DeniedUsersIDS:            opts.DeniedUsersIds.Value(),
		AllowedUsersIDS:           opts.AllowedUsersIds.Value(),
		RolesIDS:                  opts.RolesIds.Value(),
		OnlineOnly:                opts.OnlineOnly.Value(),
		CooldownRolesIDs:          opts.CooldownRolesIds.Value(),
		EnabledCategories:         opts.EnabledCategories.Value(),
		RequiredWatchTime:         opts.RequiredWatchTime.Value(),
		RequiredMessages:          opts.RequiredMessages.Value(),
		RequiredUsedChannelPoints: opts.RequiredUsedChannelPoints.Value(),
		GroupID:                   nil,
		ExpiresAt:                 nil,
		ExpiresType:               nil,
		Responses:                 nil, // should be nil
	}

	if opts.GroupID.IsSet() {
		if opts.GroupID.Value() == nil {
			updateInput.GroupID = nil
		} else {
			parsedGroupId, err := uuid.Parse(*opts.GroupID.Value())
			if err != nil {
				return false, err
			}

			updateInput.GroupID = &parsedGroupId
		}
	}

	if opts.ExpiresAt.IsSet() {
		if opts.ExpiresAt.Value() == nil {
			updateInput.ExpiresAt = nil
		} else {
			updateInput.ExpiresAt = lo.ToPtr(time.UnixMilli(int64(*opts.ExpiresAt.Value())))

		}
	}

	if opts.ExpiresType.IsSet() {
		if opts.ExpiresType.Value() == nil {
			updateInput.ExpiresType = nil
		} else {
			updateInput.ExpiresType = lo.ToPtr(opts.ExpiresType.Value().String())
		}
	}

	for idx, res := range opts.Responses.Value() {
		updateInput.Responses = append(
			updateInput.Responses,
			commands_with_groups_and_responses.UpdateInputResponse{
				Text:              &res.Text,
				Order:             idx,
				TwitchCategoryIDs: res.TwitchCategoriesIds,
			},
		)
	}

	if _, err := r.commandsWithGroupsAndResponsesService.Update(
		ctx,
		parsedID,
		updateInput,
	); err != nil {
		return false, fmt.Errorf("cannot update command: %w", err)
	}

	return true, nil
}

// CommandsRemove is the resolver for the commandsRemove field.
func (r *mutationResolver) CommandsRemove(ctx context.Context, id string) (bool, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, err
	}

	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("wrong uuid: %w", err)
	}

	err = r.commandsService.Delete(
		ctx, commands.DeleteInput{
			ChannelID: dashboardId,
			ActorID:   user.ID,
			ID:        parsedID,
		},
	)
	if err != nil {
		return false, fmt.Errorf("cannot delete command: %w", err)
	}

	return true, nil
}

// Commands is the resolver for the commands field.
func (r *queryResolver) Commands(ctx context.Context) ([]gqlmodel.Command, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	cmds, err := r.commandsWithGroupsAndResponsesService.GetManyByChannelID(ctx, dashboardId)
	if err != nil {
		return nil, err
	}

	converted := make([]gqlmodel.Command, 0, len(cmds))
	for _, c := range cmds {
		command := mappers.CommandEntityTo(c.Command)
		converted = append(converted, command)
	}

	return converted, nil
}

// CommandsPublic is the resolver for the commandsPublic field.
func (r *queryResolver) CommandsPublic(ctx context.Context, channelID string) ([]gqlmodel.PublicCommand, error) {
	if channelID == "" {
		return nil, fmt.Errorf("channelID is required")
	}

	entities, err := r.commandsWithGroupsAndResponsesService.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	channelRoles, err := r.rolesService.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	filteredCommands := make([]entity.CommandWithGroupAndResponses, 0, len(entities))
	for _, cmd := range entities {
		if cmd.Command.Visible && cmd.Command.Enabled {
			filteredCommands = append(filteredCommands, cmd)
		}
	}

	convertedCommands := make([]gqlmodel.PublicCommand, 0, len(entities))
	for _, cmd := range entities {
		var description string
		if cmd.Command.Description != nil {
			description = *cmd.Command.Description
		}

		var cooldown int
		if cmd.Command.Cooldown != nil {
			cooldown = *cmd.Command.Cooldown
		}

		converted := gqlmodel.PublicCommand{
			Name:         cmd.Command.Name,
			Description:  description,
			Aliases:      cmd.Command.Aliases,
			Responses:    make([]string, 0, len(cmd.Responses)),
			Cooldown:     cooldown,
			CooldownType: cmd.Command.CooldownType,
			Module:       cmd.Command.Module,
			Permissions:  make([]gqlmodel.PublicCommandPermission, 0),
		}

		for _, response := range cmd.Responses {
			var text string
			if response.Text != nil {
				text = *response.Text
			}
			converted.Responses = append(converted.Responses, text)
		}

		if len(cmd.Command.RolesIDS) > 0 {
			for _, role := range channelRoles {
				if slices.Contains(cmd.Command.RolesIDS, role.ID) {
					continue
				}

				converted.Permissions = append(
					converted.Permissions,
					gqlmodel.PublicCommandPermission{
						Name: role.Name,
						Type: role.Type.String(),
					},
				)
			}
		}

		convertedCommands = append(convertedCommands, converted)
	}

	return convertedCommands, nil
}

// Command returns graph.CommandResolver implementation.
func (r *Resolver) Command() graph.CommandResolver { return &commandResolver{r} }

// CommandResponse returns graph.CommandResponseResolver implementation.
func (r *Resolver) CommandResponse() graph.CommandResponseResolver {
	return &commandResponseResolver{r}
}

type commandResolver struct{ *Resolver }
type commandResponseResolver struct{ *Resolver }