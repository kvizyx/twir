package custom_var

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/eval"
)

var CustomVar = &types.Variable{
	Name:        "customvar",
	Description: lo.ToPtr("Custom variable"),
	Visible:     lo.ToPtr(false),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if variableData.Params == nil {
			return result, nil
		}

		v := &model.ChannelsCustomvars{}
		err := parseCtx.Services.Gorm.Where(`"name" = ?`, variableData.Params).WithContext(ctx).Find(v).Error
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		if v.ID == "" || (v.Response == "" && v.EvalValue == "") {
			return result, nil
		}

		if v.Type == model.CustomVarScript {
			req, err := parseCtx.Services.GrpcClients.Eval.Process(context.Background(), &eval.Evaluate{
				Script: v.EvalValue,
			})

			if err != nil {
				parseCtx.Services.Logger.Sugar().Error(err)

				return nil, errors.New(
					"cannot evaluate variable. This is internal error, please report this bug",
				)
			}

			result.Result = req.Result
		}

		if v.Type == model.CustomVarText || v.Type == model.CustomVarNumber {
			result.Result = v.Response
		}

		return result, nil
	},
}