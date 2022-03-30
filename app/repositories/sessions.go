package repositories

type SessionsRepository interface {
	Add(id int, additionalData string) (string, error)
	Get(secret string) (int, string, error)
	Remove(secret string) (err error)
}
