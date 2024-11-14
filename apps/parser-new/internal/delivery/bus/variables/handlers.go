package variables

import (
	"context"

	"github.com/twirapp/twir/libs/bus-core/parser"
)

func (bl *BusListener) GetBuiltInVariables(
	ctx context.Context,
	_ struct{},
) []parser.BuiltInVariable {
	// TODO: implement me
	return nil
}
