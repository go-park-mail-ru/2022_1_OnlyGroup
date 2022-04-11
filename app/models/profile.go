package models

const BirthdaySize = 10
const InterestSize = 32
const BirthdayRexexp = "d{2}.d{2}.d{4}"

type Profile struct {
	FirstName string   `validate:"min=0,max=40,regexp=^[a-zA-Z]*$"`
	LastName  string   `validate:"min=0,max=40,regexp=^[a-zA-Z]*$"`
	Birthday  string   `validate:"birthday"`
	City      string   `validate:"min=0,max=32,regexp=^[a-zA-Z]*$"`
	Interests []string `validate:"interests"`
	AboutUser string   `validate:"min=0,max=256"`
	UserId    int      `validate:"min=0"`
	Gender    int      `validate:"min=0, max=1"`
	Height    int      `validate:"min=0, max=250"`
}

type ShortProfile struct {
	FirstName string
	LastName  string
	City      string
}

type VectorCandidate struct {
	Candidates []int
}
