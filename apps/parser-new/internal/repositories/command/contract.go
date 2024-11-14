package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/parser-new/internal/entity"
)

// RepositoryContract is a contract for command repository.
type RepositoryContract interface {
	GetResponsesByID(ctx context.Context, commandID uuid.UUID) ([]entity.CommandResponse, error)
}
