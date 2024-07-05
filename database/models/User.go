package models

type User struct {
	GivenFields

	Email        string
	Password     string
	PasswordSalt string
}
