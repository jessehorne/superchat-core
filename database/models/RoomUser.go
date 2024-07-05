package models

type RoomUser struct {
	GivenFields

	RoomID string
	UserID string
	Muted  bool
}
