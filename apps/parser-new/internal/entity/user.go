package entity

type UserRole int

const (
	RoleBroadcaster UserRole = iota
	RoleModerator
	RoleVIP
	RoleSubscriber
)
