package models

type Filters struct {
	AgeFilter    [2]int `json:",omitempty" validate:"ageFilter"`
	GenderFilter int    `json:",omitempty" validate:"min=0, max=1"`
	HeightFilter [2]int `json:",omitempty" validate:"heightFilter"`
}
