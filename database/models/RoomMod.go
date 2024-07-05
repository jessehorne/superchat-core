package models

const (
	RoomModRoleOwner = iota
	RoomModRoleMod
)

type RoomMod struct {
	GivenFields

	UserID string
	RoomID string
	Role   int
}
