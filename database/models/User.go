package models

type User struct {
	GivenFields

	Email        string
	Name         string
	Password     string
	PasswordSalt string
}
