package models

type UserID struct {
	ID int
}

type UserAuthInfo struct {
	Email    string
	Password string
}

type UserAuthProfile struct {
	Email       string
	OldPassword string
	NewPassword string
}
