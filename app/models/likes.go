package models

type Likes struct {
	Id     int
	Action int `validate:"min=1,max=2"`
}

type LikesMatched struct {
	VectorId []int
}
