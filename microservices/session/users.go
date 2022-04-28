package repositories

type UsersRepository interface {
	AddUser(email string, password string) (id int, err error)
	Authorize(email string, password string) (id int, err error)
	ChangePassword(id int, newPassword string) (err error)
}
