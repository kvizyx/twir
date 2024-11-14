package variables

import (
	"fmt"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
)

const (
	queueGroup = "parser"
)

type BusListener struct {
	bus *buscore.Bus
}

type Params struct {
	fx.In

	Bus *buscore.Bus
}

func NewBusListener(params Params) BusListener {
	return BusListener{
		bus: params.Bus,
	}
}

func (bl *BusListener) Subscribe() error {
	if err := bl.bus.Parser.GetBuiltInVariables.SubscribeGroup(
		queueGroup, bl.GetBuiltInVariables,
	); err != nil {
		return fmt.Errorf("subscribe on get built-in variables: %w", err)
	}

	return nil
}

func (bl *BusListener) Unsubscribe() {
	bl.bus.Parser.GetBuiltInVariables.Unsubscribe()
}
