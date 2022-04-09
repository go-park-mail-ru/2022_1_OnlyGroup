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
	Email        string `validate:"regexp=^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$"`
	Password     string `validate:"password"`
	CaptchaToken string
}

type UserAuthProfile struct {
	Email       string
	OldPassword string
	NewPassword string
}
