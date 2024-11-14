package commands

import (
	"context"

	"go.uber.org/fx"
)

func NewBusListenerFx(params Params, lc fx.Lifecycle) BusListener {
	busListener := NewBusListener(params)

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return busListener.Subscribe()
			},
			OnStop: func(ctx context.Context) error {
				busListener.Unsubscribe()
				return nil
			},
		},
	)

	return busListener
}
