package taskqueue

import (
	"context"
)

type HandlerFunc func(context.Context, Task) error

type Task interface {
	Type() string
	Payload() []byte
}

type TaskDistributor interface {
	DistributeModUser(context.Context, *TaskModUserPayload, ...Option) error
}

type TaskProcessor interface {
	ProcessDistributeModUser(context.Context, Task) error
}

type TaskProcessorListener interface {
	TaskProcessor

	Start() error
	Stop() error
}
