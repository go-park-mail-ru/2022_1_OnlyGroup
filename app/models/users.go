package models

type UserID struct {
	ID int
}

type UserAuthInfo struct {
	Email        string
	Password     string
	CaptchaToken string
}

type ErrorLogin struct {
	ErrorCode    int
	ErrorMessage string
	CaptchaToken string
}

type UserAuthProfile struct {
	Email       string
	OldPassword string
	NewPassword string
}
