package util

import "testing"

func Test_Password_All(t *testing.T) {
	password := "my cabbages"

	salt, hash := ProcessPassword(password)

	shouldBeTrue := ComparePassword(password, salt, hash)

	if shouldBeTrue != true {
		t.Errorf("Your crypto library is trash. Valid passwords aren't valid??")
	}

	shouldBeFalse := ComparePassword("definitely not valid", salt, hash)
	if shouldBeFalse != false {
		t.Errorf("If this happens in production...people are going to be upset.")
	}
}

func Test_Token_All(t *testing.T) {
	token, tokenHash := CreateToken()

	shouldBeTrue := ValidateToken(token, tokenHash)
	if shouldBeTrue != true {
		t.Error("This shouldn't happen. The tokens isn't valid??")
	}

	shouldBeFalse := ValidateToken("invalid token", tokenHash)
	if shouldBeFalse != false {
		t.Error("This definitely shouldn't happen in production. This token is invalid but it says it isn't.")
	}
}
