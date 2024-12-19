package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	helix "github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
)

// RewardsRedemptionsHistory is the resolver for the rewardsRedemptionsHistory field.
func (r *queryResolver) RewardsRedemptionsHistory(ctx context.Context, opts gqlmodel.TwitchRedemptionsOpts) (*gqlmodel.TwitchRedemptionResponse, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	channelIdForRequest := dashboardId
	if opts.ByChannelID.IsSet() {
		channelIdForRequest = *opts.ByChannelID.Value()
	}

	page := 0
	perPage := 20
	if opts.Page.IsSet() {
		page = *opts.Page.Value()
	}
	if opts.PerPage.IsSet() {
		perPage = *opts.PerPage.Value()
	}

	if perPage > 500 {
		return nil, fmt.Errorf("perPage must be less than or equal to 500")
	}

	query := r.gorm.
		WithContext(ctx).
		Order("redeemed_at DESC")

	var foundTwitchChannels []helix.Channel
	if opts.UserSearch.IsSet() {
		channels, err := r.cachedTwitchClient.SearchChannels(ctx, *opts.UserSearch.Value())
		if err != nil {
			return nil, err
		}

		foundTwitchChannels = channels
	}
	if len(foundTwitchChannels) > 0 {
		var ids []string
		for _, user := range foundTwitchChannels {
			ids = append(ids, user.ID)
		}

		query = query.Where(`"channel_redemptions_history"."user_id" IN ?`, ids)
	}

	query = query.Where(`"channel_redemptions_history"."channel_id" = ?`, channelIdForRequest)

	if opts.RewardsIds.IsSet() && len(opts.RewardsIds.Value()) > 0 {
		query = query.Where(`"channel_redemptions_history"."reward_id" IN ?`, opts.RewardsIds.Value())
	}

	var entities []model.ChannelRedemption
	if err := query.
		Limit(perPage).
		Offset(page * perPage).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	if len(entities) == 0 {
		return &gqlmodel.TwitchRedemptionResponse{
			Redemptions: nil,
			Total:       0,
		}, nil
	}

	rewards, err := r.TwitchRewards(ctx, &channelIdForRequest)
	if err != nil {
		return nil, err
	}

	res := make([]gqlmodel.TwitchRedemption, 0, len(entities))
	for _, entity := range entities {
		reward := gqlmodel.TwitchReward{
			ID:        entity.RewardID.String(),
			Title:     entity.RewardTitle,
			Cost:      entity.RewardCost,
			ImageUrls: nil,
		}

		twitchReward, twitchRewardFound := lo.Find(
			rewards, func(r gqlmodel.TwitchReward) bool {
				return r.ID == entity.RewardID.String()
			},
		)
		if twitchRewardFound {
			reward.Title = twitchReward.Title
			reward.ImageUrls = twitchReward.ImageUrls
			reward.UsedTimes = twitchReward.UsedTimes
			reward.Enabled = twitchReward.Enabled
		}

		redemption := gqlmodel.TwitchRedemption{
			ID:         entity.ID.String(),
			ChannelID:  entity.ChannelID,
			RedeemedAt: entity.RedeemedAt,
			User:       &gqlmodel.TwirUserTwitchInfo{ID: entity.UserID},
			Reward:     &reward,
			Prompt:     entity.RewardPrompt.Ptr(),
		}

		res = append(res, redemption)
	}

	return &gqlmodel.TwitchRedemptionResponse{
		Redemptions: res,
		Total:       0,
	}, nil
}

// User is the resolver for the user field.
func (r *twitchRedemptionResolver) User(ctx context.Context, obj *gqlmodel.TwitchRedemption) (*gqlmodel.TwirUserTwitchInfo, error) {
	return data_loader.GetHelixUserById(ctx, obj.User.ID)
}

// TwitchRedemption returns graph.TwitchRedemptionResolver implementation.
func (r *Resolver) TwitchRedemption() graph.TwitchRedemptionResolver {
	return &twitchRedemptionResolver{r}
}

type twitchRedemptionResolver struct{ *Resolver }