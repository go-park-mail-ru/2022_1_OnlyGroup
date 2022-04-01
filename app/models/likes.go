package models

type Likes struct {
	Id     int
	Action int
}

type LikesMatched struct {
	VectorId []int
}
