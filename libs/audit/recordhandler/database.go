package recordhandler

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/audit"
	auditlogs "github.com/twirapp/twir/libs/repositories/audit_logs"
	"github.com/twirapp/twir/libs/repositories/audit_logs/model"
)

// Database is an [audit.RecorderHandler] implementation with [auditlogs.Repository].
type Database struct {
	repository auditlogs.Repository
}

var _ audit.RecorderHandler = (*Database)(nil)

func NewDatabase(repository auditlogs.Repository) Database {
	return Database{
		repository: repository,
	}
}

func (d Database) HandleRecord(
	ctx context.Context,
	system audit.System,
	operation audit.Operation,
) error {
	auditLog := auditlogs.CreateInput{
		// System:         metadata.System,
		OperationType: mapAuditOperationType(operation.Action),
		ObjectID:      operation.Metadata.ObjectID,
		ChannelID:     operation.Metadata.ChannelID,
		UserID:        operation.Metadata.ActorID,
	}

	if operation.NewValue != nil {
		newValueBytes, err := json.Marshal(operation.NewValue)
		if err != nil {
			return fmt.Errorf("marshal operation new value: %w", err)
		}

		auditLog.NewValue = lo.ToPtr(string(newValueBytes))
	}

	if operation.OldValue != nil {
		oldValueBytes, err := json.Marshal(operation.OldValue)
		if err != nil {
			return fmt.Errorf("marshal operation old value: %w", err)
		}

		auditLog.OldValue = lo.ToPtr(string(oldValueBytes))
	}

	if err := d.repository.Create(ctx, auditLog); err != nil {
		return fmt.Errorf("create audit log in database: %w", err)
	}

	return nil
}

func mapAuditOperationType(action audit.Action) model.AuditOperationType {
	switch action {
	case audit.ActionCreate:
		return model.AuditOperationCreate
	case audit.ActionDelete:
		return model.AuditOperationDelete
	case audit.ActionUpdate:
		return model.AuditOperationUpdate
	}

	return model.AuditOperationTypeUnknown
}
