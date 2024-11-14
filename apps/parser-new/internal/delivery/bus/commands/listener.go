package commands

import (
	"fmt"

	"github.com/satont/twir/apps/parser-new/internal/commands/parser"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
)

const (
	queueGroup = "parser"
)

type BusListener struct {
	bus *buscore.Bus

	logger        logger.Logger
	commandParser parser.CommandParser
}

type Params struct {
	fx.In

	Bus           *buscore.Bus
	Logger        logger.Logger
	CommandParser parser.CommandParser
}

func NewBusListener(params Params) BusListener {
	return BusListener{
		bus:           params.Bus,
		logger:        params.Logger,
		commandParser: params.CommandParser,
	}
}

func (bl *BusListener) Subscribe() error {
	if err := bl.bus.Parser.GetCommandResponse.SubscribeGroup(
		queueGroup, bl.GetCommandResponse,
	); err != nil {
		return fmt.Errorf("subscribe on get command response: %w", err)
	}

	if err := bl.bus.Parser.ParseVariablesInText.SubscribeGroup(
		queueGroup, bl.ParseVariablesInText,
	); err != nil {
		return fmt.Errorf("subscribe on parse variables in text: %w", err)
	}

	if err := bl.bus.Parser.ProcessMessageAsCommand.SubscribeGroup(
		queueGroup, bl.ProcessMessageAsCommand,
	); err != nil {
		return fmt.Errorf("subscribe on process message as command: %w", err)
	}

	return nil
}

func (bl *BusListener) Unsubscribe() {
	bl.bus.Parser.GetCommandResponse.Unsubscribe()
	bl.bus.Parser.ParseVariablesInText.Unsubscribe()
	bl.bus.Parser.ProcessMessageAsCommand.Unsubscribe()
}
