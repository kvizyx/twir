package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/date"
)

var Uptime = &types.Variable{
	Name:                "stream.uptime",
	Description:         lo.ToPtr("Prints uptime of stream"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := parseCtx.Cacher.GetChannelStream(ctx)
		if stream == nil {
			result.Result = "offline"
			return &result, nil
		}

		result.Result = date.Duration(
			stream.StartedAt, &date.DurationOpts{
				UseUtc: true,
				Hide:   date.DurationOptsHide{},
			},
		)

		return &result, nil
	},
}
