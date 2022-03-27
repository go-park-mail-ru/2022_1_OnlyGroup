package repositories

type SessionsRepository interface {
	AddSession(id int, additionalData string) (string, error)
	GetIdBySession(secret string) (int, string, error)
	RemoveSession(secret string) (err error)
}
