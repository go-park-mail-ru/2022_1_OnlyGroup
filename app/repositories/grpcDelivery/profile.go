package grpcDelivery

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/microservices/profile/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileGrpc struct {
	delivery proto.ProfileRepositoryClient
}

func NewProfileGrpc(delivery proto.ProfileRepositoryClient) *ProfileGrpc {
	return &ProfileGrpc{delivery: delivery}
}

func GRPCErrToHttpErr(err error, errInf string) error {
	if status.Code(err) == codes.InvalidArgument {
		return http.ErrBadRequest
	}
	if status.Code(err) == codes.Unavailable {
		return http.ErrServiceUnavailable.Wrap(err, "connection error")
	}

	return http.ErrBaseApp.Wrap(err, "failed: "+errInf)
}

func (repo *ProfileGrpc) Get(profileId int) (models.Profile, error) {
	msg, err := repo.delivery.Get(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return models.Profile{}, GRPCErrToHttpErr(err, "get profile Grpc")
	}
	return models.GRPCToModelProfile(msg), nil
}

func (repo *ProfileGrpc) GetShort(profileId int) (models.ShortProfile, error) {
	msg, err := repo.delivery.GetShort(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return models.ShortProfile{}, GRPCErrToHttpErr(err, "get shortProfile Grpc")
	}
	return models.GRPCToModelShortProfile(msg), nil
}

func (repo *ProfileGrpc) Change(profileId int, profile models.Profile) error {
	_, err := repo.delivery.Change(context.Background(), models.ModelProfileToGRPC(&profile))
	if err != nil {
		return GRPCErrToHttpErr(err, "change profile Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) Delete(profileId int) error {
	_, err := repo.delivery.Delete(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return GRPCErrToHttpErr(err, "delete profile Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) Add(profile models.Profile) error {
	_, err := repo.delivery.Add(context.Background(), models.ModelProfileToGRPC(&profile))
	if err != nil {
		return GRPCErrToHttpErr(err, "add profile Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) AddEmpty(profileId int) error {
	_, err := repo.delivery.AddEmpty(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return GRPCErrToHttpErr(err, "add emptyProfile Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) FindCandidate(profileId int) (models.VectorCandidate, error) {
	msg, err := repo.delivery.FindCandidate(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return models.VectorCandidate{}, GRPCErrToHttpErr(err, "findCandidate Grpc")
	}
	return models.GRPCToModelVectorCandidate(msg), nil
}

func (repo *ProfileGrpc) CheckFiled(profileId int) error {
	_, err := repo.delivery.CheckFiled(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return GRPCErrToHttpErr(err, "checkFiled Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) GetFilters(userId int) (models.Filters, error) {
	filtersModel, err := repo.delivery.GetFilters(context.Background(), models.ModelUserIdToGRPC(userId))
	if err != nil {
		return models.Filters{}, GRPCErrToHttpErr(err, "getFilters Grpc")
	}
	return models.GRPCToModelFilters(filtersModel), nil
}

func (repo *ProfileGrpc) ChangeFilters(userId int, filters models.Filters) error {
	_, err := repo.delivery.ChangeFilters(context.Background(), models.ModelFiltersToGRPC(userId, &filters))
	if err != nil {
		return GRPCErrToHttpErr(err, "changeFilters Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) GetInterests() ([]models.Interest, error) {
	msg, err := repo.delivery.GetInterests(context.Background(), &proto.Nothing{})
	if err != nil {
		return nil, GRPCErrToHttpErr(err, "getInterests Grpc")
	}
	return models.GRPCToModelInterests(msg), nil
}

func (repo *ProfileGrpc) GetDynamicInterest(interest string) ([]models.Interest, error) {
	msg, err := repo.delivery.GetDynamicInterest(context.Background(), models.ModelStrToGRPC(interest))
	if err != nil {
		return nil, GRPCErrToHttpErr(err, "getDynamicInterest Grpc")
	}
	return models.GRPCToModelInterests(msg), nil
}

func (repo *ProfileGrpc) CheckInterests(interests []models.Interest) error {
	_, err := repo.delivery.CheckInterests(context.Background(), models.ModelInterestsToGRPC(interests))
	if err != nil {
		return GRPCErrToHttpErr(err, "checkInterests Grpc")
	}
	return nil
}

func (repo *ProfileGrpc) SetAction(profileId int, likes models.Likes) error {
	_, err := repo.delivery.SetAction(context.Background(), models.ModelLikesToGRPC(profileId, likes))
	if err != nil {
		return GRPCErrToHttpErr(err, "SetAction Grpc")
	}
	return err
}

func (repo *ProfileGrpc) GetMatched(profileId int) (models.LikesMatched, error) {
	msg, err := repo.delivery.GetMatched(context.Background(), models.ModelUserIdToGRPC(profileId))
	if err != nil {
		return models.LikesMatched{}, GRPCErrToHttpErr(err, "failed getMatched Grpc")
	}
	return models.GRPCToModelLikesMatched(msg), nil
}
