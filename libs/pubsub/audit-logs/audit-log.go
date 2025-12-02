package auditlog

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/audit"
	auditlog "github.com/twirapp/twir/libs/bus-core/audit-logs"
)

type AuditLog struct {
	ID                uuid.UUID
	System            audit.System
	OperationType     auditlog.AuditOperationType
	OperationMetadata audit.OperationMetadata
	OldValue          null.String
	NewValue          null.String
	CreatedAt         time.Time
}
