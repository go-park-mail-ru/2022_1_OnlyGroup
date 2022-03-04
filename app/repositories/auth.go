package repositories

type AuthRepository interface {
	AddUser(email string, password string) (id int, err error)
	Authorize(email string, password string) (id int, err error)
	ChangePassword(id int, newPassword string) (err error)

	AddSession(id int) (secret string, err error)
	GetIdBySession(secret string) (id int, err error)
	RemoveSession(secret string) (err error)
}
