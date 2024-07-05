package models

type RoomMessage struct {
	GivenFields

	RoomID  string
	UserID  string
	Message string
}
