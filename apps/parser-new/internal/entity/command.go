package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type CommandID = uuid.UUID

type (
	CommandCooldownType int
	CommandExpireType   int
)

const (
	CmdCooldownTypeGlobal CommandCooldownType = iota
	CmdCooldownTypePerUser
)

const (
	CmdExpireTypeDisable CommandExpireType = iota
	CmdExpireTypeDelete
)

type Command struct {
	ID                        CommandID
	Name                      string
	Cooldown                  null.Int
	CooldownType              CommandCooldownType
	IsEnabled                 bool
	Description               null.String
	IsVisible                 bool
	IsDefault                 bool
	DefaultName               null.String
	ChannelID                 string
	Aliases                   []string
	Module                    string
	IsReply                   bool
	IsKeepResponseOrder       bool
	GroupID                   null.Value[uuid.UUID]
	RolesIDS                  []string
	DeniedUsersIDS            []string
	AllowedUsersIDS           []string
	IsOnlineOnly              bool
	RequiredWatchTime         int
	RequiredMessages          int
	RequiredUsedChannelPoints int
	CooldownRolesIDs          []string
	EnabledCategories         []string
	ExpiresAt                 null.Time
	ExpiresType               CommandExpireType
}

func (c *Command) IsExpired() bool {
	isBefore := c.ExpiresAt.Time.Before(time.Now())
	return c.ExpiresAt.Valid && !isBefore
}
