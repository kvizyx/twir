package recordhandler

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/audit"
	busauditlog "github.com/twirapp/twir/libs/bus-core/audit-logs"
	auditlogs "github.com/twirapp/twir/libs/pubsub/audit-logs"
)

// PubSub is an [audit.Recorder] implementation with [auditlogs.PubSub].
type PubSub struct {
	pubSub auditlogs.PubSub
}

var _ audit.RecorderHandler = (*PubSub)(nil)

func NewPubSub(pubSub auditlogs.PubSub) PubSub {
	return PubSub{
		pubSub: pubSub,
	}
}

func (p PubSub) HandleCreateRecord(
	ctx context.Context,
	metadata audit.RecordMetadata,
	newValue json.RawMessage,
) error {
	return p.publishOperation(ctx, busauditlog.AuditOperationTypeCreate, metadata, newValue, nil)
}

func (p PubSub) HandleDeleteRecord(
	ctx context.Context,
	metadata audit.RecordMetadata,
	oldValue json.RawMessage,
) error {
	return p.publishOperation(ctx, busauditlog.AuditOperationTypeDelete, metadata, nil, oldValue)
}

func (p PubSub) HandleUpdateRecord(
	ctx context.Context,
	metadata audit.RecordMetadata,
	oldValue json.RawMessage,
	newValue json.RawMessage,
) error {
	return p.publishOperation(ctx, busauditlog.AuditOperationTypeUpdate, metadata, newValue, oldValue)
}

func (p PubSub) publishOperation(
	ctx context.Context,
	operationType busauditlog.AuditOperationType,
	metadata audit.RecordMetadata,
	newValue json.RawMessage,
	oldValue json.RawMessage,
) error {
	auditLog := auditlogs.AuditLog{
		ID:                uuid.New(),
		System:            metadata.System,
		OperationType:     operationType,
		OperationMetadata: metadata.OperationMetadata,
		CreatedAt:         time.Now(),
	}

	if newValue != nil {
		auditLog.NewValue = null.StringFrom(string(newValue))
	}

	if oldValue != nil {
		auditLog.OldValue = null.StringFrom(string(oldValue))
	}

	if err := p.pubSub.Publish(ctx, auditLog); err != nil {
		return fmt.Errorf("publish audit log to pubsub: %w", err)
	}

	return nil
}
