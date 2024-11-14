package entity

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type CommandResponse struct {
	ID                string
	CommandID         uuid.UUID
	Text              null.String
	Order             int
	TwitchCategoryIDs []string
}

type FullCommandResponse struct {
	Responses []CommandResponse
	IsReply   bool
}
