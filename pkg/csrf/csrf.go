package csrf

type CsrfGenerator interface {
	Create(session string, id int, url string) (string, error)
	Check(session string, id int, url string, inputToken string) error
}
