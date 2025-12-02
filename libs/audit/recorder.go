package audit

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
)

type OperationMetadata struct {
	ActorID   string
	ChannelID string
	ObjectID  string
}

type Recorder[Value any] struct {
	system  System
	handler RecorderHandler
}

func NewRecorder[Value any](system System, handler RecorderHandler) Recorder[Value] {
	return Recorder[Value]{
		system:  system,
		handler: handler,
	}
}

func (r Recorder[Value]) RecordCreateOperation(
	ctx context.Context,
	metadata OperationMetadata,
	newValue Value,
) error {
	rawNewValue, err := json.Marshal(newValue)
	if err != nil {
		return fmt.Errorf("marshal new value: %w", err)
	}

	return r.handler.HandleCreateRecord(ctx, r.newRecordMetadata(metadata), rawNewValue)
}

func (r Recorder[Value]) RecordDeleteOperation(
	ctx context.Context,
	metadata OperationMetadata,
	oldValue Value,
) error {
	rawOldValue, err := json.Marshal(oldValue)
	if err != nil {
		return fmt.Errorf("marshal old value: %w", err)
	}

	return r.handler.HandleDeleteRecord(ctx, r.newRecordMetadata(metadata), rawOldValue)
}

func (r Recorder[Value]) RecordUpdateOperation(
	ctx context.Context,
	metadata OperationMetadata,
	oldValue Value,
	newValue Value,
) error {
	rawOldValue, err := json.Marshal(oldValue)
	if err != nil {
		return fmt.Errorf("marshal old value: %w", err)
	}

	rawNewValue, err := json.Marshal(newValue)
	if err != nil {
		return fmt.Errorf("marshal new value: %w", err)
	}

	return r.handler.HandleUpdateRecord(ctx, r.newRecordMetadata(metadata), rawOldValue, rawNewValue)
}

func (r Recorder[Value]) newRecordMetadata(metadata OperationMetadata) RecordMetadata {
	return RecordMetadata{
		System:            r.system,
		OperationMetadata: metadata,
	}
}
