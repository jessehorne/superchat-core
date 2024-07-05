package util

import "testing"

func Test_Password_All(t *testing.T) {
	password := "my cabbages"

	salt, hash := ProcessPassword(password)

	isValid := ComparePassword(password, salt, hash)

	if !isValid {
		t.Errorf("Your crypto library is trash. Valid passwords aren't valid??")
	}

	shouldBeInvalid := ComparePassword("definitely not valid", salt, hash)
	if shouldBeInvalid {
		t.Errorf("If this happens in production...people are going to be upset.")
	}
}
