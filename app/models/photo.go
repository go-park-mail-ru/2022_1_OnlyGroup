package models

type PhotoID struct {
	ID int
}

type PhotoParams struct {
	LeftTop     [2]int
	RightBottom [2]int
}

type UserPhotos struct {
	Photos []int
}

type UserAvatar struct {
	Avatar int
	Params PhotoParams
}
