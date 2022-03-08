package models

type Profile struct {
	FirstName string
	LastName  string
	Birthday  string
	City      string
	Interests []string
	AboutUser string
	UserId    int
	Gender    string
}

type ShortProfile struct {
	FirstName string
	LastName  string
	City      string
}

type VectorCandidate struct {
	Candidates []int
}
