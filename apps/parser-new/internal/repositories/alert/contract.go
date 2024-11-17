package alert

import (
	"context"

	"github.com/satont/twir/apps/parser-new/internal/entity"
)

// RepositoryContract is a contract for alert repository.
type RepositoryContract interface {
	GetByCommandIDs(
		ctx context.Context,
		channelID string,
		commandIDs ...entity.CommandID,
	) (entity.Alert, error)
}
