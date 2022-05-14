package models

import "2022_1_OnlyGroup_back/microservices/profile/proto"

type Likes struct {
	Id     int
	Action int `validate:"min=1,max=2"`
}

type LikesMatched struct {
	VectorId []int
}

func ModelLikesToGRPC(userId int, model Likes) *proto.Likes {
	return &proto.Likes{Action: int64(model.Action), WhoId: int64(userId), WhomId: int64(model.Id)}
}

func GRPCToModelLikes(grpcModel *proto.Likes) Likes {
	return Likes{Id: int(grpcModel.WhomId), Action: int(grpcModel.Action)}
}

func ModelLikesMatchedToGRPC(model LikesMatched) *proto.LikesMatched {
	var likesMatched []int64
	for _, val := range model.VectorId {
		likesMatched = append(likesMatched, int64(val))
	}
	return &proto.LikesMatched{VectorId: likesMatched}
}

func GRPCToModelLikesMatched(grpcModel *proto.LikesMatched) LikesMatched {
	var likesMatched []int
	for _, val := range grpcModel.VectorId {
		likesMatched = append(likesMatched, int(val))
	}
	return LikesMatched{VectorId: likesMatched}
}
