package customvar

import (
	"errors"
	"fmt"
	"time"
	model "tsuwari/models"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/golang/protobuf/proto"
	eval "github.com/satont/tsuwari/libs/nats/eval"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "customvar",
	Description: lo.ToPtr("Custom variable"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if data.Params == nil {
			return result, nil
		}

		v := getVarByName(ctx, *data.Params)

		if v == nil {
			return result, nil
		}

		if v.Type == "SCRIPT" {
			bytes, _ := proto.Marshal(&eval.Evaluate{
				Script: v.EvalValue.String,
			})

			msg, err := ctx.Services.Nats.Request("eval", bytes, 3*time.Second)
			if err != nil {
				return nil, errors.New(
					"cannot evaluate variable. This is internal error, please report this bug",
				)
			}

			response := eval.EvaluateResult{}

			if err := proto.Unmarshal(msg.Data, &response); err != nil {
				return nil, errors.New(
					"cannot unwrap response. This is internal error, please report this bug",
				)
			}

			result.Result = response.Result
		} else {
			result.Result = v.Response.String
		}

		return result, nil
	},
}

type CustomVar struct {
	Type      *string `json:"type"`
	EvalValue *string `json:"evalValue"`
	Response  *string `json:"response"`
}

func getVarByName(
	ctx *variables_cache.VariablesCacheService,
	name string,
) *model.ChannelsCustomvars {
	variable := &model.ChannelsCustomvars{}
	err := ctx.Services.Db.Where(`"name" = ?`, name).First(variable).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return variable
}
