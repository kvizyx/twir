package defaults

import (
	"github.com/satont/twir/apps/parser-new/internal/variables"
	"go.uber.org/fx"
)

type DefaultVariables struct {
	defaultVariables []variables.DefaultVariable
}

type Params struct {
	fx.In
}

func NewDefaultVariables(params Params) DefaultVariables {
	return DefaultVariables{}
}

// DefaultVariables is a getter for default variables.
func (dv *DefaultVariables) DefaultVariables() []variables.DefaultVariable {
	return dv.defaultVariables
}
