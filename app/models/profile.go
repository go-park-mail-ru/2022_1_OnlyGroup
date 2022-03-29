package models

type Profile struct {
	FirstName string   `db:"firstname"`
	LastName  string   `db:"lastname"`
	Birthday  string   `db:"birthday"`
	City      string   `db:"city"`
	Interests []string `db:"interests"`
	AboutUser string   `db:"aboutuser"`
	UserId    int      `db:"userid"`
	Gender    string   `db:"gender"`
}

type ShortProfile struct {
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	City      string `db:"city"`
}

type VectorCandidate struct {
	Candidates []int
}
