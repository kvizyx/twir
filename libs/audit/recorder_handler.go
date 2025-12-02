package audit

import (
	"context"

	"github.com/goccy/go-json"
)

type (
	RecordValue    = json.RawMessage
	RecordMetadata struct {
		System            System
		OperationMetadata OperationMetadata
	}
)

type RecorderHandler interface {
	HandleCreateRecord(ctx context.Context, metadata RecordMetadata, newValue RecordValue) error
	HandleDeleteRecord(ctx context.Context, metadata RecordMetadata, oldValue RecordValue) error
	HandleUpdateRecord(
		ctx context.Context,
		metadata RecordMetadata,
		oldValue RecordValue,
		newValue RecordValue,
	) error
}
