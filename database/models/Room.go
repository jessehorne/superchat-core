package models

type Room struct {
	GivenFields

	Name              string
	PasswordProtected bool
	Password          string
	PasswordSalt      string
}
