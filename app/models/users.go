package models

const PasswordPatternLowerCase = `[a-z]+`
const PasswordPatternUpperCase = `[A-Z]+`
const PasswordPatternNumber = `[0-9]+`
const PasswordMinLength = 6
const PasswordMaxLength = 32

type UserID struct {
	ID int
}

type UserAuthInfo struct {
	Email        string `validate:"regexp=^[0-9a-zA-Z!#$%&'*+/=?^_{|}~-]+@[0-9a-zA-Z+/=?^_{|}~-]+(\\.[0-9a-zA-Z+/=?^_{|}~-]+)+$"`
	Password     string `validate:"password"`
	CaptchaToken string
}

type UserAuthProfile struct {
	Email       string
	OldPassword string
	NewPassword string
}
