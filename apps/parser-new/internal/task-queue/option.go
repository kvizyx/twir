package taskqueue

import (
	"time"
)

type (
	processInOption time.Duration
)

type OptionType int

const (
	OptionTypeProcessIn OptionType = iota
)

type Option interface {
	Type() OptionType
	Value() any
}

func ProcessIn(processIn time.Duration) Option { return processInOption(processIn) }

func (opt processInOption) Type() OptionType { return OptionTypeProcessIn }
func (opt processInOption) Value() any       { return time.Duration(opt) }
