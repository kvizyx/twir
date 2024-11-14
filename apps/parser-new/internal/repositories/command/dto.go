package command

import (
	"github.com/guregu/null/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/satont/twir/apps/parser-new/internal/entity"
)

type ResponseDTO struct {
	ID                string               `db:"id"`
	CommandID         pgtype.UUID          `db:"commandId"`
	Text              null.String          `db:"text"`
	Order             int                  `db:"order"`
	TwitchCategoryIDs pgtype.Array[string] `db:"twitch_category_id"`
}

func fromCommandResponses(from []ResponseDTO) []entity.CommandResponse {
	to := make([]entity.CommandResponse, len(from))
	for index, response := range from {
		to[index] = entity.CommandResponse{
			ID:                response.ID,
			CommandID:         response.CommandID,
			Text:              response.Text,
			Order:             response.Order,
			TwitchCategoryIDs: response.TwitchCategoryIDs.Elements,
		}
	}

	return to
}
