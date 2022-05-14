package models

import (
	"2022_1_OnlyGroup_back/microservices/profile/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const BirthdayTopLimit = 100
const BirthdayBottomLimit = 18

type Profile struct {
	FirstName string     `json:",omitempty" validate:"min=0,max=40,regexp=^[a-zA-Z]*$"`
	LastName  string     `json:",omitempty" validate:"min=0,max=40,regexp=^[a-zA-Z]*$"`
	Birthday  *time.Time `json:",omitempty" validate:"birthday"`
	City      string     `json:",omitempty" validate:"min=0,max=32,regexp=^[a-zA-Z]*$"`
	Interests []Interest `json:",omitempty"`
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

type Interest struct {
	Id    int    `json:",omitempty"`
	Title string `json:",omitempty"`
}

func ModelProfileToGRPC(model *Profile) *proto.Profile {
	var interests []*proto.Interest
	for _, val := range model.Interests {
		interests = append(interests, &proto.Interest{Id: int64(val.Id), Title: val.Title})
	}
	birthdayTime := timestamppb.New(*model.Birthday)
	return &proto.Profile{Firstname: model.FirstName,
		LastName:  model.LastName,
		Birthday:  birthdayTime,
		City:      model.City,
		Interests: interests,
		AboutUser: model.AboutUser,
		UserId:    int64(model.UserId),
		Gender:    int64(model.Gender),
		Height:    int64(model.Height),
		Age:       model.Age,
	}
}

func GRPCToModelProfile(grpcModel *proto.Profile) Profile {
	birthday := grpcModel.Birthday.AsTime()
	var interests []Interest
	for _, val := range grpcModel.Interests {
		interests = append(interests, Interest{Id: int(val.Id), Title: val.Title})
	}
	return Profile{FirstName: grpcModel.Firstname,
		LastName:  grpcModel.LastName,
		Birthday:  &birthday,
		City:      grpcModel.City,
		Interests: interests,
		AboutUser: grpcModel.AboutUser,
		UserId:    int(grpcModel.UserId),
		Gender:    int(grpcModel.Gender),
		Height:    int(grpcModel.Height),
		Age:       grpcModel.Age,
	}
}

func ModelShortProfileToGRPC(model *ShortProfile) *proto.ShortProfile {
	return &proto.ShortProfile{FirstName: model.FirstName,
		LastName: model.LastName,
		City:     model.City,
	}
}

func GRPCToModelShortProfile(grpcModel *proto.ShortProfile) ShortProfile {
	return ShortProfile{FirstName: grpcModel.FirstName,
		LastName: grpcModel.LastName,
		City:     grpcModel.City,
	}
}

func ModelVectorCandidateToGRPC(model *VectorCandidate) *proto.VectorCandidate {
	var candidates []int64
	for _, val := range model.Candidates {
		candidates = append(candidates, int64(val))
	}
	return &proto.VectorCandidate{Candidates: candidates}
}

func GRPCToModelVectorCandidate(grpcModel *proto.VectorCandidate) VectorCandidate {
	var candidates []int
	for _, val := range grpcModel.Candidates {
		candidates = append(candidates, int(val))
	}
	return VectorCandidate{Candidates: candidates}
}

func ModelInterestsToGRPC(model []Interest) *proto.Interests {
	var interests []*proto.Interest
	for _, val := range model {
		interests = append(interests, &proto.Interest{Id: int64(val.Id), Title: val.Title})
	}

	return &proto.Interests{Interest: interests}
}

func GRPCToModelInterests(grpcModel *proto.Interests) []Interest {
	var interests []Interest
	for _, val := range grpcModel.Interest {
		interests = append(interests, Interest{Id: int(val.Id), Title: val.Title})
	}
	return interests
}

func ModelUserIdToGRPC(userId int) *proto.UserID {
	return &proto.UserID{Id: int64(userId)}
}

func ModelStrToGRPC(str string) *proto.StrInterest {
	return &proto.StrInterest{StrInterest: str}
}
