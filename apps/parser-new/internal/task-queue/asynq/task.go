package asynq

import (
	"github.com/satont/twir/apps/parser-new/internal/task-queue"
)

type Task struct {
	taskType string
	payload  []byte
}

var _ taskqueue.Task = (*Task)(nil)

type TaskParams struct {
	Type    string
	Payload []byte
}

func NewTask(params TaskParams) *Task {
	return &Task{
		taskType: params.Type,
		payload:  params.Payload,
	}
}

func (t *Task) Type() string {
	return t.taskType
}

func (t *Task) Payload() []byte {
	return t.payload
}
