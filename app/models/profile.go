package models

import (
	"time"
)

const BirthdayTopLimit = 100
const BirthdayBottomLimit = 18
const InterestSize = 32

type Profile struct {
	FirstName string     `json:",omitempty" validate:"min=0,max=40,regexp=^[a-zA-Z]*$"`
	LastName  string     `json:",omitempty" validate:"min=0,max=40,regexp=^[a-zA-Z]*$"`
	Birthday  *time.Time `json:",omitempty" validate:"birthday"`
	City      string     `json:",omitempty" validate:"min=0,max=32,regexp=^[a-zA-Z]*$"`
	Interests []string   `json:",omitempty" validate:"interests"`
	AboutUser string     `json:",omitempty" validate:"min=0,max=256"`
	UserId    int        `validate:"min=0"`
	Gender    int        `validate:"min=0, max=1"`
	Height    int        `validate:"min=0, max=250"`
	Age       string     `json:",omitempty"`
}

type ShortProfile struct {
	FirstName string
	LastName  string
	City      string
}

type VectorCandidate struct {
	Candidates []int
}
